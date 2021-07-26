/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	workflowv1 "github.com/salaboy/knative-workflow/api/v1"
	eventingapi "knative.dev/eventing/pkg/apis/eventing/v1"
	knativeEventingClient "knative.dev/eventing/pkg/client/clientset/versioned"
	servingapi "knative.dev/serving/pkg/apis/serving/v1"
	knativeServingClient "knative.dev/serving/pkg/client/clientset/versioned"
)

var RUNNER_IMAGE = os.Getenv("RUNNER_IMAGE")

// WorkflowRunReconciler reconciles a WorkflowRun object
type WorkflowRunReconciler struct {
	client.Client
	knativeServingClient  *knativeServingClient.Clientset
	knativeEventingClient *knativeEventingClient.Clientset
	Scheme                *runtime.Scheme
}

type WorkflowRunCreatedResponse struct {
	Id string `json:"id"`
}

//+kubebuilder:rbac:groups=workflow.knative.dev,resources=workflowruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workflow.knative.dev,resources=workflowruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=workflow.knative.dev,resources=workflowruns/finalizers,verbs=update

// +kubebuilder:rbac:groups=serving.knative.dev,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=eventing.knative.dev,resources=triggers,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WorkflowRun object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *WorkflowRunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx).WithValues(">>> Reconcile: workflowrun", req.NamespacedName)
	var workflowRun workflowv1.WorkflowRun

	if err := r.Get(ctx, req.NamespacedName, &workflowRun); err != nil {
		// it might be not found if this is a delete request
		if ignoreNotFound(err) == nil {
			log.Info("Hey there.. deleting workflowrun happened: " + req.NamespacedName.Name)

			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch workflowrun")

		return ctrl.Result{}, err
	}

	if workflowRun.Spec.WorkflowRef != "" {
		var workflow workflowv1.Workflow
		if err := r.Get(ctx, types.NamespacedName{
			Namespace: "default",
			Name:      workflowRun.Spec.WorkflowRef,
		}, &workflow); err != nil {
			// it might be not found if this is a delete request

			return ctrl.Result{}, err
		}

		yamlStates, err := yaml.Marshal(workflow.Spec.WorkflowDefinition.WorkflowStates)
		if err != nil {
			log.Error(err, "failed to parse yaml from workflow definition states")
			return ctrl.Result{}, err
		}
		if RUNNER_IMAGE == "" {
			RUNNER_IMAGE = "kind.local/knative-workflow-runner-ddfac3ccbf87482f858add851df61835:5a7b7aa766d0e97c76431442d225d28fe72908b69f2216fa49fecb46ab0c7b8b"
		}
		service := &servingapi.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kservice-" + workflow.Name,
				Namespace: "default",
			},
			Spec: servingapi.ServiceSpec{
				ConfigurationSpec: servingapi.ConfigurationSpec{

					Template: servingapi.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								"autoscaling.knative.dev/minScale": "1",
							},
						},
						Spec: servingapi.RevisionSpec{
							PodSpec: v1.PodSpec{

								Containers: []v1.Container{
									v1.Container{
										Name:  "knative-workflow-runner",
										Image: RUNNER_IMAGE,
										Env: []v1.EnvVar{
											v1.EnvVar{
												Name:  "WORKFLOW",
												Value: fmt.Sprintf("%s", yamlStates),
											},
											v1.EnvVar{
												Name:  "EVENT_SINK",
												Value: workflowRun.Spec.Sink,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		serviceExist, err := r.knativeServingClient.ServingV1().Services("default").Get(ctx, service.Name, metav1.GetOptions{})
		if err != nil {
			if ignoreNotFound(err) == nil {
				log.Info("KService doesn't exist, so creating KService: " + service.Name)
				_, err := ctrl.CreateOrUpdate(ctx, r.Client, service, func() error {
					return ctrl.SetControllerReference(&workflowRun, service, r.Scheme)
				})
				if err != nil {
					log.Error(err, "Error Creating or Updating and Setting Controller References to Knative Service: "+service.Name)
				}
			} else {
				log.Error(err, "Error Fetching Knative Service: "+service.Name)
			}
		} else if serviceExist.Name != "" {
			log.Info("KService exist, so checking the Status URL: " + service.Name)
			if serviceExist.Status.URL != nil {

				log.Info("> Created KService URL for subscriber : " + serviceExist.Status.URL.String())
				parsedURL, err := apis.ParseURL("http://" + serviceExist.Name + ".default.svc.cluster.local" + "/workflows/events")
				if err != nil {
					log.Error(err, "Error Parsing URl for: "+serviceExist.Status.URL.String())
					return ctrl.Result{}, err
				}

				// Create Triggers for Workflow definition
				for stateType, _ := range workflow.Spec.WorkflowDefinition.WorkflowStates.States {
					log.Info("> Looking for Events in State : " + string(stateType))
					// Create triggers for events that the workflow is waiting for
					for eventName, _ := range workflow.Spec.WorkflowDefinition.WorkflowStates.States[stateType].Events {
						log.Info("> Creating trigger for Event: " + string(eventName) + " in State : " + string(stateType))
						trigger := &eventingapi.Trigger{
							ObjectMeta: metav1.ObjectMeta{
								Name:      strings.ToLower("t-" + workflow.Name + "-" + string(eventName)),
								Namespace: "default",
							},
							Spec: eventingapi.TriggerSpec{
								Broker: workflowRun.Spec.Broker,
								Filter: &eventingapi.TriggerFilter{
									Attributes: map[string]string{
										"type": string(eventName),
									},
								},
								Subscriber: duckv1.Destination{
									URI: parsedURL,
								},
							},
						}
						_, err := ctrl.CreateOrUpdate(ctx, r.Client, trigger, func() error {
							return ctrl.SetControllerReference(&workflowRun, trigger, r.Scheme)
						})
						if err != nil {
							log.Error(err, "Error Creating or Updating and Setting Controller References to Knative Trigger: "+trigger.Name)
						}

					}
				}

				for _, condition := range serviceExist.Status.Conditions {
					if condition.Type == apis.ConditionReady {
						if condition.Status == "True" && workflowRun.Status.WorkflowId == "" {
							// Create instance in runner
							var jsonStr = []byte(`{}`)

							newInstanceUrl, _ := apis.ParseURL("http://" + serviceExist.Name + ".default.127.0.0.1.nip.io" + "/workflows")

							resp, err := http.Post(newInstanceUrl.String(), "application/json", bytes.NewBuffer(jsonStr))

							//Handle Error
							if err != nil {
								log.Error(err, "Something failed sending a request to the runner")
							}
							log.Info("response Status:" + fmt.Sprintf("%v", resp.Status))
							log.Info("response Headers:" + fmt.Sprintf("%v", resp.Header))
							body, _ := ioutil.ReadAll(resp.Body)

							log.Info("response Body:" + string(body))

							var workflowRunCreatedResponse WorkflowRunCreatedResponse
							if err := json.Unmarshal(body, &workflowRunCreatedResponse); err != nil {
								log.Error(err, "Error Unmarshaling workflowRunCreatedResponse")
							}

							workflowRun.Status.WorkflowId = workflowRunCreatedResponse.Id
							workflowRun.Status.RunnerUrl = "http://" + serviceExist.Name + ".default.127.0.0.1.nip.io"

							if err := r.Status().Update(ctx, &workflowRun); err != nil {
								log.Error(err, "unable to update WorkflowRun status")
								return ctrl.Result{}, err
							}
							log.Info("> WorkflowRun Updated: " + workflowRun.Name + " Workflow Run: " + workflowRun.Status.WorkflowId)

						}
					}
				}

			} else {
				log.Info("KService exist, but Status URL is nil")
			}
		}

	}

	return ctrl.Result{}, nil
}

func ignoreNotFound(err error) error {
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkflowRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.knativeServingClient = knativeServingClient.NewForConfigOrDie(mgr.GetConfig())
	r.knativeEventingClient = knativeEventingClient.NewForConfigOrDie(mgr.GetConfig())
	return ctrl.NewControllerManagedBy(mgr).
		For(&workflowv1.WorkflowRun{}).
		Owns(&servingapi.Service{}).
		Owns(&eventingapi.Trigger{}).
		Watches(&source.Kind{Type: &servingapi.Service{}},
			&handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    &workflowv1.WorkflowRun{}}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(e event.DeleteEvent) bool {
				// The reconciler adds a finalizer so we perform clean-up
				// when the delete timestamp is added
				// Suppress Delete events to avoid filtering them out in the Reconcile function
				return false
			},
		}).
		Complete(r)
}

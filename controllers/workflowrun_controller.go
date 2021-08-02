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
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
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

	// I need to find the appropiate runner based on the workflowRef and then create an instance:
	// - If there is no runner for the def.. create a new WorkflowRunner first

	//if condition.Status == "True" && workflowRunner.Status.WorkflowId == "" {
	//	// Create instance in runner
	//	var jsonStr = []byte(`{}`)
	//
	//	newInstanceUrl, _ := apis.ParseURL("http://" + serviceExist.Name + ".default.127.0.0.1.nip.io" + "/workflows")
	//
	//	resp, err := http.Post(newInstanceUrl.String(), "application/json", bytes.NewBuffer(jsonStr))
	//
	//	//Handle Error
	//	if err != nil {
	//		log.Error(err, "Something failed sending a request to the runner")
	//	}
	//	log.Info("response Status:" + fmt.Sprintf("%v", resp.Status))
	//	log.Info("response Headers:" + fmt.Sprintf("%v", resp.Header))
	//	body, _ := ioutil.ReadAll(resp.Body)
	//
	//	log.Info("response Body:" + string(body))
	//
	//	var workflowRunCreatedResponse WorkflowRunCreatedResponse
	//	if err := json.Unmarshal(body, &workflowRunCreatedResponse); err != nil {
	//		log.Error(err, "Error Unmarshaling workflowRunCreatedResponse")
	//	}
	//
	//	workflowRun.Status.WorkflowId = workflowRunCreatedResponse.Id
	//	workflowRun.Status.RunnerUrl = "http://" + serviceExist.Name + ".default.127.0.0.1.nip.io"
	//
	//	if err := r.Status().Update(ctx, &workflowRun); err != nil {
	//		log.Error(err, "unable to update WorkflowRun status")
	//		return ctrl.Result{}, err
	//	}
	//	log.Info("> WorkflowRun Updated: " + workflowRun.Name + " Workflow Run: " + workflowRun.Status.WorkflowId)
	//
	//}

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

package statemachinerunner

import (
	"context"
	"fmt"
	statev1 "github.com/salaboy/knative-state/pkg/apis/state/v1"
	statemachinerunnerreconciler "github.com/salaboy/knative-state/pkg/client/injection/reconciler/state/v1/statemachinerunner"
	listers "github.com/salaboy/knative-state/pkg/client/listers/state/v1"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	eventingapi "knative.dev/eventing/pkg/apis/eventing/v1"
	eventingClientSet "knative.dev/eventing/pkg/client/clientset/versioned"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	servingapi "knative.dev/serving/pkg/apis/serving/v1"
	servingClientSet "knative.dev/serving/pkg/client/clientset/versioned"
	"os"
	"strings"
)

type Reconciler struct {
	servingClientSet  servingClientSet.Interface
	eventingClientSet eventingClientSet.Interface

	dynamicClientSet   dynamic.Interface
	statemachineLister listers.StateMachineLister
}

// Check that our Reconciler implements Interface
var _ statemachinerunnerreconciler.Interface = (*Reconciler)(nil)

// ReconcilerArgs are the arguments needed to create a broker.Reconciler.
type ReconcilerArgs struct {
}

var RUNNER_IMAGE = os.Getenv("RUNNER_IMAGE")

// Reconcile StateMachineRunners
// A StateMachineRunner needs:
//   - A StateMachine reference that is valid and can be parsed to obtain which events are we interested in
//   - A ready broker to receive these events
//   - A Service that can be used to create new StateMachine instances and query state

// Lifecycle for the reconcilation:
//  - Check StateMachineRef and parse states, fail if there is no StateMachine available or if it cannot be parsed
//  - Check if a runner exists for this statemachine runner,
//     - if Not Create Runner (KService), get the Service URL to create triggers
//  - Check if a broker was specified in the StateMachineRunner Spec
//      - if it was specified: check to see if the broker exists and get the status
//      - if the broker wasn't specified, create a new broker
//  - Once we know which broker we will be using
//    - Create Triggers to the associated broker, using the KService URL

//  The StateMachineRunner state should be changed to Ready only when:
//    - The KService is up and Ready
//    - The Selected Broker has a URL and is Ready

func (r *Reconciler) ReconcileKind(ctx context.Context, smr *statev1.StateMachineRunner) pkgreconciler.Event {
	logging.FromContext(ctx).Infow("Reconciling", zap.Any("StateMachineRunner", smr))

	// @TODO: remove from here
	if RUNNER_IMAGE == "" {
		RUNNER_IMAGE = "kind.local/knative-statemachine-runner-7a3c815d2bf3ebf9af9650f7624a29c9:93b9adcd6af50be3ba7f7b4848c79da214c0b4dcca39709c98c28905eb91b6a0"
	}

	var states []byte
	var stateMachine *statev1.StateMachine
	var err error
	// Check if the StateMachine Reference is ok and parse the states
	if smr.Spec.StateMachineRef != "" {
		stateMachine, err = r.GetStateMachine(ctx, smr)
		if err != nil {
			return fmt.Errorf("failed to find statemachine definition: %w", err)
		}
		logging.FromContext(ctx).Infow("StateMachine reference found: ", zap.Any("stateMachine", stateMachine))
		states, err = yaml.Marshal(stateMachine.Spec.StateMachineDefinition.StateMachineStates.States)
		if err != nil {
			return fmt.Errorf("failed to parse statemachine definition: %w", err)
		}
	}

	serviceName := "kservice-" + stateMachine.Name
	serviceExist, err := r.servingClientSet.ServingV1().Services(smr.Namespace).Get(ctx, serviceName, metav1.GetOptions{})
	var serviceUrl *apis.URL

	if serviceExist.Name == "" && err != nil { // the Service doesn't exist we need to create it
		logging.FromContext(ctx).Infow("Knative Service doesn't exist: "+serviceName+". Let's create it!", zap.Any("service", serviceExist))
		service := makeService(serviceName, states, stateMachine, smr)
		if _, err := r.servingClientSet.ServingV1().Services(smr.Namespace).Create(ctx, service, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}

	} else { // The service exist so get the URL and the Status
		logging.FromContext(ctx).Infow("Knative Service already exist: "+serviceName+".", zap.Any("service", serviceExist))
		if serviceExist.Status.URL != nil {
			logging.FromContext(ctx).Infow("Created KService URL for trigger", zap.Any("serviceUrl",  serviceExist.Status.URL.String()))
			serviceUrl, err = apis.ParseURL("http://" + serviceExist.Name + "." + smr.Namespace + ".svc.cluster.local" + "/statemachines/events")
			if err != nil {
				return fmt.Errorf("failed to parse URL for Service: %w"+serviceExist.Status.URL.String(), err)
			}
		}
	}

	var createBroker bool = false
	var brokerName string = ""
	if smr.Spec.Broker != "" { // if a broker name is provided we need to use that one
		brokerName = smr.Spec.Broker
	}else{ // if not we need to create one
		createBroker = true
		brokerName = "broker" + stateMachine.Name
	}

	if _, err := r.eventingClientSet.EventingV1().Brokers(smr.Namespace).Get(ctx, brokerName, metav1.GetOptions{}); err != nil {
		logging.FromContext(ctx).Infow("The Broker Doesn't exist", zap.Any("broker",  smr.Spec.Broker))
		//@TODO: get Broker Status:  broker.Status.Status
	}else{
		// There is a broker get the status, don't recreate
		createBroker = false
	}

	logging.FromContext(ctx).Infow("Do I need to create a Broker? ", zap.Any("createBroker",  createBroker))
	if createBroker { // If there is no broker name provided we need to create a new broker

		logging.FromContext(ctx).Infow("Knative Broker doesn't exist: "+brokerName+". Let's create it!")


		broker := makeBroker(brokerName, smr)
		if _, err := r.eventingClientSet.EventingV1().Brokers(smr.Namespace).Create(ctx, broker, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("failed to create broker: %w", err)
		}

	}

	// Once we have a broker let's create Triggers for StateMachine definition
	for stateType, _ := range stateMachine.Spec.StateMachineDefinition.StateMachineStates.States {
		logging.FromContext(ctx).Infow("Looking for Events in State: ", string(stateType))

		// Create triggers for events that the workflow is waiting for
		for eventName, _ := range stateMachine.Spec.StateMachineDefinition.StateMachineStates.States[stateType].Events {

			logging.FromContext(ctx).Infow("Creating trigger for Event: ", string(eventName), " in State : ", string(stateType))
			triggerName := strings.ToLower("t-" + stateMachine.Name + "-" + string(eventName))

			trigger, err := r.eventingClientSet.EventingV1().Triggers(smr.Namespace).Get(ctx, triggerName, metav1.GetOptions{})
			if trigger.Name == "" && err != nil { // The trigger doesn't exist lets create it
				trigger := makeTrigger(triggerName, string(eventName), brokerName, serviceUrl, smr)

				if _, err := r.eventingClientSet.EventingV1().Triggers(smr.Namespace).Create(ctx, trigger, metav1.CreateOptions{}); err != nil {
					return fmt.Errorf("failed to create trigger: %w", err)
				}
			}else {
				// DO nothing the trigger exist, maybe check that the URLs are the required ones, or log :)
			}
		}
	}

	for _, condition := range serviceExist.Status.Conditions {
		if condition.Type == apis.ConditionReady {
			logging.FromContext(ctx).Infow("StateMachineRunner Ready!")
			smr.Status.RunnerUrl = "http://" + serviceExist.Name + "." + smr.Namespace + ".127.0.0.1.nip.io"
			smr.Status.RunnerId = ""  // Need to fetch the ID from the Info endpoint
			smr.Status.BrokerUrl = "" // Need to check if the broker is up and add the URL here

		}
	}

	return nil
}

// isPodReady returns whether or not the given pod is ready.
func isServiceReady(s *servingapi.Service) bool {

	return false
}

func makeBroker(brokerName string, stateMachineRunner *statev1.StateMachineRunner) *eventingapi.Broker {
	return &eventingapi.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:            brokerName,
			Namespace:       stateMachineRunner.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(stateMachineRunner)},
		},
		Spec: eventingapi.BrokerSpec{},
	}
}

func makeTrigger(triggerName string, eventName string, brokerName string, url *apis.URL, stateMachineRunner *statev1.StateMachineRunner) *eventingapi.Trigger {
	return &eventingapi.Trigger{
		ObjectMeta: metav1.ObjectMeta{
			Name:            triggerName,
			Namespace:       stateMachineRunner.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(stateMachineRunner)},
		},
		Spec: eventingapi.TriggerSpec{
			Broker: brokerName,
			Filter: &eventingapi.TriggerFilter{
				Attributes: map[string]string{
					"type": eventName,
				},
			},
			Subscriber: duckv1.Destination{
				URI: url,
			},
		},
	}
}

func makeService(serviceName string, states []byte, stateMachine *statev1.StateMachine, stateMachineRunner *statev1.StateMachineRunner) *servingapi.Service {
	return &servingapi.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            serviceName,
			Namespace:       stateMachineRunner.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(stateMachineRunner)},
		},
		Spec: servingapi.ServiceSpec{
			ConfigurationSpec: servingapi.ConfigurationSpec{

				Template: servingapi.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							// Remove if we have redis in the runner
							"autoscaling.knative.dev/minScale": "1",
						},
					},
					Spec: servingapi.RevisionSpec{
						PodSpec: v1.PodSpec{

							Containers: []v1.Container{
								v1.Container{
									Name:  "knative-statemachine-runner",
									Image: RUNNER_IMAGE,
									Env: []v1.EnvVar{
										v1.EnvVar{
											Name:  "STATEMACHINE_NAME",
											Value: stateMachine.Name,
										},
										v1.EnvVar{
											Name:  "STATEMACHINE_VERSION",
											Value: stateMachine.Spec.StateMachineDefinition.Version,
										},
										v1.EnvVar{
											Name:  "STATEMACHINE_DEF",
											Value: fmt.Sprintf("%s", states),
										},
										v1.EnvVar{
											Name:  "EVENT_SINK",
											Value: stateMachineRunner.Spec.Sink,
										},
										v1.EnvVar{
											Name:  "REDIS_HOST",
											Value: stateMachineRunner.Spec.RedisHost,
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
}

func (r *Reconciler) GetStateMachine(ctx context.Context, smr *statev1.StateMachineRunner) (*statev1.StateMachine, pkgreconciler.Event) {

	stateMachine, err := r.statemachineLister.StateMachines(smr.Namespace).Get(smr.Spec.StateMachineRef)
	if err != nil {
		return nil, err
	}

	return stateMachine, nil
}

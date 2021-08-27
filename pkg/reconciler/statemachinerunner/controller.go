package statemachinerunner

import (
	"context"
	statemachineinformer "github.com/salaboy/knative-state/pkg/client/injection/informers/state/v1/statemachine"
	statemachinerunnerinformer "github.com/salaboy/knative-state/pkg/client/injection/informers/state/v1/statemachinerunner"
	statemachinerunnerreconciler "github.com/salaboy/knative-state/pkg/client/injection/reconciler/state/v1/statemachinerunner"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/logging"
	eventingclient "knative.dev/eventing/pkg/client/injection/client"
	servingclient "knative.dev/serving/pkg/client/injection/client"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)
	stateMachineRunnerInformer := statemachinerunnerinformer.Get(ctx)
	stateMachineInformer := statemachineinformer.Get(ctx)

	r := &Reconciler{
		dynamicClientSet:   dynamicclient.Get(ctx),
		statemachineLister: stateMachineInformer.Lister(),
		eventingClientSet:  eventingclient.Get(ctx),
		servingClientSet:  servingclient.Get(ctx),
	}
	impl := statemachinerunnerreconciler.NewImpl(ctx, r)

	logger.Info("Setting up event handlers for StateMachineRunners")

	stateMachineRunnerInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))


	return impl
}

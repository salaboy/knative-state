package statemachine

import (
	"context"
	statemachineinformer "github.com/salaboy/knative-state/pkg/client/injection/informers/state/v1/statemachine"
	statemachinereconciler "github.com/salaboy/knative-state/pkg/client/injection/reconciler/state/v1/statemachine"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/logging"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)
	stateMachineInformer := statemachineinformer.Get(ctx)
	r := &Reconciler{

		dynamicClientSet:   dynamicclient.Get(ctx),

	}
	impl := statemachinereconciler.NewImpl(ctx, r)

	logger.Info("Setting up event handlers")

	stateMachineInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))


	return impl
}
package statemachine

import (
	"context"
	statev1 "github.com/salaboy/knative-state/pkg/apis/state/v1"
	statemachinereconciler "github.com/salaboy/knative-state/pkg/client/injection/reconciler/state/v1/statemachine"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
)

type Reconciler struct {
	dynamicClientSet  dynamic.Interface

}

// Check that our Reconciler implements Interface
var _ statemachinereconciler.Interface = (*Reconciler)(nil)

// ReconcilerArgs are the arguments needed to create a broker.Reconciler.
type ReconcilerArgs struct {

}

func (r *Reconciler) ReconcileKind(ctx context.Context, b *statev1.StateMachine) pkgreconciler.Event {
	logging.FromContext(ctx).Infow("Reconciling", zap.Any("StateMachine", b))

	// @TODO: logic here


	return nil
}


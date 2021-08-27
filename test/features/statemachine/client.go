package statemachine

import (
	"context"
	statev1 "github.com/salaboy/knative-state/pkg/apis/state/v1"
	"knative.dev/reconciler-test/pkg/feature"
)

func GetStateMachine(ctx context.Context, t feature.T) *statev1.StateMachine {
	//c := Client(ctx)
	//name := state.GetStringOrFail(ctx, t, StateMachineNameKey)

	//stateMachine, err := c.StateMachines.Get(ctx, name, metav1.GetOptions{})
	//if err != nil {
	//	t.Errorf("failed to get StateMachine, %v", err)
	//}
	//return stateMachine
	return nil
}

type StateClient struct {
	//StateMachines  stateclient.StateMachineInterface
}

func Client(ctx context.Context) *StateClient {
	//sc := stateclient.Get(ctx).StateV1()
	//env := environment.FromContext(ctx)
	//
	//return &StateClient{
	//	StateMachines:  sc.StateMachines(env.Namespace()),
	//
	//}
	return nil
}

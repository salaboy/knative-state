package statemachine

import (
	"context"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/state"
)

const (
	StateMachineNameKey = "stateMachineName"
)

func ControlPlaneStateMachine(stateMachineName string) *feature.Feature {
	f := feature.NewFeatureNamed("StateMachine")

	f.Setup("Set StateMachine Name", setStateMachineName(stateMachineName))

	f.Stable("Conformance").
		Should("Broker objects SHOULD include a Ready condition in their status",
			stateMachineIsParseable)

	return f
}

func setStateMachineName(name string) feature.StepFn {
	return func(ctx context.Context, t feature.T) {
		state.SetOrFail(ctx, t, StateMachineNameKey, name)
	}
}

func stateMachineIsParseable(ctx context.Context, t feature.T) {
	//trigger := triggerfeatures.GetTrigger(ctx, t)
	//if trigger.Spec.Broker == "" {
	//	t.Error("broker is empty")
	//}
	//if strings.Contains(trigger.Spec.Broker, ",") {
	//	t.Errorf("more than one broker specified: %q", trigger.Spec.Broker)
	//}
}

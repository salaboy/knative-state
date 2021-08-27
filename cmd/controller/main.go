package main

import (
	statemachine "github.com/salaboy/knative-state/pkg/reconciler/statemachine"
	"github.com/salaboy/knative-state/pkg/reconciler/statemachinerunner"
	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	sharedmain.Main("controller",
		// StateMachine definition
		statemachine.NewController,
		// StateMachineRunner definition
		statemachinerunner.NewController,
	)
}

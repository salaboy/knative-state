// +build e2e

package test

import (
	"testing"
)



// TestStateMachine is an example simple test.
func TestStateMachine(t *testing.T) {
	// Signal to the go test framework that this test can be run in parallel
	// with other tests.
	t.Parallel()

	// Create an instance of an environment. The environment will be configured
	// with any relevant configuration and settings based on the global
	// environment settings. Using environment.Managed(t) will call env.Finish()
	// on test completion. If this option is not used, the test should call
	// env.Finish() to perform cleanup at the end of the test. Additional options
	// can be passed to Environment() if customization is required.
	//ctx, env := global.Environment(
	//	knative.WithKnativeNamespace(system.Namespace()),
	//	knative.WithLoggingConfig,
	//	knative.WithTracingConfig,
	//	k8s.WithEventListener,
	//)

	// With the instance of an Environment, perform one or more calls to Test().
	// Note: env.Test() is blocking until the feature completes.
	//env.Test(ctx, t, FooFeature1())

	//env.Finish()
}

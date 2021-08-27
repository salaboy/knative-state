package statemachine

import (
	"embed"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/manifest"
)

//go:embed *.yaml
var yaml embed.FS

func init() {
	// Process EventingGlobal.
	//if err := envconfig.Process("", &EnvCfg); err != nil {
	//	log.Fatal("Failed to process env var", err)
	//}
}

func GVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "flow.knative.dev", Version: "v1", Resource: "statemachines"}
}

// Install will create a Broker resource, augmented with the config fn options.
func Install(name string, opts ...manifest.CfgFn) feature.StepFn {
	//cfg := map[string]interface{}{
	//	"name": name,
	//}
	//for _, fn := range opts {
	//	fn(cfg)
	//}
	//
	//if dir, ok := cfg["__brokerTemplateDir"]; ok {
	//	return func(ctx context.Context, t feature.T) {
	//		if _, err := manifest.InstallYamlFS(ctx, os.DirFS(dir.(string)), cfg); err != nil {
	//			t.Fatal(err)
	//		}
	//	}
	//}
	//return func(ctx context.Context, t feature.T) {
	//	if _, err := manifest.InstallYamlFS(ctx, yaml, cfg); err != nil {
	//		t.Fatal(err)
	//	}
	//}
	return nil
}

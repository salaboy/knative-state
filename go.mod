module github.com/salaboy/knative-state

go 1.16

require (
	github.com/ahmetb/gen-crd-api-reference-docs v0.3.1-0.20210609063737-0067dc6dcea2
	go.uber.org/zap v1.18.1
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	knative.dev/eventing v0.25.0
	knative.dev/hack v0.0.0-20210622141627-e28525d8d260
	knative.dev/hack/schema v0.0.0-20210622141627-e28525d8d260
	knative.dev/pkg v0.0.0-20210803160015-21eb4c167cc5
	knative.dev/reconciler-test v0.0.0-20210803183715-b61cc77c06f6
	knative.dev/serving v0.25.0
)

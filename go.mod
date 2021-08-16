module github.com/salaboy/knative-state

go 1.15

require (
	github.com/ghodss/yaml v1.0.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.4
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	knative.dev/eventing v0.23.1
	knative.dev/pkg v0.0.0-20210510175900-4564797bf3b7
	knative.dev/serving v0.23.0
	sigs.k8s.io/controller-runtime v0.8.3
)

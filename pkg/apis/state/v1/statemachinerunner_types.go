/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)


// StateMachineRunnerSpec defines the desired state of WorkflowRunner
type StateMachineRunnerSpec struct {

	// Add connection URL for Redis here??
	// Here we should have something like RunnerConfig

	StateMachineRef string `json:"stateMachineRef,omitempty"`

	// +optional
	Broker string `json:"broker,omitempty"`

	// +optional
	Sink string `json:"sink,omitempty"`

	// +optional
	RedisHost string `json:"redisHost,omitempty"`
}

// StateMachineRunnerStatus defines the observed state of WorkflowRunner
type StateMachineRunnerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	RunnerUrl string `json:"runnerUrl,omitempty"`

	RunnerId string `json:"runnerId,omitempty"`

	BrokerUrl string `json:"brokerUrl,omitempty"`

	// inherits duck/v1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the StateMachineRunner that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1.Status `json:",inline"`
}

// +genreconciler
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StateMachineRunner is the Schema for the statemachinerunners API
type StateMachineRunner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StateMachineRunnerSpec   `json:"spec,omitempty"`
	Status StateMachineRunnerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StateMachineRunnerList contains a list of StateMachineRunner
type StateMachineRunnerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StateMachineRunner `json:"items"`
}

// GetGroupVersionKind returns GroupVersionKind for StateMachines
func (t *StateMachineRunner) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("StateMachineRunner")
}

// GetStatus retrieves the status of the StateMachine. Implements the KRShaped interface.
func (t *StateMachineRunner) GetStatus() *duckv1.Status {
	return &t.Status.Status
}
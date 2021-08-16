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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StateMachineRunnerSpec defines the desired state of WorkflowRunner
type StateMachineRunnerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Add connection URL for Redis here??
	// Here we should have something like RunnerConfig

	StateMachineRef string `json:"stateMachineRef,omitempty"`

	Broker string `json:"broker,omitempty"`

	Sink string `json:"sink,omitempty"`

	RedisHost string `json:"redisHost,omitempty"`
}

// StateMachineRunnerStatus defines the observed state of WorkflowRunner
type StateMachineRunnerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	RunnerUrl string `json:"runnerUrl,omitempty"`

	RunnerId string `json:"runnerId,omitempty"`

	BrokerUrl string `json:"brokerUrl,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StateMachineRunner is the Schema for the statemachinerunners API
type StateMachineRunner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StateMachineRunnerSpec   `json:"spec,omitempty"`
	Status StateMachineRunnerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StateMachineRunnerList contains a list of StateMachineRunner
type StateMachineRunnerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StateMachineRunner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StateMachineRunner{}, &StateMachineRunnerList{})
}

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

// @TODO: delete from here This comes from StateMachineRunner
type StateMachineDefinition struct {
	Id                 string             `json:"id,omitempty"`
	Name               string             `json:"name,omitempty"`
	Version            string             `json:"version,omitempty"`
	StateMachineStates StateMachineStates `json:"stateMachineStates"`
}

type StateMachineStates struct {
	States States `json:"states"`
}

// States represents a mapping of states and their implementations.
type States map[StateType]State

// Events represents a mapping of events and states.
type Events map[EventType]StateType

// State binds a state with an action and a set of events it can handle.
type State struct {
	Events Events `json:"events,omitempty"`
}

// StateType represents an extensible state type in the state machine.
type StateType string

// EventType represents an extensible event type in the state machine.
type EventType string

// @TODO: delete from here This comes from StateMachineRunner

// StateMachineSpec defines the desired state of Workflow
type StateMachineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	StateMachineDefinition StateMachineDefinition `json:"stateMachine,omitempty"`
}

// StateMachineStatus defines the observed state of Workflow
type StateMachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StateMachine is the Schema for the workflows API
type StateMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StateMachineSpec   `json:"spec,omitempty"`
	Status StateMachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StateMachineList contains a list of Workflow
type StateMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StateMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StateMachine{}, &StateMachineList{})
}

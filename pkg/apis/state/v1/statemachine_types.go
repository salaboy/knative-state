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


// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// @TODO: delete from here This comes from StateMachineRunner
type StateMachineDefinition struct {
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

// StateMachineSpec defines the desired state of StateMachine
type StateMachineSpec struct {

	StateMachineDefinition StateMachineDefinition `json:"stateMachine,omitempty"`
}

// StateMachineStatus defines the observed state of Workflow
type StateMachineStatus struct {
	// inherits duck/v1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the StateMachine that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1.Status `json:",inline"`
}

// +genreconciler
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StateMachine is the Schema for the StateMachine API
type StateMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StateMachineSpec   `json:"spec,omitempty"`
	Status StateMachineStatus `json:"status,omitempty"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StateMachineList contains a list of StateMachine
type StateMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StateMachine `json:"items"`
}


// GetGroupVersionKind returns GroupVersionKind for StateMachines
func (t *StateMachine) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("StateMachine")
}

// GetStatus retrieves the status of the StateMachine. Implements the KRShaped interface.
func (t *StateMachine) GetStatus() *duckv1.Status {
	return &t.Status.Status
}
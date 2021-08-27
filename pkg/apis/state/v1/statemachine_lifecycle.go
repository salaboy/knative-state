/*
Copyright 2020 The Knative Authors

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
	"sync"

	"knative.dev/pkg/apis"
)

const (
	StateMachineConditionReady = apis.ConditionReady
)

var stateMachineCondSet = apis.NewLivingConditionSet(
	StateMachineConditionReady,
)
var stateMachineCondSetLock = sync.RWMutex{}

// RegisterAlternateStateMachineConditionSet register a apis.ConditionSet for the given statemachine class.
func RegisterAlternateStateMachineConditionSet(conditionSet apis.ConditionSet) {
	stateMachineCondSetLock.Lock()
	defer stateMachineCondSetLock.Unlock()

	stateMachineCondSet = conditionSet
}

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (sm *StateMachine) GetConditionSet() apis.ConditionSet {
	stateMachineCondSetLock.RLock()
	defer stateMachineCondSetLock.RUnlock()

	return stateMachineCondSet
}

// GetConditionSet retrieves the condition set for this resource.
func (sms *StateMachineStatus) GetConditionSet() apis.ConditionSet {
	stateMachineCondSetLock.RLock()
	defer stateMachineCondSetLock.RUnlock()

	return stateMachineCondSet
}

// GetTopLevelCondition returns the top level Condition.
func (sms *StateMachineStatus) GetTopLevelCondition() *apis.Condition {
	return sms.GetConditionSet().Manage(sms).GetTopLevelCondition()
}

// SetAddress makes this Broker addressable by setting the URI. It also
// sets the BrokerConditionAddressable to true.
func (sms *StateMachineStatus) SetAddress(url *apis.URL) {
	//sms.Address.URL = url
	//if url != nil {
	//	sms.GetConditionSet().Manage(sms).MarkTrue(BrokerConditionAddressable)
	//} else {
	//	sms.GetConditionSet().Manage(sms).MarkFalse(BrokerConditionAddressable, "nil URL", "URL is nil")
	//}
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (sms *StateMachineStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return sms.GetConditionSet().Manage(sms).GetCondition(t)
}

// IsReady returns true if the resource is ready overall.
func (sms *StateMachineStatus) IsReady() bool {
	return sms.GetConditionSet().Manage(sms).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (sms *StateMachineStatus) InitializeConditions() {
	sms.GetConditionSet().Manage(sms).InitializeConditions()
}

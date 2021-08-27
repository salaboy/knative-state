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
	StateMachineRunnerConditionReady                             = apis.ConditionReady
	StateMachineRunnerAddressable    apis.ConditionType = "Addressable"
)

var stateMachineRunnerCondSet = apis.NewLivingConditionSet(
	StateMachineRunnerConditionReady,
	StateMachineRunnerAddressable,

)
var stateMachineRunnerCondSetLock = sync.RWMutex{}

// RegisterAlternateStateMachineRunnerConditionSet register a apis.ConditionSet for the given broker class.
func RegisterAlternateStateMachineRunnerConditionSet(conditionSet apis.ConditionSet) {
	stateMachineRunnerCondSetLock.Lock()
	defer stateMachineRunnerCondSetLock.Unlock()

	stateMachineRunnerCondSet = conditionSet
}

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (smr *StateMachineRunner) GetConditionSet() apis.ConditionSet {
	stateMachineRunnerCondSetLock.RLock()
	defer stateMachineRunnerCondSetLock.RUnlock()

	return stateMachineRunnerCondSet
}

// GetConditionSet retrieves the condition set for this resource.
func (smrs *StateMachineRunnerStatus) GetConditionSet() apis.ConditionSet {
	stateMachineRunnerCondSetLock.RLock()
	defer stateMachineRunnerCondSetLock.RUnlock()

	return stateMachineRunnerCondSet
}

// GetTopLevelCondition returns the top level Condition.
func (smrs *StateMachineRunnerStatus) GetTopLevelCondition() *apis.Condition {
	return smrs.GetConditionSet().Manage(smrs).GetTopLevelCondition()
}

// SetAddress makes this Broker addressable by setting the URI. It also
// sets the BrokerConditionAddressable to true.
func (smrs *StateMachineRunnerStatus) SetAddress(url *apis.URL) {
	//smrs.Address.URL = url
	if url != nil {
		smrs.GetConditionSet().Manage(smrs).MarkTrue(StateMachineRunnerAddressable)
	} else {
		smrs.GetConditionSet().Manage(smrs).MarkFalse(StateMachineRunnerAddressable, "nil URL", "URL is nil")
	}
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (smrs *StateMachineRunnerStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return smrs.GetConditionSet().Manage(smrs).GetCondition(t)
}

// IsReady returns true if the resource is ready overall.
func (smrs *StateMachineRunnerStatus) IsReady() bool {
	return smrs.GetConditionSet().Manage(smrs).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (smrs *StateMachineRunnerStatus) InitializeConditions() {
	smrs.GetConditionSet().Manage(smrs).InitializeConditions()
}

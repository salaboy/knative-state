/*
Copyright 2018 The Knative Authors

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

package state

import "k8s.io/apimachinery/pkg/runtime/schema"

const (
	GroupName = "flow.knative.dev"


)

var (
	// StateMachinesResource represents a Knative State StateMachine
	StateMachinesResource = schema.GroupResource{
		Group:    GroupName,
		Resource: "statemachines",
	}
	// StateMachineRunnersResource represents a Knative State StateMachineRunner
	StateMachineRunnersResource = schema.GroupResource{
		Group:    GroupName,
		Resource: "statemachinerunners",
	}

)

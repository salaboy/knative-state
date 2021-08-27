/*
Copyright 2021 The Knative Authors

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	statev1 "github.com/salaboy/knative-state/pkg/apis/state/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeStateMachineRunners implements StateMachineRunnerInterface
type FakeStateMachineRunners struct {
	Fake *FakeFlowV1
	ns   string
}

var statemachinerunnersResource = schema.GroupVersionResource{Group: "flow.knative.dev", Version: "v1", Resource: "statemachinerunners"}

var statemachinerunnersKind = schema.GroupVersionKind{Group: "flow.knative.dev", Version: "v1", Kind: "StateMachineRunner"}

// Get takes name of the stateMachineRunner, and returns the corresponding stateMachineRunner object, and an error if there is any.
func (c *FakeStateMachineRunners) Get(ctx context.Context, name string, options v1.GetOptions) (result *statev1.StateMachineRunner, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(statemachinerunnersResource, c.ns, name), &statev1.StateMachineRunner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*statev1.StateMachineRunner), err
}

// List takes label and field selectors, and returns the list of StateMachineRunners that match those selectors.
func (c *FakeStateMachineRunners) List(ctx context.Context, opts v1.ListOptions) (result *statev1.StateMachineRunnerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(statemachinerunnersResource, statemachinerunnersKind, c.ns, opts), &statev1.StateMachineRunnerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &statev1.StateMachineRunnerList{ListMeta: obj.(*statev1.StateMachineRunnerList).ListMeta}
	for _, item := range obj.(*statev1.StateMachineRunnerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested stateMachineRunners.
func (c *FakeStateMachineRunners) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(statemachinerunnersResource, c.ns, opts))

}

// Create takes the representation of a stateMachineRunner and creates it.  Returns the server's representation of the stateMachineRunner, and an error, if there is any.
func (c *FakeStateMachineRunners) Create(ctx context.Context, stateMachineRunner *statev1.StateMachineRunner, opts v1.CreateOptions) (result *statev1.StateMachineRunner, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(statemachinerunnersResource, c.ns, stateMachineRunner), &statev1.StateMachineRunner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*statev1.StateMachineRunner), err
}

// Update takes the representation of a stateMachineRunner and updates it. Returns the server's representation of the stateMachineRunner, and an error, if there is any.
func (c *FakeStateMachineRunners) Update(ctx context.Context, stateMachineRunner *statev1.StateMachineRunner, opts v1.UpdateOptions) (result *statev1.StateMachineRunner, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(statemachinerunnersResource, c.ns, stateMachineRunner), &statev1.StateMachineRunner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*statev1.StateMachineRunner), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeStateMachineRunners) UpdateStatus(ctx context.Context, stateMachineRunner *statev1.StateMachineRunner, opts v1.UpdateOptions) (*statev1.StateMachineRunner, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(statemachinerunnersResource, "status", c.ns, stateMachineRunner), &statev1.StateMachineRunner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*statev1.StateMachineRunner), err
}

// Delete takes name of the stateMachineRunner and deletes it. Returns an error if one occurs.
func (c *FakeStateMachineRunners) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(statemachinerunnersResource, c.ns, name), &statev1.StateMachineRunner{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeStateMachineRunners) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(statemachinerunnersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &statev1.StateMachineRunnerList{})
	return err
}

// Patch applies the patch and returns the patched stateMachineRunner.
func (c *FakeStateMachineRunners) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *statev1.StateMachineRunner, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(statemachinerunnersResource, c.ns, name, pt, data, subresources...), &statev1.StateMachineRunner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*statev1.StateMachineRunner), err
}
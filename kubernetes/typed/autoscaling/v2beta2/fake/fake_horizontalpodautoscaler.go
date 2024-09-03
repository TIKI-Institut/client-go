/*
Copyright The Kubernetes Authors.

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
	context "context"
	json "encoding/json"
	fmt "fmt"

	v2beta2 "k8s.io/api/autoscaling/v2beta2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	autoscalingv2beta2 "k8s.io/client-go/applyconfigurations/autoscaling/v2beta2"
	testing "k8s.io/client-go/testing"
)

// FakeHorizontalPodAutoscalers implements HorizontalPodAutoscalerInterface
type FakeHorizontalPodAutoscalers struct {
	Fake *FakeAutoscalingV2beta2
	ns   string
}

var horizontalpodautoscalersResource = v2beta2.SchemeGroupVersion.WithResource("horizontalpodautoscalers")

var horizontalpodautoscalersKind = v2beta2.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler")

// Get takes name of the horizontalPodAutoscaler, and returns the corresponding horizontalPodAutoscaler object, and an error if there is any.
func (c *FakeHorizontalPodAutoscalers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2beta2.HorizontalPodAutoscaler, err error) {
	emptyResult := &v2beta2.HorizontalPodAutoscaler{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(horizontalpodautoscalersResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2beta2.HorizontalPodAutoscaler), err
}

// List takes label and field selectors, and returns the list of HorizontalPodAutoscalers that match those selectors.
func (c *FakeHorizontalPodAutoscalers) List(ctx context.Context, opts v1.ListOptions) (result *v2beta2.HorizontalPodAutoscalerList, err error) {
	emptyResult := &v2beta2.HorizontalPodAutoscalerList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(horizontalpodautoscalersResource, horizontalpodautoscalersKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2beta2.HorizontalPodAutoscalerList{ListMeta: obj.(*v2beta2.HorizontalPodAutoscalerList).ListMeta}
	for _, item := range obj.(*v2beta2.HorizontalPodAutoscalerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested horizontalPodAutoscalers.
func (c *FakeHorizontalPodAutoscalers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(horizontalpodautoscalersResource, c.ns, opts))

}

// Create takes the representation of a horizontalPodAutoscaler and creates it.  Returns the server's representation of the horizontalPodAutoscaler, and an error, if there is any.
func (c *FakeHorizontalPodAutoscalers) Create(ctx context.Context, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, opts v1.CreateOptions) (result *v2beta2.HorizontalPodAutoscaler, err error) {
	emptyResult := &v2beta2.HorizontalPodAutoscaler{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(horizontalpodautoscalersResource, c.ns, horizontalPodAutoscaler, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2beta2.HorizontalPodAutoscaler), err
}

// Update takes the representation of a horizontalPodAutoscaler and updates it. Returns the server's representation of the horizontalPodAutoscaler, and an error, if there is any.
func (c *FakeHorizontalPodAutoscalers) Update(ctx context.Context, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, opts v1.UpdateOptions) (result *v2beta2.HorizontalPodAutoscaler, err error) {
	emptyResult := &v2beta2.HorizontalPodAutoscaler{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(horizontalpodautoscalersResource, c.ns, horizontalPodAutoscaler, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2beta2.HorizontalPodAutoscaler), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeHorizontalPodAutoscalers) UpdateStatus(ctx context.Context, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, opts v1.UpdateOptions) (result *v2beta2.HorizontalPodAutoscaler, err error) {
	emptyResult := &v2beta2.HorizontalPodAutoscaler{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(horizontalpodautoscalersResource, "status", c.ns, horizontalPodAutoscaler, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2beta2.HorizontalPodAutoscaler), err
}

// Delete takes name of the horizontalPodAutoscaler and deletes it. Returns an error if one occurs.
func (c *FakeHorizontalPodAutoscalers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(horizontalpodautoscalersResource, c.ns, name, opts), &v2beta2.HorizontalPodAutoscaler{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeHorizontalPodAutoscalers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(horizontalpodautoscalersResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v2beta2.HorizontalPodAutoscalerList{})
	return err
}

// Patch applies the patch and returns the patched horizontalPodAutoscaler.
func (c *FakeHorizontalPodAutoscalers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2beta2.HorizontalPodAutoscaler, err error) {
	emptyResult := &v2beta2.HorizontalPodAutoscaler{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(horizontalpodautoscalersResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2beta2.HorizontalPodAutoscaler), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied horizontalPodAutoscaler.
func (c *FakeHorizontalPodAutoscalers) Apply(ctx context.Context, horizontalPodAutoscaler *autoscalingv2beta2.HorizontalPodAutoscalerApplyConfiguration, opts v1.ApplyOptions) (result *v2beta2.HorizontalPodAutoscaler, err error) {
	if horizontalPodAutoscaler == nil {
		return nil, fmt.Errorf("horizontalPodAutoscaler provided to Apply must not be nil")
	}
	data, err := json.Marshal(horizontalPodAutoscaler)
	if err != nil {
		return nil, err
	}
	name := horizontalPodAutoscaler.Name
	if name == nil {
		return nil, fmt.Errorf("horizontalPodAutoscaler.Name must be provided to Apply")
	}
	emptyResult := &v2beta2.HorizontalPodAutoscaler{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(horizontalpodautoscalersResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions()), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2beta2.HorizontalPodAutoscaler), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeHorizontalPodAutoscalers) ApplyStatus(ctx context.Context, horizontalPodAutoscaler *autoscalingv2beta2.HorizontalPodAutoscalerApplyConfiguration, opts v1.ApplyOptions) (result *v2beta2.HorizontalPodAutoscaler, err error) {
	if horizontalPodAutoscaler == nil {
		return nil, fmt.Errorf("horizontalPodAutoscaler provided to Apply must not be nil")
	}
	data, err := json.Marshal(horizontalPodAutoscaler)
	if err != nil {
		return nil, err
	}
	name := horizontalPodAutoscaler.Name
	if name == nil {
		return nil, fmt.Errorf("horizontalPodAutoscaler.Name must be provided to Apply")
	}
	emptyResult := &v2beta2.HorizontalPodAutoscaler{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(horizontalpodautoscalersResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions(), "status"), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2beta2.HorizontalPodAutoscaler), err
}

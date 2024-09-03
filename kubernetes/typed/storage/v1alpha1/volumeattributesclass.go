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

package v1alpha1

import (
	context "context"

	storagev1alpha1 "k8s.io/api/storage/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	applyconfigurationsstoragev1alpha1 "k8s.io/client-go/applyconfigurations/storage/v1alpha1"
	gentype "k8s.io/client-go/gentype"
	scheme "k8s.io/client-go/kubernetes/scheme"
)

// VolumeAttributesClassesGetter has a method to return a VolumeAttributesClassInterface.
// A group's client should implement this interface.
type VolumeAttributesClassesGetter interface {
	VolumeAttributesClasses() VolumeAttributesClassInterface
}

// VolumeAttributesClassInterface has methods to work with VolumeAttributesClass resources.
type VolumeAttributesClassInterface interface {
	Create(ctx context.Context, volumeAttributesClass *storagev1alpha1.VolumeAttributesClass, opts v1.CreateOptions) (*storagev1alpha1.VolumeAttributesClass, error)
	Update(ctx context.Context, volumeAttributesClass *storagev1alpha1.VolumeAttributesClass, opts v1.UpdateOptions) (*storagev1alpha1.VolumeAttributesClass, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*storagev1alpha1.VolumeAttributesClass, error)
	List(ctx context.Context, opts v1.ListOptions) (*storagev1alpha1.VolumeAttributesClassList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *storagev1alpha1.VolumeAttributesClass, err error)
	Apply(ctx context.Context, volumeAttributesClass *applyconfigurationsstoragev1alpha1.VolumeAttributesClassApplyConfiguration, opts v1.ApplyOptions) (result *storagev1alpha1.VolumeAttributesClass, err error)
	VolumeAttributesClassExpansion
}

// volumeAttributesClasses implements VolumeAttributesClassInterface
type volumeAttributesClasses struct {
	*gentype.ClientWithListAndApply[*storagev1alpha1.VolumeAttributesClass, *storagev1alpha1.VolumeAttributesClassList, *applyconfigurationsstoragev1alpha1.VolumeAttributesClassApplyConfiguration]
}

// newVolumeAttributesClasses returns a VolumeAttributesClasses
func newVolumeAttributesClasses(c *StorageV1alpha1Client) *volumeAttributesClasses {
	return &volumeAttributesClasses{
		gentype.NewClientWithListAndApply[*storagev1alpha1.VolumeAttributesClass, *storagev1alpha1.VolumeAttributesClassList, *applyconfigurationsstoragev1alpha1.VolumeAttributesClassApplyConfiguration](
			"volumeattributesclasses",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *storagev1alpha1.VolumeAttributesClass { return &storagev1alpha1.VolumeAttributesClass{} },
			func() *storagev1alpha1.VolumeAttributesClassList { return &storagev1alpha1.VolumeAttributesClassList{} }),
	}
}

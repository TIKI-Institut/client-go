/*
Copyright 2015 The Kubernetes Authors.

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

package testing

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

func getArbitraryResource(s schema.GroupVersionResource, name, namespace string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       s.Resource,
			"apiVersion": s.Version,
			"metadata": map[string]interface{}{
				"name":            name,
				"namespace":       namespace,
				"generateName":    "test_generateName",
				"uid":             "test_uid",
				"resourceVersion": "test_resourceVersion",
			},
			"data": strconv.Itoa(rand.Int()),
		},
	}
}

func TestWatchCallNonNamespace(t *testing.T) {
	testResource := schema.GroupVersionResource{Group: "", Version: "test_version", Resource: "test_kind"}
	testObj := getArbitraryResource(testResource, "test_name", "test_namespace")
	accessor, err := meta.Accessor(testObj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ns := accessor.GetNamespace()
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	o := NewObjectTracker(scheme, codecs.UniversalDecoder())
	watch, err := o.Watch(testResource, ns)
	if err != nil {
		t.Fatalf("test resource watch failed in %s: %v ", ns, err)
	}
	go func() {
		err := o.Create(testResource, testObj, ns)
		if err != nil {
			t.Errorf("test resource creation failed: %v", err)
		}
	}()
	out := <-watch.ResultChan()
	assert.Equal(t, testObj, out.Object, "watched object mismatch")
}

func TestWatchCallAllNamespace(t *testing.T) {
	testResource := schema.GroupVersionResource{Group: "", Version: "test_version", Resource: "test_kind"}
	testObj := getArbitraryResource(testResource, "test_name", "test_namespace")
	accessor, err := meta.Accessor(testObj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ns := accessor.GetNamespace()
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	o := NewObjectTracker(scheme, codecs.UniversalDecoder())
	w, err := o.Watch(testResource, "test_namespace")
	if err != nil {
		t.Fatalf("test resource watch failed in test_namespace: %v", err)
	}
	wAll, err := o.Watch(testResource, "")
	if err != nil {
		t.Fatalf("test resource watch failed in all namespaces: %v", err)
	}
	go func() {
		err := o.Create(testResource, testObj, ns)
		assert.NoError(t, err, "test resource creation failed")
	}()
	out := <-w.ResultChan()
	outAll := <-wAll.ResultChan()
	assert.Equal(t, watch.Added, out.Type, "watch event mismatch")
	assert.Equal(t, watch.Added, outAll.Type, "watch event mismatch")
	assert.Equal(t, testObj, out.Object, "watched created object mismatch")
	assert.Equal(t, testObj, outAll.Object, "watched created object mismatch")
	go func() {
		err := o.Update(testResource, testObj, ns)
		assert.NoError(t, err, "test resource updating failed")
	}()
	out = <-w.ResultChan()
	outAll = <-wAll.ResultChan()
	assert.Equal(t, watch.Modified, out.Type, "watch event mismatch")
	assert.Equal(t, watch.Modified, outAll.Type, "watch event mismatch")
	assert.Equal(t, testObj, out.Object, "watched updated object mismatch")
	assert.Equal(t, testObj, outAll.Object, "watched updated object mismatch")
	go func() {
		err := o.Delete(testResource, "test_namespace", "test_name")
		assert.NoError(t, err, "test resource deletion failed")
	}()
	out = <-w.ResultChan()
	outAll = <-wAll.ResultChan()
	assert.Equal(t, watch.Deleted, out.Type, "watch event mismatch")
	assert.Equal(t, watch.Deleted, outAll.Type, "watch event mismatch")
	assert.Equal(t, testObj, out.Object, "watched deleted object mismatch")
	assert.Equal(t, testObj, outAll.Object, "watched deleted object mismatch")
}

func TestWatchCallMultipleInvocation(t *testing.T) {
	cases := []struct {
		name string
		op   watch.EventType
		ns   string
	}{
		{
			"foo",
			watch.Added,
			"test_namespace",
		},
		{
			"bar",
			watch.Added,
			"test_namespace",
		},
		{
			"baz",
			watch.Added,
			"",
		},
		{
			"bar",
			watch.Modified,
			"test_namespace",
		},
		{
			"baz",
			watch.Modified,
			"",
		},
		{
			"foo",
			watch.Deleted,
			"test_namespace",
		},
		{
			"bar",
			watch.Deleted,
			"test_namespace",
		},
		{
			"baz",
			watch.Deleted,
			"",
		},
	}

	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	testResource := schema.GroupVersionResource{Group: "", Version: "test_version", Resource: "test_kind"}

	o := NewObjectTracker(scheme, codecs.UniversalDecoder())
	watchNamespaces := []string{
		"",
		"",
		"test_namespace",
		"test_namespace",
	}
	var wg sync.WaitGroup
	wg.Add(len(watchNamespaces))
	for idx, watchNamespace := range watchNamespaces {
		i := idx
		watchNamespace := watchNamespace
		w, err := o.Watch(testResource, watchNamespace)
		if err != nil {
			t.Fatalf("test resource watch failed in %s: %v", watchNamespace, err)
		}
		go func() {
			assert.NoError(t, err, "watch invocation failed")
			for _, c := range cases {
				if watchNamespace == "" || c.ns == watchNamespace {
					fmt.Printf("%#v %#v\n", c, i)
					event := <-w.ResultChan()
					accessor, err := meta.Accessor(event.Object)
					if err != nil {
						t.Errorf("unexpected error: %v", err)
						break
					}
					assert.Equal(t, c.op, event.Type, "watch event mismatched")
					assert.Equal(t, c.name, accessor.GetName(), "watched object mismatch")
					assert.Equal(t, c.ns, accessor.GetNamespace(), "watched object mismatch")
				}
			}
			wg.Done()
		}()
	}
	for _, c := range cases {
		switch c.op {
		case watch.Added:
			obj := getArbitraryResource(testResource, c.name, c.ns)
			o.Create(testResource, obj, c.ns)
		case watch.Modified:
			obj := getArbitraryResource(testResource, c.name, c.ns)
			o.Update(testResource, obj, c.ns)
		case watch.Deleted:
			o.Delete(testResource, c.ns, c.name)
		}
	}
	wg.Wait()
}

func TestWatchAddAfterStop(t *testing.T) {
	testResource := schema.GroupVersionResource{Group: "", Version: "test_version", Resource: "test_kind"}
	testObj := getArbitraryResource(testResource, "test_name", "test_namespace")
	accessor, err := meta.Accessor(testObj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ns := accessor.GetNamespace()
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	o := NewObjectTracker(scheme, codecs.UniversalDecoder())
	watch, err := o.Watch(testResource, ns)
	if err != nil {
		t.Errorf("watch creation failed: %v", err)
	}

	// When the watch is stopped it should ignore later events without panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Watch panicked when it should have ignored create after stop: %v", r)
		}
	}()

	watch.Stop()
	err = o.Create(testResource, testObj, ns)
	if err != nil {
		t.Errorf("test resource creation failed: %v", err)
	}
}

func TestPatch(t *testing.T) {

	t.Run("WithMissingObject", func(t *testing.T) {
		nodesResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}

		scheme := runtime.NewScheme()
		codecs := serializer.NewCodecFactory(scheme)
		o := NewObjectTracker(scheme, codecs.UniversalDecoder())
		reaction := ObjectReaction(o)
		action := NewRootPatchSubresourceAction(nodesResource, "node-1", types.StrategicMergePatchType, []byte(`{}`))
		handled, node, err := reaction(action)
		assert.True(t, handled)
		assert.Nil(t, node)
		assert.EqualError(t, err, `nodes "node-1" not found`)
	})

	t.Run("PatchingProperty", func(t *testing.T) {
		configmapsResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}
		configmapsGvk := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}

		targetObject := &v1.ConfigMap{}

		scheme := runtime.NewScheme()
		scheme.AddKnownTypeWithName(configmapsGvk, targetObject)

		codecs := serializer.NewCodecFactory(scheme)

		targetMeta, err := strategicpatch.NewPatchMetaFromStruct(targetObject)
		assert.NoError(t, err)

		schemeWithPatchMeta := NewSchemeWithPatchMeta(
			scheme,
			map[schema.GroupVersionKind]strategicpatch.LookupPatchMeta{
				configmapsGvk: targetMeta,
			},
		)

		o := NewObjectTracker(schemeWithPatchMeta, codecs.UniversalDecoder())

		unstructuredObj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "ConfigMap",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name":      "cm-1",
					"namespace": "cm-namespace",
				},
				"data": map[string]interface{}{"field-1": "A", "field-2": "B"},
			},
		}
		expected := unstructuredObj.DeepCopy()
		expected.Object["data"] = map[string]interface{}{"field-1": "C", "field-2": "B"}

		err = o.Add(unstructuredObj)
		assert.NoError(t, err)

		reaction := ObjectReaction(o)
		action := NewPatchAction(configmapsResource, "cm-namespace", "cm-1", types.StrategicMergePatchType, []byte(`{"data":{"field-1":"C"}}`))
		handled, changedCm, err := reaction(action)
		assert.NoError(t, err)
		assert.True(t, handled)
		assert.EqualValues(t, expected, changedCm)

		trackedObject, err := o.Get(configmapsResource, "cm-namespace", "cm-1")
		assert.NoError(t, err)
		assert.EqualValues(t, expected, trackedObject)

	})
}

func TestGetWithExactMatch(t *testing.T) {
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)

	constructObject := func(s schema.GroupVersionResource, name, namespace string) (*unstructured.Unstructured, schema.GroupVersionResource) {
		obj := getArbitraryResource(s, name, namespace)
		gvks, _, err := scheme.ObjectKinds(obj)
		assert.NoError(t, err)
		gvr, _ := meta.UnsafeGuessKindToResource(gvks[0])
		return obj, gvr
	}

	var err error
	// Object with empty namespace
	o := NewObjectTracker(scheme, codecs.UniversalDecoder())
	nodeResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "node"}
	node, gvr := constructObject(nodeResource, "node", "")

	assert.Nil(t, o.Add(node))

	// Exact match
	_, err = o.Get(gvr, "", "node")
	assert.NoError(t, err)

	// Unexpected namespace provided
	_, err = o.Get(gvr, "ns", "node")
	assert.Error(t, err)
	errNotFound := errors.NewNotFound(gvr.GroupResource(), "node")
	assert.EqualError(t, err, errNotFound.Error())

	// Object with non-empty namespace
	o = NewObjectTracker(scheme, codecs.UniversalDecoder())
	podResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pod"}
	pod, gvr := constructObject(podResource, "pod", "default")
	assert.Nil(t, o.Add(pod))

	// Exact match
	_, err = o.Get(gvr, "default", "pod")
	assert.NoError(t, err)

	// Missing namespace
	_, err = o.Get(gvr, "", "pod")
	assert.Error(t, err)
	errNotFound = errors.NewNotFound(gvr.GroupResource(), "pod")
	assert.EqualError(t, err, errNotFound.Error())
}

func Test_resourceCovers(t *testing.T) {
	type args struct {
		resource string
		action   Action
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				resource: "*",
				action:   ActionImpl{},
			},
			want: true,
		},
		{
			args: args{
				resource: "serviceaccounts",
				action:   ActionImpl{},
			},
			want: false,
		},
		{
			args: args{
				resource: "serviceaccounts",
				action: ActionImpl{
					Resource: schema.GroupVersionResource{
						Resource: "serviceaccounts",
					},
				},
			},
			want: true,
		},
		{
			args: args{
				resource: "serviceaccounts/token",
				action: ActionImpl{
					Resource: schema.GroupVersionResource{},
				},
			},
			want: false,
		},
		{
			args: args{
				resource: "serviceaccounts/token",
				action: ActionImpl{
					Resource: schema.GroupVersionResource{
						Resource: "serviceaccounts",
					},
				},
			},
			want: false,
		},
		{
			args: args{
				resource: "serviceaccounts/token",
				action: ActionImpl{
					Resource:    schema.GroupVersionResource{},
					Subresource: "token",
				},
			},
			want: false,
		},
		{
			args: args{
				resource: "serviceaccounts/token",
				action: ActionImpl{
					Resource: schema.GroupVersionResource{
						Resource: "serviceaccounts",
					},
					Subresource: "token",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceCovers(tt.args.resource, tt.args.action); got != tt.want {
				t.Errorf("resourceCovers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategicMergePatch(t *testing.T) {
	//StrategicMergePatch is not supposed to be called with an *unstructured.Unstructured as reference for the underlying schema

	//local entries in dynamic fake client will always be of type *unstructured.Unstructured (see: dynamic/fake/simple.go:50 "objects, err := convertObjectsToUnstructured(scheme, objects)")
	unstructuredObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "ConfigMap",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "some-config-map",
				"namespace": "some-namespace",
			},
			"data": map[string]interface{}{
				"foo":  "bar",
				"john": "doe",
			},
		},
	}

	//the expected object shall no longer have the 'data' property
	expectedObj := unstructuredObj.DeepCopy()
	expectedObj.Object["data"] = map[string]interface{}{"foo": "foobar", "john": "doe"}
	expectedObjBytes, err := json.Marshal(expectedObj)
	assert.NoError(t, err)

	unstructuredObjBytes, err := json.Marshal(unstructuredObj)
	assert.NoError(t, err)

	t.Run("using no struct", func(t *testing.T) {
		outObj, err := strategicpatch.StrategicMergePatch(unstructuredObjBytes, []byte(`{"data":{"foo": "foobar"}}`), unstructuredObj)
		assert.EqualError(t, err, `unable to find api field in struct Unstructured for the json field "data"`)
		assert.Nil(t, outObj)
	})

	t.Run("using struct", func(t *testing.T) {
		outObj, err := strategicpatch.StrategicMergePatch(unstructuredObjBytes, []byte(`{"data":{"foo": "foobar"}}`), &v1.ConfigMap{})
		assert.NoError(t, err)
		assert.Equal(t, expectedObjBytes, outObj, "object mismatch")
	})

}

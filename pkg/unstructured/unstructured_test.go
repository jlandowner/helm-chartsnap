package unstructured

import (
	"reflect"
	"testing"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructuredscheme" // Required for DeepCopy
	"github.com/stretchr/testify/assert"
	"github.com/jlandowner/helm-chartsnap/pkg/api/v1alpha1"
)

func TestEncode(t *testing.T) {
	type args struct {
		arr []metaV1.Unstructured
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				arr: []metaV1.Unstructured{
					{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Pod",
							"metadata": map[string]interface{}{
								"name": "pod1",
							},
						},
					},
					{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Service",
							"metadata": map[string]interface{}{
								"name": "service1",
							},
						},
					},
				},
			},
			want: `apiVersion: v1
kind: Pod
metadata:
  name: pod1
---
apiVersion: v1
kind: Service
metadata:
  name: service1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.arr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Encode() = %s, want %v", got, tt.want)
			}
		})
	}
}

func newTestUnstructuredConfigMap(name string, data map[string]interface{}) *metaV1.Unstructured {
	return &metaV1.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"data": data,
		},
	}
}

func TestApplyFixedValue(t *testing.T) {
	tests := []struct {
		name              string
		initialManifests  []metaV1.Unstructured
		config            v1alpha1.SnapshotConfig
		expectedManifests []metaV1.Unstructured // For full object comparison if needed, or use expectedFields
		expectedFields    map[int]map[string]string // map[manifestIndex]map[jsonPath]expectedValue
		wantErr           bool
	}{
		{
			name: "Scenario 1: JSONPathList with paths only",
			initialManifests: []metaV1.Unstructured{
				*newTestUnstructuredConfigMap("my-cm", map[string]interface{}{"key1": "value1", "key2": "value2"}),
			},
			config: v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "my-cm",
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1"},
							{Path: "/metadata/name"},
						},
					},
				},
			},
			expectedFields: map[int]map[string]string{
				0: {
					"/data/key1":     v1alpha1.DynamicValue,
					"/metadata/name": v1alpha1.DynamicValue,
				},
			},
		},
		{
			name: "Scenario 1.1: JSONPathList with paths only and Base64",
			initialManifests: []metaV1.Unstructured{
				*newTestUnstructuredConfigMap("my-secret", map[string]interface{}{"secretKey": "c2VjcmV0VmFsdWU="}), // Assume data is already base64 encoded as per k8s
			},
			config: v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap", // Testing with ConfigMap for simplicity, field would be stringData for Secret
						Name:       "my-secret",
						Base64:     true,
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/secretKey"},
						},
					},
				},
			},
			expectedFields: map[int]map[string]string{
				0: {
					"/data/secretKey": v1alpha1.Base64DynamicValue,
				},
			},
		},
		{
			name: "Scenario 2: JSONPathList with paths and values",
			initialManifests: []metaV1.Unstructured{
				*newTestUnstructuredConfigMap("my-cm-paths-values", map[string]interface{}{"key1": "value1", "key2": "value2"}),
			},
			config: v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "my-cm-paths-values",
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1", Value: "fixed-value1"},
							{Path: "/metadata/name", Value: "fixed-name-cm"},
						},
					},
				},
			},
			expectedFields: map[int]map[string]string{
				0: {
					"/data/key1":     "fixed-value1",
					"/metadata/name": "fixed-name-cm",
				},
			},
		},
		{
			name: "Scenario 3: JSONPathList with mixed paths-only and paths-with-values",
			initialManifests: []metaV1.Unstructured{
				*newTestUnstructuredConfigMap("my-cm-mixed", map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"}),
			},
			config: v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "my-cm-mixed",
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1", Value: "fixed-value-mixed"},
							{Path: "/data/key2"}, // Fallback to DynamicValue
							{Path: "/metadata/name", Value: "fixed-name-mixed-cm"},
						},
					},
				},
			},
			expectedFields: map[int]map[string]string{
				0: {
					"/data/key1":     "fixed-value-mixed",
					"/data/key2":     v1alpha1.DynamicValue,
					"/metadata/name": "fixed-name-mixed-cm",
				},
			},
		},
		{
			name: "No matching manifest",
			initialManifests: []metaV1.Unstructured{
				*newTestUnstructuredConfigMap("another-cm", map[string]interface{}{"key1": "value1"}),
			},
			config: v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "non-existent-cm", // This name does not match
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1", Value: "should-not-apply"},
						},
					},
				},
			},
			expectedFields: map[int]map[string]string{
				0: { // Original values should remain
					"/data/key1":     "value1",
					"/metadata/name": "another-cm",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Deep copy initial manifests to avoid modification across test cases
			currentManifests := make([]metaV1.Unstructured, len(tt.initialManifests))
			for i, m := range tt.initialManifests {
				currentManifests[i] = *m.DeepCopyObject().(*metaV1.Unstructured)
			}

			err := ApplyFixedValue(tt.config, currentManifests)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				for i, expectedFieldsForManifest := range tt.expectedFields {
					manifest := currentManifests[i]
					for jsonPath, expectedValue := range expectedFieldsForManifest {
						pathParts := strings.Split(strings.TrimPrefix(jsonPath, "/"), "/")
						actualValue, found, _ := metaV1.NestedString(manifest.Object, pathParts...)
						assert.True(t, found, "Path %s not found in manifest %d", jsonPath, i)
						assert.Equal(t, expectedValue, actualValue, "Value mismatch for path %s in manifest %d", jsonPath, i)
					}
				}
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    []metaV1.Unstructured
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				source: `
apiVersion: v1
kind: Pod
metadata:
  name: pod1
---
apiVersion: v1
kind: Service
metadata:
   name: service1
`,
			},
			want: []metaV1.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name": "pod1",
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Service",
						"metadata": map[string]interface{}{
							"name": "service1",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, errs := Decode(tt.args.source)
			if (len(errs) != 0) != tt.wantErr {
				t.Errorf("Decode() errorcount = %d, wantErr %v", len(errs), tt.wantErr)
				for _, err := range errs {
					t.Errorf("Decode() error = %v", err)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	type args struct {
		obj   metaV1.Unstructured
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *metaV1.Unstructured
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				obj: metaV1.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name": "pod1",
						},
					},
				},
				key:   "/metadata/name",
				value: "pod2",
			},
			want: &metaV1.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]interface{}{
						"name": "pod2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Replace(tt.args.obj, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

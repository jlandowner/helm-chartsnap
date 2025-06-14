package unstructured

import (
	"encoding/json"
	"reflect"
	"testing"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func TestSetLogger(t *testing.T) {
	// Test SetLogger with nil
	SetLogger(nil)
	
	// Test that log() returns a default logger when logger is nil
	logger := log()
	if logger == nil {
		t.Error("Expected log() to return a logger even when SetLogger(nil) is called")
	}
	
	// Reset logger
	SetLogger(nil)
}

func TestStringToUnstructured(t *testing.T) {
	yamlStr := `apiVersion: v1
kind: Pod
metadata:
  name: test-pod`
	
	gvk, result, err := StringToUnstructured(yamlStr)
	if err != nil {
		t.Errorf("StringToUnstructured() error = %v", err)
		return
	}
	
	if gvk == nil {
		t.Error("Expected gvk to be non-nil")
		return
	}
	
	if result.GetAPIVersion() != "v1" {
		t.Errorf("Expected apiVersion 'v1', got '%s'", result.GetAPIVersion())
	}
	if result.GetKind() != "Pod" {
		t.Errorf("Expected kind 'Pod', got '%s'", result.GetKind())
	}
	if result.GetName() != "test-pod" {
		t.Errorf("Expected name 'test-pod', got '%s'", result.GetName())
	}
}

func TestBytesToUnstructured(t *testing.T) {
	yamlBytes := []byte(`apiVersion: v1
kind: Service
metadata:
  name: test-service`)
	
	gvk, result, err := BytesToUnstructured(yamlBytes)
	if err != nil {
		t.Errorf("BytesToUnstructured() error = %v", err)
		return
	}
	
	if gvk == nil {
		t.Error("Expected gvk to be non-nil")
		return
	}
	
	if result.GetKind() != "Service" {
		t.Errorf("Expected kind 'Service', got '%s'", result.GetKind())
	}
	if result.GetName() != "test-service" {
		t.Errorf("Expected name 'test-service', got '%s'", result.GetName())
	}
}

func TestUnstructuredToJSONBytes(t *testing.T) {
	obj := &metaV1.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": "test-config",
			},
		},
	}
	
	result, err := UnstructuredToJSONBytes(obj)
	if err != nil {
		t.Errorf("UnstructuredToJSONBytes() error = %v", err)
		return
	}
	
	// Should be valid JSON
	var jsonObj map[string]interface{}
	if err := json.Unmarshal(result, &jsonObj); err != nil {
		t.Errorf("Result is not valid JSON: %v", err)
	}
	
	// Check content
	if jsonObj["kind"] != "ConfigMap" {
		t.Errorf("Expected kind 'ConfigMap', got '%v'", jsonObj["kind"])
	}
}


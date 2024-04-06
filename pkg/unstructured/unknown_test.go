package unstructured

import (
	"reflect"
	"testing"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestUnknownError_Unstructured(t *testing.T) {
	raw := "some raw data"
	err := NewUnknownError(raw)

	obj := err.Unstructured()
	expectedObj := &metaV1.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "helm-chartsnap.jlandowner.dev/v1alpha1",
			"kind":       "Unknown",
			"raw":        "some raw data",
		},
	}

	if !reflect.DeepEqual(obj, expectedObj) {
		t.Errorf("Expected obj to be %v, but got %v", err, obj)
	}

}

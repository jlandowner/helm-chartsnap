package unstructured

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "helm-chartsnap.jlandowner.dev", Version: "v1alpha1"}
)

func NewUnknownError(raw string) *UnknownError {
	return &UnknownError{Raw: raw}
}

type UnknownError struct {
	Raw string
}

func (e *UnknownError) Error() string {
	out, err := yaml.Marshal(e.Unstructured())
	if err != nil {
		return "xxx"
	}
	return fmt.Sprintf("WARN: failed to recognize a resource. snapshot as Unknown: \n---\n%s\n---", out)
}

func (e *UnknownError) Unstructured() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": GroupVersion.String(),
			"kind":       "Unknown",
			"raw":        e.Raw,
		},
	}
}

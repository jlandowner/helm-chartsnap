package charts

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	unstructuredutil "github.com/jlandowner/helm-chartsnap/pkg/unstructured"
)

type SnapshotValues struct {
	TestSpec TestSpec `yaml:"testSpec,omitempty"`
}

type TestSpec struct {
	DynamicFields []ManifestPath `yaml:"dynamicFields,omitempty"`
}

type ManifestPath struct {
	Kind       string   `yaml:"kind,omitempty"`
	APIVersion string   `yaml:"apiVersion,omitempty"`
	Name       string   `yaml:"name,omitempty"`
	JSONPath   []string `yaml:"jsonPath,omitempty"`
}

func (t *TestSpec) ApplyFixedValue(manifests []unstructured.Unstructured) error {
	for _, v := range t.DynamicFields {
		for i, obj := range manifests {
			if v.APIVersion == obj.GetAPIVersion() &&
				v.Kind == obj.GetKind() &&
				v.Name == obj.GetName() {
				for _, p := range v.JSONPath {
					newObj, err := unstructuredutil.Replace(manifests[i], p, "###DYNAMIC_FIELD###")
					if err != nil {
						return fmt.Errorf("failed to replace json path: %w", err)
					}
					manifests[i] = *newObj
				}
			}
		}
	}
	return nil
}

package charts

import (
	"encoding/base64"
	"fmt"
	"os"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yaml "sigs.k8s.io/yaml/goyaml.v3"

	unst "github.com/jlandowner/helm-chartsnap/pkg/unstructured"
)

func LoadSnapshotConfig[T SnapshotValues | SnapshotConfig](filePath string, out *T) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(out)
	if err != nil {
		return fmt.Errorf("failed to decode file '%s': %w", filePath, err)
	}
	return nil
}

type SnapshotValues struct {
	TestSpec SnapshotConfig `yaml:"testSpec,omitempty"`
}

type SnapshotConfig struct {
	DynamicFields []ManifestPath `yaml:"dynamicFields,omitempty"`
}

type ManifestPath struct {
	Kind       string   `yaml:"kind,omitempty"`
	APIVersion string   `yaml:"apiVersion,omitempty"`
	Name       string   `yaml:"name,omitempty"`
	JSONPath   []string `yaml:"jsonPath,omitempty"`
	Base64     bool     `yaml:"base64,omitempty"`
}

const DynamicValue = "###DYNAMIC_FIELD###"

var Base64DynamicValue = base64.StdEncoding.EncodeToString([]byte(DynamicValue))

func (t *SnapshotConfig) ApplyFixedValue(manifests []metaV1.Unstructured) error {
	for _, v := range t.DynamicFields {
		for i, obj := range manifests {
			if v.APIVersion == obj.GetAPIVersion() &&
				v.Kind == obj.GetKind() &&
				v.Name == obj.GetName() {
				for _, p := range v.JSONPath {
					var value string
					if v.Base64 {
						value = Base64DynamicValue
					} else {
						value = DynamicValue
					}
					newObj, err := unst.Replace(manifests[i], p, value)
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

// Merge merges the snapshot configs into the current snapshot config
// The current snapshot config has higher priority than the given snapshot config
func (t *SnapshotConfig) Merge(cfg SnapshotConfig) {
	// For DynamicFields, it doesn't matter if the same field is replaced with a fixed value several times
	// But the current snapshot config has higher priority than the given snapshot config
	t.DynamicFields = append(cfg.DynamicFields, t.DynamicFields...)
}

package charts

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yaml "sigs.k8s.io/yaml/goyaml.v3"

	unstructuredutil "github.com/jlandowner/helm-chartsnap/pkg/unstructured"
)

func LoadSnapshotConfig(file string) (SnapshotConfig, error) {
	cfg := SnapshotConfig{}
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // ignore not found
		} else {
			return cfg, fmt.Errorf("failed to open config file '%s': %w", file, err)
		}
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to decode config file '%s': %w", file, err)
	}
	return cfg, nil
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
}

func (t *SnapshotConfig) ApplyFixedValue(manifests []unstructured.Unstructured) error {
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

func (t *SnapshotConfig) Merge(cfg SnapshotConfig) {
	// dynamic fields
	// It doesn't matter if the same field is replaced with a fixed value several times, so just append and not consider duplication.
	t.DynamicFields = append(t.DynamicFields, cfg.DynamicFields...)
}

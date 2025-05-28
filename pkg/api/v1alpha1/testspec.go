package v1alpha1

import (
	"encoding/base64"
	"fmt"
	"os"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

func FromFile[T SnapshotValues | SnapshotConfig](filePath string, out *T) error {
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
	DynamicFields   []ManifestPath `yaml:"dynamicFields,omitempty"`
	SnapshotFileExt string         `yaml:"snapshotFileExt,omitempty"`
	SnapshotVersion string         `yaml:"snapshotVersion,omitempty"`
}

type ManifestPath struct {
	Kind       string   `yaml:"kind,omitempty"`
	APIVersion string       `yaml:"apiVersion,omitempty"`
	Name       string       `yaml:"name,omitempty"`
	JSONPath   JSONPathList `yaml:"jsonPath,omitempty"`
	Base64     bool         `yaml:"base64,omitempty"`
}

// JSONPathItem is a struct that holds a JSON path and a value.
type JSONPathItem struct {
	Path  string `yaml:"path"`
	Value string `yaml:"value,omitempty"`
}

// JSONPathList is a slice of JSONPathItem.
type JSONPathList []JSONPathItem

// UnmarshalYAML implements the yaml.Unmarshaler interface.
// It allows JSONPathList to be unmarshalled from a list of JSONPathItem objects
// or a list of strings.
func (j *JSONPathList) UnmarshalYAML(value *yaml.Node) error {
	var items []JSONPathItem
	if err := value.Decode(&items); err == nil {
		*j = items
		return nil
	}

	var paths []string
	if err := value.Decode(&paths); err == nil {
		*j = make([]JSONPathItem, len(paths))
		for i, path := range paths {
			(*j)[i] = JSONPathItem{Path: path}
		}
		return nil
	}

	return fmt.Errorf("failed to unmarshal JSONPathList: expected a list of JSONPathItem objects or a list of strings")
}

const DynamicValue = "###DYNAMIC_FIELD###"

var Base64DynamicValue = base64.StdEncoding.EncodeToString([]byte(DynamicValue))

func (v *ManifestPath) DynamicValue() string {
	if v.Base64 {
		return Base64DynamicValue
	} else {
		return DynamicValue
	}
}

// Merge merges the snapshot configs into the current snapshot config
// The current snapshot config has higher priority than the given snapshot config
func (t *SnapshotConfig) Merge(cfg SnapshotConfig) {
	// For DynamicFields, it doesn't matter if the same field is replaced with a fixed value several times
	// But the current snapshot config has higher priority than the given snapshot config
	t.DynamicFields = append(cfg.DynamicFields, t.DynamicFields...)
	if cfg.SnapshotFileExt != "" {
		t.SnapshotFileExt = cfg.SnapshotFileExt
	}
	if cfg.SnapshotVersion != "" {
		t.SnapshotVersion = cfg.SnapshotVersion
	}
}

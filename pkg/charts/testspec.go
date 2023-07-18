package charts

type SnapshotValues struct {
	TestSpec TestSpec `yaml:"testSpec,omitempty"`
}

type TestSpec struct {
	Description   string         `yaml:"desc,omitempty"`
	DynamicFields []ManifestPath `yaml:"dynamicFields,omitempty"`
}

type ManifestPath struct {
	Kind       string   `yaml:"kind,omitempty"`
	APIVersion string   `yaml:"apiVersion,omitempty"`
	Name       string   `yaml:"name,omitempty"`
	JSONPath   []string `yaml:"jsonPath,omitempty"`
}

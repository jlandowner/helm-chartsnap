package v1alpha1

import (
	"testing"

	. "github.com/jlandowner/helm-chartsnap/pkg/snap/gomega"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCharts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Charts Suite")
}

var _ = Describe("TestSpec", func() {
	var _ = Describe("FromFile", func() {
		Context("when values.yaml has testSpec", func() {
			It("should load config", func() {
				var v SnapshotValues
				err := FromFile("testdata/testspec_values.yaml", &v)
				Expect(err).NotTo(HaveOccurred())
				Expect(v).To(MatchSnapShot())
			})
		})

		Context("when loading .chartsnap.yaml", func() {
			It("should load config", func() {
				var v SnapshotConfig
				err := FromFile("testdata/.chartsnap.yaml", &v)
				Expect(err).NotTo(HaveOccurred())
				Expect(v).To(MatchSnapShot())
			})
		})

		Context("when loading invalid yaml", func() {
			It("should not load config", func() {
				var v SnapshotValues
				err := FromFile("testdata/testspec_values_invalid.yaml", &v)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(MatchSnapShot())
			})
		})

		Context("when loading not found", func() {
			It("should not load config", func() {
				var v SnapshotValues
				err := FromFile("testdata/notfound.yaml", &v)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(MatchSnapShot())
			})
		})
	})

	var _ = Describe("Merge", func() {
		It("should merge dynamic fields", func() {
			cfg1 := SnapshotConfig{
				DynamicFields: []ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "service",
						Name:       "chartsnap-app1",
						JSONPath: []string{
							"/spec/ports/0/targetPort",
						},
					},
					{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "chartsnap-app1",
						JSONPath: []string{
							"/spec/ports/1/targetPort",
						},
					},
				},
			}
			cfg2 := SnapshotConfig{
				DynamicFields: []ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "service",
						Name:       "chartsnap-app1",
						JSONPath: []string{
							"/spec/ports/0/targetPort",
						},
					},
					{
						APIVersion: "v1",
						Kind:       "Pod",
						Name:       "chartsnap-app1-test-connection",
						JSONPath: []string{
							"/metadata/name",
						},
					},
				},
			}
			cfg1.Merge(cfg2)
			Expect(cfg1).To(MatchSnapShot())
		})

		It("should merge snapshot file extensions", func() {
			base := SnapshotConfig{
				SnapshotFileExt: "json",
			}
			overwrite := SnapshotConfig{
				SnapshotFileExt: "yaml",
			}
			base.Merge(overwrite)
			Expect(base.SnapshotFileExt).To(Equal("yaml"))
		})

		It("should override snapshot file extension if not set in current config", func() {
			base := SnapshotConfig{}
			overwrite := SnapshotConfig{
				SnapshotFileExt: "yaml",
			}
			base.Merge(overwrite)
			Expect(base.SnapshotFileExt).To(Equal("yaml"))
		})

		It("should merge snapshot version when target is empty", func() {
			base := SnapshotConfig{}
			overwrite := SnapshotConfig{
				SnapshotVersion: "v1",
			}
			base.Merge(overwrite)
			Expect(base.SnapshotVersion).To(Equal("v1"))
		})

		It("should override snapshot version when target has value", func() {
			base := SnapshotConfig{
				SnapshotVersion: "v2",
			}
			overwrite := SnapshotConfig{
				SnapshotVersion: "v1",
			}
			base.Merge(overwrite)
			Expect(base.SnapshotVersion).To(Equal("v1")) // overwrite の値で上書きされることを期待
		})
	})

})

func TestManifestPath_DynamicValue(t *testing.T) {
	type fields struct {
		Kind       string
		APIVersion string
		Name       string
		JSONPath   []string
		Base64     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test DynamicValue",
			fields: fields{
				Base64: false,
			},
			want: "###DYNAMIC_FIELD###",
		},
		{
			name: "Test DynamicValue Base64",
			fields: fields{
				Base64: true,
			},
			want: "IyMjRFlOQU1JQ19GSUVMRCMjIw==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ManifestPath{
				Kind:       tt.fields.Kind,
				APIVersion: tt.fields.APIVersion,
				Name:       tt.fields.Name,
				JSONPath:   tt.fields.JSONPath,
				Base64:     tt.fields.Base64,
			}
			if got := v.DynamicValue(); got != tt.want {
				t.Errorf("ManifestPath.DynamicValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

var _ = Describe("JSONPathList UnmarshalYAML", func() {
	// Helper struct to trigger the custom unmarshaler
	type TestHostStruct struct {
		JPList JSONPathList `yaml:"jsonPath"`
	}

	Context("when unmarshalling from a list of strings", func() {
		It("should correctly populate the JSONPathList", func() {
			yamlInput := `jsonPath: ["/data/foo", "/data/bar"]`
			var host TestHostStruct
			err := yaml.Unmarshal([]byte(yamlInput), &host)
			Expect(err).NotTo(HaveOccurred())
			Expect(host.JPList).To(Equal(JSONPathList{
				{Path: "/data/foo", Value: ""},
				{Path: "/data/bar", Value: ""},
			}))
		})
	})

	Context("when unmarshalling from a list of objects", func() {
		It("should correctly populate the JSONPathList", func() {
			yamlInput := `jsonPath: [{path: "/data/foo", value: "val1"}, {path: "/data/bar", value: "val2"}]`
			var host TestHostStruct
			err := yaml.Unmarshal([]byte(yamlInput), &host)
			Expect(err).NotTo(HaveOccurred())
			Expect(host.JPList).To(Equal(JSONPathList{
				{Path: "/data/foo", Value: "val1"},
				{Path: "/data/bar", Value: "val2"},
			}))
		})
	})

	Context("when unmarshalling from an empty list", func() {
		It("should result in an empty JSONPathList", func() {
			yamlInput := `jsonPath: []`
			var host TestHostStruct
			err := yaml.Unmarshal([]byte(yamlInput), &host)
			Expect(err).NotTo(HaveOccurred())
			Expect(host.JPList).To(BeEmpty())
		})
	})

	Context("when unmarshalling from a mixed list (string first)", func() {
		It("should return an error", func() {
			yamlInput := `jsonPath: ["/data/foo", {path: "/data/bar", value: "val2"}]`
			var host TestHostStruct
			err := yaml.Unmarshal([]byte(yamlInput), &host)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to unmarshal JSONPathList: expected a list of JSONPathItem objects or a list of strings"))
		})
	})

	Context("when unmarshalling from a mixed list (object first)", func() {
		It("should return an error", func() {
			yamlInput := `jsonPath: [{path: "/data/foo", value: "val1"}, "/data/bar"]`
			var host TestHostStruct
			err := yaml.Unmarshal([]byte(yamlInput), &host)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to unmarshal JSONPathList: expected a list of JSONPathItem objects or a list of strings"))
		})
	})

	Context("when input is a map instead of a list", func() {
		It("should return an error", func() {
			yamlInput := `jsonPath: {path: "/data/foo"}`
			var host TestHostStruct
			err := yaml.Unmarshal([]byte(yamlInput), &host)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to unmarshal JSONPathList: expected a list of JSONPathItem objects or a list of strings"))
		})
	})

	Context("when input is a string instead of a list", func() {
		It("should return an error", func() {
			yamlInput := `jsonPath: "/data/foo"`
			var host TestHostStruct
			err := yaml.Unmarshal([]byte(yamlInput), &host)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to unmarshal JSONPathList: expected a list of JSONPathItem objects or a list of strings"))
		})
	})
})

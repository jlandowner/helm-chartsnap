package yaml

import (
	"io"
	"os"
	"testing"

	. "github.com/jlandowner/helm-chartsnap/pkg/snap/gomega"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/jlandowner/helm-chartsnap/pkg/api/v1alpha1"
)

func TestYAML(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YAML Suite")
}

var _ = Describe("Decode & Encode", func() {
	It("should success", func() {
		// Decode
		manifests := load("testdata/input.yaml")

		// Encode again
		buf, err := Encode(manifests)
		Expect(err).NotTo(HaveOccurred())
		Ω(buf).Should(MatchSnapShot())
	})

	It("should success with converting invalid YAML format", func() {
		// Decode
		manifests := load("testdata/input2.yaml")

		// Encode again
		buf, err := Encode(manifests)
		Expect(err).NotTo(HaveOccurred())
		Ω(buf).Should(MatchSnapShot())
	})

	It("should success with converting ScalerNode", func() {
		// Decode
		manifests := load("testdata/input3.yaml")

		// Encode again
		buf, err := Encode(manifests)
		Expect(err).NotTo(HaveOccurred())
		Ω(buf).Should(MatchSnapShot())
	})
})

var _ = Describe("ApplyFixedValueToDynamicFieleds", func() {
	It("should replace specified fields", func() {
		cfg := v1alpha1.SnapshotConfig{
			DynamicFields: []v1alpha1.ManifestPath{
				{
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "chartsnap-app2",
					JSONPath: []string{
						"/spec/ports/0/targetPort",
					},
				},
				{
					APIVersion: "v2",
					Kind:       "Service",
					Name:       "chartsnap-app1",
					JSONPath: []string{
						"/spec/ports/0/targetPort",
					},
				},
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
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "chartsnap-app1",
					JSONPath: []string{
						"/metadata/labels/app.kubernetes.io~1version",
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
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "chartsnap-app1",
					JSONPath: []string{
						"/spec/template/spec/serviceAccountName",
					},
					Base64: true,
				},
			},
		}
		manifests := load("testdata/input.yaml")
		err := ApplyFixedValueToDynamicFieleds(cfg, manifests)
		Expect(err).NotTo(HaveOccurred())

		// Encode
		buf, err := Encode(manifests)
		Expect(err).NotTo(HaveOccurred())
		Ω(buf).Should(MatchSnapShot())
	})
})

var _ = Describe("ApplyFixedValueToDynamicFields with JSONPathList scenarios", func() {
	var (
		baseConfigMapYAML = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-cm
data:
  key1: value1
  key2: value2
  key3: value3
`
	)

	Context("Scenario 1: JSONPathList with paths only", func() {
		It("should replace specified fields with DynamicValue", func() {
			manifests, err := Decode([]byte(baseConfigMapYAML))
			Expect(err).NotTo(HaveOccurred())
			Expect(manifests).To(HaveLen(1))

			cfg := v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "test-cm",
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1"},
							{Path: "/metadata/name"},
						},
					},
				},
			}

			err = ApplyFixedValueToDynamicFieleds(cfg, manifests)
			Expect(err).NotTo(HaveOccurred())

			// Assertions
			dataKey1, err := manifests[0].Pipe(yaml.Lookup("data", "key1"))
			Expect(err).NotTo(HaveOccurred())
			Expect(dataKey1.YNode().Value).To(Equal(v1alpha1.DynamicValue))

			metadataName, err := manifests[0].Pipe(yaml.Lookup("metadata", "name"))
			Expect(err).NotTo(HaveOccurred())
			Expect(metadataName.YNode().Value).To(Equal(v1alpha1.DynamicValue))

			// Ensure other fields are untouched
			dataKey2, err := manifests[0].Pipe(yaml.Lookup("data", "key2"))
			Expect(err).NotTo(HaveOccurred())
			Expect(dataKey2.YNode().Value).To(Equal("value2"))
		})
	})

	Context("Scenario 1.1: JSONPathList with paths only and Base64", func() {
		It("should replace specified fields with Base64DynamicValue", func() {
			manifests, err := Decode([]byte(baseConfigMapYAML)) // Using CM for simplicity, imagine it's a Secret
			Expect(err).NotTo(HaveOccurred())
			Expect(manifests).To(HaveLen(1))

			cfg := v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "test-cm",
						Base64:     true,
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1"},
						},
					},
				},
			}

			err = ApplyFixedValueToDynamicFieleds(cfg, manifests)
			Expect(err).NotTo(HaveOccurred())

			dataKey1, err := manifests[0].Pipe(yaml.Lookup("data", "key1"))
			Expect(err).NotTo(HaveOccurred())
			Expect(dataKey1.YNode().Value).To(Equal(v1alpha1.Base64DynamicValue))
		})
	})

	Context("Scenario 2: JSONPathList with paths and values", func() {
		It("should replace specified fields with fixed values", func() {
			manifests, err := Decode([]byte(baseConfigMapYAML))
			Expect(err).NotTo(HaveOccurred())

			cfg := v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "test-cm",
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1", Value: "fixed-key1-value"},
							{Path: "/metadata/name", Value: "fixed-cm-name"},
						},
					},
				},
			}

			err = ApplyFixedValueToDynamicFieleds(cfg, manifests)
			Expect(err).NotTo(HaveOccurred())

			dataKey1, err := manifests[0].Pipe(yaml.Lookup("data", "key1"))
			Expect(err).NotTo(HaveOccurred())
			Expect(dataKey1.YNode().Value).To(Equal("fixed-key1-value"))

			metadataName, err := manifests[0].Pipe(yaml.Lookup("metadata", "name"))
			Expect(err).NotTo(HaveOccurred())
			Expect(metadataName.YNode().Value).To(Equal("fixed-cm-name"))
		})
	})

	Context("Scenario 3: JSONPathList with mixed paths-only and paths-with-values", func() {
		It("should replace fields correctly based on item definition", func() {
			manifests, err := Decode([]byte(baseConfigMapYAML))
			Expect(err).NotTo(HaveOccurred())

			cfg := v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "test-cm",
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1", Value: "fixed-key1-mixed"}, // Fixed value
							{Path: "/data/key2"},                            // Dynamic placeholder
							{Path: "/metadata/name", Value: "fixed-name-mixed"}, // Fixed value
						},
					},
				},
			}

			err = ApplyFixedValueToDynamicFieleds(cfg, manifests)
			Expect(err).NotTo(HaveOccurred())

			dataKey1, err := manifests[0].Pipe(yaml.Lookup("data", "key1"))
			Expect(err).NotTo(HaveOccurred())
			Expect(dataKey1.YNode().Value).To(Equal("fixed-key1-mixed"))

			dataKey2, err := manifests[0].Pipe(yaml.Lookup("data", "key2"))
			Expect(err).NotTo(HaveOccurred())
			Expect(dataKey2.YNode().Value).To(Equal(v1alpha1.DynamicValue))
			
			metadataName, err := manifests[0].Pipe(yaml.Lookup("metadata", "name"))
			Expect(err).NotTo(HaveOccurred())
			Expect(metadataName.YNode().Value).To(Equal("fixed-name-mixed"))

			// Check untouched field
			dataKey3, err := manifests[0].Pipe(yaml.Lookup("data", "key3"))
			Expect(err).NotTo(HaveOccurred())
			Expect(dataKey3.YNode().Value).To(Equal("value3"))
		})
	})

	Context("No matching manifest", func() {
		It("should not modify the manifest", func() {
			manifests, err := Decode([]byte(baseConfigMapYAML))
			Expect(err).NotTo(HaveOccurred())
			originalYAML, err := Encode(manifests) // Get original state
			Expect(err).NotTo(HaveOccurred())

			cfg := v1alpha1.SnapshotConfig{
				DynamicFields: []v1alpha1.ManifestPath{
					{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "non-existent-cm", // Does not match "test-cm"
						JSONPath: v1alpha1.JSONPathList{
							{Path: "/data/key1", Value: "should-not-apply"},
						},
					},
				},
			}

			err = ApplyFixedValueToDynamicFieleds(cfg, manifests)
			Expect(err).NotTo(HaveOccurred())

			modifiedYAML, err := Encode(manifests)
			Expect(err).NotTo(HaveOccurred())
			Expect(modifiedYAML).To(Equal(originalYAML)) // Content should be unchanged
		})
	})
})

var _ = Describe("ConvertToUnknownNode", func() {
	DescribeTable("should convert nodes appropriately based on apiVersion and kind",
		func(name string, input string, expected bool) {
			node, err := yaml.Parse(input)
			Expect(err).NotTo(HaveOccurred())

			docs := []*yaml.RNode{node}
			err = convertToUnknownNode(docs)
			Expect(err).NotTo(HaveOccurred())

			// Check if the document was converted to Unknown
			if expected {
				// Should be converted to Unknown
				Expect(docs[0].GetKind()).To(Equal(v1alpha1.UnknownKind))
				Expect(docs[0].GetApiVersion()).To(Equal(v1alpha1.GroupVersion.String()))
				Expect(docs[0].ShouldKeep).To(BeTrue())
			} else {
				// Should remain unchanged
				Expect(docs[0].GetKind()).NotTo(Equal(v1alpha1.UnknownKind))
				Expect(docs[0].GetApiVersion()).NotTo(Equal(v1alpha1.GroupVersion.String()))
			}
		},
		Entry("normal yaml with apiVersion and kind", "normal yaml with apiVersion and kind", `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
`, false),
		Entry("yaml without apiVersion and kind", "yaml without apiVersion and kind", `
foo: bar
baz: qux
`, true),
		Entry("string value", "string value", `
just a string
`, true),
	)
})

func load(filePath string) []*yaml.RNode {
	f, err := os.Open(filePath)
	Expect(err).NotTo(HaveOccurred())
	defer f.Close()

	buf, err := io.ReadAll(f)
	Expect(err).NotTo(HaveOccurred())

	manifests, err := Decode(buf)
	Expect(err).NotTo(HaveOccurred())

	return manifests
}

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

var _ = Describe("SetLogger", func() {
	It("should set the logger", func() {
		// Test that SetLogger sets the logger
		SetLogger(nil)
		// logger should be nil when passing nil
	})
})

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

	DescribeTable("wildcard matching",
		func(df []v1alpha1.ManifestPath) {
			cfg := v1alpha1.SnapshotConfig{
				DynamicFields: df,
			}
			manifests := load("testdata/input4.yaml")
			err := ApplyFixedValueToDynamicFieleds(cfg, manifests)
			Expect(err).NotTo(HaveOccurred())

			// Encode
			buf, err := Encode(manifests)
			Expect(err).NotTo(HaveOccurred())
			Ω(buf).Should(MatchSnapShot())
		},
		Entry("wildcard for all matches any resource", []v1alpha1.ManifestPath{
			{
				// APIVersion: "",
				// Kind: "",
				// Name: "",
				JSONPath: []string{
					"/data/ca.crt",
					"/data/tls.crt",
					"/data/tls.key",
				},
			},
			{
				JSONPath: []string{
					"/metadata/labels/app.kubernetes.io~1version",
				},
			},
		}),
		Entry("to all resources which kind='DummySecret'", []v1alpha1.ManifestPath{
			{
				// APIVersion: "",
				Kind: "DummySecret",
				// Name: "",
				JSONPath: []string{
					"/data/ca.crt",
					"/data/tls.crt",
					"/data/tls.key",
				},
			},
			{
				JSONPath: []string{
					"/metadata/labels/app.kubernetes.io~1version",
				},
			},
		}),
		Entry("to all name is 'test-secret'", []v1alpha1.ManifestPath{
			{
				// APIVersion: "",
				// Kind: "",
				Name: "test-secret",
				JSONPath: []string{
					"/data/ca.crt",
					"/data/tls.crt",
					"/data/tls.key",
				},
			},
			{
				JSONPath: []string{
					"/metadata/labels/app.kubernetes.io~1version",
				},
			},
		}),
		Entry("to all resources which apiVersion='dummy/v1'", []v1alpha1.ManifestPath{
			{
				APIVersion: "dummy/v1",
				// Kind: "",
				// Name: "",
				JSONPath: []string{
					"/data/ca.crt",
					"/data/tls.crt",
					"/data/tls.key",
				},
			},
			{
				JSONPath: []string{
					"/metadata/labels/app.kubernetes.io~1version",
				},
			},
		}),
		Entry("to all resources which apiVersion='v1' and kind='Secret'", []v1alpha1.ManifestPath{
			{
				APIVersion: "v1",
				Kind:       "Secret",
				// Name: "",
				JSONPath: []string{
					"/data/ca.crt",
					"/data/tls.crt",
					"/data/tls.key",
				},
			},
		}),
		Entry("to all resources which apiVersion='dummy/v1' and name='test-secret'", []v1alpha1.ManifestPath{
			{
				APIVersion: "dummy/v1",
				// Kind: "",
				Name: "test-secret",
				JSONPath: []string{
					"/data/ca.crt",
					"/data/tls.crt",
					"/data/tls.key",
				},
			},
			{
				JSONPath: []string{
					"/metadata/labels/app.kubernetes.io~1version",
				},
			},
		}),
		Entry("to all resources which kind='Secret' and name='test-secret'", []v1alpha1.ManifestPath{
			{
				// APIVersion: "",
				Kind: "Secret",
				Name: "test-secret",
				JSONPath: []string{
					"/data/ca.crt",
					"/data/tls.crt",
					"/data/tls.key",
				},
			},
			{
				JSONPath: []string{
					"/metadata/labels/app.kubernetes.io~1version",
				},
			},
		}),
	)
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

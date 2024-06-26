package unstructured

import (
	"io"
	"os"
	"testing"

	. "github.com/jlandowner/helm-chartsnap/pkg/snap/gomega"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/jlandowner/helm-chartsnap/pkg/api/v1alpha1"
)

func TestUnstructured(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unstructured Suite")
}

var _ = Describe("Diff", func() {
	Context("DiffContextLineN is 3", func() {
		It("should return the extracted diff with previous/next 3 lines", func() {
			expectedSnap := mustReadFile("testdata/expected.snap")
			actualSnap := mustReadFile("testdata/actual.snap")

			d := DiffOptions{
				ContextLineN: 3,
			}
			diff := d.Diff(expectedSnap, actualSnap)
			Ω(diff).To(MatchSnapShot())
		})
	})

	Context("DiffContextLineN is 0", func() {
		It("should return all diff", func() {
			expectedSnap := mustReadFile("testdata/expected.snap")
			actualSnap := mustReadFile("testdata/actual.snap")

			d := DiffOptions{
				ContextLineN: 0,
			}
			diff := d.Diff(expectedSnap, actualSnap)
			Ω(diff).To(MatchSnapShot())
		})
	})
})

var _ = Describe("Unknown", func() {
	Context("OK", func() {
		It("report unknown as warning", func() {
			raw := `some: raw data
raw:
  data: here`
			err := v1alpha1.NewUnknownError(raw)

			Ω(err.Error()).To(MatchSnapShot())
		})
	})
})

var _ = Describe("ApplyDynamicFields", func() {
	load := func(filePath string) []metaV1.Unstructured {
		f, err := os.Open(filePath)
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()

		buf, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())

		manifests, errs := Decode(string(buf))
		Expect(len(errs)).To(BeZero())

		return manifests
	}

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
		manifests := load("testdata/testspec_test.yaml")
		err := ApplyFixedValue(cfg, manifests)
		Expect(err).NotTo(HaveOccurred())
		Expect(manifests).To(MatchSnapShot())
	})
})

func mustReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

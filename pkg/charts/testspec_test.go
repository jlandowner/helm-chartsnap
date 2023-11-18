package charts

import (
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jlandowner/helm-chartsnap/pkg/snap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	unstructuredutil "github.com/jlandowner/helm-chartsnap/pkg/unstructured"
)

func TestSnapshot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testspec Suite")
}

var _ = Describe("ApplyDynamicFields", func() {
	load := func(filePath string) []unstructured.Unstructured {
		f, err := os.Open(filePath)
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()

		buf, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())

		manifests, err := unstructuredutil.Decode(string(buf))
		Expect(err).NotTo(HaveOccurred())

		return manifests
	}

	It("should replace specified fields", func() {
		testSpec := &TestSpec{
			DynamicFields: []ManifestPath{
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
			},
		}
		manifests := load("testdata/testspec_test.yaml")
		err := testSpec.ApplyFixedValue(manifests)
		Expect(err).NotTo(HaveOccurred())
		Expect(manifests).To(snap.MatchSnapShot())
	})
})

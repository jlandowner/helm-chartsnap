package gomega

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jlandowner/helm-chartsnap/pkg/unstructured"
)

func TestSnap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snap Suite")
}

var _ = Describe("Snap", func() {
	It("takes a full snapshot", func() {
		b, err := os.ReadFile("testdata/pod.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(b)).To(MatchSnapShot())
	})

	It("takes a snapshot without dynamic values", func() {
		b, err := os.ReadFile("testdata/pod.yaml")
		Expect(err).NotTo(HaveOccurred())

		_, obj, err := unstructured.BytesToUnstructured(b)
		Expect(err).NotTo(HaveOccurred())

		Expect(ObjectSnapshot(obj)).To(MatchSnapShot())
	})
})

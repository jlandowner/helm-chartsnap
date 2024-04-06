package unstructured

import (
	"os"
	"testing"

	. "github.com/jlandowner/helm-chartsnap/pkg/snap/gomega"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
			err := NewUnknownError(raw)

			Ω(err.Error()).To(MatchSnapShot())
		})
	})
})

func mustReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

package snap

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Snapshot", func() {
	AfterEach(func() {
		os.Remove("__snapshots__/snapshot_create_test.snap")
	})

	It("should match", func() {
		testdata := `
Helm is a tool for managing Charts. Charts are packages of pre-configured Kubernetes resources.

Use Helm to:
		
- Find and use popular software packaged as Helm Charts to run in Kubernetes
- Share your own applications as Helm Charts
- Create reproducible builds of your Kubernetes applications
- Intelligently manage your Kubernetes manifest files
- Manage releases of Helm packages
`
		// create snapshot file
		Expect(testdata).To(MatchSnapShot())
	})
})

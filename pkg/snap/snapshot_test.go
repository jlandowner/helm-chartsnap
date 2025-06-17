package snap

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	pwd string
)

func TestSnap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snap Suite")
}

var _ = Describe("Snap", func() {
	BeforeEach(func() {
		_pwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		pwd = _pwd
	})

	AfterEach(func() {
		logger = nil
	})

	Describe("SetLogger", func() {
		It("should set the logger", func() {
			// Test that SetLogger sets the logger
			SetLogger(nil)
			// Just test that the function can be called without error
		})
	})

	Describe("SnapshotMatcher", func() {
		Context("snapshot file does not exist", func() {
			It("should pass and create snapshot file", func() {
				var (
					snapFile     = "not_found.snap"
					snapFilePath = filepath.Join(pwd, "__snapshot__", snapFile)
					snapData     = `
Beautiful is better than ugly.
Explicit is better than implicit.
Simple is better than complex.
Complex is better than complicated.
Flat is better than nested.
Sparse is better than dense.
Readability counts.
Special cases aren't special enough to break the rules.
Although practicality beats purity.
Errors should never pass silently.
Unless explicitly silenced.
In the face of ambiguity, refuse the temptation to guess.
There should be one-- and preferably only one --obvious way to do it.
Although that way may not be obvious at first unless you're Dutch.
Now is better than never.
Although never is often better than *right* now.
If the implementation is hard to explain, it's a bad idea.
If the implementation is easy to explain, it may be a good idea.
Namespaces are one honking great idea -- let's do more of those!
`
				)
				defer os.Remove(snapFilePath)

				_, err := os.Stat(snapFilePath)
				Expect(os.IsNotExist(err)).To(BeTrue())

				matcher := SnapshotMatcher(snapFilePath)
				success, err := matcher.Match(snapData)
				Expect(err).NotTo(HaveOccurred())
				Expect(success).To(BeTrue())

				fileContent, err := os.ReadFile(snapFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(fileContent)).To(Equal(snapData))
			})
		})

		Context("snapshot file exist", func() {
			It("match snapshot", func() {
				var (
					snapFile     = "single.snap"
					snapFilePath = filepath.Join(pwd, "__snapshot__", snapFile)
				)

				_, err := os.Stat(snapFilePath)
				Expect(err).NotTo(HaveOccurred())

				fileContent, err := os.ReadFile(snapFilePath)
				Expect(err).NotTo(HaveOccurred())

				matcher := SnapshotMatcher(snapFilePath)
				success, err := matcher.Match(fileContent)
				Expect(err).NotTo(HaveOccurred())
				Expect(success).To(BeTrue())
			})
		})

		Context("multi-formatted snapshot file exist", func() {
			It("match snapshot", func() {
				var (
					singleSnapFile     = "single.snap"
					singleSnapFilePath = filepath.Join(pwd, "__snapshot__", singleSnapFile)

					snapFile     = "multi.snap"
					snapFilePath = filepath.Join(pwd, "__snapshot__", snapFile)
				)

				_, err := os.Stat(snapFilePath)
				Expect(err).NotTo(HaveOccurred())

				// file content is the same as single.snap
				fileContent, err := os.ReadFile(singleSnapFilePath)
				Expect(err).NotTo(HaveOccurred())

				matcher := SnapshotMatcher(snapFilePath, WithSnapshotID("default"))
				success, err := matcher.Match(fileContent)
				Expect(err).NotTo(HaveOccurred())
				Expect(success).To(BeTrue())
			})
		})

		Context("multi-formatted snapshot file", func() {
			It("append new snapID", func() {
				snapFilePath := filepath.Join(pwd, "__snapshot__", "multi-not-found.snap")
				existingFileData := `[xxx]
SnapShot = """
apiVersion: v1
kind: Namespace
metadata:
  name: default
"""
`
				yyySnapshot := `apiVersion: v1
kind: Pod
metadata:
  name: nginx
`
				expectedFileData := `[xxx]
SnapShot = """
apiVersion: v1
kind: Namespace
metadata:
  name: default
"""

[yyy]
SnapShot = """
apiVersion: v1
kind: Pod
metadata:
  name: nginx
"""
`

				_, err := os.Stat(snapFilePath)
				Expect(os.IsNotExist(err)).To(BeTrue())

				os.WriteFile(snapFilePath, []byte(existingFileData), 0644)
				defer os.Remove(snapFilePath)

				matcher := SnapshotMatcher(snapFilePath, WithSnapshotID("yyy"))
				success, err := matcher.Match(yyySnapshot)
				Expect(err).NotTo(HaveOccurred())
				Expect(success).To(BeTrue())

				fileContent, err := os.ReadFile(snapFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(fileContent)).To(Equal(expectedFileData))
			})
		})
	})

	Context("json snapshot file", func() {
		It("match snapshot", func() {
			testStruct := struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{Name: "John", Age: 30}

			snapFilePath := filepath.Join(pwd, "__snapshot__", "json.snap")

			// _, err := os.Stat(snapFilePath)
			// Expect(err).NotTo(HaveOccurred())

			matcher := SnapshotMatcher(snapFilePath, WithSnapshotID("json"))
			success, err := matcher.Match(testStruct)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(matcher.FailureMessage(nil))
			Expect(success).To(BeTrue())
		})
	})

	Context("NegatedFailureMessage", func() {
		It("should return negated failure message", func() {
			snapFilePath := filepath.Join(pwd, "__snapshot__", "test.snap")
			matcher := SnapshotMatcher(snapFilePath)

			message := matcher.NegatedFailureMessage(nil)
			Expect(message).To(ContainSubstring("Expected"))
		})
	})
})

func TestIsMultiSnapshots(t *testing.T) {
	// table test
	tests := []struct {
		name     string
		snapFile string
		want     bool
	}{
		{
			name:     "single snapshot",
			snapFile: "single.snap",
			want:     false,
		},
		{
			name:     "multi snapshot",
			snapFile: "multi.snap",
			want:     true,
		},
	}
	for _, tt := range tests {
		got := IsMultiSnapshots(filepath.Join(pwd, "__snapshot__", tt.snapFile))
		if got != tt.want {
			t.Errorf("Expected file to contain multiple snapshots got %v, want %v", got, tt.want)
		}
	}
}

func TestDiffFunc(t *testing.T) {
	// Set up test data
	x := "test1"
	y := "test2"

	// Define the expected diff function
	expectedDiffFunc := func(x, y string) string {
		return "diff"
	}

	// Create a snapshot matcher with the diff function
	matcher := SnapshotMatcher("test_snap.toml", WithDiffFunc(expectedDiffFunc))

	// Get the actual diff function
	actualDiffFunc := matcher.diffFunc

	// Check the result
	if actualDiffFunc(x, y) != expectedDiffFunc(x, y) {
		t.Errorf("Diff function does not match expected")
	}
}

package snap

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CacheFS", func() {
	Context("CacheFS", func() {
		It("WriteFile", func() {
			snapTempDir := filepath.Join(os.TempDir(), "__snapshot__")
			os.RemoveAll(snapTempDir)

			filePath := filepath.Join(snapTempDir, generateRandomFileName())
			defer os.RemoveAll(snapTempDir)

			data := []byte("Hello, World! 0")

			err := WriteFile(filePath, data)
			Expect(err).NotTo(HaveOccurred())

			fileContent, err := os.ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileContent).To(Equal(data))
		})

		It("ReadFile", func() {
			tmpDir := createTempDir()
			defer os.RemoveAll(tmpDir)

			filePath := filepath.Join(tmpDir, "test.txt")
			data := []byte("Hello, World! 1")

			err := os.WriteFile(filePath, data, 0644)
			Expect(err).NotTo(HaveOccurred())

			fileContent, err := ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileContent).To(Equal(data))

			// Remove the file directly
			err = os.Remove(filePath)
			Expect(err).NotTo(HaveOccurred())

			// Verify that the file has been removed
			_, err = os.Stat(filePath)
			Expect(os.IsNotExist(err)).To(BeTrue())

			// Verify that cacheFs still keeps the file
			fileContent, err = ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileContent).To(Equal(data))
		})

		It("RemoveFile", func() {
			tmpDir := createTempDir()
			defer os.RemoveAll(tmpDir)

			filePath := filepath.Join(tmpDir, "test.txt")
			data := []byte("Hello, World! 2")

			err := os.WriteFile(filePath, data, 0644)
			Expect(err).NotTo(HaveOccurred())

			fileContent, err := ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileContent).To(Equal(data))

			// Remove the file via cacheFs
			err = RemoveFile(filePath)
			Expect(err).NotTo(HaveOccurred())

			// Verify that the file has been removed
			_, err = os.Stat(filePath)
			Expect(os.IsNotExist(err)).To(BeTrue())

			// Verify that the file has been removed from cacheFs
			_, err = ReadFile(filePath)
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

		It("WriteFile with permission error", func() {
			// Try to write to a directory that doesn't exist and can't be created
			invalidPath := "/root/nonexistent/test.txt"
			data := []byte("test data")

			err := WriteFile(invalidPath, data)
			// Should handle error gracefully even if directory creation fails
			// The exact behavior depends on the system, but we expect it to handle errors
			Expect(err).To(HaveOccurred())
		})

		It("ReadFile with non-existent file", func() {
			nonExistentPath := "/tmp/non-existent-file.txt"

			_, err := ReadFile(nonExistentPath)
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

	})
})

// Helper function to create a temporary directory for testing
func createTempDir() string {
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		panic(err)
	}
	return tmpDir
}

func generateRandomFileName() string {
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(randBytes)
}

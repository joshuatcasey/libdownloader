package libdownloader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joshuatcasey/libdownloader"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testSimpleDownloadedFile(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		dir, path string

		downloadedFile libdownloader.DownloadedFile
	)

	it.Before(func() {
		dir = t.TempDir()
		path = filepath.Join(dir, "original_name.txt")
		err := os.WriteFile(path, []byte("contents"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		downloadedFile = libdownloader.NewSimpleDownloadedFile(path, "contents")
	})

	context("SimpleDownloadedFile", func() {
		context("Path", func() {
			it("will return the path", func() {
				Expect(downloadedFile.Path()).To(Equal(path))
			})
		})

		context("Contents", func() {
			it("will return the contents", func() {
				Expect(downloadedFile.Contents()).To(Equal("contents"))
			})

			context("when given content is empty", func() {
				it.Before(func() {
					path = filepath.Join(dir, "contents.txt")
					err := os.WriteFile(path, []byte("contents of file"), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())

					downloadedFile = libdownloader.NewSimpleDownloadedFile(path, "")
				})

				it("will read from the file", func() {
					Expect(downloadedFile.Contents()).To(Equal("contents of file"))
				})
			})
		})

		context("SHA256", func() {
			it("will retrieve the SHA256", func() {
				sha256, err := downloadedFile.Sha256()
				Expect(err).NotTo(HaveOccurred())
				Expect(sha256).To(Equal("d1b2a59fbea7e20077af9f91b27e95e865061b270be03ff539ab3b73587882e8"))
			})

			it("will cache the SHA256", func() {
				sha256, err := downloadedFile.Sha256()
				Expect(err).NotTo(HaveOccurred())
				Expect(sha256).To(Equal("d1b2a59fbea7e20077af9f91b27e95e865061b270be03ff539ab3b73587882e8"))

				err = os.RemoveAll(dir)
				Expect(err).NotTo(HaveOccurred())

				sha256, err = downloadedFile.Sha256()
				Expect(err).NotTo(HaveOccurred())
				Expect(sha256).To(Equal("d1b2a59fbea7e20077af9f91b27e95e865061b270be03ff539ab3b73587882e8"))
			})

			it("will use the file contents, not the given contents", func() {
				newFile := filepath.Join(dir, "new_file.txt")
				err := os.WriteFile(newFile, []byte("file contents"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				downloadedFile = libdownloader.NewSimpleDownloadedFile(newFile, "given contents")

				sha256, err := downloadedFile.Sha256()
				Expect(err).NotTo(HaveOccurred())

				Expect(sha256).To(Equal("7bb6f9f7a47a63e684925af3608c059edcc371eb81188c48c9714896fb1091fd"))
			})

			context("failure cases", func() {
				it("will error if the file does not exist", func() {
					downloadedFile = libdownloader.NewSimpleDownloadedFile("", "")

					_, err := downloadedFile.Sha256()
					Expect(err).To(MatchError(os.ErrNotExist))
				})
			})
		})

		context("Cleanup", func() {
			it("will cleanup the file and local variables", func() {
				err := downloadedFile.Cleanup()
				Expect(err).NotTo(HaveOccurred())

				Expect(path).NotTo(BeAnExistingFile())

				Expect(downloadedFile.Path()).To(Equal(""))

				_, err = downloadedFile.Contents()
				Expect(err).To(MatchError(os.ErrNotExist))

				_, err = downloadedFile.Sha256()
				Expect(err).To(MatchError(os.ErrNotExist))
			})
		})
	})
}

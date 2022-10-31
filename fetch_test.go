package libdownloader_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/joshuatcasey/libdownloader"
	"github.com/joshuatcasey/libdownloader/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testFetch(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		server    *httptest.Server
		serverURL *url.URL
	)

	it.Before(func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method != http.MethodGet {
				http.Error(w, "NotFound", http.StatusNotFound)
				return
			}

			switch req.URL.Path {
			case "/filename":
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "contents of download")
			default:
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "unknown download")
			}
		}))

		var err error
		serverURL, err = serverURL.Parse(server.URL)
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		server.Close()
	})

	context("Fetch", func() {
		it("will fetch /filename and return its checksum", func() {
			file, err := libdownloader.Fetch(serverURL.JoinPath("filename").String())
			Expect(err).NotTo(HaveOccurred())

			Expect(file.Contents()).To(Equal("contents of download"))
			Expect(file.Sha256()).To(Equal("b7d7817b26898ad5f5fada279bd34b9d7b29a04bc5e42d41e6151333a2be8c2b"))
			Expect(filepath.Base(file.Path())).To(Equal("filename"))
		})

		it("will fetch /unknown and return its checksum", func() {
			file, err := libdownloader.Fetch(serverURL.JoinPath("unknown").String())
			Expect(err).NotTo(HaveOccurred())

			Expect(file.Contents()).To(Equal("unknown download"))
			Expect(file.Sha256()).To(Equal("48b0c3096c357b5800e09a4420a2b0c20076a12364e1586d4ebc7d954ec04214"))
			Expect(file.Path()).To(HaveSuffix("unknown"))
		})

		context("With Options", func() {
			context("WithFilename", func() {
				it("will use the given filename", func() {
					file, err := libdownloader.Fetch(
						serverURL.JoinPath("filename").String(),
						libdownloader.WithFilename("hello"))
					Expect(err).NotTo(HaveOccurred())

					Expect(filepath.Base(file.Path())).To(Equal("hello"))
				})
			})

			context("WithHttpClient", func() {
				var httpClient *fakes.SimpleHttpClient

				it.Before(func() {
					httpClient = &fakes.SimpleHttpClient{}
					httpClient.GetCall.Returns.Error = errors.New("bad client")
				})

				it("will use the given client", func() {
					_, err := libdownloader.Fetch(
						serverURL.JoinPath("redirect").String(),
						libdownloader.WithHttpClient(httpClient))
					Expect(err.Error()).To(Equal("could not get url: bad client"))
				})
			})

			context("WithDownloadDirectory", func() {
				var downloadDir string

				it.Before(func() {
					downloadDir = t.TempDir()
				})

				it("will use the given download directory", func() {
					file, err := libdownloader.Fetch(
						serverURL.JoinPath("filename").String(),
						libdownloader.WithDownloadDirectory(downloadDir))
					Expect(err).NotTo(HaveOccurred())

					Expect(file.Path()).To(Equal(filepath.Join(downloadDir, "filename")))
				})

				context("failure cases", func() {
					context("when the given dir is not writeable", func() {
						it.Before(func() {
							err := os.Chmod(downloadDir, 0000)
							Expect(err).NotTo(HaveOccurred())
						})

						it("returns error", func() {
							_, err := libdownloader.Fetch(
								serverURL.JoinPath("filename").String(),
								libdownloader.WithDownloadDirectory(downloadDir))

							Expect(err.Error()).To(ContainSubstring("could not write to file: open %s: permission denied", filepath.Join(downloadDir, "filename")))
						})
					})
				})
			})
		})
	})

	context("failure cases", func() {
		it("returns error when file not found", func() {
			_, err := libdownloader.Fetch("https://4877be5d01cc44c5b2ea236f476da950.com")
			Expect(err.Error()).To(ContainSubstring(`could not get url: Get "https://4877be5d01cc44c5b2ea236f476da950.com": dial tcp: lookup 4877be5d01cc44c5b2ea236f476da950.com: no such host`))
		})
	})
}

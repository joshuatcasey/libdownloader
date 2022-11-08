package libdownloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// SimpleHttpClient is implemented by http.Client.
// It exists for purposes of testing.
//
//go:generate faux --interface simpleHttpClient --output fakes/simple_http_client.go
type simpleHttpClient interface {
	Get(url string) (*http.Response, error)
}

type fetchConfig struct {
	downloadDir string
	filename    string
	httpClient  simpleHttpClient
}

type Option func(*fetchConfig)

// WithHttpClient allows the consumer to provide an HTTP client for fetching the remote resource.
// Potentially useful for timeout and redirect specification.
// Defaults to http.DefaultClient.
func WithHttpClient(httpClient simpleHttpClient) Option {
	return func(c *fetchConfig) {
		c.httpClient = httpClient
	}
}

// WithFilename allows the consumer to specify just the filename.
// Defaults to filepath.Base of the given url.
func WithFilename(filename string) Option {
	return func(c *fetchConfig) {
		c.filename = filename
	}
}

// WithDownloadDirectory allows specification of where the file should be placed.
// Defaults to os.MkdirTemp.
func WithDownloadDirectory(downloadDir string) Option {
	return func(c *fetchConfig) {
		c.downloadDir = downloadDir
	}
}

// Fetch will retrieve a file from the internet and download it to the local filesystem.
// It returns a DownloadedFile through which consumers can access the path, contents, and checksum.
// Use the Option list to modify where and how the file is downloaded.
func Fetch(url string, options ...Option) (DownloadedFile, error) {
	config := fetchConfig{
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(&config)
	}

	response, err := config.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get url: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("could not get url %s, with status code %d", url, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}

	if config.downloadDir == "" {
		config.downloadDir, err = os.MkdirTemp("", "")
		if err != nil {
			return nil, errors.New("could not create a temp dir")
		}
	}

	filename := config.filename
	if filename == "" {
		filename = filepath.Base(url)
	}

	filePath := filepath.Join(config.downloadDir, filename)

	err = os.WriteFile(filePath, body, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not write to file: %w", err)
	}

	return NewSimpleDownloadedFile(filePath, string(body)), nil
}

package libdownloader

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type DownloadedFile interface {
	// Path returns the fully-qualified path to the downloaded file
	Path() string

	// Contents returns the entirety of the file contents.
	// Implementations can choose whether to cache the contents for easy retrieval.
	Contents() (string, error)

	// Sha256 returns the SHA256 for the file contents.
	Sha256() (string, error)

	// Cleanup is meant to perform any desired operations before releasing memory.
	// This could include removing cached contents or deleting the file.
	Cleanup() error
}

type simpleDownloadedFile struct {
	path     string
	contents string
	sha256   string
}

func NewSimpleDownloadedFile(path, contents string) *simpleDownloadedFile {
	return &simpleDownloadedFile{
		path:     path,
		contents: contents,
	}
}

func (simpleDownloadedFile *simpleDownloadedFile) Contents() (string, error) {
	if simpleDownloadedFile.contents != "" {
		return simpleDownloadedFile.contents, nil
	}

	bytes, err := os.ReadFile(simpleDownloadedFile.path)
	if err != nil {
		return "", err
	}
	simpleDownloadedFile.contents = string(bytes)
	return simpleDownloadedFile.contents, nil
}

func (simpleDownloadedFile *simpleDownloadedFile) Sha256() (string, error) {
	if simpleDownloadedFile.sha256 != "" {
		return simpleDownloadedFile.sha256, nil
	}

	file, err := os.Open(simpleDownloadedFile.path)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	simpleDownloadedFile.sha256 = hex.EncodeToString(hash.Sum(nil))

	return simpleDownloadedFile.sha256, nil
}

func (simpleDownloadedFile *simpleDownloadedFile) Path() string {
	return simpleDownloadedFile.path
}

func (simpleDownloadedFile *simpleDownloadedFile) Cleanup() error {
	err := os.RemoveAll(simpleDownloadedFile.path)
	if err != nil {
		return err
	}

	simpleDownloadedFile.path = ""
	simpleDownloadedFile.sha256 = ""
	simpleDownloadedFile.contents = ""
	return nil
}

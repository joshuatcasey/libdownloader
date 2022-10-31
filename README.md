# Libdownloader

Exists to provide a convenient way to download files from the internet.

## Usage

Download to a temporary location, with an arbitrary filename based on the url.

```golang
package sample

import (
	"fmt"
	
	"github.com/joshuatcasey/libdownloader"
)

func sample() error {
	downloadedFile, err := libdownloader.Fetch("https://example.com/")
	if err != nil {
		return err
	}

	// downloadedFile.Path() is a string path to a file containing the contents.
	fmt.Printf("Downloaded to %s\n", downloadedFile.Path())
	
	// downloadedFile.Contents() is the contents of the file as a string.
	contents, err := downloadedFile.Contents()
	if err != nil {
		return err
	}
	fmt.Printf("Found contents %s\n", contents)

	// downloadedFile.Sha256() is a string containing the SHA256.
	sha256, err := downloadedFile.Sha256()
	if err != nil {
		return err
	}
	fmt.Printf("Found SHA256 %s\n", sha256)
	
	return nil
}
```

Download to a known location, with a known filename.

```golang
package sample

import (
	"fmt"
	
	"github.com/joshuatcasey/libdownloader"
)

func sample() error {
	downloadedFile, err := libdownloader.Fetch("https://example.com/", 
		libdownloader.WithDownloadDirectory("/tmp/dir"),
		libdownloader.WithFilename("filename.txt"))
	if err != nil {
		return err
	}

	// downloadedFile.Path() should return "/tmp/dir/filename.txt"
	fmt.Printf("Downloaded to %s\n", downloadedFile.Path())

	return nil
}
```

## Unit Tests

```shell
./scripts/unit.sh
```
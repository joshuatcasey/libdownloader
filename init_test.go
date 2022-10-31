package libdownloader_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitRetrieve(t *testing.T) {
	suite := spec.New("libdownloader", spec.Report(report.Terminal{}))
	suite("Fetch", testFetch, spec.Sequential())
	suite("SimpleDownloadedFile", testSimpleDownloadedFile, spec.Sequential())
	suite.Run(t)
}

package dependency

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/build-info-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
)

func TestDetailsString(t *testing.T) {
	dependency := Details{Path: "github.com/jfrog/jfrog/jfrog-cli", Version: "v1.0.0"}
	depName, err := dependency.String()
	assert.NoError(t, err)
	assert.Equal(t, "jfrog-cli v1.0.0", depName)
}

func TestToString(t *testing.T) {
	dependencies := []Details{{Path: "github.com/jfrog/jfrog/jfrog-cli", Version: "v1.0.0"}, {Path: "github.com/jfrog/jfrog/jfrog-cli-core", Version: "v1.2.0"}}
	depsDetails, err := ToString(dependencies)
	assert.NoError(t, err)
	assert.Equal(t, "jfrog-cli v1.0.0, jfrog-cli-core v1.2.0", depsDetails)
}

func TestUpgrade(t *testing.T) {
	tempDirPath := t.TempDir()
	wd, err := os.Getwd()
	assert.NoError(t, err)

	assert.NoError(t, utils.CopyFile(tempDirPath, filepath.Join(wd, "testdata", "gomod")))
	assert.NoError(t, utils.MoveFile(filepath.Join(tempDirPath, "gomod"), filepath.Join(tempDirPath, "go.mod")))
	assert.NoError(t, Upgrade(tempDirPath, []Details{{Path: "github.com/jfrog/jfrog-cli-core", Version: "v1.2.6"}, {Path: "github.com/jfrog/jfrog-client-go", Version: "v0.18.0"}}))
	fileDetails, err := fileutils.GetFileDetails(filepath.Join(tempDirPath, "go.mod"), true)
	assert.NoError(t, err)
	assert.Equal(t, fileDetails.Checksum.Md5, "393573bda8c6f6a10dee023785165ee1")
}

func TestIsCoreV1DepIncluded(t *testing.T) {
	dependency := []Details{{Path: "github.com/jfrog/jfrog-cli-core/v2", Version: "v1.11.4"}}
	assert.True(t, IsCoreV1DepIncluded(dependency))
	dependency = []Details{{Path: "github.com/jfrog/jfrog-cli-core/v2", Version: "v2.11.4"}}
	assert.False(t, IsCoreV1DepIncluded(dependency))
}

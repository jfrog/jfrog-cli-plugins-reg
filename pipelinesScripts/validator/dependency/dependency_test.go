package dependency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetailsString(t *testing.T) {
	dependency := Details{Path: "github.com/jfrog/jfrog/jfrog-cli", Version: "v1.0.0"}
	assert.Equal(t, "jfrog-cli v1.0.0", dependency.String())
}

func TestToString(t *testing.T) {
	dependencies := []Details{{Path: "github.com/jfrog/jfrog/jfrog-cli", Version: "v1.0.0"}, {Path: "github.com/jfrog/jfrog/jfrog-cli-core", Version: "v1.2.0"}}
	assert.Equal(t, "jfrog-cli v1.0.0, jfrog-cli-core v1.2.0", ToString(dependencies))
}

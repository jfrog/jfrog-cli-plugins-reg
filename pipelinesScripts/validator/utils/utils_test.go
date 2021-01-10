package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRepoDetails(t *testing.T) {
	owner, repo := GetRepoDetails("https://github.com/JFrog/Jfrog-ClI-plugins")
	assert.Equal(t, "jfrog", owner)
	assert.Equal(t, "jfrog-cli-plugins", repo)
}

func TestGetPluginsDescriptor(t *testing.T) {
	results, err := GetPluginsDescriptor()
	assert.NoError(t, err)
	assert.NotZero(t, len(results))
}

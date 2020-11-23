package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// Creates a new plugin descriptor in the plugins folder.
func CreatePluginYaml(data *PluginDescriptor, descriptorName string) error {
	dataBytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	pluginPath := filepath.Join(currentDir, "plugins", descriptorName)
	return ioutil.WriteFile(pluginPath, dataBytes, 0644)
}

// Commit all the modified files.
func CommitAllFiles() error {
	if err := RunCommand("git", "add", "."); err != nil {
		return err
	}
	return RunCommand("git", "commit", "-m", "plugin_tests")
}

func CreatePlaygroundForJfrogCliTest() (string, string, error) {
	tempDirPath, err := ioutil.TempDir("", "out")
	if err != nil {
		return "", "", err
	}
	playgroundPath, err := CloneRepository(tempDirPath, JfrogCliPluginRegUrl, "", JfrogCliPluginRegBranch, "")
	if err != nil {
		return "", "", err
	}
	return tempDirPath, playgroundPath, nil
}

func CleanupTestPlayground(t *testing.T, tempDirPath string, oldCW string) {
	assert.NoError(t, os.Chdir(oldCW))
	assert.NoError(t, os.RemoveAll(tempDirPath))
}

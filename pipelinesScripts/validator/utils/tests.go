package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

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
	if err := RunCommand("git", "commit", "-m", "plugin_tests"); err != nil {
		return err
	}
	return nil
}

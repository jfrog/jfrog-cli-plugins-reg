package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-cli-plugins-reg/git"
	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const testPluginRepo = "https://github.com/jfrog/jfrog-cli-plugins"

func TestValidateExtension(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := git.CreatePlaygroundForJfrogCliTest()
	require.NoError(t, err)
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, git.CleanupTestPlayground(tempDirPath, oldCW))
	}()
	// cd to the cloned project
	assert.NoError(t, os.Chdir(playgroundPath))
	descriptorName := "test_extention_plugin"
	assert.NoError(t, execValidator(&utils.PluginDescriptor{}, descriptorName+".yml", validateExtension))
	assert.Error(t, execValidator(&utils.PluginDescriptor{}, descriptorName, validateExtension))
}

func TestValidateDescriptorStructure(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := git.CreatePlaygroundForJfrogCliTest()
	require.NoError(t, err)
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, git.CleanupTestPlayground(tempDirPath, oldCW))
	}()

	// cd to the cloned project
	assert.NoError(t, os.Chdir(playgroundPath))

	pluginDescriptor := &utils.PluginDescriptor{
		PluginName:  "My Plugin",
		Version:     "v1.0.0",
		Repository:  "www.myrepo.com",
		Maintainers: []string{"First_maintainer", "Second_maintainer"},
	}
	descriptorName := "test_extention_plugin.yml"
	assert.NoError(t, execValidator(pluginDescriptor, descriptorName, validateDescriptor))

	// Validate mandatory fields
	pluginDescriptorCopy := *pluginDescriptor
	pluginDescriptorCopy.PluginName = ""
	assert.EqualError(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor), "Errors detected in the yml descriptor file:\n* 'name' is missing\n")
	pluginDescriptorCopy = *pluginDescriptor
	pluginDescriptorCopy.Version = ""
	assert.EqualError(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor), "Errors detected in the yml descriptor file:\n* 'version' is missing\n")

	pluginDescriptorCopy = *pluginDescriptor
	pluginDescriptorCopy.Repository = ""
	assert.EqualError(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor), "Errors detected in the yml descriptor file:\n* 'repository' is missing\n")

	pluginDescriptorCopy = *pluginDescriptor
	pluginDescriptorCopy.Maintainers = nil
	assert.EqualError(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor), "Errors detected in the yml descriptor file:\n* 'maintainers' is missing\n")

	pluginDescriptorCopy = *pluginDescriptor
	pluginDescriptorCopy.Branch = "my-branch"
	pluginDescriptorCopy.Tag = "my-tag"
	assert.EqualError(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor), "Errors detected in the yml descriptor file:\n* Plugin descriptor yml cannot include both 'tag' and 'branch'.\n")
}

func TestValidateDescriptorTests(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := git.CreatePlaygroundForJfrogCliTest()
	require.NoError(t, err)
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, git.CleanupTestPlayground(tempDirPath, oldCW))
	}()

	// cd to the cloned project
	assert.NoError(t, os.Chdir(playgroundPath))
	pluginDescriptor := &utils.PluginDescriptor{
		PluginName:   "My Plugin",
		Version:      "v1.0.0",
		Repository:   testPluginRepo,
		Maintainers:  []string{"First_maintainer", "Second_maintainer"},
		RelativePath: "build-deps-info",
	}
	descriptorName := "test_extention_plugin.yml"
	assert.NoError(t, execValidator(pluginDescriptor, descriptorName, runTests))
}

func execValidator(pluginDescriptor *utils.PluginDescriptor, descriptorName string, validatorFunc func() error) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := CreatePluginYaml(pluginDescriptor, descriptorName, currentDir); err != nil {
		return err
	}
	if err := git.CommitAllFiles(currentDir); err != nil {
		return err
	}
	// Mimic jfrog pipelines entry point.
	defer os.Chdir(currentDir)
	os.Chdir(filepath.Join("pipelinesScripts", "validator"))
	return validatorFunc()
}

// Creates a new plugin descriptor in the plugins folder.
func CreatePluginYaml(data *utils.PluginDescriptor, descriptorName, currentDir string) error {
	dataBytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	pluginPath := filepath.Join(currentDir, "plugins", descriptorName)
	return ioutil.WriteFile(pluginPath, dataBytes, 0644)
}

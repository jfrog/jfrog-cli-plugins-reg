package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPluginRepo = "https://github.com/jfrog/jfrog-cli-plugins"

func TestValidateExtension(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := utils.CreatePlaygroundForJfrogCliTest()
	require.NoError(t, err)
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer utils.CleanupTestPlayground(t, tempDirPath, oldCW)

	// cd to the cloned project
	assert.NoError(t, os.Chdir(playgroundPath))
	descriptorName := "test_extention_plugin"
	assert.NoError(t, execValidator(&utils.PluginDescriptor{}, descriptorName+".yml", validateExtension))
	assert.Error(t, execValidator(&utils.PluginDescriptor{}, descriptorName, validateExtension))
}

func TestValidateDescriptorStructure(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := utils.CreatePlaygroundForJfrogCliTest()
	require.NoError(t, err)
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer utils.CleanupTestPlayground(t, tempDirPath, oldCW)

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
	tempDirPath, playgroundPath, err := utils.CreatePlaygroundForJfrogCliTest()
	require.NoError(t, err)
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer utils.CleanupTestPlayground(t, tempDirPath, oldCW)

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
	if err := utils.CreatePluginYaml(pluginDescriptor, descriptorName); err != nil {
		return err
	}
	if err := utils.CommitAllFiles(); err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	// Mimic jfrog pipelines entry point.
	defer os.Chdir(currentDir)
	os.Chdir(filepath.Join("pipelinesScripts", "validator"))
	return validatorFunc()
}

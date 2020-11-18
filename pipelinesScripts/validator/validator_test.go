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
	tempDir, err := utils.CreatePLaygroundForJfrogCliTest()
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	descriptorName := "test_extention_plugin"
	assert.NoError(t, execValidator(&utils.PluginDescriptor{}, descriptorName+".yml", validateExtension))
	assert.Error(t, execValidator(&utils.PluginDescriptor{}, descriptorName, validateExtension))
}

func TestValidateDescriptorStructure(t *testing.T) {
	tempDir, err := utils.CreatePLaygroundForJfrogCliTest()
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

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
	assert.Error(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor))

	pluginDescriptorCopy = *pluginDescriptor
	pluginDescriptorCopy.Version = ""
	assert.Error(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor))

	pluginDescriptorCopy = *pluginDescriptor
	pluginDescriptorCopy.Repository = ""
	assert.Error(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor))

	pluginDescriptorCopy = *pluginDescriptor
	pluginDescriptorCopy.Maintainers = nil
	assert.Error(t, execValidator(&pluginDescriptorCopy, descriptorName, validateDescriptor))
}

func TestValidateDescriptorTests(t *testing.T) {
	tempDir, err := utils.CreatePLaygroundForJfrogCliTest()
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

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
	// Mimice jfrog pipelines entry point.
	defer os.Chdir(currentDir)
	os.Chdir(filepath.Join("pipelinesScripts", "validator"))
	return validatorFunc()
}

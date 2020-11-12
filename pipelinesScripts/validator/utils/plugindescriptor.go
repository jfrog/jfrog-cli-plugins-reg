package utils

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// PluginDescriptor describes the plugin descriptor yml.
type PluginDescriptor struct {
	// Mandatory fields:
	PluginName  string   `yaml:"pluginName"`  // Example: RT-FS
	Version     string   `yaml:"version"`     // Example: 1.0.0
	Repository  string   `yaml:"repository"`  // Example: https://github.com/jfrog/jfrog-cli-plugins-reg.git
	Maintainers []string `yaml:"maintainers"` // Example: ['frog1', 'frog2']

	// Optional fields:
	RelativePath string `yaml:"relativePath"` // Example: rt-fs
	Branch       string `yaml:"branch"`       // Example: rel-1.0.0
	Tag          string `yaml:"tag"`          // Example: 1.0.0
}

func ReadDescriptor(filePath string) (*PluginDescriptor, error) {
	rootPath, err := getRootPath()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadFile(filepath.Join(rootPath, filePath))
	if err != nil {
		return nil, errors.New("Fail to read '" + filePath + "'. Error: " + err.Error())
	}
	var descriptor PluginDescriptor
	if err := yaml.UnmarshalStrict(content, &descriptor); err != nil {
		return nil, errors.New("Fail to unmarshal '" + filePath + "'. Error: " + err.Error())
	}
	return &descriptor, nil
}

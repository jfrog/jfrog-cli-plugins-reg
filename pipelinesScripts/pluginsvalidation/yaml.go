package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Or-Gabay/jfrog-cli-plugins-reg/utils"
	"gopkg.in/yaml.v2"
)

// This program validate a new jfrog cli plugin register
func main() {
	arg := os.Args[1]
	if len(arg) == 0 {
		fmt.Println("No args was specify")
		os.Exit(1)
	}

	var err error
	switch strings.ToLower(arg) {
	case "extension":
		err = extensionCheck()
	case "structure":
		err = structureCheck()
	case "tests":
		err = pluginTests()
	default:
		err = errors.New("unknown command: " + arg)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// In order to add a plugin to the registry,
// the maintainer should create a pull request to the registry.
// The pull request should include the plugin(s) YAML.
// If other files extentions exists, return an error.
func extensionCheck() error {
	files, err := utils.GetModifiedFiles()
	if err != nil {
		return err
	}
	forbiddenFiles := ""
	for _, committedFilePath := range files {
		if !strings.HasSuffix(committedFilePath, ".yml") {
			forbiddenFiles += committedFilePath + "\n"
		}
	}
	if forbiddenFiles != "" {
		return errors.New("Failed, only .yml files are permitted to be in the pull request. Please remove: " + forbiddenFiles)
	}
	return nil
}

// PluginsYAMLFile describes the plugin YAML.
type PluginsYAMLFile struct {
	// Mandatory
	PluginName      string `yaml:"pluginName"`
	Version         string `yaml:"version"`
	Repository      string `yaml:"repository"`
	MaintainerName  string `yaml:"maintainerName"`
	MaintainerEmail string `yaml:"maintainerEmail"`
	// Optionals
	RelativePath string `yaml:"relativePath"`
	Branch       string `yaml:"branch"`
	Tag          string `yaml:"tag"`
}

// Check the plugin YAML file format. if one of the mandatory fields are missing, return an error.
func structureCheck() error {
	files, err := utils.GetModifiedFiles()
	if err != nil {
		return err
	}
	for _, yamlPath := range files {
		content, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			return errors.New("Fail to ReadFile yaml, error:" + err.Error())
		}
		var pluginsYAML PluginsYAMLFile
		if err := yaml.UnmarshalStrict(content, &pluginsYAML); err != nil {
			return errors.New("Fail to unmarshal yaml, error:" + err.Error())
		}
		fmt.Println("Analyzing:" + yamlPath)
		if err := validateContent(pluginsYAML); err != nil {
			return err
		}
	}
	return nil
}

// Verifies the plugin and run the plugin tests using 'go test ./...'.
func pluginTests() error {
	files, err := utils.GetModifiedFiles()
	if err != nil {
		return err
	}
	for _, yamlPath := range files {
		content, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			return errors.New("Fail to ReadFile yaml, error:" + err.Error())
		}
		var pluginsYAML PluginsYAMLFile
		if err := yaml.Unmarshal(content, &pluginsYAML); err != nil {
			return errors.New("Fail to unmarshal yaml, error:" + err.Error())
		}
		fmt.Println("Analyzing:" + yamlPath)
		tempDir, err := ioutil.TempDir("", "pluginRepo")
		if err != nil {
			return errors.New("Fail to create temp dir, error:" + err.Error())
		}
		defer os.RemoveAll(filepath.Join(tempDir))
		projectPath, err := utils.CloneProject(tempDir, pluginsYAML.Repository, pluginsYAML.RelativePath, pluginsYAML.Branch, pluginsYAML.Tag)
		if err != nil {
			return err
		}
		if err := runProjectTests(projectPath); err != nil {
			return err
		}
	}
	return nil
}

func runProjectTests(projectPath string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.New("Fail to get current directory, error:" + err.Error())
	}
	defer os.Chdir(currentDir)
	err = os.Chdir(projectPath)
	if err != nil {
		return errors.New("Fail to get change directory to" + projectPath + ", error:" + err.Error())
	}
	cmd := exec.Command("go", "test", "./...")
	if _, err := cmd.Output(); err != nil {
		return errors.New("Failed plugin tests for " + projectPath + ", error:" + err.Error())
	}
	return nil
}

func validateContent(pluginsYAML PluginsYAMLFile) error {
	missingfields := ""
	if pluginsYAML.PluginName == "" {
		missingfields += "name\n"
	}
	if pluginsYAML.Version == "" {
		missingfields += "version\n"
	}
	if pluginsYAML.Repository == "" {
		missingfields += "repository\n"
	}
	if pluginsYAML.MaintainerName == "" {
		missingfields += "maintainer name\n"
	}
	if pluginsYAML.MaintainerEmail == "" {
		missingfields += "maintainer email\n"
	}
	if missingfields != "" {
		return errors.New("YAML is missing the following:\n" + missingfields)
	}
	return nil
}

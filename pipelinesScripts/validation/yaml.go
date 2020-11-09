package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	arg := os.Args[1]
	if len(arg) == 0 {
		fmt.Println("No args was specify")
		os.Exit(1)
	}
	var err error
	switch strings.ToLower(arg) {
	case "containyamls":
		err = containYamls()
	case "yamlcontent":
		err = yamlContent()
	default:
		err = errors.New("unknown command: " + arg)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func containYamls() error {
	pathToResource, commitSha := os.Getenv("res_validatePluginCriteria_resourcePath"), os.Getenv("res_validatePluginCriteria_commitSha")
	if pathToResource == "" || commitSha == "" {
		return errors.New("Failed to parse env vars: res_validatePluginCriteria_resourcePath & res_validatePluginCriteria_commitSha")
	}
	os.Chdir(pathToResource)
	cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", commitSha)
	output, err := cmd.Output()
	if err != nil {
		return errors.New("Failed to run git cmd, error:" + err.Error())
	}
	outputStr := strings.Trim(string(output), "\n")
	for _, commitFile := range strings.Split(outputStr, "\n") {
		if !strings.HasSuffix(commitFile, ".yml") {
			return errors.New("Failed, only .yml files are permitted to be in the pull request.")
		}
		fmt.Print(pathToResource + "/" + commitFile + ";")
	}
	return nil
}

// PluginsYAMLFile describes a plugin for jfrog in order to register on 'jfrog-cli-plugins-reg'.
type PluginsYAMLFile struct {
	PluginName      string
	Version         string
	Repository      string
	MaintainerName  string
	MaintainerEmail string
	// Optionals
	RelativePath string
	Branch       string
	Tag          string
}

func yamlContent() error {
	yamlsStr := os.Getenv("yaml_files_path")
	if yamlsStr == "" {
		return errors.New("No YAML is found in the plugins directory.")
	}
	yamlsPath := strings.Split(yamlsStr, ";")
	for _, yamlPath := range yamlsPath {
		content, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			return errors.New("Fail to ReadFile yaml, error:" + err.Error())
		}
		var pluginsYAML PluginsYAMLFile
		if err := yaml.Unmarshal(content, &pluginsYAML); err != nil {
			return errors.New("Fail to unmarshal yaml, error:" + err.Error())
		}
		if err := validateContent(pluginsYAML); err != nil {
			return err
		}
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
		return errors.New("YAML is missing the following: " + missingfields)
	}
	return nil
}

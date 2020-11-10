package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	case "extension":
		err = extensionCheck()
	case "structure":
		err = structureCheck()
	case "tests":
		err = structureTests()
	default:
		err = errors.New("unknown command: " + arg)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func extensionCheck() error {
	files, err := getModifiedFiles()
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

// PluginsYAMLFile describes a plugin for jfrog in order to register on 'jfrog-cli-plugins-reg'.
type PluginsYAMLFile struct {
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

func structureCheck() error {
	files, err := getModifiedFiles()
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
		if err := validateContent(pluginsYAML); err != nil {
			return err
		}
	}
	return nil
}
func structureTests() error {
	files, err := getModifiedFiles()
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
		projectPath, err := cloneProject(pluginsYAML.Repository)
		runProjectTests(projectPath)
		defer os.RemoveAll(filepath.Join(projectPath, ".."))
	}
	return nil
}

func cloneProject(projectRepository string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", errors.New("Fail to get current directory, error:" + err.Error())
	}
	defer os.Chdir(currentDir)
	tempDir, err := ioutil.TempDir(currentDir, "pluginRepo")
	if err != nil {
		return "", errors.New("Fail to create temp dir, error:" + err.Error())
	}
	err = os.Chdir(tempDir)
	if err != nil {
		return "", errors.New("Fail to get change directory to" + tempDir + ", error:" + err.Error())
	}
	projectRepository = strings.TrimSuffix(projectRepository, ".git")
	cmd := exec.Command("git", "clone", projectRepository+".git")
	if _, err := cmd.Output(); err != nil {
		return "", errors.New("Failed to run git clone for " + projectRepository + ", error:" + err.Error())
	}
	repositoryName := projectRepository[strings.LastIndex(projectRepository, "/")+1:]
	return filepath.Join(tempDir, repositoryName), nil
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

// Return the paths to the modified files for all affected files since master's commit.
func getModifiedFiles() ([]string, error) {
	pathToResource, commitSha := os.Getenv("res_jfrog_cli_plugins_reg_resourcePath"), os.Getenv("res_jfrog_cli_plugins_reg_commitSha")
	if pathToResource == "" || commitSha == "" {
		return nil, errors.New("Failed to parse env vars: res_jfrog_cli_plugins_reg_resourcePath & res_jfrog_cli_plugins_reg_commitSha")
	}
	os.Chdir(pathToResource)
	cmd := exec.Command("git", "diff", "--no-commit-id", "--name-only", "-r", "master..."+commitSha)
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("Failed to run git cmd, error:" + err.Error())
	}
	var fullPathCommittedFiles []string
	for _, committedFile := range strings.Split(string(output), "\n") {
		if committedFile != "" {
			fullPathCommittedFiles = append(fullPathCommittedFiles, pathToResource+"/"+committedFile)
		}
	}
	return fullPathCommittedFiles, nil
}

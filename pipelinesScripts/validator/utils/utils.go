package utils

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

type ValidationType string

const (
	Extension           ValidationType = "extension"
	Structure                          = "structure"
	Tests                              = "tests"
	UpgradeJfrogPlugins                = "upgradejfrogplugins"
	PluginDescriptorDir                = "plugins"
)

func PrintUsageAndExit() {
	fmt.Printf("Usage: `go run validator.go <command>`\nPossible commands: '%s', '%s' or '%s'\n", Extension, Structure, Tests)
	os.Exit(1)
}

func RunCommand(dir string, getOutput bool, name string, args ...string) (output string, err error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if getOutput {
		var data []byte
		data, err = cmd.Output()
		output = string(data)
	} else {
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		err = cmd.Run()
	}
	return
}

func getRootPath() (string, error) {
	rootPath := filepath.Join("..", "..")
	absRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return "", errors.New("Failed to convert path to Abs path for " + rootPath + ". Error:" + err.Error())
	}
	if _, err := os.Stat(filepath.Join(absRootPath, "plugins")); os.IsNotExist(err) {
		return "", errors.New("Failed to find 'plugin' folder in:" + rootPath + ". Error:" + err.Error())
	}
	return absRootPath, nil
}

func GetPluginsDescriptor() ([]*PluginDescriptor, error) {
	rootPath, err := getRootPath()
	if err != nil {
		return nil, err
	}
	fileInfos, err := ioutil.ReadDir(filepath.Join(rootPath, PluginDescriptorDir))
	if err != nil {
		return nil, err
	}
	var resutls []*PluginDescriptor
	for _, file := range fileInfos {
		fDescriptor, err := ReadDescriptor(filepath.Join(PluginDescriptorDir, file.Name()))
		if err != nil {
			return nil, err
		}
		resutls = append(resutls, fDescriptor)
	}
	return resutls, nil
}

// Returns the plugins repository owner and name. e.g.: https://github.com/JFrog/jfrog-cli-plugins => (jfrog, jfrog-cli-plugins)
func GetRepoDetails(pluginRepository string) (owner string, repo string) {
	pluginRepository = strings.Replace(pluginRepository, "https://github.com/", "", 1)
	splitted := strings.Split(pluginRepository, "/")
	return strings.ToLower(splitted[0]), strings.ToLower(splitted[1])
}

func UpdateGoDependency(runAt, DepName, depVersion string) (err error) {
	_, err = RunCommand(runAt, false, "go", "get", DepName+"@"+depVersion)
	if err != nil {
		fmt.Println("Go Get failed for at" + runAt)
	}
	return
}

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
		return nil, errors.New("Failed to read '" + filePath + "'. Error: " + err.Error())
	}
	var descriptor PluginDescriptor
	if err := yaml.UnmarshalStrict(content, &descriptor); err != nil {
		return nil, errors.New("Failed to unmarshal '" + filePath + "'. Error: " + err.Error())
	}
	return &descriptor, nil
}

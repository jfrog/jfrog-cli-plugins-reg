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
	Structure           ValidationType = "structure"
	Tests               ValidationType = "tests"
	UpgradeJfrogPlugins ValidationType = "upgradejfrogplugins"
	PluginDescriptorDir string         = "plugins"
	rootDirName                        = "jfrog-cli-plugins-reg"
)

func PrintUsageAndExit() {
	fmt.Printf("Usage: `go run validator.go <command>`\nPossible commands: '%s', '%s', '%s' or '%s'\n", Extension, Structure, Tests, UpgradeJfrogPlugins)
	os.Exit(1)
}

func RunCommand(dir string, getOutput bool, name string, args ...string) (output string, err error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if getOutput {
		var data []byte
		data, err = cmd.CombinedOutput()
		output = string(data)
	} else {
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		err = cmd.Run()
	}
	return
}

// Gets the root directory of `jfrog-cli-plugins-reg` project, where the plugins descriptors directory located.
func getRootPath() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if !strings.Contains(pwd, rootDirName) {
		return "", errors.New("Failed to find 'plugin' folder in:" + pwd + ".")
	}
	return strings.Split(pwd, rootDirName)[0] + rootDirName, nil
}

func GetPluginsDescriptors() ([]*PluginDescriptor, error) {
	rootPath, err := getRootPath()
	if err != nil {
		return nil, err
	}
	fileInfos, err := ioutil.ReadDir(filepath.Join(rootPath, PluginDescriptorDir))
	if err != nil {
		return nil, err
	}
	var results []*PluginDescriptor
	for _, file := range fileInfos {
		fDescriptor, err := ReadDescriptor(filepath.Join(PluginDescriptorDir, file.Name()))
		if err != nil {
			return nil, err
		}
		results = append(results, fDescriptor)
	}
	return results, nil
}

// Returns the plugins repository owner and name lowecase. e.g.: https://github.com/JFrog/jfrog-CLI-plugins => (jfrog, jfrog-cli-plugins)
func ExtractRepoDetails(pluginRepository string) (owner string, repo string) {
	pluginRepository = strings.Replace(pluginRepository, "https://github.com/", "", 1)
	splitted := strings.Split(pluginRepository, "/")
	return strings.ToLower(splitted[0]), strings.ToLower(splitted[1])
}

func UpdateGoDependency(runAt, depName, depVersion string) (err error) {
	dependency := depName + "@" + depVersion
	fmt.Println(fmt.Sprintf("Running command 'go get %v' at '%v'", dependency, runAt))
	var output string
	output, err = RunCommand(runAt, true, "go", "get", dependency)
	if err != nil {
		fmt.Println(fmt.Sprintf("Go Get failed at %v, output:'%v'", runAt, output))
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

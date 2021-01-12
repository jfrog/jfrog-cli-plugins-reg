package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jfrog/jfrog-cli-plugins-reg/dependency"
	"github.com/jfrog/jfrog-cli-plugins-reg/git"
	"github.com/jfrog/jfrog-cli-plugins-reg/github"
	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
)

const (
	GitHubIssueTitle = "Failed upgrading dependencies"
	GitHubIssueBody  = "This issue opened by the JFrog CLI plugins bot. I attended to upgrade the following plugin(s) to\n%s:%s\nThe following commands failed after upgrading:\n go ver ./...\ngo test -v ./...\nThe upgrade was therefore aborted. Please fix and upgrade manually."
)

// This program runs a series of validations on a new JFrog CLI plugin, following a pull request to register it in the public registry.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: Wrong number of arguments.")
		utils.PrintUsageAndExit()
	}
	command := os.Args[1]
	var err error
	switch strings.ToLower(command) {
	case string(utils.Extension):
		err = validateExtension()
	case string(utils.Structure):
		err = validateDescriptor()
	case string(utils.Tests):
		err = runTests()
	case string(utils.UpgradeJfrogPlugins):
		err = upgradeJfrogPlugins()
	default:
		err = errors.New("Unknown command: " + command)
	}
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}

// In order to add a plugin to the registry,
// the maintainer should create a pull request to the registry.
// The pull request should include the plugin(s) YAML.
// If the pull request includes other files, return an error.
func validateExtension() error {
	prFiles, err := git.GetModifiedFiles()
	if err != nil {
		return err
	}
	forbiddenFiles := ""
	for _, committedFilePath := range prFiles {
		if !strings.HasSuffix(committedFilePath, ".yml") || !strings.HasPrefix(committedFilePath, utils.PluginDescriptorDir+"/") {
			forbiddenFiles += committedFilePath + "\n"
		}
	}
	if forbiddenFiles != "" {
		return errors.New("Only .yml files are permitted to be included in the pull request. Please remove: " + forbiddenFiles)
	}
	return nil
}

// Check the plugin YAML file format. if one of the mandatory fields are missing, return an error.
func validateDescriptor() error {
	files, err := git.GetModifiedFiles()
	if err != nil {
		return err
	}
	for _, yamlPath := range files {
		fmt.Println("Validating:" + yamlPath)

		descriptor, err := utils.ReadDescriptor(yamlPath)
		if err != nil {
			return err
		}

		if err := validateContent(descriptor); err != nil {
			return err
		}
	}
	return nil
}

// Verifies the plugin and run the plugin tests using 'go test ./...'.
func runTests() error {
	files, err := git.GetModifiedFiles()
	if err != nil {
		return err
	}
	for _, yamlPath := range files {
		fmt.Println("Analyzing:" + yamlPath)

		descriptor, err := utils.ReadDescriptor(yamlPath)
		if err != nil {
			return err
		}
		tempDir, err := ioutil.TempDir("", "pluginRepo")
		if err != nil {
			return errors.New("ERROR: Failed to create temp dir: " + err.Error())
		}
		defer func() {
			if deferErr := os.RemoveAll(tempDir); deferErr != nil {
				fmt.Println("ERROR: Failed to remove temp dir. Error:" + deferErr.Error())
			}
		}()
		projectPath, err := git.CloneRepository(tempDir, descriptor.Repository, descriptor.RelativePath, descriptor.Branch, descriptor.Tag)
		if err != nil {
			return err
		}
		if err := runProjectTests(projectPath); err != nil {
			return err
		}
	}
	return nil
}

func upgradeJfrogPlugins() error {
	if len(os.Args) < 3 {
		return errors.New("missing cli plugin path.")
	}
	cliPluginPath := os.Args[2]
	fileInfo, err := os.Stat(cliPluginPath)
	if os.IsNotExist(err) || !fileInfo.IsDir() {
		return errors.New("ERROR: " + cliPluginPath + " is not a directory.")
	}
	descriptors, err := utils.GetPluginsDescriptor()
	if err != nil {
		return err
	}
	fmt.Println("Starting to upgrade JFrog plugins...")
	token := os.Getenv("issue_token")
	if token == "" {
		return errors.New("issue_token was not found.")
	}
	depToUpgrade, err := dependency.GetJfrogLatest()
	if err != nil {
		return err
	}
	var failedPlugins []string
	for _, descriptor := range descriptors {
		// Filter plugins that are not owned by JFrog.
		owner, _ := utils.ExtractRepoDetails(descriptor.Repository)
		if owner != "jfrog" {
			continue
		}
		fmt.Println("Upgrading: " + descriptor.PluginName)
		projectPath := filepath.Join(cliPluginPath, descriptor.RelativePath)
		if err := dependency.Upgrade(projectPath, depToUpgrade); err != nil {
			return err
		}
		fmt.Println("Running tests after upgrade...")
		if err := runValidation(projectPath); err != nil {
			fmt.Println("ERROR: Go test/vert failed, skipping upgrade " + descriptor.PluginName + ".")
			failedPlugins = append(failedPlugins, descriptor.PluginName)
			continue
		}
		fmt.Println("Stage go.mod and go.sum")
		stagedCount, err := git.StageModifiedFiles(projectPath, "go.mod", "go.sum")
		if err != nil {
			return err
		}
		if stagedCount == 0 {
			fmt.Println("No file were changed due to upgrade for plugin: " + descriptor.PluginName)
		} else {
			fmt.Println(fmt.Sprintf("%v files were staged.", stagedCount))
		}
	}
	if len(failedPlugins) > 0 {
		pluginsSummary := ""
		for _, pluginName := range failedPlugins {
			pluginsSummary += "\n" + pluginName
		}
		req := github.IssuesReq{
			Title: GitHubIssueTitle,
			Body:  fmt.Sprintf(GitHubIssueBody, dependency.ToString(depToUpgrade), pluginsSummary),
		}
		if err := github.OpenIssue("jfrog", "jfrog-cli-plugins", token, req); err != nil {
			return err
		}
	}
	return nil
}

func runProjectTests(projectPath string) error {
	var currentDir string
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.New("Failed to get current directory: " + err.Error())
	}
	defer func() {
		if deferErr := os.Chdir(currentDir); deferErr != nil {
			fmt.Println("ERROR: Failed to change dir to " + currentDir + ". Error:" + deferErr.Error())
		}
	}()
	err = os.Chdir(projectPath)
	if err != nil {
		return errors.New("Failed to get change directory to" + projectPath + ": " + err.Error())
	}
	return runValidation(projectPath)
}

func validateContent(descriptor *utils.PluginDescriptor) error {
	missingfields := ""
	if descriptor.PluginName == "" {
		missingfields += "* 'name' is missing\n"
	}
	if descriptor.Version == "" {
		missingfields += "* 'version' is missing\n"
	}
	if descriptor.Repository == "" {
		missingfields += "* 'repository' is missing\n"
	}
	if len(descriptor.Maintainers) == 0 {
		missingfields += "* 'maintainers' is missing\n"
	}
	if descriptor.Tag != "" && descriptor.Branch != "" {
		missingfields += "* Plugin descriptor yml cannot include both 'tag' and 'branch'.\n"
	}
	if missingfields != "" {
		return errors.New("Errors detected in the yml descriptor file:\n" + missingfields)
	}
	return nil
}

func runValidation(projectPath string) (err error) {
	var output string
	if output, err = utils.RunCommand(projectPath, true, "go", "vet", "-v", "./..."); err != nil {
		fmt.Println("Failed to Lint plugin source code, located at " + projectPath + ". Error:\n" + output)
		return
	}
	if output, err = utils.RunCommand(projectPath, true, "go", "test", "-v", "./..."); err != nil {
		fmt.Println("Plugin Tests failed at " + projectPath + ". Error:\n" + output)
	}
	return
}

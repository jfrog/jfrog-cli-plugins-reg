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
	GitHubIssueBody  = "This issue was opened by the JFrog CLI Plugins Registry bot. We attended to upgrade the following plugin(s) to\n%s:%s\nThe following commands failed after upgrading:\n go vet -v ./...\ngo test -v ./...\nThe upgrade commit and push were therefore aborted. Please fix the issue and upgrade manually."
)

// This program runs a series of validations and upgrades on JFrog CLI plugins, following a pull request to register it in the public registry.
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
// the maintainer required to create a pull request to the registry.
// The pull request must include the plugin(s) YAML.
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

// Run the plugin tests using 'go test ./...'.
func runTests() (err error) {
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
			return err
		}
		defer func() {
			if deferErr := os.RemoveAll(tempDir); err == nil {
				err = deferErr
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
	return err
}

func upgradeJfrogPlugins() error {
	cliPluginPath, token, err := getUpgradeArgs()
	if err != nil {
		return err
	}
	descriptors, err := utils.GetPluginsDescriptors()
	if err != nil {
		return err
	}
	depToUpgrade, err := dependency.GetJfrogLatest()
	if err != nil {
		return err
	}
	if dependency.IsCoreV1DepIncluded(depToUpgrade) {
		fmt.Println("The JFrog-cli-core release 1.x.x was detected. Skipping the upgrade process...")
		return nil
	}
	fmt.Println("Starting to upgrade JFrog plugins...")
	failedPlugins, err := doUpgrade(descriptors, depToUpgrade, cliPluginPath)
	if err != nil {
		return err
	}
	if len(failedPlugins) > 0 {
		if err := openIssue(failedPlugins, depToUpgrade, token); err != nil {
			return err
		}
	}
	return nil
}

// Upgrade jfrog plugins dependencies.
// Returns a list of plugins that failed in the upgrade process.
func doUpgrade(descriptors []*utils.PluginDescriptor, depToUpgrade []dependency.Details, pluginsRoot string) ([]string, error) {
	var failedPlugins []string
	for _, descriptor := range descriptors {
		// Filter out plugins which are not owned by JFrog.
		owner, repo := utils.ExtractRepoDetails(descriptor.Repository)
		if owner != "jfrog" || !strings.HasPrefix(repo, "jfrog-cli-plugins") {
			continue
		}
		fmt.Println("Upgrading: " + descriptor.PluginName)
		projectPath := filepath.Join(pluginsRoot, descriptor.RelativePath)
		if err := dependency.Upgrade(projectPath, depToUpgrade); err != nil {
			return nil, err
		}
		fmt.Println("Running tests after upgrade...")
		if err := runValidation(projectPath); err != nil {
			fmt.Println(err.Error() + ". Skipping upgrade " + descriptor.PluginName + ".")
			failedPlugins = append(failedPlugins, descriptor.PluginName)
			continue
		}
		fmt.Println("Stage go.mod and go.sum")
		stagedCount, err := git.StageModifiedFiles(projectPath, "go.mod", "go.sum")
		if err != nil {
			return nil, err
		}
		if stagedCount == 0 {
			fmt.Println("No file were changed due to upgrade for plugin: " + descriptor.PluginName)
		} else {
			fmt.Println(fmt.Sprintf("%v files were staged.", stagedCount))
		}
	}
	return failedPlugins, nil
}

// Open a new GitHub issue for plugins that failed in the upgrade process.
func openIssue(failedPlugins []string, depToUpgrade []dependency.Details, token string) error {
	pluginsSummary := ""
	for _, pluginName := range failedPlugins {
		pluginsSummary += "\n" + pluginName
	}
	depsDetails, err := dependency.ToString(depToUpgrade)
	if err != nil {
		return err
	}
	req := github.IssuesReq{
		Title: GitHubIssueTitle,
		Body:  fmt.Sprintf(GitHubIssueBody, depsDetails, pluginsSummary),
	}
	if err := github.OpenIssue("jfrog", "jfrog-cli-plugins", token, req); err != nil {
		return err
	}
	return nil
}

// Returns the necessary arguments to run the upgrade process.
func getUpgradeArgs() (string, string, error) {
	if len(os.Args) < 3 {
		return "", "", errors.New("missing cli plugin path.")
	}
	cliPluginPath := os.Args[2]
	fileInfo, err := os.Stat(cliPluginPath)
	if os.IsNotExist(err) || !fileInfo.IsDir() {
		return "", "", errors.New(cliPluginPath + " is not a directory.")
	}
	token := os.Getenv("issue_token")
	if token == "" {
		return "", "", errors.New("issue_token env was not found.")
	}
	return cliPluginPath, token, nil
}

func runProjectTests(projectPath string) (err error) {
	var currentDir string
	currentDir, err = os.Getwd()
	if err != nil {
		return err
	}
	defer func() {
		if deferErr := os.Chdir(currentDir); err == nil {
			err = deferErr
		}
	}()
	if err = os.Chdir(projectPath); err != nil {
		return err
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
	// Golang 1.16 or above requires to run 'go mod tidy' in order to update go.mod and go.sum files after the upgrade.
	if output, err = utils.RunCommand(projectPath, true, "go", "mod", "tidy"); err != nil {
		err = errors.New("Failed to run 'go mod tidy' located at " + projectPath + ". Error:\n" + output)
		return
	}
	if output, err = utils.RunCommand(projectPath, true, "go", "vet", "-v", "./..."); err != nil {
		err = errors.New("Failed to run 'go vet -v ./...' located at " + projectPath + ". Error:\n" + output)
		return
	}
	if output, err = utils.RunCommand(projectPath, true, "go", "test", "-v", "./..."); err != nil {
		err = errors.New("Plugin Tests failed at " + projectPath + ". Error:\n" + output)
	}
	return
}

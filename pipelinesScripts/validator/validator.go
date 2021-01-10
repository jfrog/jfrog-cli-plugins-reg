package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jfrog/jfrog-cli-plugins-reg/dependency"
	"github.com/jfrog/jfrog-cli-plugins-reg/git"
	"github.com/jfrog/jfrog-cli-plugins-reg/github"
	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
)

const (
	GitHubIssueTitle = "Failed upgrading dependencies"
	GitHubIssueBody  = "This issue opened by the JFrog CLI plugins bot. I attended to upgrade '%s' plugin to:\n %s.\n The following commands failed after upgrading:\n go ver ./...\ngo test -v ./...\nThe upgrade was therefore aborted. Please fix and upgrade manually."
)

// This program runs a series of validations on a new JFrog CLI plugin, following a pull request to register it in the public registry.
func main() {
	if len(os.Args) != 2 {
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
	descriptors, err := utils.GetPluginsDescriptor()
	if err != nil {
		return err
	}
	var localErr error
	fmt.Println("Starting to upgrade JFrog plugins...")
	token := os.Getenv("int_github_generic_token")
	if token == "" {
		return errors.New("missing Token to open an issue")
	}
	for _, descriptor := range descriptors {
		// Filter plugins that are not owned by JFrog.
		owner, repo := utils.GetRepoDetails(descriptor.Repository)
		if owner != "jfrog" {
			continue
		}
		fmt.Println("Upgrading: " + descriptor.PluginName)
		tempDir, err := ioutil.TempDir("", "pluginRepo")
		if err != nil {
			return errors.New("ERROR: Failed to create temp dir: " + err.Error())
		}
		defer func() {
			if deferErr := os.RemoveAll(tempDir); deferErr != nil {
				fmt.Println("ERROR: Failed to remove temp dir. Error:" + deferErr.Error())
			}
		}()
		fmt.Println("Cloning repository: " + descriptor.Repository)
		projectPath, err := git.CloneRepository(tempDir, descriptor.Repository, descriptor.RelativePath, descriptor.Branch, descriptor.Tag)
		if err != nil {
			return err
		}
		depToUpgrade, err := dependency.GetJfrogLatest(projectPath)
		if err != nil {
			return err
		}
		fmt.Println("Upgrading dependencies...")
		if err := dependency.Upgrade(projectPath, depToUpgrade); err != nil {
			fmt.Println("ERROR: " + descriptor.PluginName + " failed to upgrade. Error" + err.Error())
			localErr = err
		}
		stagedCount, err := git.StageModifiedFiles(projectPath, "go.mod", "go.sum")
		if err != nil {
			return err
		}
		if stagedCount == 0 {
			fmt.Println("No file were changed due to upgrade. Skipping commit step ")
			continue
		}
		fmt.Println("Running tests after upgrade...")
		if err := runValidation(projectPath); err != nil {
			fmt.Println("ERROR: Go test/vert failed, skipping upgrading " + descriptor.PluginName + ". Opening a GitHub issue ")
			req := github.IssuesReq{
				Title: GitHubIssueTitle,
				Body:  fmt.Sprintf(GitHubIssueBody, descriptor.PluginName, dependency.ToString(depToUpgrade)),
			}
			github.OpenIssue(owner, repo, token, req)
			continue
		}
		fmt.Println("Commiting changes...")
		if err := git.CommitStagedFiles(projectPath, "Upgrade dependencies of plugin "+descriptor.PluginName); err != nil {
			return err
		}
		fmt.Println("Pushing changes...")
		if err := git.Push(projectPath, descriptor.Repository, token, descriptor.Branch); err != nil {
			return err
		}
		fmt.Println(descriptor.PluginName + " plugin upgraded successfully")
	}
	return localErr
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
	utils.RunCommand(projectPath, false, "go", "vet", "-v", "./...")
	if _, err = utils.RunCommand(projectPath, false, "go", "vet", "-v", "./..."); err != nil {
		fmt.Println("Failed to Lint plugin source code, located at " + projectPath)
		return
	}
	if _, err = utils.RunCommand(projectPath, false, "go", "test", "-v", "./..."); err != nil {
		fmt.Println("Plugin Tests failed at " + projectPath)
	}
	return
}

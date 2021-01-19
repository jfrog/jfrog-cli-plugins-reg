package git

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
)

const (
	JfrogCliPluginsRegUrl    = "https://github.com/jfrog/jfrog-cli-plugins-reg.git"
	JfrogCliPluginsRegBranch = "master"
)

// Clones the plugin's repository to a local temp directory and returns the full path of the plugin's source code.
// 'tempDir' - Temporary dir to which the project will be cloned.
// 'gitRepository' - The GitHub repository to clone.
// 'relativePath' - Relative path inside the repository.
// 'branch' - If specified, override the default branch with the input branch.
// 'tag' - If specified, checkout to the input tag.
// returns: (project-path, error)
func CloneRepository(destination, gitRepository, relativePath, branch, tag string) (string, error) {
	gitRepository = strings.TrimSuffix(gitRepository, ".git")
	if err := cloneRepository(destination, branch, tag, gitRepository); err != nil {
		return "", errors.New("Failed to run git clone for " + gitRepository + ", error:" + err.Error())
	}
	repositoryName := gitRepository[strings.LastIndex(gitRepository, "/")+1:]
	return filepath.Join(destination, repositoryName, relativePath), nil
}

func GetModifiedFiles() (modifiedFiles []string, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.New("Couldn't get current directory: " + err.Error())
	}
	// Create unique upstream and branch names
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	uniqueUpstream := "remote-origin-" + timestamp
	if err := addRemote(currentDir, uniqueUpstream, JfrogCliPluginsRegUrl); err != nil {
		return nil, err
	}
	defer func() {
		if deferErr := removeRemote(currentDir, uniqueUpstream); err == nil {
			err = deferErr
		}
	}()
	// Fetch from upstream
	if err := fetch(currentDir, uniqueUpstream, JfrogCliPluginsRegBranch); err != nil {
		return nil, err
	}
	return runGitDiff(currentDir, uniqueUpstream+"/master")
}

// Runs the cmd 'git add FILE -v' over all 'files' and returns the total number of staged files.
func StageModifiedFiles(runAt string, files ...string) (stagedCount int, err error) {
	var cmdOutput string
	for _, file := range files {
		if cmdOutput, err = utils.RunCommand(runAt, true, "git", "add", file, "-v"); err != nil {
			return
		}
		if cmdOutput != "" {
			stagedCount++
		}
	}
	return
}

func addRemote(runAt, remoteName, remoteUrl string) (err error) {
	if _, err = utils.RunCommand(runAt, false, "git", "remote", "add", remoteName, remoteUrl); err != nil {
		err = errors.New("Failed to add git remote for " + remoteName + " upstream and" + remoteUrl + " branch. Error:" + err.Error())
	}
	return
}

func removeRemote(runAt, remoteName string) (err error) {
	if _, err = utils.RunCommand(runAt, false, "git", "remote", "remove", remoteName); err != nil {
		err = errors.New("Failed to remove remote upstream " + remoteName + ". Error:" + err.Error())
	}
	return
}

func fetch(runAt, remoteName, branch string) (err error) {
	if _, err = utils.RunCommand(runAt, false, "git", "fetch", remoteName, branch); err != nil {
		err = errors.New("Failed to fetch from " + remoteName + ", branch " + branch + ". Error:" + err.Error())
	}
	return
}

func cloneRepository(runAt, branch, tag, repo string) (err error) {
	flags := []string{"clone"}
	if branch != "" {
		flags = append(flags, "--branch="+branch)
	}
	if tag != "" {
		flags = append(flags, "--branch="+tag)
	}
	flags = append(flags, repo+".git")
	_, err = utils.RunCommand(runAt, false, "git", flags...)
	return
}

func runGitDiff(runAt, compareTo string) ([]string, error) {
	output, err := utils.RunCommand(runAt, true, "git", "diff", "--no-commit-id", "--name-only", "-r", compareTo+"...HEAD")
	if err != nil {
		return nil, errors.New("Failed to run git diff command: " + err.Error())
	}
	var fullPathCommittedFiles []string
	for _, file := range strings.Split(string(output), "\n") {
		if file != "" {
			fullPathCommittedFiles = append(fullPathCommittedFiles, file)
		}
	}
	return fullPathCommittedFiles, nil
}

// Commit all the modified files.
func CommitAllFiles(runAt string) (err error) {
	if _, err = utils.RunCommand(runAt, false, "git", "add", "."); err != nil {
		return
	}
	_, err = utils.RunCommand(runAt, false, "git", "commit", "-m", "plugin_tests")
	return
}

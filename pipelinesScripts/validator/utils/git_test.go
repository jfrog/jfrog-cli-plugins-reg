package utils

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModifiedFiles(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := CreatePLaygroundForJfrogCliTest()
	require.NoError(t, err)
	// cd to the cloned project
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(playgroundPath))

	// Create new file.
	assert.NoError(t, ioutil.WriteFile("file.txt", []byte("test"), 0644))
	assert.NoError(t, CommitAllFiles())

	files, err := GetModifiedFiles()
	require.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, files[0], "file.txt")

	// Cleanup
	assert.NoError(t, os.RemoveAll(tempDirPath))
	assert.NoError(t, os.Chdir(oldCW))
}

func TestGetModifiedFilesCleanupBranches(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := CreatePLaygroundForJfrogCliTest()
	require.NoError(t, err)
	// cd to the cloned project
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(playgroundPath))

	cmd := exec.Command("git", "branch")
	branchesBefore, err := cmd.Output()
	assert.NoError(t, err)

	modifiedFiles, err := GetModifiedFiles()
	assert.NoError(t, err)
	assert.Empty(t, modifiedFiles)

	cmd = exec.Command("git", "branch")
	branchesAfter, err := cmd.Output()
	assert.NoError(t, err)
	assert.Equal(t, branchesBefore, branchesAfter)

	// Cleanup
	assert.NoError(t, os.RemoveAll(tempDirPath))
	assert.NoError(t, os.Chdir(oldCW))
}

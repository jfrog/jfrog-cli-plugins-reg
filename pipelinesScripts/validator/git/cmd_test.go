package git

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModifiedFiles(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := CreatePlaygroundForJfrogCliTest(t)
	require.NoError(t, err)
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, CleanupTestPlayground(tempDirPath, oldCW))
	}()
	// CD to the cloned project
	require.NoError(t, os.Chdir(playgroundPath))
	// Create new file.
	assert.NoError(t, os.WriteFile("file.txt", []byte("test"), 0600))
	assert.NoError(t, CommitAllFiles(playgroundPath))
	files, err := GetModifiedFiles()
	require.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, files[0], "file.txt")
}

func TestStageModifiedFiles(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := CreatePlaygroundForJfrogCliTest(t)
	require.NoError(t, err)
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, CleanupTestPlayground(tempDirPath, oldCW))
	}()
	// CD to the cloned project
	require.NoError(t, os.Chdir(playgroundPath))
	// Create new file.
	assert.NoError(t, os.WriteFile("file", []byte("test"), 0600))
	assert.NoError(t, os.WriteFile("file2", []byte("test"), 0600))
	count, err := StageModifiedFiles(playgroundPath, "file", "file2")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestGetModifiedFilesCleanupBranches(t *testing.T) {
	// Init playground
	tempDirPath, playgroundPath, err := CreatePlaygroundForJfrogCliTest(t)
	require.NoError(t, err)

	// CD to the cloned project
	oldCW, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, CleanupTestPlayground(tempDirPath, oldCW))
	}()

	// CD to the cloned project
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
}

package git

import (
	"os"
	"testing"
)

// Clones the 'jfrog-cli-plugin-reg' repo to a temp dir for tests purposes.
func CreatePlaygroundForJfrogCliTest(t *testing.T) (string, string, error) {
	tempDirPath := t.TempDir()
	playgroundPath, err := CloneRepository(tempDirPath, JfrogCliPluginsRegUrl, "", JfrogCliPluginsRegBranch, "")
	if err != nil {
		return "", "", err
	}
	return tempDirPath, playgroundPath, nil
}

func CleanupTestPlayground(tempDirPath string, oldCW string) (err error) {
	if err = os.Chdir(oldCW); err != nil {
		return
	}
	err = os.RemoveAll(tempDirPath)
	return
}

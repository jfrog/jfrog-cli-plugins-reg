package git

import (
	"io/ioutil"
	"os"
)

// Clones the 'jfrog-cli-plugin-reg' repo to a temp dir for tests purposes.
func CreatePlaygroundForJfrogCliTest() (string, string, error) {
	tempDirPath, err := ioutil.TempDir("", "out")
	if err != nil {
		return "", "", err
	}
	playgroundPath, err := CloneRepository(tempDirPath, JfrogCliPluginRegUrl, "", JfrogCliPluginRegBranch, "")
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

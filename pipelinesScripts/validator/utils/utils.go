package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ValidationType string

const (
	Extension             ValidationType = "extension"
	Structure                            = "structure"
	Tests                                = "tests"
	PluginDescriptoPrefix                = "plugins/"
)

func PrintUsageAndExit() {
	fmt.Printf("Usage: `go run validator.go <command>`\nPossible commands: '%s', '%s' or '%s'\n", Extension, Structure, Tests)
	os.Exit(1)
}

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
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

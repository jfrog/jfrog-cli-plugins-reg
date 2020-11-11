package utils

import (
	"fmt"
	"os"
	"os/exec"
)

type ValidationType string

const (
	Extension ValidationType = "extension"
	Structure                = "structure"
	Tests                    = "tests"
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

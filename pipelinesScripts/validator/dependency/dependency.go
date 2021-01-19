// Package dependency collects common JFrog dependencies for plugins to update.
package dependency

import (
	"errors"
	"strings"

	"github.com/jfrog/jfrog-cli-plugins-reg/github"
	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
)

// `Details` struct identifies the golang dependency.
type Details struct {
	// Dependency Import path as listed in go.mod.
	Path    string
	Version string
}

// String returns the name and version of the dependency, omitting the path prefix.
func (d *Details) String() (string, error) {
	// Starting position of dependency name.
	idx := strings.LastIndex(d.Path, "/")
	if idx == -1 {
		return "", errors.New("failed to locate dependency name")
	}
	return d.Path[strings.LastIndex(d.Path, "/")+1:] + " " + d.Version, nil
}

var (
	jfrogDependencies = [...]string{"jfrog-cli-core", "jfrog-client-go"}
)

// Returns thr latest jfrog dependencies version in order to upgrade plugins dependencies.
func GetJfrogLatest() (dependencies []Details, err error) {
	for _, dep := range jfrogDependencies {
		latest, err := github.GetLatestRelease("jfrog", dep)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, Details{Path: "github.com/jfrog/" + dep, Version: latest})
	}
	return
}

func Upgrade(projectPath string, dependencies []Details) (err error) {
	for _, dependency := range dependencies {
		if err = utils.UpdateGoDependency(projectPath, dependency.Path, dependency.Version); err != nil {
			return
		}
	}
	return
}

// Generates a string of all the dependencies' details.
func ToString(deps []Details) (results string, err error) {
	var depName string
	for i, dep := range deps {
		if i != 0 {
			results += ", "
		}
		depName, err = dep.String()
		if err != nil {
			return
		}
		results += depName
	}
	return
}

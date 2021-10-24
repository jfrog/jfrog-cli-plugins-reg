package dependency

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jfrog/jfrog-cli-plugins-reg/github"
	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
)

// This struct represents a go dependency in the go.mod file.
type Details struct {
	Path                    string
	Version                 string
	MajorVersionsModulePath string
}

type JFrogDependencyDetails struct {
	Name                    string
	MajorVersionsModulePath string
}

// String returns the name and version of the dependency, omitting the path prefix.
func (d *Details) String() (string, error) {
	idx := strings.LastIndex(d.Path, "/")
	if idx == -1 {
		return "", errors.New("failed to locate dependency name")
	}
	return d.Path[strings.LastIndex(d.Path, "/")+1:] + " " + d.Version, nil
}

var (
	jfrogDependencies = [...]JFrogDependencyDetails{{Name: "jfrog-cli-core", MajorVersionsModulePath: "/v2"}, {Name: "jfrog-client-go", MajorVersionsModulePath: ""}}
)

// Returns the latest jfrog dependencies version in order to upgrade plugins dependencies.
func GetJfrogLatest() (dependencies []Details, err error) {
	for _, dep := range jfrogDependencies {
		latest, err := github.GetLatestRelease("jfrog", dep.Name)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, Details{Path: "github.com/jfrog/" + dep.Name + dep.MajorVersionsModulePath, Version: latest})
	}
	return
}

func Upgrade(projectPath string, dependencies []Details) (err error) {
	for _, dependency := range dependencies {
		fmt.Println("Updating: " + dependency.Path + " to version " + dependency.Version)
		if err = utils.UpdateGoDependency(projectPath, dependency.Path, dependency.Version); err != nil {
			return
		}
		fmt.Println("Successfully updated!")
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

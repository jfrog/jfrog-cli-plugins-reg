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
	Path    string
	Version string
}

type JFrogDependencyDetails struct {
	Name                   string
	MajorVersionModulePath string
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
	jfrogDependencies = [...]JFrogDependencyDetails{{Name: "jfrog-cli-core", MajorVersionModulePath: "/v2"}, {Name: "jfrog-client-go", MajorVersionModulePath: ""}}
)

// Returns the latest jfrog dependencies version in order to upgrade plugins dependencies.
func GetJfrogLatest() (dependencies []Details, err error) {
	for _, dep := range jfrogDependencies {
		latest, err := github.GetLatestRelease("jfrog", dep.Name)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, Details{Path: "github.com/jfrog/" + dep.Name + dep.MajorVersionModulePath, Version: latest})
	}
	return
}

func IsCoreVersionOneReleased(dependencies []Details) bool{
	for _, dep := range dependencies {
		if strings.Contains(dep.Path,"jfrog-cli-core") {
			return strings.HasPrefix(dep.Version,"v1.")
		}
	}
	return false
}

func Upgrade(projectPath string, dependencies []Details) (err error) {
	for _, dependency := range dependencies {
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

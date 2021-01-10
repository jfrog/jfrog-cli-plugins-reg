package dependency

import (
	"strings"

	"github.com/jfrog/jfrog-cli-plugins-reg/github"
	"github.com/jfrog/jfrog-cli-plugins-reg/utils"
)

type Details struct {
	Path    string
	Version string
}

func (d *Details) String() string {
	return d.Path[strings.LastIndex(d.Path, "/")+1:] + " " + d.Version
}

var (
	JfrogDependencies = [...]string{"jfrog-cli-core", "jfrog-client-go"}
)

// Returns thr latest jfrog dependencies version in order to upgrade plugings dependencies.
func GetJfrogLatest(projectPath string) (dependencies []Details, err error) {
	for _, dep := range JfrogDependencies {
		latest, err := github.GetLatest("jfrog", dep)
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
func ToString(deps []Details) (results string) {
	for i, dep := range deps {
		if i != 0 {
			results += ", "
		}
		results += dep.String()
	}
	return
}

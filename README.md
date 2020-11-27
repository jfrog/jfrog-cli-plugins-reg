# JFrog CLI plugins registry
## General
**JFrog CLI Plugins** allow enhancing the functionality of [JFrog CLI](https://www.jfrog.com/confluence/display/CLI/JFrog+CLI) to meet the specific user and organization needs. The source code of a plugin is maintained as an open source Go project on GitHub. All public plugins are registered in **JFrog CLI's Plugins Registry**. The Registry is hosted in this GitHub repository. The [plugins](plugins) directory includes a descriptor file for each plugin included in the Registry. 

## Installing a plugin 
After a plugin is included in this Registry, it becomes available for installation using JFrog CLI. JFrog CLI version 1.41.1 or above is required. To install a plugin included in this registry, run the following JFrog CLI command -  `jfrog plugin install plugin-name`. 

## The list of available plugins
* [build-deps-info](https://github.com/jfrog/jfrog-cli-plugins/tree/main/build-deps-info) - The build-deps-info plugin prints the dependencies' details of a specific build, which has been previously published to Artifactory.

* [build-report](https://github.com/jfrog/jfrog-cli-plugins/tree/main/build-report) - This JFrog CLI plugin prints a report of a published build info in Artifactory, or the diff between two builds.

* [file-spec-gen](https://github.com/jfrog/jfrog-cli-plugins/tree/main/file-spec-gen) - This plugin provides an easy way for generating file-specs json.

* [keyring](https://github.com/jfrog/jfrog-cli-plugins/tree/main/keyring) - This plugin allows using the OS keyring for managing the Artifactory connection details.

* [rb-gen](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rb-gen) - This plugin is designed to simplify interaction with release bundles, by generating them from other formats. Currently, it can generate release bundles from Helm charts.

* [rt-cleanup](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rt-cleanup) - This plugin is a simple Artifactory cleanup plugin. It can be used to delete all artifacts that have not been downloaded for the past n time units (both can bu configured) from a given repository.

* [rt-fs](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rt-fs) - This plugin executes file system commands in Artifactory. It is designed to mimic the functionality of the Linux/Unix 'ls' and 'cat' commands.

## Developing and publishing plugins
We encourage you, as developers, to create plugins and share them publicly with the rest of the community. Read the [JFrog CLI Plugins Developer Guide](https://github.com/jfrog/jfrog-cli/blob/master/guides/jfrog-cli-plugins-developer-guide.md) for information about developing and publishing JFrog CLI Plugins.

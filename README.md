# JFrog CLI plugins registry
## General
**JFrog CLI Plugins** allow enhancing the functionality of [JFrog CLI](https://www.jfrog.com/confluence/display/CLI/JFrog+CLI) to meet the specific user and organization needs. The source code of a plugin is maintained as an open source Go project on GitHub. All public plugins are registered in **JFrog CLI's Plugins Registry**. The Registry is hosted in this GitHub repository. The [plugins](plugins) directory includes a descriptor file for each plugin included in the Registry. 

## Installing a plugin 
After a plugin is included in this Registry, it becomes available for installation using JFrog CLI. JFrog CLI version 1.41.1 or above is required. To install a plugin included in this registry, run the following JFrog CLI command -  `jfrog plugin install plugin-name`. 

## The list of available plugins
* [build-deps-info](https://github.com/jfrog/jfrog-cli-plugins/tree/main/build-deps-info)
* [build-report](https://github.com/jfrog/jfrog-cli-plugins/tree/main/build-report)
* [file-spec-gen](https://github.com/jfrog/jfrog-cli-plugins/tree/main/file-spec-gen)
* [keyring](https://github.com/jfrog/jfrog-cli-plugins/tree/main/keyring)
* [release-bundle-generator](https://github.com/jfrog/jfrog-cli-plugins/tree/main/release-bundle-generator)
* [rt-cleanup](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rt-cleanup)
* [rt-fs](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rt-fs)

## Developing and publishing plugins
We encourage you, as developers, to create plugins and share them publicly with the rest of the community. Read the [JFrog CLI Plugins Developer Guide](https://github.com/jfrog/jfrog-cli/blob/master/guides/jfrog-cli-plugins-developer-guide.md) for information about developing and publishing JFrog CLI Plugins.

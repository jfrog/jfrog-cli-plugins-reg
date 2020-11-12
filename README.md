# JFrog CLI plugins registry
## General
JFrog CLI plugins support enhancing the functionality of [JFrog CLI](https://www.jfrog.com/confluence/display/CLI/JFrog+CLI) to meet the specific user and organization needs. The source code of a plugin is maintained as an open source Go project on GitHub. All public plugins are registered in JFrog CLI's Plugins Registry. The registry is hosted in this GitHub repository. The registry includes information about all the public JFrog CLI plugins, along with installation instructions.

## Installing a plugin 
To install a plugin which is included in this registry, run the following JFrog CLI command from your machine -  `jfrog plugin install plugin-name`. JFrog CLI version 1.41.0 or above is required. 

## The list of available plugins
* [build-report](https://github.com/jfrog/jfrog-cli-plugins/tree/main/build-report)
* [keyring](https://github.com/jfrog/jfrog-cli-plugins/tree/main/keyring)
* [release-bundle-generator](https://github.com/jfrog/jfrog-cli-plugins/tree/main/release-bundle-generator)
* [rt-cleanup](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rt-cleanup)
* [rt-ls](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rt-fs)

## Developing and publishing plugins
We encourage you, as developers, to create plugins and share them publicly with the rest of your community. Read the [JFrog CLI Plugins Developer Guide](https://github.com/jfrog/jfrog-cli/blob/master/guides/jfrog-cli-plugins-developer-guide.md) for information about developing and publishing JFrog CLI Plugins.

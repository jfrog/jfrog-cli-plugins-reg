# JFrog CLI plugins registry
## General
**JFrog CLI Plugins** allow enhancing the functionality of [JFrog CLI](https://www.jfrog.com/confluence/display/CLI/JFrog+CLI) to meet the specific user and organization needs. The source code of a plugin is maintained as an open source Go project on GitHub. All public plugins are registered in **JFrog CLI's Plugins Registry**. The Registry is hosted in this GitHub repository. The [plugins](plugins) directory includes a descriptor file for each plugin included in the Registry. 

## Installing a plugin 
After a plugin is included in this Registry, it becomes available for installation using JFrog CLI. JFrog CLI version 1.41.1 or above is required. To install a plugin included in this registry, run the following JFrog CLI command -  `jfrog plugin install plugin-name`. 

## The list of available plugins
* [build-deps-info](https://github.com/jfrog/jfrog-cli-plugins/tree/main/build-deps-info) - The build-deps-info plugin prints the dependencies' details of a specific build, which has been previously published to Artifactory.

* [build-report](https://github.com/jfrog/jfrog-cli-plugins/tree/main/build-report) - This JFrog CLI plugin prints a report of a published build info in Artifactory, or the diff between two builds.

* [file-spec-gen](https://github.com/jfrog/jfrog-cli-plugins/tree/main/file-spec-gen) - This plugin provides an easy way for generating file-specs json.

* [JCheck](https://github.com/rdar-lab/JCheck) - A Micro-UTP, plug-able sanity checker for any on-prem JFrog platform instance.

* [jfrog-yocto](https://github.com/rdar-lab/jfrog-cli-yocto-plugin) - This plugin allows integrating Yocto builds with the JFrog Platform.

* [live-logs](https://github.com/jfrog/live-logs) - The JFrog Platform includes an integrated Live Logs plugin, which allows customers to get the JFrog product logs (Artifactory, Xray, Mission Control, Distribution, and Pipelines) using the JFrog CLI Plugin. The plugin also provides the ability to cat and tail -f any log on any product node.

* [metrics-viewer](https://github.com/eldada/metrics-viewer/tree/master) - A plugin or standalone binary to show open-metrics formatted data in a terminal based graph.

* [repostats](https://github.com/chanti529/repostats) - This plugin can help find out the most popularly downlaoded artifacts in a given repository, Artifacts that are consuming the most space in a given repository with various levels of customization available. Results obtained can also be sorted and filtered.

* [rm-empty](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rm-empty) - This plugin deletes all the empty folders under a specific path in Artifactory.

* [rt-cleanup](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rt-cleanup) - This plugin is a simple Artifactory cleanup plugin. It can be used to delete all artifacts that have not been downloaded for the past n time units (both can be configured) from a given repository.

* [rt-fs](https://github.com/jfrog/jfrog-cli-plugins/tree/main/rt-fs) - This plugin executes file system commands in Artifactory. It is designed to mimic the functionality of the Linux/Unix 'ls' and 'cat' commands.

## Developing and publishing plugins
We encourage you, as developers, to create plugins and share them publicly with the rest of the community. Read the [JFrog CLI Plugins Developer Guide](https://github.com/jfrog/jfrog-cli/blob/master/guides/jfrog-cli-plugins-developer-guide.md) for information about developing and publishing JFrog CLI Plugins.

[![rt-fs-plugin](images/rt-fs-plugin.png)](https://youtu.be/zQ1JV83frFI)

[![build-report-plugin](images/build-report-plugin.png)](https://youtu.be/_oPNuiDm04g)
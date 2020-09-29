# JFrog CLI plugins registry
This repository contains information about all shared JFrog CLI plugins, and how to install them.

## Installation
To install a plugin, run `jfrog plugin install plugin-name` from your machine. jfrog-cli version TODO or above is required. 

## Adding plugin to registry
To add a plugin, create a PR to the `plugins` directory in this repo, with a `yml` file with the name of your plugin.

The yml file should have the following structure:
```
pluginName: hello-frog
repository: https://github.com/jfrog/jfrog-cli-plugin-template
version: v0.1.0
```

NOTE: plugin deployment to registry is currently manual.

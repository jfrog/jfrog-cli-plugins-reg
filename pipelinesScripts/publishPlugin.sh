#!/bin/bash

# Publish plugin -build, verify unique version and upload for all architectures.
publishPlugin () {
  jf plugin p "$JFROG_CLI_PLUGIN_PLUGIN_NAME" "$JFROG_CLI_PLUGIN_VERSION"
}

publishPlugin
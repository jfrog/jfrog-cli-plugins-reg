#!/bin/bash
set -eu

#function verifyPluginVersionMatching()
verifyPluginVersionMatching () {
  echo "Verifying provided plugin version matches built version..."
  res=$(eval "./$JFROG_CLI_PLUGIN_PLUGIN_NAME -v")
  exitCode=$?
  if [ $exitCode -ne 0 ]; then
    echo "Error: Failed verifying version matches"
    exit $exitCode
  fi

  # Get the version which is after the last space. (expected output to -v for example: "plugin-name version v1.0.0")
  echo "Output: $res"
  builtVersion="${res##* }"
  # Compare versions
  if [ "$builtVersion" != "$JFROG_CLI_PLUGIN_VERSION" ]; then
    echo "Versions dont match. Provided: $JFROG_CLI_PLUGIN_VERSION, Actual: $builtVersion"
    exit 1
  fi
  echo "Versions match."
}

#function build(pkg, goos, goarch, exeName)
build () {
  pkg="$1"
  export GOOS="$2"
  export GOARCH="$3"
  exeName="$4"
  echo "Building $exeName for $GOOS-$GOARCH ..."

  CGO_ENABLED=0 go build -o "$exeName" -ldflags '-w -extldflags "-static"' main.go

  # Run verification after building plugin for the correct platform of this image.
  if [ "$pkg" = "linux-386" ]; then
    verifyPluginVersionMatching
  fi
}

#function verifyUniqueVersion()
verifyUniqueVersion () {
  echo "Verifying version uniqueness..."
  versionFolderUrl="$JFROG_CLI_PLUGINS_REGISTRY_URL/$JFROG_CLI_PLUGINS_REGISTRY_REPO/$JFROG_CLI_PLUGIN_PLUGIN_NAME/$JFROG_CLI_PLUGIN_VERSION/"

  echo "Checking existence of $versionFolderUrl"
  res=$(curl -o /dev/null -s --head --fail "$versionFolderUrl")
  exitCode=$?
  if [ $exitCode -ne 0 ]; then
    echo "Error: Failed verifying uniqueness of the plugin's version"
    exit $exitCode
  fi

  echo "Artifactory response: $res"
  if [ $res -eq 200 ]; then
    echo "Error: Version already exists in registry"
    exit 1
  fi
}

#function downloadJfrogCli()
downloadJfrogCli () {
  echo "Downloading the latest version of JFrog CLI..."
  curl -sSfL https://getcli.jfrog.io | sh
  # Verify CLI was downloaded
  if [ ! -f ./jfrog ]; then
      echo "Error: JFrog CLI downloaded failed."
      exit 1
  fi
}

#function buildAndUpload(pkg, goos, goarch, fileExtension)
buildAndUpload () {
  pkg="$1"
  goos="$2"
  goarch="$3"
  fileExtension="$4"
  exeName="$JFROG_CLI_PLUGIN_PLUGIN_NAME$fileExtension"

  build $pkg $goos $goarch $exeName

  destPath="$JFROG_CLI_PLUGINS_REGISTRY_REPO/$JFROG_CLI_PLUGIN_PLUGIN_NAME/$JFROG_CLI_PLUGIN_VERSION/$pkg/$exeName"
  echo "Uploading $exeName to $JFROG_CLI_PLUGINS_REGISTRY_URL/$destPath ..."

  ./jfrog rt u "./$exeName" "$destPath" --url="$JFROG_CLI_PLUGINS_REGISTRY_URL" --access-token=$int_releases_jfrog_token
  exitCode=$?
  if [ $exitCode -ne 0 ]; then
    echo "Error: Failed uploading plugin"
    exit $exitCode
  fi
}

#function copyToLatestDir()
copyToLatestDir () {
  pluginPath="$JFROG_CLI_PLUGINS_REGISTRY_REPO/$JFROG_CLI_PLUGIN_PLUGIN_NAME"
  echo "Copy version to latest dir: $pluginPath"

  ./jfrog rt cp "$pluginPath/$JFROG_CLI_PLUGIN_VERSION/(*)" "$pluginPath/latest/{1}" --flat --url="$JFROG_CLI_PLUGINS_REGISTRY_URL" --access-token=$int_releases_jfrog_token
  exitCode=$?
  if [ $exitCode -ne 0 ]; then
    echo "Error: Failed uploading plugin"
    exit $exitCode
  fi
}

# Verify uniqueness of the requested plugin's version
verifyUniqueVersion

# Download JFrog CLI
downloadJfrogCli

# Build and upload for every architecture.
# Keep 'linux-386' first to prevent unnecessary uploads in case the built version doesn't match the provided one.
buildAndUpload 'linux-386' 'linux' '386' ''
buildAndUpload 'linux-amd64' 'linux' 'amd64' ''
buildAndUpload 'linux-arm64' 'linux' 'arm64' ''
buildAndUpload 'linux-arm' 'linux' 'arm' ''
buildAndUpload 'mac-386' 'darwin' 'amd64' ''
buildAndUpload 'windows-amd64' 'windows' 'amd64' '.exe'

# Copy the uploaded version to override latest dir
copyToLatestDir

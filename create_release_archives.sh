#!/bin/bash

read -p "Release Version> " version

re="^v?[0-9]+\.[0-9]+\.[0-9]+$"

# Ensure that the version conforms to version naming requirements
if ! [[ $version =~ $re ]]; then
  echo "Error: version number must fit semver naming requirements"
  exit 1
fi

echo "Creating release archives for Fu $version..."

if [[ -d "build" && -x "build" ]]; then
  rm -rf build
fi

echo "Compiling..."
if hash gox 2>/dev/null; then
  gox -output "build/fu-$version-{{.OS}}-{{.Arch}}/fu" -os "linux darwin" -arch "amd64" -ldflags "-s -w"
else
  echo "Error: gox is not installed, can install it by running `go get github.com/mitchellh/gox`"
  exit 1
fi

cd build

echo -n "Compressing... "
tar czf "fu-$version-darwin-amd64.tar.gz" "fu-$version-darwin-amd64/"
tar czf "fu-$version-linux-amd64.tar.gz" "fu-$version-linux-amd64/"
echo "done."

echo -n "Cleaning up... "
rm -rf "fu-$version-darwin-amd64"
rm -rf "fu-$version-linux-amd64"
echo "done."

echo "All done."

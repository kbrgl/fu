#!/bin/bash
if [[ -d "build" && -x "build" ]]; then
    rm -r build
fi
echo "Compiling..."
if hash gox 2>/dev/null; then
    gox -output "build/fu_{{.OS}}_{{.Arch}}/fu" -os "linux darwin" -arch "amd64" -ldflags "-s -w"
else
    echo "gox is not installed, can install it by running `go get github.com/mitchellh/gox`."
    exit 1
fi
cd build
echo -n "Compressing... "
tar czf fu_darwin_amd64.tar.gz fu_darwin_amd64/
tar czf fu_linux_amd64.tar.gz fu_linux_amd64/
echo "done."
echo -n "Cleaning up... "
rm -rf fu_darwin_amd64
rm -rf fu_linux_amd64
echo "done."
echo "All done."

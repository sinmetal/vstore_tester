#!/bin/sh -eux

targets=`find . -type f \( -name '*.go' -and -not -iwholename '*vendor*' \)`
packages=$(go list ./...)

# Apply tools
export PATH=$(pwd)/build-cmd:$PATH
goimports -w $targets
go tool vet $targets
golint $packages

#!/usr/bin/env bash

GOPATH="/root/go"
export PATH=$PATH:$GOPATH/bin
GINKGO_OLD_VERSION="Ginkgo Version 1.16.5"

if ! which ginkgo || ginkgo version -eq "$GINKGO_OLD_VERSION"; then {  
	echo "Downloading ginkgo tool"
	go install github.com/onsi/ginkgo/v2/ginkgo@v2.8.1
} fi

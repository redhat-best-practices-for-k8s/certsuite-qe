#!/usr/bin/env bash

GOPATH="${GOPATH:-~/go}"
export PATH=$PATH:$GOPATH/bin

if ! which ginkgo ; then
	echo "Downloading ginkgo tool"
	go install github.com/onsi/ginkgo/ginkgo@v1.16.4
fi

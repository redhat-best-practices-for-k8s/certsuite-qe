#!/usr/bin/env bash

GOPATH="${GOPATH:-~/go}"
export PATH=$PATH:$GOPATH/bin

if ! which ginkgo ; then
	echo "Downloading ginkgo tool"
	go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo@latest
fi

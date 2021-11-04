#!/usr/bin/env bash

GOPATH="${GOPATH:-~/go}"
export GOFLAGS="${GOFLAGS:-"-mod=vendor"}"
export PATH=$PATH:$GOPATH/bin

if ! which ginkgo ; then
	echo "Downloading ginkgo tool"
	go install github.com/onsi/ginkgo/ginkgo
fi

#!/usr/bin/env bash
set -e

. $(dirname "$0")/common.sh

if which golangci-lint; then
	echo "golint installed"
else
	echo "Downloading golint tool"
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
fi

golangci-lint run --skip-dirs-use-default

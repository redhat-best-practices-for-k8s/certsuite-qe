#!/usr/bin/env bash
set -e

# shellcheck disable=SC1091 # Not following.
. "$(dirname "$0")"/common.sh

if which golangci-lint; then
	echo "golint installed"
else
	echo "Downloading golint tool"
	if [[ -z "${GOPATH}" ]]; then
		DEPLOY_PATH=/tmp/
	else
		DEPLOY_PATH=${GOPATH}/bin
	fi
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${DEPLOY_PATH}" v1.59.1
fi

PATH=${PATH}:${DEPLOY_PATH} golangci-lint run -v

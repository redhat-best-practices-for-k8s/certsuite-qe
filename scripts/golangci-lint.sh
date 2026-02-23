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
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/2691aacde61559ce53cfda05513ee7c70677ea31/install.sh | sh -s -- -b "${DEPLOY_PATH}" v2.5.0
fi

PATH=${PATH}:${DEPLOY_PATH} golangci-lint run -v

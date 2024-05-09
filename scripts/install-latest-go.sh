#!/usr/bin/env bash

REQUIRED_GO_VERSION=1.22.3
INSTALLED_GO_VERSION="$(
	go version | {
		read -r _ _ v _
		echo "${v#go}"
	}
)"

function version {
	echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4); }'
}

if [ "$(version "${INSTALLED_GO_VERSION}")" -ge "$(version ${REQUIRED_GO_VERSION})" ]; then
	echo "Version is up to date"
else
	unameOut="$(uname -s)"
	case "${unameOut}" in
	Linux*)
		wget https://go.dev/dl/go"${REQUIRED_GO_VERSION}".linux-amd64.tar.gz
		rm -rf /usr/local/go && tar -C /usr/local -xzf go"${REQUIRED_GO_VERSION}".linux-amd64.tar.gz
		rm go"${REQUIRED_GO_VERSION}".linux-amd64.tar.gz
		;;
	*)
		echo "Please install go with version ${REQUIRED_GO_VERSION}"
		;;
	esac
fi

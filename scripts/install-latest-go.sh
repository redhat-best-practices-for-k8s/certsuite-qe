#!/usr/bin/env bash

EXPECTED_GO_VERSION="1.20.2"
INSTALLED_GO_VERSION="$(go version | { read -r _ _ v _; echo "${v#go}"; })"

function version { echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4); }'; }

if [ "$(version "${INSTALLED_GO_VERSION}")" -ge "$(version ${EXPECTED_GO_VERSION})" ]; then
    echo "Version is up to date"
else
    unameOut="$(uname -s)"
    case "${unameOut}" in
        Linux*)
                    wget https://go.dev/dl/"${EXPECTED_GO_VERSION}".linux-amd64.tar.gz
                    rm -rf /usr/local/go && tar -C /usr/local -xzf "${EXPECTED_GO_VERSION}".linux-amd64.tar.gz
                    rm "${LATEST_GO_VERSION}".linux-amd64.tar.gz
                    ;;
        *)          echo "Install go with version ${EXPECTED_GO_VERSION}"
    esac
fi
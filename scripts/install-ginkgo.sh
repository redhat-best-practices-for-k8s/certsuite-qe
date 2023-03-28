#!/usr/bin/env bash

REQUIRED_GINKGO_VERSION="$(grep '/ginkgo' go.mod | cut -d ' ' -f2 | cut -c 2-)"
INSTALLED_GINKGO_VERSION="$(ginkgo version | { read -r _ _ v _; echo "${v}"; })"

GOPATH="${GOPATH:-/root/go}"
export PATH=$PATH:$GOPATH/bin

function version { echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4); }'; }

if [ "$(version "${INSTALLED_GINKGO_VERSION}")" -ge "$(version "${REQUIRED_GINKGO_VERSION}")" ]; then
    echo "Version is up to date"
else
        GINKGO_TMP_DIR=$(mktemp -d)
        cd "$GINKGO_TMP_DIR" || exit
        go mod init tmp
        GOFLAGS=-mod=mod go install github.com/onsi/ginkgo/v2/ginkgo@v"$REQUIRED_GINKGO_VERSION"
        rm -rf "$GINKGO_TMP_DIR"
        echo "Downloading ginkgo tool"
        cd - || exit
fi
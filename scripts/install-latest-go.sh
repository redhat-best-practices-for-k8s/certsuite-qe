#!/usr/bin/env bash


LATEST_GO_VERSION=$(curl https://go.dev/VERSION?m=text)
INSTALLED_GO_VERSION=$(go version | grep -oP "go\d+\.\d+\.\d+")

if [ "$INSTALLED_GO_VERSION" != "$LATEST_GO_VERSION" ]; then {
    wget https://go.dev/dl/${LATEST_GO_VERSION}.linux-amd64.tar.gz
    rm -rf /usr/local/go && tar -C /usr/local -xzf ${LATEST_GO_VERSION}.linux-amd64.tar.gz
    rm ${LATEST_GO_VERSION}.linux-amd64.tar.gz
} fi



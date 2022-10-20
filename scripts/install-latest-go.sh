#!/usr/bin/env bash
set -e

LATEST=$(curl https://go.dev/VERSION?m=text)
wget https://go.dev/dl/${LATEST}.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf ${LATEST}.linux-amd64.tar.gz
rm ${LATEST}.linux-amd64.tar.gz


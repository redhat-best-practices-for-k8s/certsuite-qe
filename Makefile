# Export GO111MODULE=on to enable project to be built from within GOPATH/src
export GO111MODULE=on
GO_PACKAGES=$(shell go list ./... | grep -v vendor)

.PHONY: lint \
		deps-update \
		gofmt \
		test-all \
		test-features \
		install \
		vet \

vet:
	go vet ${GO_PACKAGES}

lint:
	@echo "Running go lint"
	scripts/golangci-lint.sh

update-go:
	scripts/install-latest-go.sh

deps-update:
	go mod tidy && \
	go mod vendor

test: unit-tests

gofmt:
	@echo "Running gofmt"
	gofmt -s -l `find . -path ./vendor -prune -o -type f -name '*.go' -print`

test-all: update-go install-ginkgo download-unstable
	./scripts/run-tests.sh all

test-features: update-go install-ginkgo download-unstable
	FEATURES="$(FEATURES)" ./scripts/run-tests.sh features

download-unstable:
	./scripts/download-unstable.sh

install-ginkgo:
	go install "$$(awk '/ginkgo/ {printf "%s/ginkgo@%s", $$1, $$2}' go.mod)"

install: deps-update install-ginkgo
	@echo "Installing needed dependencies"

unit-tests:
	UNIT_TEST=true go test ./... -tags=utest -coverprofile=cover.out

coverage-html: test
	go tool cover -html cover.out



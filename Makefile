# Export GO111MODULE=on to enable project to be built from within GOPATH/src
export GO111MODULE=on
GO_PACKAGES=$(shell go list ./... | grep -v vendor)

.PHONY: lint \
		deps-update \
		gofmt \
		test-all \
		test-features \
		install \
		vet

vet:
	go vet ${GO_PACKAGES}

lint:
	@echo "Running go lint"
	scripts/golangci-lint.sh

deps-update:
	go mod tidy && \
	go mod vendor

gofmt:
	@echo "Running gofmt"
	gofmt -s -l `find . -path ./vendor -prune -o -type f -name '*.go' -print`

test-all:
	./scripts/run-tests.sh all

test-features:
	FEATURES="$(FEATURES)" ./scripts/run-tests.sh features

install-ginkgo:
	scripts/install-ginkgo.sh

install: deps-update install-ginkgo
	@echo "Installing needed dependencies"

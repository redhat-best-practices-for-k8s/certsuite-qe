# Export GO111MODULE=on to enable project to be built from within GOPATH/src
export GO111MODULE=on

.PHONY: govet \
		golint \
		deps-update \
		gofmt \
		test-all \
		test-features \
		install

govet:
	@echo "Running go vet"
	# Disabling GO111MODULE just for go vet execution
	GO111MODULE=off go vet github.com/test-network-function/cnfcert-tests-verification/test...

golint:
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

install: deps-update
	@echo "Installing needed dependencies"
	scripts/install-ginkgo.sh

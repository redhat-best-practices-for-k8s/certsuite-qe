# Export GO111MODULE=on to enable project to be built from within GOPATH/src
export GO111MODULE=on
GO_PACKAGES=$(shell go list ./... | grep -v vendor)

# Minimum total statement coverage enforced by `unit-tests`. Kept one point
# under the current total so ordinary churn passes but a whole untested package
# does not. Override on the command line: make unit-tests COVERAGE_FLOOR=95
COVERAGE_FLOOR ?= 92.0

# Color definitions for better output formatting
RED     := \033[31m
GREEN   := \033[32m
YELLOW  := \033[33m
BLUE    := \033[34m
MAGENTA := \033[35m
CYAN    := \033[36m
WHITE   := \033[37m
BOLD    := \033[1m
RESET   := \033[0m

# Default target - show help when make is run without arguments
.DEFAULT_GOAL := help

.PHONY: help \
		lint \
		lint-ci \
		compile-tests \
		verify \
		deps-update \
		gofmt \
		fmt \
		test-all \
		test-features \
		install \
		vet \
		test \
		unit-tests \
		coverage-html \
		download-unstable \
		install-ginkgo

help: ## Display this help message with available targets
	@echo "$(BOLD)$(BLUE)╔════════════════════════════════════════════════════════════╗$(RESET)"
	@echo "$(BOLD)$(BLUE)║                    $(WHITE)Available Make Targets$(BLUE)                    ║$(RESET)"
	@echo "$(BOLD)$(BLUE)╚════════════════════════════════════════════════════════════╝$(RESET)"
	@echo ""
	@echo "$(BOLD)$(YELLOW)📋 Development and Code Quality:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "; section=""} /^# [A-Z]/ {section=$$0; gsub(/^# /, "", section)} /^[a-zA-Z_-]+:.*?## / && (section=="Development and Code Quality Targets") {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)$(YELLOW)📦 Dependency Management:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "; section=""} /^# [A-Z]/ {section=$$0; gsub(/^# /, "", section)} /^[a-zA-Z_-]+:.*?## / && (section=="Dependency Management") {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)$(YELLOW)🧪 Testing:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "; section=""} /^# [A-Z]/ {section=$$0; gsub(/^# /, "", section)} /^[a-zA-Z_-]+:.*?## / && (section=="Testing Targets") {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)$(YELLOW)🔧 Utilities:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "; section=""} /^# [A-Z]/ {section=$$0; gsub(/^# /, "", section)} /^[a-zA-Z_-]+:.*?## / && (section=="Utility Targets") {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)$(GREEN)💡 Examples:$(RESET)"
	@echo "  $(WHITE)make help$(RESET)           $(MAGENTA)# Show this help$(RESET)"
	@echo "  $(WHITE)make test$(RESET)           $(MAGENTA)# Run unit tests$(RESET)"
	@echo "  $(WHITE)make lint$(RESET)           $(MAGENTA)# Run linting$(RESET)"
	@echo "  $(WHITE)make install$(RESET)        $(MAGENTA)# Install dependencies$(RESET)"
	@echo ""

# Development and Code Quality Targets
vet: ## Run go vet to examine Go source code and report suspicious constructs
	@echo "$(BOLD)$(BLUE)🔍 Running go vet...$(RESET)"
	@go vet ${GO_PACKAGES} && echo "$(GREEN)✅ go vet completed successfully$(RESET)" || (echo "$(RED)❌ go vet failed$(RESET)" && exit 1)

lint: ## Run golangci-lint to check code quality and style
	@echo "$(BOLD)$(BLUE)🔧 Running golangci-lint...$(RESET)"
	@scripts/golangci-lint.sh && echo "$(GREEN)✅ Linting completed successfully$(RESET)" || (echo "$(RED)❌ Linting failed$(RESET)" && exit 1)
	@echo "$(BOLD)$(BLUE)🔧 Running shfmt on shell scripts...$(RESET)"
	@shfmt -d scripts/*.sh && echo "$(GREEN)✅ Shell scripts formatted correctly$(RESET)" || (echo "$(RED)❌ Shell script formatting issues found. Run 'shfmt -w scripts/*.sh' to fix$(RESET)" && exit 1)

# `lint-ci` mirrors .github/workflows/pre-main.yml. Its tools are installed out
# of band, so a missing tool is reported as SKIP rather than failing the target.
# markdownlint's MD060 (table column style) ships in 0.38+, newer than the
# version CI pins, so the call below disables it to match CI's result.
lint-ci: lint compile-tests ## Run every static check CI runs (shellcheck, yamllint, typos, markdownlint, ginkgo build)
	@echo "$(BOLD)$(BLUE)🔎 Running CI-only static checks...$(RESET)"
	@if ! command -v shellcheck >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠ SKIP: shellcheck not installed (brew install shellcheck)$(RESET)"; \
	elif ! shellcheck scripts/*.sh rebase_all.sh; then \
		echo "$(RED)❌ shellcheck reported issues$(RESET)"; exit 1; \
	fi
	@if ! command -v yamllint >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠ SKIP: yamllint not installed (pip install yamllint)$(RESET)"; \
	elif ! yamllint .github/workflows/*.yml config/*.yaml; then \
		echo "$(RED)❌ yamllint reported issues$(RESET)"; exit 1; \
	fi
	@if ! command -v typos >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠ SKIP: typos not installed (brew install typos-cli)$(RESET)"; \
	elif ! typos; then \
		echo "$(RED)❌ typos reported misspellings$(RESET)"; exit 1; \
	fi
	@if ! command -v markdownlint >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠ SKIP: markdownlint not installed (npm install -g markdownlint-cli)$(RESET)"; \
	elif ! markdownlint --disable MD060 -- README.md; then \
		echo "$(RED)❌ markdownlint reported issues$(RESET)"; exit 1; \
	fi

compile-tests: ## Compile every Ginkgo suite (the only local check that builds the integration packages)
	@echo "$(BOLD)$(BLUE)🏗  Compiling Ginkgo suites...$(RESET)"
	@ginkgo build -r ./tests && status=0 || status=$$?; \
	find ./tests -type f -name '*.test' -delete; \
	if [ $$status -ne 0 ]; then echo "$(RED)❌ Ginkgo suites failed to compile$(RESET)"; exit 1; fi; \
	echo "$(GREEN)✅ Ginkgo suites compiled successfully$(RESET)"

verify: vet gofmt lint-ci unit-tests ## Run the full pre-push suite (static checks first, then unit tests)
	@echo "$(BOLD)$(GREEN)🎉 All pre-push checks passed!$(RESET)"

gofmt: ## Check Go code formatting (use 'gofmt -w' to fix issues)
	@echo "$(BOLD)$(BLUE)📝 Checking Go code formatting...$(RESET)"
	@files=$$(gofmt -s -l `find . -path ./vendor -prune -o -type f -name '*.go' -print`); \
	if [ -n "$$files" ]; then \
		echo "$(RED)❌ The following files need formatting:$(RESET)"; \
		echo "$$files"; \
		echo "$(YELLOW)💡 Run 'gofmt -w .' to fix formatting$(RESET)"; \
		exit 1; \
	else \
		echo "$(GREEN)✅ All Go files are properly formatted$(RESET)"; \
	fi

fmt: ## Format Go code using gofmt
	@echo "$(BOLD)$(BLUE)✨ Formatting Go code...$(RESET)"
	@gofmt -s -w `find . -path ./vendor -prune -o -type f -name '*.go' -print` && echo "$(GREEN)✅ Go code formatted successfully$(RESET)" || (echo "$(RED)❌ Failed to format Go code$(RESET)" && exit 1)

# Dependency Management
deps-update: ## Update and vendor Go module dependencies
	@echo "$(BOLD)$(BLUE)📦 Updating dependencies...$(RESET)"
	@go mod tidy && go mod vendor && echo "$(GREEN)✅ Dependencies updated successfully$(RESET)" || (echo "$(RED)❌ Failed to update dependencies$(RESET)" && exit 1)

install-ginkgo: ## Install Ginkgo testing framework
	@echo "$(BOLD)$(BLUE)⚙️  Installing Ginkgo...$(RESET)"
	@go install "$$(awk '/ginkgo/ {printf "%s/ginkgo@%s", $$1, $$2}' go.mod)" && echo "$(GREEN)✅ Ginkgo installed successfully$(RESET)" || (echo "$(RED)❌ Failed to install Ginkgo$(RESET)" && exit 1)

install: deps-update install-ginkgo ## Install all project dependencies and tools
	@echo "$(BOLD)$(GREEN)🎉 All dependencies installed successfully!$(RESET)"

# Testing Targets
test: unit-tests ## Run unit tests (alias for unit-tests)

unit-tests: ## Run unit tests with race detection and the coverage floor
	@echo "$(BOLD)$(BLUE)🧪 Running unit tests...$(RESET)"
	@UNIT_TEST=true go test ./... -tags=utest -race -count=1 -coverprofile=cover.out && echo "$(GREEN)✅ Unit tests completed successfully$(RESET)" || (echo "$(RED)❌ Unit tests failed$(RESET)" && exit 1)
	@coverage=$$(go tool cover -func=cover.out | awk '/^total:/ {gsub(/%/, "", $$NF); print $$NF}'); \
	echo "$(BOLD)$(BLUE)📊 Total coverage:$(RESET) $${coverage}% (floor $(COVERAGE_FLOOR)%)"; \
	awk -v actual="$$coverage" -v floor="$(COVERAGE_FLOOR)" 'BEGIN { if (actual + 0 < floor + 0) { printf "$(RED)❌ Coverage is below the %.1f%% floor$(RESET)\n", floor; exit 1 } }' || exit 1

test-all: install-ginkgo download-unstable ## Run all test suites (requires cluster setup)
	@echo "$(BOLD)$(BLUE)🚀 Running all test suites...$(RESET)"
	@echo "$(YELLOW)⚠️  This requires a properly configured cluster$(RESET)"
	@./scripts/run-tests.sh all && echo "$(GREEN)✅ All tests completed successfully$(RESET)" || (echo "$(RED)❌ Some tests failed$(RESET)" && exit 1)

test-features: install-ginkgo download-unstable ## Run specific feature tests (set FEATURES env var)
	@echo "$(BOLD)$(BLUE)🎯 Running feature tests...$(RESET)"
	@if [ -z "$(FEATURES)" ]; then \
		echo "$(YELLOW)⚠️  No FEATURES specified. Set FEATURES environment variable.$(RESET)"; \
		echo "$(CYAN)Example: FEATURES=networking make test-features$(RESET)"; \
	else \
		echo "$(CYAN)📋 Running tests for features: $(BOLD)$(FEATURES)$(RESET)"; \
		FEATURES="$(FEATURES)" ./scripts/run-tests.sh features && echo "$(GREEN)✅ Feature tests completed successfully$(RESET)" || (echo "$(RED)❌ Feature tests failed$(RESET)" && exit 1); \
	fi

coverage-html: test ## Generate HTML coverage report from test results
	@echo "$(BOLD)$(BLUE)📊 Generating HTML coverage report...$(RESET)"
	@go tool cover -html cover.out && echo "$(GREEN)✅ Coverage report generated: $(CYAN)cover.out$(RESET)" || (echo "$(RED)❌ Failed to generate coverage report$(RESET)" && exit 1)

# Utility Targets
download-unstable: ## Download unstable certsuite image for testing
	@echo "$(BOLD)$(BLUE)⬇️  Downloading unstable certsuite...$(RESET)"
	@./scripts/download-unstable.sh && echo "$(GREEN)✅ Unstable certsuite downloaded successfully$(RESET)" || (echo "$(RED)❌ Failed to download unstable certsuite$(RESET)" && exit 1)



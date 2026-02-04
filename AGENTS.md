# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

CertSuite QE is a Go-based test harness that validates the [certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite) project. It runs certsuite test scenarios against various pre-configured OCP (OpenShift Container Platform) environments and verifies that tests produce expected results (pass/fail/skip).

The project uses Ginkgo/Gomega BDD testing framework. Tests deploy resources to an OCP cluster, execute certsuite (via container image or binary), then parse the resulting claim.json to verify expected outcomes.

## Commands

### Installation and Setup
```bash
make install              # Install all dependencies (runs deps-update + install-ginkgo)
make install-ginkgo       # Install Ginkgo testing framework
make deps-update          # Update and vendor Go module dependencies
```

### Running Tests
```bash
# Unit tests (no cluster required)
make test                 # Run unit tests with coverage
make unit-tests           # Same as above
make coverage-html        # Generate HTML coverage report

# Integration tests (requires OCP cluster)
KUBECONFIG=$HOME/.kube/config CERTSUITE_REPO_PATH=/path/to/certsuite make test-all

# Run specific feature tests
FEATURES=networking KUBECONFIG=$HOME/.kube/config make test-features
FEATURES=accesscontrol,lifecycle make test-features

# Debug mode
DEBUG_CERTSUITE=true CERTSUITE_LOG_LEVEL=debug make test-features
```

### Running a Single Test
```bash
# Using ginkgo focus
ginkgo -v --focus="test name pattern" ./tests/networking/...

# Or run a specific suite
ginkgo -v ./tests/accesscontrol/
```

### Linting
```bash
make lint                 # Run golangci-lint and shfmt
make gofmt                # Check Go code formatting
make fmt                  # Format Go code
make vet                  # Run go vet
```

## Architecture

### Test Flow
1. **Setup**: Create random namespace, generate config/report directories
2. **Deploy**: Create test resources (deployments, daemonsets, operators, etc.)
3. **Execute**: Run certsuite via container image or binary with specific test labels
4. **Verify**: Parse claim.json and assert expected test status (passed/failed/skipped)
5. **Cleanup**: Delete namespace and temporary directories

### Directory Structure

**tests/** - Main test code organized by certsuite test suite:
- `accesscontrol/`, `networking/`, `lifecycle/`, `operator/`, etc. - Each maps to a certsuite test suite
- Each suite contains:
  - `<suite>_suite_test.go` - Ginkgo suite entry point
  - `tests/` - Individual test cases (`.go` files with `_test.go` stubs)
  - `parameters/` - Suite-specific constants and config
  - `helper/` - Suite-specific helper functions

**tests/globalhelper/** - Shared test utilities:
- `init.go` - Client initialization and namespace setup helpers
- `runhelper.go` - Functions to launch certsuite (via image or binary)
- `reporthelper.go` - Claim.json parsing and validation
- Resource helpers: `deployment.go`, `daemonset.go`, `pod.go`, `namespaces.go`, etc.

**tests/utils/** - Low-level Kubernetes resource builders:
- `client/` - Kubernetes client wrapper
- `config/` - Configuration loading from YAML and env vars
- Individual packages for each resource type (pod, deployment, service, etc.)

**tests/globalparameters/** - Shared constants and type definitions

### Test Pattern
Tests follow a consistent pattern using Ginkgo's BDD syntax:

```go
var _ = Describe("Feature", func() {
    BeforeEach(func() {
        // Create random namespace and config
        randomNamespace, randomReportDir, randomCertsuiteConfigDir =
            globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNamespace)
        // Define certsuite config
        err = globalhelper.DefineCertsuiteConfig(...)
    })

    AfterEach(func() {
        // Cleanup
        globalhelper.AfterEachCleanupWithRandomNamespace(...)
    })

    It("test case description", func() {
        // 1. Deploy resources
        err := tshelper.DefineAndCreateDeploymentOnCluster(3, randomNamespace)
        Expect(err).ToNot(HaveOccurred())

        // 2. Run certsuite
        err = globalhelper.LaunchTests(
            tsparams.CertsuiteTestCaseName,
            globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
            randomReportDir, randomCertsuiteConfigDir)
        Expect(err).ToNot(HaveOccurred())

        // 3. Verify result
        err = globalhelper.ValidateIfReportsAreValid(
            tsparams.CertsuiteTestCaseName,
            globalparameters.TestCasePassed, randomReportDir)
        Expect(err).ToNot(HaveOccurred())
    })
})
```

### Unit Tests vs Integration Tests
- **Unit tests**: Use build tag `//go:build utest`, mock Kubernetes clients, run without cluster
- **Integration tests**: Use build tag `//go:build !utest`, require live OCP cluster

### Configuration

**config/config.yaml** - Default configuration values
**Environment Variables**:
- `KUBECONFIG` - Path to cluster kubeconfig (required)
- `CERTSUITE_REPO_PATH` - Path to certsuite repo (required for binary mode)
- `FEATURES` - Comma-separated test features to run
- `CERTSUITE_IMAGE` / `CERTSUITE_IMAGE_TAG` - Custom certsuite image
- `DEBUG_CERTSUITE` - Enable debug logging
- `USE_BINARY` - Use local binary instead of container
- `DISABLE_INTRUSIVE_TESTS` - Skip intrusive tests
- `NON_LINUX_ENV` - Set to any value (including empty string) for macOS development
- `CERTSUITE_CONTAINER_CLIENT` - Container runtime (`docker` or `podman`)
- `DOCKER_CONFIG_DIR` - Docker config directory (macOS: `$HOME/.docker`)

### Key Dependencies
- `github.com/onsi/ginkgo/v2` / `github.com/onsi/gomega` - Testing framework
- `github.com/redhat-best-practices-for-k8s/certsuite-claim` - Claim file parsing
- `github.com/openshift-kni/eco-goinfra` - OpenShift resource builders
- `k8s.io/client-go` - Kubernetes client

## Test Features

Available test features (set via `FEATURES` env var):
- `accesscontrol` - Security context and privilege tests
- `affiliatedcertification` - Container/operator certification
- `lifecycle` - Pod lifecycle management
- `manageability` - Configuration management
- `networking` - Network policies and connectivity
- `observability` - Logging and monitoring
- `operator` - Operator installation and status
- `performance` - CPU pinning and scheduling
- `platformalteration` - Node and kernel configuration
- `preflight` - Red Hat preflight certification checks

### Kind Cluster Limitations
The following tests are known to fail on local Kind clusters:
- `access-control-security-context`
- `affiliated-certification-container-is-certified-digest`
- `affiliated-certification-operator-is-certified`
- `platform-alteration-tainted-node-kernel`

<!-- markdownlint-disable line-length no-bare-urls -->
# cnfcert-tests-verification

[![makefile ci](https://github.com/test-network-function/cnfcert-tests-verification/actions/workflows/makefile.yml/badge.svg)](https://github.com/test-network-function/cnfcert-tests-verification/actions/workflows/makefile.yml)
[![red hat](https://img.shields.io/badge/red%20hat---?color=gray&logo=redhat&logoColor=red&style=flat)](https://www.redhat.com) [![openshift](https://img.shields.io/badge/openshift---?color=gray&logo=redhatopenshift&logoColor=red&style=flat)](https://www.redhat.com/en/technologies/cloud-computing/openshift)
[![license](https://img.shields.io/github/license/test-network-function/cnfcert-tests-verification?color=blue&labelColor=gray&logo=apache&logoColor=lightgray&style=flat)](https://github.com/test-network-function/cnf-certification-test-partner/blob/master/LICENSE)

## Objective

> The repository contains a set of test cases that run different test scenarios from [cnf-certification-test](https://github.com/test-network-function/cnf-certification-test) project and verifies if these scenarios behave correctly under different environment conditions.

The cnfcert-tests-verification project is based on golang+[ginkgo](https://onsi.github.io/ginkgo) framework.

`cnfcert-tests-verification` project triggers the same test scenario from
[cnf-certification-test](https://github.com/test-network-function/cnf-certification-test)
several times using different pre-configured OCP environment.

Once the triggered scenario is completed, the test case processes the report and verifies that the scenario is completed with the excepted result: skip/fail/pass.

## Requirements

The tests are run on the OCP cluster with certain requirements that are listed below.

|  | Conditions | Mandatory |
| ------ | ------ | ------ |
| OCP Cluster | Version: >= 4.7, Node Count >= 3 with 2 cnf-worker nodes | Yes
| Installed Operators | Performance Addon, Machine-config-operator | Yes
|  | Machine config pool, PTP operator, SR-IOV operator| No

> Bare-minimum requirements consists of a OCP cluster with 3 nodes where 2 are cnf-worker nodes and 1 worker node.

## Overview

The following test features are can run selectively or altogether.

* *accesscontrol*
* *affiliatedcertification*
* *lifecycle*
* *manageability*
* *networking*
* *observability*
* *platformalteration*
* *performance*
* *operator*
* *preflight*

Choose the variant that suits you best:

> **`make test-features`** - will only run tests for the features that were defined in the `FEATURES` environment variable
> **`make test-all`** - will run the test suite for all features

### Environment variables

The following environment variables are used to configure the test setup.
| Env Variable Name | Purpose |
| ------ | ------ |
| FEATURES | To select the test scenarios that you are going to test, comma separated
| TNF_REPO_PATH | Points to the absolute path to  [cnf-certification-test](https://github.com/test-network-function/cnf-certification-test) on your machine
| TNF_IMAGE | Links to the TNF image. Default is quay.io/testnetworkfunction/cnf-certification-test
| TNF_IMAGE_TAG | image tag that is going to be tested. Default is latest
| TEST_IMAGE | Test image that is going to be used for all test resources such as deployments, daemonsets and so on. Default is quay.io/testnetworkfunction/cnf-test-partner
| DEBUG_TNF | Generate `Debug` folder that will contain TNF suites folders with TNF logs for each test.
| TNF_LOG_LEVEL | Log level. Default is 4
| DISABLE_INTRUSIVE_TESTS | Turns off the intrusive tests for faster execution. Default is `false`.
| ENABLE_PARALLEL | Enable ginkgo -p parallel flags (experimental). Default is `false`.

## Steps to run the tests

### Pre-requisites

Make sure docker or podman is running on the local machine. You could consider using [Colima - container runtime on macOS (and Linux) with minimal setup](https://github.com/abiosoft/colima).

#### Clone the repo and change directory to the cloned repo

```sh
git clone https://github.com/test-network-function/cnfcert-tests-verification.git
cd cnfcert-tests-verification
```

#### Download and install needed dependencies

```sh
make install
```

#### Set environment variables

* `testconfig.yaml` inside the `config` directory stores the local environment related information.
Update `tnf_config_dir` and `tnf_report_dir`, and `docker_config_dir` as specific to your local workspace.

Optionally, update `tnf_image`, `test_image`, and `tnf_image_tag` as per needs.

```yaml
# Sample configurations snippet
general:
  tnf_config_dir: "/Users/bmandal/rhdev/github.com/cnfcert-tests-verification/tnf_config"
  tnf_report_dir: "/Users/bmandal/rhdev/github.com/cnfcert-tests-verification/tnf_report"
  tnf_image: "quay.io/testnetworkfunction/cnf-certification-test"
  tnf_image_tag: "unstable"
  docker_config_dir: "/tmp"
```

* To use this test config file, you need to set `LOCAL_TESTING` environment variable while running the test.
* If you need to force the download of the `unstable` image, set the `FORCE_DOWNLOAD_UNSTABLE=true` environment variable.

>**Mac Users**:
Set `NON_LINUX_ENV=` to signal the repo code that the suite is run against the non Linux local env.

#### Execute tests

* To run all tests

```sh
# Mac user
 \
  export TNF_CONTAINER_CLIENT=docker &&
  DOCKER_CONFIG_DIR=$HOME/.docker \
  KUBECONFIG=$HOME/.kube/config \
  NON_LINUX_ENV= \
  TNF_REPO_PATH=$HOME/path/to/cnf-certification-test \
  make test-all
```

```sh
# Linux user
 \
  KUBECONFIG=$HOME/.kube/config \
  LOCAL_TESTING= \
  TNF_REPO_PATH=$HOME/path/to/cnf-certification-test \
  make test-all
```

```sh
# Linux user with force download unstable image
 \
  FORCE_DOWNLOAD_UNSTABLE=true \
  KUBECONFIG=$HOME/.kube/config \
  LOCAL_TESTING= \
  TNF_REPO_PATH=$HOME/path/to/cnf-certification-test \
  make test-all
```

* To run a specific feature

```sh
# Mac user
 \
  export TNF_CONTAINER_CLIENT=docker &&
  DOCKER_CONFIG_DIR=$HOME/.docker \
  FEATURES=platformalteration \
  KUBECONFIG=$HOME/.kube/config \
  NON_LINUX_ENV= \
  TNF_REPO_PATH=$HOME/path/to/cnf-certification-test \
  make test-features
```

```sh
# Linux user
 \
  FEATURES=platformalteration \
  KUBECONFIG=$HOME/.kube/config \
  LOCAL_TESTING= \
  TNF_REPO_PATH=$HOME/path/to/cnf-certification-test \
  make test-features
```

* To debug

Use `DEBUG_TNF=true` and `TNF_LOG_LEVEL=trace` while running the above commands.
This would create a `Debug` folder containing suites folders with TNF logs for each of the tests.

```sh
# Mac user
 \
  export TNF_CONTAINER_CLIENT=docker &&
  DEBUG_TNF=true \
  DOCKER_CONFIG_DIR=$HOME/.docker \
  FEATURES=platformalteration \
  KUBECONFIG=$HOME/.kube/config \
  NON_LINUX_ENV= \
  TNF_LOG_LEVEL=trace \
  TNF_REPO_PATH=$HOME/path/to/cnf-certification-test \
  make test-features
```

```sh
# Linux user
 \
  DEBUG_TNF=true \
  FEATURES=platformalteration \
  KUBECONFIG=$HOME/.kube/config \
  LOCAL_TESTING= \
  TNF_LOG_LEVEL=trace \
  TNF_REPO_PATH=$HOME/path/to/cnf-certification-test \
  make test-features
```

## Running the unit tests

To execute the unit tests in the repository, run the following:

```sh
make test
```

## Test exceptions on local kind cluster

* access-control-security-context
* affiliated-certification-container-is-certified-digest
* affiliated-certification-operator-is-certified

## Contribution Guidelines

Fork the repo, create a new branch and create a PR with your changes.

## License

CNF Certification Test Partner is copyright [Red Hat, Inc.](https://www.redhat.com) and available
under an
[Apache 2 license](https://github.com/test-network-function/cnfcert-tests-verification/blob/main/LICENSE).

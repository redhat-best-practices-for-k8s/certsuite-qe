<!-- markdownlint-disable line-length no-bare-urls -->
# certsuite-qe

[![Test Incoming Changes](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/actions/workflows/pre-main.yml/badge.svg)](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/actions/workflows/pre-main.yml)
[![red hat](https://img.shields.io/badge/red%20hat---?color=gray&logo=redhat&logoColor=red&style=flat)](https://www.redhat.com) [![openshift](https://img.shields.io/badge/openshift---?color=gray&logo=redhatopenshift&logoColor=red&style=flat)](https://www.redhat.com/en/technologies/cloud-computing/openshift)
[![license](https://img.shields.io/github/license/redhat-best-practices-for-k8s/certsuite-qe?color=blue&labelColor=gray&logo=apache&logoColor=lightgray&style=flat)](https://github.com/redhat-best-practices-for-k8s/certsuite-partner/blob/master/LICENSE)

## Objective

The repository contains a set of test cases that run different test scenarios from [certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite) project and verifies if these scenarios behave correctly under different environment conditions.

The certsuite-qe project is based on golang+[ginkgo](https://onsi.github.io/ginkgo) framework.

`certsuite-qe` project triggers the same test scenario from
[certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite)
several times using different pre-configured OCP environment.

Once the triggered scenario is completed, the test case processes the report and verifies that the scenario is completed with the excepted result: skip/fail/pass.

## Requirements

The tests are run on the OCP cluster with certain requirements that are listed below.

|  | Conditions | Mandatory |
| ------ | ------ | ------ |
| OCP Cluster | Version: >= 4.12, Node Count >= 3 with 2 cnf-worker nodes | Yes
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
| CERTSUITE_REPO_PATH | Points to the absolute path to  [certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite) on your machine
| CERTSUITE_IMAGE | Links to the Certsuite image. Default is quay.io/redhat-best-practices-for-k8s/certsuite
| CERTSUITE_IMAGE_TAG | image tag that is going to be tested. Default is latest
| TEST_IMAGE | Test image that is going to be used for all test resources such as deployments, daemonsets and so on. Default is quay.io/testnetworkfunction/k8s-best-practices-debug
| DEBUG_CERTSUITE | Generate `Debug` folder that will contain Certsuite suites folders with Certsuite logs for each test.
| CERTSUITE_LOG_LEVEL | Log level. Default is 4
| DISABLE_INTRUSIVE_TESTS | Turns off the intrusive tests for faster execution. Default is `false`.
| ENABLE_PARALLEL | Enable ginkgo -p parallel flags (experimental). Default is `false`.
| FORCE_DOWNLOAD_UNSTABLE | Force download the unstable image. Default is `false`.
| NON_LINUX_ENV | Allow the test suites to run in a non Linux environment. Default is `false`.

## Steps to run the tests

### Pre-requisites

Make sure [docker](https://www.docker.com/) or [podman](https://podman.io/) is running on the local machine.

Set your local container runtime to your environment with:

```sh
export CERTSUITE_CONTAINER_CLIENT=docker
```

#### Clone the repo and change directory to the cloned repo

```sh
git clone https://github.com/redhat-best-practices-for-k8s/certsuite-qe.git
cd certsuite-qe
```

#### Download and install needed dependencies

```sh
make install
```

#### Execute tests

* To run all tests

```sh
# Mac user
  DOCKER_CONFIG_DIR=$HOME/.docker \
  KUBECONFIG=$HOME/.kube/config \
  NON_LINUX_ENV= \
  CERTSUITE_REPO_PATH=$HOME/path/to/certsuite \
  make test-all
```

```sh
# Linux user
  KUBECONFIG=$HOME/.kube/config \
  CERTSUITE_REPO_PATH=$HOME/path/to/certsuite \
  make test-all
```

```sh
# Linux user with force download unstable image
 \
  FORCE_DOWNLOAD_UNSTABLE=true \
  KUBECONFIG=$HOME/.kube/config \
  CERTSUITE_REPO_PATH=$HOME/path/to/certsuite \
  make test-all
```

* To run a specific test-suite:

```sh
# Mac user
  DOCKER_CONFIG_DIR=$HOME/.docker \
  FEATURES=platformalteration \
  KUBECONFIG=$HOME/.kube/config \
  NON_LINUX_ENV= \
  CERTSUITE_REPO_PATH=$HOME/path/to/certsuite \
  make test-features
```

```sh
# Linux user
  FEATURES=platformalteration \
  KUBECONFIG=$HOME/.kube/config \
  DOCKER_CONFIG_DIR=$HOME/.docker \
  CERTSUITE_REPO_PATH=$HOME/path/to/certsuite \
  make test-features
```

* To debug

Use `DEBUG_CERTSUITE=true` and `CERTSUITE_LOG_LEVEL=debug` while running the above commands.
This would create a `Debug` folder containing suites folders with Certsuite logs for each of the tests.

```sh
# Mac user
  DEBUG_CERTSUITE=true \
  DOCKER_CONFIG_DIR=$HOME/.docker \
  FEATURES=platformalteration \
  KUBECONFIG=$HOME/.kube/config \
  NON_LINUX_ENV= \
  CERTSUITE_LOG_LEVEL=debug \
  CERTSUITE_REPO_PATH=$HOME/path/to/certsuite \
  make test-features
```

```sh
# Linux user
  DEBUG_CERTSUITE=true \
  FEATURES=platformalteration \
  KUBECONFIG=$HOME/.kube/config \
  CERTSUITE_LOG_LEVEL=debug \
  CERTSUITE_REPO_PATH=$HOME/path/to/certsuite \
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
* platform-alteration-tainted-node-kernel

## Nightly Runs Against Various Environments

The QE repo is being used in nightly automated runs in the following files:

* [QE via Kind (Github Hosted)](https://github.com/redhat-best-practices-for-k8s/certsuite/blob/main/.github/workflows/qe-hosted.yml)
* [QE via OCP (Intrusive)](https://github.com/redhat-best-practices-for-k8s/certsuite/blob/main/.github/workflows/qe-ocp-intrusive.yaml)
* [QE via OCP (Non-Intrusive)](https://github.com/redhat-best-practices-for-k8s/certsuite/blob/main/.github/workflows/qe-ocp.yaml)

## Contribution Guidelines

Fork the repo, create a new branch and create a PR with your changes.

## License

CNFCert Tests Verification is copyright [Red Hat, Inc.](https://www.redhat.com) and available
under an
[Apache 2 license](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/blob/main/LICENSE).

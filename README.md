<!-- markdownlint-disable line-length no-bare-urls -->
# CertSuite QE

[![Test Incoming Changes](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/actions/workflows/pre-main.yml/badge.svg)](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/actions/workflows/pre-main.yml)
[![red hat](https://img.shields.io/badge/red%20hat---?color=gray&logo=redhat&logoColor=red&style=flat)](https://www.redhat.com) [![openshift](https://img.shields.io/badge/openshift---?color=gray&logo=redhatopenshift&logoColor=red&style=flat)](https://www.redhat.com/en/technologies/cloud-computing/openshift)
[![license](https://img.shields.io/github/license/redhat-best-practices-for-k8s/certsuite-qe?color=blue&labelColor=gray&logo=apache&logoColor=lightgray&style=flat)](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/blob/main/LICENSE)

## Objective

The repository contains a set of test cases that run different test scenarios from [certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite) project and verifies if these scenarios behave correctly under different environment conditions.

The certsuite-qe project is based on golang+[ginkgo](https://onsi.github.io/ginkgo) framework.

`certsuite-qe` project triggers the same test scenario from
[certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite)
several times using different pre-configured OCP environment.

Once the triggered scenario is completed, the test case processes the report and verifies that the scenario is completed with the expected result: skip/fail/pass.

## Requirements

The tests are run on the OCP cluster with certain requirements that are listed below.

|  | Conditions | Mandatory |
| ------ | ------ | ------ |
| OCP Cluster | Version: >= 4.12, Node Count >= 3 with 2 cnf-worker nodes | Yes |
| Installed Operators | Performance Addon, Machine-config-operator | Yes |
|  | Machine config pool, PTP operator, SR-IOV operator| No |

> Bare-minimum requirements consist of an OCP cluster with 3 nodes where 2 are cnf-worker nodes and 1 worker node.

## Overview

The following test features can run selectively or altogether.

* *accesscontrol*
* *affiliatedcertification*
* *lifecycle*
* *manageability*
* *networking*
* *observability*
* *platformalteration*
* *performance*
* *operator*

Choose the variant that suits you best:

> **`make test-features`** - will only run tests for the features that were defined in the `FEATURES` environment variable
> **`make test-all`** - will run the test suite for all features

### Environment variables

The following environment variables are used to configure the test setup.

| Env Variable Name | Purpose |
| ------ | ------ |
| KUBECONFIG | Path to cluster kubeconfig (required) |
| FEATURES | Select test scenarios to run, comma separated |
| CERTSUITE_REPO_PATH | Absolute path to [certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite) on your machine |
| CERTSUITE_IMAGE | Certsuite image. Default is `quay.io/redhat-best-practices-for-k8s/certsuite` |
| CERTSUITE_IMAGE_TAG | Image tag to test. Default is `latest` |
| USE_BINARY | Use local certsuite binary instead of container image. Default is `false` |
| DEBUG_CERTSUITE | Generate a `Debug` folder with Certsuite logs for each test |
| CERTSUITE_LOG_LEVEL | Log level when debugging. Set to `debug` with `DEBUG_CERTSUITE=true` |
| DISABLE_INTRUSIVE_TESTS | Skip intrusive tests for faster execution. Default is `false` |
| ENABLE_PARALLEL | Enable ginkgo parallel execution via `--procs=16` (experimental). Default is `false` |
| FORCE_DOWNLOAD_UNSTABLE | Force download the unstable image. Default is `false` |
| NON_LINUX_ENV | Set to any value (including empty string) to run on macOS. Unset on Linux |
| DOCKER_CONFIG_DIR | Docker config directory (required on macOS; example: `$HOME/.docker`) |
| CONTAINER_ENGINE | Container runtime to use (`docker` or `podman`). Default is `docker` |

## Steps to run the tests

### Pre-requisites

Make sure [docker](https://www.docker.com/) or [podman](https://podman.io/) is running on the local machine.

Set your local container runtime to your environment with:

```sh
export CONTAINER_ENGINE=docker
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

## CI Workflows

Nightly and on-demand QE runs use the following workflows:

* [QE via Kind](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/blob/main/.github/workflows/qe.yml) (this repo)
* [QE via OCP](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/blob/main/.github/workflows/qe-ocp.yml) (this repo)
* [QE via Kind (certsuite-hosted)](https://github.com/redhat-best-practices-for-k8s/certsuite/blob/main/.github/workflows/qe-hosted.yml)
* [QE via OCP 4.22 (certsuite-hosted)](https://github.com/redhat-best-practices-for-k8s/certsuite/blob/main/.github/workflows/qe-ocp-422.yaml)

For agent-oriented architecture and command reference, see [AGENTS.md](AGENTS.md).

## Contribution Guidelines

Fork the repo, create a new branch and create a PR with your changes.

## License

CertSuite QE is copyright [Red Hat, Inc.](https://www.redhat.com) and available
under an
[Apache 2 license](https://github.com/redhat-best-practices-for-k8s/certsuite-qe/blob/main/LICENSE).

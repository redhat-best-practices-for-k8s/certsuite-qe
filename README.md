# cnfcert-tests-verification

## Overview
The repository contains a set of test cases that run different test scenarios from [test-network-function](https://github.com/test-network-function/cnf-certification-test) project and verifies if these scenarios behave correctly under different environment conditions.
The cnfcert-tests-verification project based on golang+[ginkgo](https://onsi.github.io/ginkgo) framework.

cnfcert-tests-verification project triggers the same test scenario from
[test-network-function](https://github.com/test-network-function/cnf-certification-test)
several times using different preconfigured OCP environment. Once the triggered scenario is completed, the test case processes the report and verifies that the scenario is completed with the excepted result: skip/fail/pass.

## cnfcert-tests-verification
The cnfcert-tests-verification is designed to test [test-network-function](https://github.com/test-network-function/cnf-certification-test) project using pre-installed OCP cluster with version 4.7 and above. In order to be able to test all test scenarios, the following requirements should be met:

Mandatory requirements:
* OCP cluster installed with version >=4.7
* Minimum 3 worker nodes where 2 of them are cnf-workers nodes
* Performance Addon Operator
* Machine-config-operator

Optional:
* Machine config pool
* PTP operator
* SR-IOV operator
* Performance Addon Operator

### Recommended environment:
* 3 master nodes
* 2 cnf-worker nodes
* 1 worker node

#### Environment variables
* `FEATURES` - select the test scenarios that you are going to test, comma separated
* `TNF_REPO_PATH` - absolute path to  [test-network-function](https://github.com/test-network-function/cnf-certification-test) on your machine
* `TNF_IMAGE` - link to tnf image. Default is quay.io/testnetworkfunction/cnf-certification-test
* `TNF_IMAGE_TAG` - image tag that is going to be tested. Default is latest
* `TNF_LOG_LEVEL` - Log level. Default is 4
* `TEST_IMAGE` - Test image that is going to be used for all test resources such as deployments, daemonsets and so on. Default is quay.io/testnetworkfunction/cnf-test-partner
* `DEBUG_TNF` - Generate `Debug` folder that will contain TNF suites folders with TNF logs for each test.

#### Available features
The list of available features:
* *networking*
* *affiliated-certification*
* *lifecycle*


#### Running the tests

Choose the variant that suits you best:

* `make test-features` - will only run tests for the features that were defined in the `FEATURES` variable
* `make test-all` - will run the test suite for all features

#### Pre-configuration

`make install` - download and install all required dependencies for the cnfcert-tests-verification project

## How to run

Below is an e2e flow example:

1. Clone the project to your local computer - `git clone https://github.com/test-network-function/cnfcert-tests-verification.git`

2. Change the folder to the project folder - `cd cnfcert-tests-verification`

3. Download and install needed dependencies - `make install`

4. Define absolute path to [test-network-function](https://github.com/test-network-function/cnf-certification-test) project on you computer by setting environment variable: `export TNF_REPO_PATH=/path/to/project/cnf-certification-test`

5. OPTIONAL: Set cnf-certification-test container tag that you are going to test. Default is latest. `export TNF_IMAGE_TAG=unstable` 

6. Run all tests - `make test-all`

# How to debug

1. export `DEBUG_TNF`=true.
2. run `make test-all` / `make test-features`.
3. under reports folder, a `Debug` folder will be generated, containing suites folders with TNF logs for each of the tests. 


# cnfcert-tests-verification - How to contribute

The project uses a development method - forking workflow
### The following is a step-by-step example of forking workflow:
1) A developer [forks](https://docs.gitlab.com/ee/user/project/repository/forking_workflow.html#creating-a-fork)
   the [cnfcert-tests-verification](https://github.com/test-network-function/cnfcert-tests-verification) project
2) A new local feature branch is created
3) The developer makes changes on the new branch.
4) New commits are created for the changes.
5) The branch gets pushed to the developer's own server-side copy.
6) Changes are tested.
7) The developer opens a pull request(`PR`) from the new branch to
   the [cnfcert-tests-verification](https://github.com/test-network-function/cnfcert-tests-verification).
8) The pull request gets approved for merge and is merged into
   the [cnfcert-tests-verification](https://github.com/test-network-function/cnfcert-tests-verification).

# cnfcert-tests-verification - Project structure
    .
    ├── config                         # Config files
    ├── scripts                        # Makefile Scripts 
    ├── tests                          # Test cases directory
    │   ├── networking                 # Networking test cases directory
    │   │   ├── nethelper       # Networking common test function
    │   │   ├── netparameters   # Networking constants and parameters 
    │   │   └── tests                  # Networking test suite directory
    |   ├── affiliatedcertification    # Affiliated Certification test cases directory   
    |   |   ├── affiliatedcerthelper         # Affiliated Certification common test function
    |   |   ├── affiliatedcertparameters     # Affiliated Certification constants and parameters 
    |   |   └── tests                        # Affiliated Certification test suite directory
    │   ├── platform                   # Platform test cases directory
    │   │   ├── platformghelper        # Platform common test function
    │   │   ├── platformparameters     # Platform constans and parameters
    │   │   └── tests                  # Platform test suite directory
    │   ├── globalhelper                     # Common test test function
    │   ├── globalparameters                 # Common test function
    │   └── utils                      # Common utils functions. These utils are based on Kubernetes api calls
    │       ├── client
    │       ├── config
    │       ├── node
    │       ├── namespace
    │       └── pod
    └── vendors                        # Dependencies folder 

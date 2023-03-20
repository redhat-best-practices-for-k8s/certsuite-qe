package helper

import (
	"fmt"
	"strings"

	"github.com/onsi/ginkgo/v2"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
)

func SetUpAndRunContainerCertTest(tcName string, containersInfo []string, expectedResult string) error {
	var err error

	ginkgo.By("Add container information to " + globalparameters.DefaultTnfConfigFileName)

	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.TestCertificationNameSpace},
		[]string{tsparams.TestPodLabel},
		containersInfo,
		[]string{})

	if err != nil {
		return fmt.Errorf("error defining tnf config file: %w", err)
	}

	ginkgo.By("Start test")

	err = globalhelper.LaunchTests(
		tsparams.TestCaseContainerAffiliatedCertName,
		tcName)

	if strings.Contains(expectedResult, globalparameters.TestCaseFailed) && err == nil {
		return fmt.Errorf("error running %s test",
			tsparams.TestCaseContainerAffiliatedCertName)
	}

	if (strings.Contains(expectedResult, globalparameters.TestCasePassed) ||
		strings.Contains(expectedResult, globalparameters.TestCaseSkipped)) && err != nil {
		return fmt.Errorf("error running %s test: %w",
			tsparams.TestCaseContainerAffiliatedCertName, err)
	}

	ginkgo.By("Verify test case status in Junit and Claim reports")

	err = globalhelper.ValidateIfReportsAreValid(
		tsparams.TestCaseContainerAffiliatedCertName,
		expectedResult)

	if err != nil {
		return fmt.Errorf("error validating test reports: %w", err)
	}

	return nil
}

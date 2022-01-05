package affiliatedcerthelper

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcertparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/nethelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
)

func SetUpAndRunContainerCertTest(containersInfo []string, expectedResult string) error {
	var err error

	By("Add container information to " + globalparameters.DefaultTnfConfigFileName)

	err = globalhelper.DefineTnfConfig(
		[]string{netparameters.TestNetworkingNameSpace},
		[]string{netparameters.TestPodLabel},
		containersInfo)

	if err != nil {
		return fmt.Errorf("error defining tnf config file: %w", err)
	}

	By("Start test")

	err = globalhelper.LaunchTests(
		[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
		affiliatedcertparameters.TestCaseContainerSkipRegEx,
	)

	if strings.Compare(expectedResult, globalparameters.TestCaseFailed) == 0 {
		if err == nil {
			return fmt.Errorf("error running %s test",
				affiliatedcertparameters.AffiliatedCertificationTestSuiteName)
		}
	} else {

		if err != nil {
			return fmt.Errorf("error running %s test: %w",
				affiliatedcertparameters.AffiliatedCertificationTestSuiteName, err)
		}
	}

	By("Verify test case status in Junit and Claim reports")

	err = nethelper.ValidateIfReportsAreValid(
		affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
		expectedResult)

	if err != nil {
		return fmt.Errorf("error validating test reports: %w", err)
	}

	return nil
}

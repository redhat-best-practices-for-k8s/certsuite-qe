package affiliatedcerthelper

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

	By("Start test")

	err = globalhelper.LaunchTests(
		[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
		affiliatedcertparameters.TestCaseContainerSkipRegEx,
	)

	if strings.Compare(expectedResult, globalparameters.TestCaseFailed) == 0 {
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")
	} else {
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")
	}

	By("Verify test case status in Junit and Claim reports")

	err = nethelper.ValidateIfReportsAreValid(
		affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
		expectedResult)
	Expect(err).ToNot(HaveOccurred(), "Error validating test reports")

	if err != nil {
		return err
	}

	return nil
}

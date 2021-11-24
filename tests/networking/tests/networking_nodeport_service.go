package tests

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/nethelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Networking custom namespace, custom deployment,", func() {

	execute.BeforeAll(func() {

		By("Clean namespace before all tests")
		err := netparameters.TestNamespace.Clean(globalhelper.ApiClient)
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, netparameters.TestNamespace.Name)
		Expect(err).ToNot(HaveOccurred())

	})

	BeforeEach(func() {

		By("Clean namespace before each test")
		err := netparameters.TestNamespace.Clean(globalhelper.ApiClient)
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())
	})

	// 45447
	It("2 custom pods, no service installed, service Should not have type of nodePort", func() {

		By("Define deployment and create it on cluster")
		err := nethelper.DefineAndCreateDeploymentOnCluster(3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseNodePortSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			netparameters.TestCaseNodePortNetworkName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45481
	It("2 custom pods, service installed without NodePort, service Should not have type of nodePort", func() {

	})

	// 45482
	It("2 custom pods, multiple services installed without NodePort, service Should not have type of nodePort",
		func() {

		})

	// 45483
	It("2 custom pods, service installed with NodePort, service Should not have type of nodePort [negative]",
		func() {

		})

	// 45484
	It("2 custom pods, multiple services installed and one has NodePort, service Should not have type of "+
		"nodePort [negative]", func() {

	})

})

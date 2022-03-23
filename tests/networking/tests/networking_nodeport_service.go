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
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Networking custom namespace, custom deployment,", func() {

	execute.BeforeAll(func() {

		By("Clean namespace before all tests")
		err := namespaces.Clean(netparameters.TestNetworkingNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, netparameters.TestNetworkingNameSpace)
		Expect(err).ToNot(HaveOccurred())

	})

	BeforeEach(func() {

		By("Clean namespace before each test")
		err := namespaces.Clean(netparameters.TestNetworkingNameSpace, globalhelper.APIClient)
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
			netparameters.NetworkingTestSuiteName,
			netparameters.TestCaseNodePortNetworkName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
			netparameters.TestCaseNodePortSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			netparameters.TestCaseNodePortNetworkName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45481
	It("2 custom pods, service installed without NodePort, service Should not have type of nodePort", func() {

		By("Define Service")
		err := nethelper.DefineAndCreateServiceOnCluster("testservice", 3022, 3022, false)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentOnCluster(3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			netparameters.NetworkingTestSuiteName,
			netparameters.TestCaseNodePortNetworkName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
			netparameters.TestCaseNodePortSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			netparameters.TestCaseNodePortNetworkName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45482
	It("2 custom pods, multiple services installed without NodePort, service Should not have type of nodePort", func() {

		By("Define multiple Services")
		err := nethelper.DefineAndCreateServiceOnCluster("testservicefirst", 3022, 3022, false)
		Expect(err).ToNot(HaveOccurred())

		err = nethelper.DefineAndCreateServiceOnCluster("testservicesecond", 3023, 3023, false)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentOnCluster(3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			netparameters.NetworkingTestSuiteName,
			netparameters.TestCaseNodePortNetworkName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
			netparameters.TestCaseNodePortSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			netparameters.TestCaseNodePortNetworkName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45483
	It("2 custom pods, service installed with NodePort, service Should not have type of nodePort [negative]", func() {

		By("Define Services with NodePort")
		err := nethelper.DefineAndCreateServiceOnCluster("testservice", 30022, 3022, true)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentOnCluster(3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			netparameters.NetworkingTestSuiteName,
			netparameters.TestCaseNodePortNetworkName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
			netparameters.TestCaseNodePortSkipRegEx,
		)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			netparameters.TestCaseNodePortNetworkName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45484
	It("2 custom pods, multiple services installed and one has NodePort, service Should not have type of "+
		"nodePort [negative]", func() {

		By("Define Services")
		err := nethelper.DefineAndCreateServiceOnCluster("testservicefirst", 30022, 3022, true)
		Expect(err).ToNot(HaveOccurred())
		err = nethelper.DefineAndCreateServiceOnCluster("testservicesecond", 3022, 3022, false)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentOnCluster(3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			netparameters.NetworkingTestSuiteName,
			netparameters.TestCaseNodePortNetworkName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
			netparameters.TestCaseNodePortSkipRegEx,
		)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			netparameters.TestCaseNodePortNetworkName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

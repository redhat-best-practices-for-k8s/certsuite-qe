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

var _ = Describe("Networking custom namespace, ", func() {

	var multusInterfaces []string

	execute.BeforeAll(func() {

		By("Clean namespace before all tests")
		err := namespaces.Clean(netparameters.TestNetworkingNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, netparameters.TestNetworkingNameSpace)
		Expect(err).ToNot(HaveOccurred())

		By("Collect list of available list of interface from cluster")
		multusInterfaces, err = nethelper.GetClusterMultusInterfaces()
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

		if len(multusInterfaces) < 1 {
			Skip("There is no enough Multus interfaces available")
		}

		By("Put interface down")
		err = nethelper.PutDownInterfaceOnNode(multusInterfaces[0], "down")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := nethelper.PutAllInterfaceUPonAllNodes(multusInterfaces[0])
		Expect(err).ToNot(HaveOccurred())
	})

	// 48346
	It("custom deployment 3 pods,1 NAD,no connectivity via Multus secondary interface[negative]",
		func() {

			By("Define and create Network-attachment-definition")
			err := nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameA, multusInterfaces[0], netparameters.TestIPamIPNetworkA)
			Expect(err).ToNot(HaveOccurred())

			By("Define deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentAName, []string{netparameters.TestNadNameA}, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseMultusSkipRegEx,
			)
			Expect(err).To(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				netparameters.TestCaseMultusConnectivityName,
				globalparameters.TestCaseFailed)
			Expect(err).ToNot(HaveOccurred())
		})

	// 48347
	It("custom deployment and daemonset 3 pods, 2 NADs, No connectivity on daemonset via Multus secondary "+
		"interface[negative]",
		func() {

			if len(multusInterfaces) < 2 {
				Skip("There is not enough secondary network interfaces to run the test case")
			}

			err := nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameA, multusInterfaces[0], netparameters.TestIPamIPNetworkA)
			Expect(err).ToNot(HaveOccurred())

			err = nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameB, multusInterfaces[1], netparameters.TestIPamIPNetworkB)
			Expect(err).ToNot(HaveOccurred())

			By("Define deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentAName, []string{netparameters.TestNadNameA}, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Define daemonset and create it on cluster")
			err = nethelper.DefineAndCreateDeamonsetWithMultusOnCluster(netparameters.TestNadNameA)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseMultusSkipRegEx,
			)
			Expect(err).To(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				netparameters.TestCaseMultusConnectivityName,
				globalparameters.TestCaseFailed)
			Expect(err).ToNot(HaveOccurred())
		})

	It("custom deployment and daemonset 3 pods, 2 NADs, multiple Multus interfaces on deployment no "+
		"connectivity via secondary interface[negative]", func() {

		if len(multusInterfaces) < 2 {
			Skip("There is no enough Multus interfaces available")
		}

		By("Define and create network-attachment-definitions")
		err := nethelper.DefineAndCreateNadOnCluster(
			netparameters.TestNadNameA, multusInterfaces[0], netparameters.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = nethelper.DefineAndCreateNadOnCluster(
			netparameters.TestNadNameB, multusInterfaces[1], netparameters.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
			netparameters.TestDeploymentAName, []string{netparameters.TestNadNameB, netparameters.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = nethelper.DefineAndCreateDeamonsetWithMultusOnCluster(netparameters.TestNadNameB)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseMultusSkipRegEx,
		)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			netparameters.TestCaseMultusConnectivityName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

package tests

import (
	"os"
	"time"

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

		By("Collect list of available interfaces from the cluster")
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
	})

	// 48328
	It("custom deployment 3 pods, 1 NAD, connectivity via Multus secondary interface",
		func() {

			By("Define and create Network-attachment-definition")
			err := nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameA, multusInterfaces[0], netparameters.TestIPamIPNetworkA)
			Expect(err).ToNot(HaveOccurred())

			By("Define deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentAName, netparameters.TestNadNameA, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseMultusSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				netparameters.TestCaseMultusConnectivityName,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())
		})

	// 48330
	It("2 custom deployments 3 pods, 1 NAD, connectivity via Multus secondary interface",
		func() {

			By("Define and create Network-attachment-definition")
			err := nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameA, multusInterfaces[0], netparameters.TestIPamIPNetworkA)
			Expect(err).ToNot(HaveOccurred())

			By("Define first deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentAName, netparameters.TestNadNameA, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Define second deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentBName, netparameters.TestNadNameA, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseMultusSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				netparameters.TestCaseMultusConnectivityName,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(30 * time.Second)
		})

	// 48331
	It("custom deployment and daemonset 3 pods, 2 NADs, connectivity via Multus secondary interfaces",
		func() {

			By("Define and create Network-attachment-definition")
			err := nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameA, multusInterfaces[0], netparameters.TestIPamIPNetworkA)
			Expect(err).ToNot(HaveOccurred())

			err = nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameB, multusInterfaces[0], netparameters.TestIPamIPNetworkB)
			Expect(err).ToNot(HaveOccurred())

			By("Define first deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentAName, netparameters.TestNadNameA, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Define second deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentBName, netparameters.TestNadNameB, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseMultusSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				netparameters.TestCaseMultusConnectivityName,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(30 * time.Second)
		})

	// 48334
	It("custom deployment 3 pods, 1 NAD missing IP, connectivity via Multus secondary interface[skip]",
		func() {
			By("Define and create Network-attachment-definition")

			err := nethelper.DefineAndCreateNadOnCluster(netparameters.TestNadNameA, multusInterfaces[0], "")
			Expect(err).ToNot(HaveOccurred())

			By("Define deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentAName, netparameters.TestNadNameA, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseMultusSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				netparameters.TestCaseMultusConnectivityName,
				globalparameters.TestCaseSkipped)
			Expect(err).ToNot(HaveOccurred())
		})

	// 48338
	It("custom deployments 3 pods and 1 pod, standalone IP, connectivity via Multus secondary interface[skip]",
		func() {
			if len(multusInterfaces) < 2 {
				Skip("There is not enough secondary network interfaces to run the test case")
			}

			By("Define and create Network-attachment-definitions")
			err := nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameA, multusInterfaces[0], netparameters.TestIPamIPNetworkA)
			Expect(err).ToNot(HaveOccurred())

			err = nethelper.DefineAndCreateNadOnCluster(netparameters.TestNadNameB, multusInterfaces[1], "")
			Expect(err).ToNot(HaveOccurred())

			By("Define deployment-a and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentAName, netparameters.TestNadNameA, 1)
			Expect(err).ToNot(HaveOccurred())

			By("Define deployment-b and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentBName, netparameters.TestNadNameB, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseMultusSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				netparameters.TestCaseMultusConnectivityName,
				globalparameters.TestCaseSkipped)
			Expect(err).ToNot(HaveOccurred())
		})

	// 48343
	It("custom deployment and daemonset 3 pods, daemonset missing ip, 2 NADs, connectivity via Multus "+
		"secondary interface",
		func() {
			if len(multusInterfaces) < 2 {
				Skip("There is not enough secondary network interfaces to run the test case")
			}

			By("Define and create network-attachment-definitions")
			err := nethelper.DefineAndCreateNadOnCluster(
				netparameters.TestNadNameA, multusInterfaces[0], netparameters.TestIPamIPNetworkA)
			Expect(err).ToNot(HaveOccurred())

			err = nethelper.DefineAndCreateNadOnCluster(netparameters.TestNadNameB, multusInterfaces[1], "")
			Expect(err).ToNot(HaveOccurred())

			By("Define deployment and create it on cluster")
			err = nethelper.DefineAndCreateDeploymentWithMultusOnCluster(
				netparameters.TestDeploymentAName, netparameters.TestNadNameA, 3)
			Expect(err).ToNot(HaveOccurred())

			By("Define daemonset and create it on cluster")

			err = nethelper.DefineAndCreateDeamonsetWithMultusOnCluster(netparameters.TestNadNameB)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseMultusSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				netparameters.TestCaseMultusConnectivityName,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())
		})
})

package tests

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/nethelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Networking custom namespace, custom deployment,", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file"))
	}
	execute.BeforeAll(func() {

		By("Clean namespace before all tests")
		err = netparameters.TestNamespace.Clean(globalhelper.ApiClient)
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

	// 45440
	It("3 custom pods on Default network networking-icmpv4-connectivity", func() {

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentOnCluster(3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseDefaultSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			netparameters.TestCaseDefaultNetworkName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45441
	It("custom daemonset, 4 custom pods on Default network", func() {

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(netparameters.TestNamespace.Name,
			configSuite.General.TestImage,
			netparameters.TestDeploymentLabels)
		daemonSet.RedefineDaemonSetWithNodeSelector(map[string]string{configSuite.General.CnfNodeLabel: ""})
		By("Create DaemonSet on cluster")
		err = nethelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet.DaemonSet, netparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseDefaultSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			netparameters.TestCaseDefaultNetworkName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45442
	It("3 custom pods on Default network networking-icmpv4-connectivity fail when one pod is "+
		"disconnected [negative]", func() {

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreatePrivilegedDeploymentOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Collect partner's pod info")
		partnerPod, err := nethelper.GetPartnerPodDefinition()
		Expect(err).ToNot(HaveOccurred())

		By("Close communication between partner and test pods")
		err = nethelper.ExecCmdOnOnePodInNamespace(
			[]string{"ip", "route", "add", partnerPod.Status.PodIP, "via", "127.0.0.1"})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseDefaultSkipRegEx,
		)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			netparameters.TestCaseDefaultNetworkName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45443
	It("2 custom pods on Default network networking-icmpv4-connectivity fail when there is no ping binary "+
		"[negative]", func() {

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreatePrivilegedDeploymentOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Remove ping binary from test pod")
		err = nethelper.ExecCmdOnOnePodInNamespace(
			[]string{"rm", "-rf", "/usr/bin/ping", "/usr/sbin/ping"})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseDefaultSkipRegEx,
		)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			netparameters.TestCaseDefaultNetworkName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45444
	It("2 custom pods on Default network networking-icmpv4-connectivity skip when label "+
		"test-network-function.com/skip_connectivity_tests is set in deployment [skip]", func() {

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentWithSkippedLabelOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Remove ping binary from test pod")
		err = nethelper.ExecCmdOnOnePodInNamespace(
			[]string{"rm", "-rf", "/usr/bin/ping", "/usr/sbin/ping"})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseDefaultSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			netparameters.TestCaseDefaultNetworkName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45445
	It("custom daemonset, 4 custom pods on Default network networking-icmpv4-connectivity pass when label "+
		"test-network-function.com/skip_connectivity_tests is set in deployment only", func() {

		By("Define deployment and create it on cluster")
		err = nethelper.DefineAndCreateDeploymentWithSkippedLabelOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Remove ping binary from test pod")
		err = nethelper.ExecCmdOnAllPodInNamespace(
			[]string{"rm", "-rf", "/usr/bin/ping", "/usr/sbin/ping"})
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(netparameters.TestNamespace.Name,
			configSuite.General.TestImage,
			netparameters.TestDeploymentLabels)
		daemonSet.RedefineDaemonSetWithNodeSelector(map[string]string{configSuite.General.CnfNodeLabel: ""})
		By("Create DaemonSet on cluster")
		err = nethelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet.DaemonSet, netparameters.WaitingTime)

		By("Create DaemonSet on cluster")
		err = nethelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet.DaemonSet, netparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseDefaultSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			netparameters.TestCaseDefaultNetworkName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45446
	It("2 custom pods on Default network networking-icmpv4-connectivity skip when there is no ip binary [skip]",
		func() {

			By("Define deployment and create it on cluster")
			err = nethelper.DefineAndCreatePrivilegedDeploymentOnCluster(2)
			Expect(err).ToNot(HaveOccurred())

			By("Remove ip binary from test pod")
			err = nethelper.ExecCmdOnAllPodInNamespace(
				[]string{"rm", "-rf", "/usr/bin/ip", "/usr/sbin/ip"})
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				[]string{netparameters.NetworkingTestSuiteName},
				netparameters.TestCaseDefaultSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = nethelper.ValidateIfReportsAreValid(
				netparameters.TestCaseDefaultNetworkName,
				globalparameters.TestCaseSkipped)
			Expect(err).ToNot(HaveOccurred())
		})
})

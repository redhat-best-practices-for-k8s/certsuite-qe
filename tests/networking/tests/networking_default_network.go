package tests

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking custom namespace, custom deployment,", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file"))
	}

	execute.BeforeAll(func() {

		By("Clean namespace before all tests")
		err = namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, tsparams.TestNetworkingNameSpace)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {

		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())

	})

	// 45440
	It("3 custom pods on Default network networking-icmpv4-connectivity", func() {
		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentOnCluster(3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDefaultNetworkTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDefaultNetworkTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45441
	It("custom daemonset, 4 custom pods on Default network", func() {
		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(tsparams.TestNetworkingNameSpace, configSuite.General.TestImage,
			tsparams.TestDeploymentLabels, "daemonsetnetworkingput")
		daemonset.RedefineDaemonSetWithNodeSelector(daemonSet, map[string]string{configSuite.General.CnfNodeLabel: ""})

		By("Create DaemonSet on cluster")
		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDefaultNetworkTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDefaultNetworkTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45442
	It("3 custom pods on Default network networking-icmpv4-connectivity fail when "+
		"one pod is disconnected [negative]", func() {
		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreatePrivilegedDeploymentOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Close communication between deployment pods")
		podsList, err := globalhelper.GetListOfPodsInNamespace(tsparams.TestNetworkingNameSpace)
		Expect(err).ToNot(HaveOccurred())

		for index := range podsList.Items {
			_, err := globalhelper.ExecCommand(
				podsList.Items[0],
				[]string{"ip", "route", "add", podsList.Items[index].Status.PodIP, "via", "127.0.0.1"},
			)
			Expect(err).ToNot(HaveOccurred())
		}

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDefaultNetworkTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDefaultNetworkTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45444
	It("2 custom pods on Default network networking-icmpv4-connectivity skip when label "+
		"test-network-function.com/skip_connectivity_tests is set in deployment [skip]", func() {
		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithSkippedLabelOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Remove ping binary from test pod")
		err = tshelper.ExecCmdOnOnePodInNamespace(
			[]string{"rm", "-rf", "/usr/bin/ping", "/usr/sbin/ping"})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDefaultNetworkTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDefaultNetworkTcName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45445
	It("custom daemonset, 4 custom pods on Default network networking-icmpv4-connectivity pass when label "+
		"test-network-function.com/skip_connectivity_tests is set in deployment only", func() {
		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithSkippedLabelOnCluster(2)
		Expect(err).ToNot(HaveOccurred())

		By("Remove ping binary from test pod")
		err = tshelper.ExecCmdOnAllPodInNamespace(
			[]string{"rm", "-rf", "/usr/bin/ping", "/usr/sbin/ping"})
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(tsparams.TestNetworkingNameSpace, configSuite.General.TestImage,
			tsparams.TestDeploymentLabels, "daemonsetnetworkingput")
		daemonset.RedefineDaemonSetWithNodeSelector(daemonSet, map[string]string{configSuite.General.CnfNodeLabel: ""})

		By("Create DaemonSet on cluster")
		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDefaultNetworkTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDefaultNetworkTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

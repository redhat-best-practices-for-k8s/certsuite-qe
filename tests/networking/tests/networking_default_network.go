package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking custom namespace, custom deployment,", func() {
	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		randomNamespace = tsparams.TestNetworkingNameSpace + "-" + globalhelper.GenerateRandomString(10)

		By(fmt.Sprintf("Create %s namespace", randomNamespace))
		err := namespaces.Create(randomNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Override default report directory")
		origReportDir = globalhelper.GetConfiguration().General.TnfReportDir
		reportDir := origReportDir + "/" + randomNamespace
		globalhelper.OverrideReportDir(reportDir)

		By("Override default TNF config directory")
		origTnfConfigDir = globalhelper.GetConfiguration().General.TnfConfigDir
		configDir := origTnfConfigDir + "/" + randomNamespace
		globalhelper.OverrideTnfConfigDir(configDir)

		By("Define TNF config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By(fmt.Sprintf("Remove %s namespace", randomNamespace))
		err := namespaces.DeleteAndWait(
			globalhelper.GetAPIClient().CoreV1Interface,
			randomNamespace,
			tsparams.WaitingTime,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Restore default report directory")
		globalhelper.GetConfiguration().General.TnfReportDir = origReportDir

		By("Restore default TNF config directory")
		globalhelper.GetConfiguration().General.TnfConfigDir = origTnfConfigDir
	})

	// 45440
	It("3 custom pods on Default network networking-icmpv4-connectivity", func() {
		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(3, randomNamespace)
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
		err := tshelper.DefineAndCreateDeploymentOnCluster(2, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, configSuite.General.TestImage,
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
		err := tshelper.DefineAndCreatePrivilegedDeploymentOnCluster(2, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Close communication between deployment pods")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
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
		err := tshelper.DefineAndCreateDeploymentWithSkippedLabelOnCluster(2, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Remove ping binary from test pod")
		err = tshelper.ExecCmdOnOnePodInNamespace(
			[]string{"rm", "-rf", "/usr/bin/ping", "/usr/sbin/ping"}, randomNamespace)
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
		err := tshelper.DefineAndCreateDeploymentWithSkippedLabelOnCluster(2, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Remove ping binary from test pod")
		err = tshelper.ExecCmdOnAllPodInNamespace(
			[]string{"rm", "-rf", "/usr/bin/ping", "/usr/sbin/ping"}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, configSuite.General.TestImage,
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

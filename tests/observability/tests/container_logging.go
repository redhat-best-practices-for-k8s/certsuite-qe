package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/observability/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"
)

var _ = Describe(tsparams.TnfContainerLoggingTcName, func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			tshelper.GetTnfTargetPodLabelsSlice(),
			[]string{},
			[]string{},
			[]string{tsparams.CrdSuffix1, tsparams.CrdSuffix2}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.CrdDeployTimeoutMins)
	})

	// 51747
	It("One deployment one pod one container that prints two log lines", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.TwoLogLines})

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51753
	It("One deployment one pod one container that prints one log line", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51754
	It("One deployment one pod with two containers, both containers print two log lines to stdout", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.TwoLogLines, tsparams.TwoLogLines})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51755
	It("One daemonset with two containers, first prints two lines, the second one line", func() {
		if globalhelper.IsKindCluster() {
			Skip("Test skipped on KIND cluster due to newline char issue")
		}

		By("Deploy daemonset in the cluster")
		daemonSet := tshelper.DefineDaemonSetWithStdoutBuffers(
			tsparams.TestDaemonSetBaseName, randomNamespace, []string{tsparams.OneLogLineWithoutNewLine})

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet,
			tsparams.DaemonSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51756
	It("Two deployments, two pods with two containers each, all printing 1 log line", func() {
		if globalhelper.IsKindCluster() {
			Skip("Test skipped on KIND cluster due to newline char issue")
		}

		By("Create deployment1 in the cluster")
		deployment1 := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName+"1", randomNamespace, 2,
			[]string{tsparams.OneLogLineWithoutNewLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment1,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment2 in the cluster")
		deployment2 := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName+"2", randomNamespace, 2,
			[]string{tsparams.OneLogLineWithoutNewLine})

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment2,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51757
	It("One deployment and one statefulset, both having one pod with one container that prints one log "+
		"line each", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create statefulset in the cluster")
		statefulset := tshelper.DefineStatefulSetWithStdoutBuffers(
			tsparams.TestStatefulSetBaseName, randomNamespace, 1,
			[]string{tsparams.OneLogLine})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset,
			tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51758
	It("One pod with one container that prints one log line to stdout", func() {
		By("Create pod in the cluster")
		pod := tshelper.DefinePodWithStdoutBuffer(
			tsparams.TestPodBaseName, randomNamespace, tsparams.OneLogLine)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, tsparams.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51759
	It("One pod with one container that prints to stdout one log line starting with a tab char", func() {
		By("Create pod in the cluster")
		pod := tshelper.DefinePodWithStdoutBuffer(tsparams.TestPodBaseName, randomNamespace,
			"\t"+tsparams.OneLogLine)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, tsparams.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51760
	It("One deployment one pod one container without any log line to stdout [negative]", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.NoLogLines})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51761
	It("One deployment one pod two containers but only one printing one log line [negative]", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.OneLogLine, tsparams.NoLogLines})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51762
	It("Two deployments one pod two containers each, first deployment passing but second fails [negative]", func() {
		By("Create deployment1 in the cluster whose containers print one line to stdout each")
		deployment1 := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName+"1", randomNamespace, 1,
			[]string{tsparams.OneLogLine, tsparams.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment1,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment2 in the cluster but only the first of its containers prints a line to stdout")
		deployment2 := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName+"2", randomNamespace, 1,
			[]string{tsparams.OneLogLine, tsparams.NoLogLines})

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment2,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51763
	It("One pod one container without any log line to stdout [negative]", func() {
		By("Create pod in the cluster")
		pod := tshelper.DefinePodWithStdoutBuffer(tsparams.TestPodBaseName, randomNamespace,
			tsparams.NoLogLines)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, tsparams.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51764
	It("One deployment and one statefulset both one container each, but only deployment prints "+
		"one log line [negative]", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Deploy statefulset in the cluster")
		statefulset := tshelper.DefineStatefulSetWithStdoutBuffers(
			tsparams.TestStatefulSetBaseName, randomNamespace, 1,
			[]string{tsparams.NoLogLines})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset,
			tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51765
	It("One deployment one pod one container printing one log line without newline char", func() {
		if globalhelper.IsKindCluster() {
			Skip("Test skipped on KIND cluster due to newline char issue")
		}

		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.OneLogLineWithoutNewLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51767
	It("One deployment one pod two containers, first prints one line, second prints "+
		"one line without newline", func() {
		if globalhelper.IsKindCluster() {
			Skip("Test skipped on KIND cluster due to newline char issue")
		}

		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]string{tsparams.OneLogLine, tsparams.OneLogLineWithoutNewLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51768
	It("One deployment with one pod and one container without TNF target labels [skip]", func() {
		By("Create deployment without TNF target labels in the cluster")
		deployment := tshelper.DefineDeploymentWithoutTargetLabels(
			tsparams.TestDeploymentBaseName, randomNamespace)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerLoggingTcName, globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

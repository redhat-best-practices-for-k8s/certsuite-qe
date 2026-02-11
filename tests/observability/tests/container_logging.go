package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/observability/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/observability/parameters"
)

var _ = Describe(tsparams.CertsuiteContainerLoggingTcName, Label("observability", "ocp-required"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			tshelper.GetCertsuiteTargetPodLabelsSlice(),
			[]string{},
			[]string{},
			[]string{tsparams.CrdSuffix1, tsparams.CrdSuffix2}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.CrdDeployTimeoutMins)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
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

		By("Assert daemonSet is ready")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDaemonSet).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
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

		By("Assert deployment1 is ready")
		runningDeployment1, err := globalhelper.GetRunningDeployment(deployment1.Namespace, deployment1.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment1).ToNot(BeNil())

		By("Create deployment2 in the cluster")
		deployment2 := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName+"2", randomNamespace, 2,
			[]string{tsparams.OneLogLineWithoutNewLine})

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment2,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment2 is ready")
		runningDeployment2, err := globalhelper.GetRunningDeployment(deployment2.Namespace, deployment2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Create statefulset in the cluster")
		statefulset := tshelper.DefineStatefulSetWithStdoutBuffers(
			tsparams.TestStatefulSetBaseName, randomNamespace, 1,
			[]string{tsparams.OneLogLine})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset,
			tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulSet is ready")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(statefulset.Namespace, statefulset.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningStatefulSet).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51758
	It("One pod with one container that prints one log line to stdout", func() {
		By("Create pod in the cluster")
		pod := tshelper.DefinePodWithStdoutBuffer(
			tsparams.TestPodBaseName, randomNamespace, tsparams.OneLogLine)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, tsparams.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is ready")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, pod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51759
	It("One pod with one container that prints to stdout one log line starting with a tab char", func() {
		By("Create pod in the cluster")
		pod := tshelper.DefinePodWithStdoutBuffer(tsparams.TestPodBaseName, randomNamespace,
			"\t"+tsparams.OneLogLine)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, tsparams.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is ready")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, pod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCaseFailed, randomReportDir)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCaseFailed, randomReportDir)
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

		By("Assert deployment1 is ready")
		runningDeployment1, err := globalhelper.GetRunningDeployment(deployment1.Namespace, deployment1.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment1).ToNot(BeNil())

		By("Create deployment2 in the cluster but only the first of its containers prints a line to stdout")
		deployment2 := tshelper.DefineDeploymentWithStdoutBuffers(
			tsparams.TestDeploymentBaseName+"2", randomNamespace, 1,
			[]string{tsparams.OneLogLine, tsparams.NoLogLines})

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment2,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment2 is ready")
		runningDeployment2, err := globalhelper.GetRunningDeployment(deployment2.Namespace, deployment2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51763
	It("One pod one container without any log line to stdout [negative]", func() {
		By("Create pod in the cluster")
		pod := tshelper.DefinePodWithStdoutBuffer(tsparams.TestPodBaseName, randomNamespace,
			tsparams.NoLogLines)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, tsparams.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is ready")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, pod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCaseFailed, randomReportDir)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Deploy statefulset in the cluster")
		statefulset := tshelper.DefineStatefulSetWithStdoutBuffers(
			tsparams.TestStatefulSetBaseName, randomNamespace, 1,
			[]string{tsparams.NoLogLines})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset,
			tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulSet is ready")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(statefulset.Namespace, statefulset.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningStatefulSet).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCaseFailed, randomReportDir)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
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

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51768
	It("One deployment with one pod and one container without Certsuite target labels [skip]", func() {
		By("Create deployment without Certsuite target labels in the cluster")
		deployment := tshelper.DefineDeploymentWithoutTargetLabels(
			tsparams.TestDeploymentBaseName, randomNamespace)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Start Certsuite " + tsparams.CertsuiteContainerLoggingTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerLoggingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerLoggingTcName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

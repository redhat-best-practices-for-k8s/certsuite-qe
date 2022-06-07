package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/observability/observabilityhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/observability/observabilityparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe(observabilityparameters.TnfContainerLoggingTcName, func() {
	const tnfTestCaseName = observabilityparameters.TnfContainerLoggingTcName
	qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

	BeforeEach(func() {
		By("Clean namespace " + observabilityparameters.TestNamespace + " before each test")
		err := namespaces.Clean(observabilityparameters.TestNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51747
	It("One deployment one pod one container that prints two log lines", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.TwoLogLines})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51753
	It("One deployment one pod one container that prints one log line", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51754
	It("One deployment one pod with two containers, both containers print two log lines to stdout", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.TwoLogLines, observabilityparameters.TwoLogLines})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51755
	It("One daemonset with two containers, first prints two lines, the second one line", func() {

		By("Deploy daemonset in the cluster")
		daemonSet := observabilityhelper.DefineDaemonSetWithStdoutBuffers(
			observabilityparameters.TestDaemonSetBaseName,
			[]string{observabilityparameters.TwoLogLines, observabilityparameters.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet,
			observabilityparameters.DaemonSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51756
	It("Two deployments, two pods with two containers each, all printing 1 log line", func() {

		By("Create deployment1 in the cluster")
		deployment1 := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName+"1", 2,
			[]string{observabilityparameters.OneLogLine, observabilityparameters.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment1,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment2 in the cluster")
		deployment2 := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName+"2", 2,
			[]string{observabilityparameters.OneLogLine, observabilityparameters.OneLogLine})

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment2,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51757
	It("One deployment and one statefulset, both having one pod with one container that prints one log "+
		"line each", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create statefulset in the cluster")
		statefulset := observabilityhelper.DefineStatefulSetWithStdoutBuffers(
			observabilityparameters.TestStatefulSetBaseName, 1,
			[]string{observabilityparameters.OneLogLine})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset,
			observabilityparameters.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51758
	It("One pod with one container that prints one log line to stdout", func() {

		By("Create pod in the cluster")
		pod := observabilityhelper.DefinePodWithStdoutBuffer(
			observabilityparameters.TestPodBaseName, observabilityparameters.OneLogLine)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, observabilityparameters.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51759
	It("One pod with one container that prints to stdout one log line starting with a tab char", func() {

		By("Create pod in the cluster")
		pod := observabilityhelper.DefinePodWithStdoutBuffer(observabilityparameters.TestPodBaseName,
			"\t"+observabilityparameters.OneLogLine)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, observabilityparameters.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51760
	It("One deployment one pod one container without any log line to stdout [negative]", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.NoLogLines})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51761
	It("One deployment one pod two containers but only one printing one log line [negative]", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.OneLogLine, observabilityparameters.NoLogLines})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51762
	It("Two deployments one pod two containers each, first deployment passing but second fails [negative]", func() {

		By("Create deployment1 in the cluster whose containers print one line to stdout each")
		deployment1 := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName+"1", 1,
			[]string{observabilityparameters.OneLogLine, observabilityparameters.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment1,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment2 in the cluster but only the first of its containers prints a line to stdout")
		deployment2 := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName+"2", 1,
			[]string{observabilityparameters.OneLogLine, observabilityparameters.NoLogLines})

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment2,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51763
	It("One pod one container without any log line to stdout [negative]", func() {

		By("Create pod in the cluster")
		pod := observabilityhelper.DefinePodWithStdoutBuffer(observabilityparameters.TestPodBaseName,
			observabilityparameters.NoLogLines)

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, observabilityparameters.PodDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51764
	It("One deployment and one statefulset both one container each, but only deployment prints "+
		"one log line [negative]", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.OneLogLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Deploy statefulset in the cluster")
		statefulset := observabilityhelper.DefineStatefulSetWithStdoutBuffers(
			observabilityparameters.TestStatefulSetBaseName, 1,
			[]string{observabilityparameters.NoLogLines})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset,
			observabilityparameters.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51765
	It("One deployment one pod one container printing one log line without newline char [negative]", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.OneLogLineWithoutNewLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51767
	It("One deployment one pod two containers, first prints one line, second prints "+
		"one line without newline [negative]", func() {

		By("Create deployment in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithStdoutBuffers(
			observabilityparameters.TestDeploymentBaseName, 1,
			[]string{observabilityparameters.OneLogLine, observabilityparameters.OneLogLineWithoutNewLine})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51768
	It("One deployment with one pod and one container without TNF target labels [skip]", func() {

		By("Create deployment without TNF target labels in the cluster")
		deployment := observabilityhelper.DefineDeploymentWithoutTargetLabels(
			observabilityparameters.TestDeploymentBaseName)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			observabilityparameters.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})

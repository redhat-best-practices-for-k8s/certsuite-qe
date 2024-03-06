package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/observability/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"
)

var _ = Describe(tsparams.TnfTerminationMsgPolicyTcName, func() {
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

	// Positive #1.
	It("One deployment one pod one container with terminationMessagePolicy set to FallbackToLogsOnError", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has terminationMessagePolicy set to FallbackToLogsOnError")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].TerminationMessagePolicy).
			To(Equal(corev1.TerminationMessageFallbackToLogsOnError))

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// // Positive #2.
	It("One deployment one pod two containers both with terminationMessagePolicy set to FallbackToLogsOnError", func() {
		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{
				corev1.TerminationMessageFallbackToLogsOnError,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has terminationMessagePolicy set to FallbackToLogsOnError")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].TerminationMessagePolicy).
			To(Equal(corev1.TerminationMessageFallbackToLogsOnError))

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Positive #3.
	It("One daemonset with two containers, both with terminationMessagePolicy "+
		"set to FallbackToLogsOnError", func() {

		By("Define daemonset")
		daemonSet := tshelper.DefineDaemonSetWithTerminationMsgPolicies(tsparams.TestDaemonSetBaseName,
			randomNamespace,
			[]corev1.TerminationMessagePolicy{
				corev1.TerminationMessageFallbackToLogsOnError,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		By("Create daemonset in the cluster")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.DaemonSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert daemonset has terminationMessagePolicy set to FallbackToLogsOnError in both containers")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		for _, container := range runningDaemonSet.Spec.Template.Spec.Containers {
			Expect(container.TerminationMessagePolicy).To(Equal(corev1.TerminationMessageFallbackToLogsOnError))
		}

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Positive #4
	It("One deployment and one statefulset, both with one pod with one container, "+
		"all with terminationMessagePolicy set to FallbackToLogsOnError", func() {

		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has terminationMessagePolicy set to FallbackToLogsOnError")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].TerminationMessagePolicy).
			To(Equal(corev1.TerminationMessageFallbackToLogsOnError))

		By("Create statefulset in the cluster")
		statefulSet := tshelper.DefineStatefulSetWithTerminationMsgPolicies(tsparams.TestStatefulSetBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulset has terminationMessagePolicy set to FallbackToLogsOnError")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(statefulSet.Namespace, statefulSet.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningStatefulSet.Spec.Template.Spec.Containers[0].TerminationMessagePolicy).
			To(Equal(corev1.TerminationMessageFallbackToLogsOnError))

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #1.
	It("One deployment one pod one container without terminationMessagePolicy [negative]", func() {

		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #2.
	It("One deployment one pod two containers, only one container with terminationMessagePolicy "+
		"set to FallbackToLogsOnError [negative]", func() {

		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{
				tsparams.UseDefaultTerminationMsgPolicy,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert the two containers have different terminationMessagePolicy")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[1].TerminationMessagePolicy).
			To(Equal(corev1.TerminationMessageFallbackToLogsOnError))

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #3.
	It("One deployment with two pods with one container each without terminationMessagePolicy set [negative]", func() {

		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 2,
			[]corev1.TerminationMessagePolicy{
				tsparams.UseDefaultTerminationMsgPolicy,
			})

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #4.
	It("One deployment and one statefulset, both with one pod with one container, "+
		"only the deployment has terminationMessagePolicy set to FallbackToLogsOnError [negative]", func() {

		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has terminationMessagePolicy set to FallbackToLogsOnError")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].TerminationMessagePolicy).
			To(Equal(corev1.TerminationMessageFallbackToLogsOnError))

		By("Create statefulset in the cluster")
		statefulSet := tshelper.DefineStatefulSetWithTerminationMsgPolicies(tsparams.TestStatefulSetBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		By("Create statefulset in the cluster")
		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #5.
	It("One deployment and one daemonset, both with one pod with one container, "+
		"only the deployment has terminationMessagePolicy set to FallbackToLogsOnError [negative]", func() {

		By("Define deployment")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		By("Create daemonset")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset")
		daemonSet := tshelper.DefineDaemonSetWithTerminationMsgPolicies(tsparams.TestDaemonSetBaseName,
			randomNamespace,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		By("Create daemonset")
		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.DaemonSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Skip #1.
	It("One deployment with one pod and one container without TNF target labels [skip]", func() {

		By("Define deployment without TNF target labels in the cluster")
		deployment := tshelper.DefineDeploymentWithoutTargetLabels(tsparams.TestDeploymentBaseName, randomNamespace)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment,
			tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseSkipped,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

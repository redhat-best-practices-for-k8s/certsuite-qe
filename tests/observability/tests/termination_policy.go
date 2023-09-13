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
	qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			tshelper.GetTnfTargetPodLabelsSlice(),
			[]string{},
			[]string{},
			[]string{tsparams.CrdSuffix1, tsparams.CrdSuffix2})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.CrdDeployTimeoutMins)
	})

	// Positive #1.
	It("One deployment one pod one container with terminationMessagePolicy set to FallbackToLogsOnError", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// // Positive #2.
	It("One deployment one pod two containers both with terminationMessagePolicy set to FallbackToLogsOnError", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName, randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{
				corev1.TerminationMessageFallbackToLogsOnError,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Positive #3.
	It("One daemonset with two containers, both with terminationMessagePolicy "+
		"set to FallbackToLogsOnError", func() {

		By("Create deployment in the cluster")
		daemonSet := tshelper.DefineDaemonSetWithTerminationMsgPolicies(tsparams.TestDaemonSetBaseName,
			randomNamespace,
			[]corev1.TerminationMessagePolicy{
				corev1.TerminationMessageFallbackToLogsOnError,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			daemonSet, tsparams.DaemonSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Positive #4
	It("One deployment and one statefulset, both with one pod with one container, "+
		"all with terminationMessagePolicy set to FallbackToLogsOnError", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create statefulset in the cluster")
		statefulSet := tshelper.DefineStatefulSetWithTerminationMsgPolicies(tsparams.TestStatefulSetBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #1.
	It("One deployment one pod one container without terminationMessagePolicy [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #2.
	It("One deployment one pod two containers, only one container with terminationMessagePolicy "+
		"set to FallbackToLogsOnError [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{
				tsparams.UseDefaultTerminationMsgPolicy,
				corev1.TerminationMessageFallbackToLogsOnError,
			})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #3.
	It("One deployment with two pods with one container each without terminationMessagePolicy set [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 2,
			[]corev1.TerminationMessagePolicy{
				tsparams.UseDefaultTerminationMsgPolicy,
			})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #4.
	It("One deployment and one statefulset, both with one pod with one container, "+
		"only the deployment has terminationMessagePolicy set to FallbackToLogsOnError [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create statefulset in the cluster")
		statefulSet := tshelper.DefineStatefulSetWithTerminationMsgPolicies(tsparams.TestStatefulSetBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #5.
	It("One deployment and one daemonset, both with one pod with one container, "+
		"only the deployment has terminationMessagePolicy set to FallbackToLogsOnError [negative]", func() {

		By("Create deployment in the cluster")
		deployment := tshelper.DefineDeploymentWithTerminationMsgPolicies(tsparams.TestDeploymentBaseName,
			randomNamespace, 1,
			[]corev1.TerminationMessagePolicy{corev1.TerminationMessageFallbackToLogsOnError})

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create daemonset in the cluster")
		daemonSet := tshelper.DefineDaemonSetWithTerminationMsgPolicies(tsparams.TestDaemonSetBaseName,
			randomNamespace,
			[]corev1.TerminationMessagePolicy{tsparams.UseDefaultTerminationMsgPolicy})

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			daemonSet, tsparams.DaemonSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Skip #1.
	It("One deployment with one pod and one container without TNF target labels [skip]", func() {

		By("Create deployment without TNF target labels in the cluster")
		deployment := tshelper.DefineDeploymentWithoutTargetLabels(tsparams.TestDeploymentBaseName, randomNamespace)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfTerminationMsgPolicyTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfTerminationMsgPolicyTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTerminationMsgPolicyTcName, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})

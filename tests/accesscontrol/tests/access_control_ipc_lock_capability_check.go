package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Access-control ipc-lock-capability-check,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})

	// 63736
	It("one deployment, one pod, one container, does not have ipc lock capability", func() {
		By("Define deployment without ipc lock")
		dep, err := tshelper.DefineDeployment(1, 1, "acdeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have IPC_LOCK capability")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlIpcLockCapability,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlIpcLockCapability,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63737
	It("one deployment, one pod, one container, does have ipc lock capability [negative]", func() {
		By("Define deployment with ipc lock")
		dep, err := tshelper.DefineDeployment(1, 1, "acdeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextIpcLock(dep)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does have IPC_LOCK capability")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].
			SecurityContext.Capabilities.Add).To(ContainElement(corev1.Capability("IPC_LOCK")))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlIpcLockCapability,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlIpcLockCapability,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63738
	It("two deployments, one pod each, one container each, does not have ipc lock capability", func() {
		By("Define deployments without ipc lock")
		dep, err := tshelper.DefineDeployment(1, 1, "acdeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have IPC_LOCK capability")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())

		dep2, err := tshelper.DefineDeployment(1, 1, "acdeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have IPC_LOCK capability")
		runningDeployment, err = globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlIpcLockCapability,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlIpcLockCapability,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63739
	It("two deployments, one pod each, one container each, one does have ipc lock capability [negative]", func() {
		By("Define deployments with varying ipc lock capabilities")
		dep, err := tshelper.DefineDeployment(1, 1, "acdeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextIpcLock(dep)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does have IPC_LOCK capability")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].
			SecurityContext.Capabilities.Add).To(ContainElement(corev1.Capability("IPC_LOCK")))

		dep2, err := tshelper.DefineDeployment(1, 1, "acdeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have IPC_LOCK capability")
		runningDeployment, err = globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlIpcLockCapability,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlIpcLockCapability,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

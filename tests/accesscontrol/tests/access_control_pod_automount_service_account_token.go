package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
)

var _ = Describe("Access-control pod-automount-service-account-token, ", func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		err = globalhelper.CreateServiceAccount(
			tsparams.ServiceAccountName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.Timeout)
	})

	// 53033
	It("one deployment, one pod, token false", func() {
		By("Define deployment with automountServiceAccountToken set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep, tsparams.ServiceAccountName)

		deployment.RedefineWithAutomountServiceAccountToken(dep, false)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment automountServiceAccountToken is false")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).To(BeFalse())
		Expect(runningDeployment.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53034
	It("one deployment, one pod, token true [negative]", func() {
		By("Define deployment with automountServiceAccountToken set to true")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep, tsparams.ServiceAccountName)

		deployment.RedefineWithAutomountServiceAccountToken(dep, true)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment automountServiceAccountToken is true")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).To(BeTrue())
		Expect(runningDeployment.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53035
	It("one deployment, one pod, token not set, service account's token false", func() {
		By("Define deployment with automountServiceAccountToken not set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep, tsparams.ServiceAccountName)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment automountServiceAccountToken is nil")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).To(BeNil())
		Expect(runningDeployment.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		By("Set namespace's default serviceaccount's automountServiceAccountToken to false")
		err = tshelper.SetServiceAccountAutomountServiceAccountToken(randomNamespace,
			tsparams.ServiceAccountName, "false")
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53036
	It("one deployment, one pod, token not set, service account's token true [negative]", func() {
		By("Define deployment with automountServiceAccountToken not set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep, tsparams.ServiceAccountName)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Set namespace's default serviceaccount's automountServiceAccountToken to true")
		err = tshelper.SetServiceAccountAutomountServiceAccountToken(randomNamespace,
			tsparams.ServiceAccountName, "true")
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment serviceaccount name is set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		// Deployment is running so this is already nil
		By("Assert deployment automountServiceAccountToken nil")
		Expect(runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).To(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53040
	It("one deployment, one pod, token not set, service account's token not set [negative]", func() {
		By("Define deployment with automountServiceAccountToken not set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep, tsparams.ServiceAccountName)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Set namespace's default serviceaccount's automountServiceAccountToken to nil")
		err = tshelper.SetServiceAccountAutomountServiceAccountToken(randomNamespace,
			tsparams.ServiceAccountName, "nil")
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment serviceaccount name is set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		By("Assert deployment automountServiceAccountToken is nil")
		Expect(runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).To(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53054
	It("one deployment, one pod, token false, service account's token true", func() {
		By("Define deployment with automountServiceAccountToken set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep, tsparams.ServiceAccountName)

		deployment.RedefineWithAutomountServiceAccountToken(dep, false)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Set namespace's default serviceaccount's automountServiceAccountToken to true")
		err = tshelper.SetServiceAccountAutomountServiceAccountToken(randomNamespace,
			tsparams.ServiceAccountName, "true")
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment serviceaccount name is set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		By("Assert deployment automountServiceAccountToken is false")
		Expect(runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).To(BeFalse())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53036
	It("two deployments, one pod each, tokens false", func() {
		By("Define deployments with automountServiceAccountTokens set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep, tsparams.ServiceAccountName)

		deployment.RedefineWithAutomountServiceAccountToken(dep, false)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment automountServiceAccountToken is false")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).To(BeFalse())
		Expect(runningDeployment.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep2, tsparams.ServiceAccountName)

		deployment.RedefineWithAutomountServiceAccountToken(dep2, false)

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment automountServiceAccountToken is false")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2.Spec.Template.Spec.AutomountServiceAccountToken).ToNot(BeNil())
		Expect(*runningDeployment2.Spec.Template.Spec.AutomountServiceAccountToken).To(BeFalse())
		Expect(runningDeployment2.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53057
	It("two deployments, one pod each, one token true [negative]", func() {
		By("Define deployments with automountServiceAccountTokens set to different values")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep, tsparams.ServiceAccountName)

		deployment.RedefineWithAutomountServiceAccountToken(dep, true)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment automountServiceAccountToken is true")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.AutomountServiceAccountToken).To(BeTrue())
		Expect(runningDeployment.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithServiceAccount(dep2, tsparams.ServiceAccountName)

		deployment.RedefineWithAutomountServiceAccountToken(dep2, false)

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment automountServiceAccountToken is false")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2.Spec.Template.Spec.AutomountServiceAccountToken).ToNot(BeNil())
		Expect(*runningDeployment2.Spec.Template.Spec.AutomountServiceAccountToken).To(BeFalse())
		Expect(runningDeployment2.Spec.Template.Spec.ServiceAccountName).To(Equal(tsparams.ServiceAccountName))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodAutomountToken,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

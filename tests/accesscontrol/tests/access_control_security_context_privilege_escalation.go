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

var _ = Describe("Access-control security-context-privilege-escalation,", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 63992
	It("one deployment, one pod, one container, does not allow privilege escalation", func() {
		By("Define deployment without privilege escalation")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment does not allow privilege escalation")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63993
	It("one deployment, one pod, one container, does allow privilege escalation [negative]", func() {
		By("Define deployment with privilege escalation")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextAllowPrivilegeEscalation(dep, true)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment does allow privilege escalation")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation).To(BeTrue())
		Expect(*runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser).To(Equal(int64(0)))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63994
	It("two deployments, one pod each, one container each, does not allow privilege escalation", func() {
		By("Define deployments without privilege escalation")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextAllowPrivilegeEscalation(dep, false)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment1 does not allow privilege escalation")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation).To(BeFalse())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment2 does not allow privilege escalation")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63995
	It("two deployments, one pod each, one container each, one does allow privilege escalation [negative]", func() {
		By("Define deployments with varying privilege escalation capabilities")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextAllowPrivilegeEscalation(dep, true)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment1 does allow privilege escalation")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation).To(BeTrue())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment2 does not allow privilege escalation")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

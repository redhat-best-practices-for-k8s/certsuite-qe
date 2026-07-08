package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
)

var _ = Describe("Access-control security-context-read-only-root-file-system,", Label("accesscontrol-readonly-root-fs"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomPrivilegedNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("one deployment, one pod, one container with read-only root filesystem", func() {
		By("Define deployment with read-only root filesystem")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithReadOnlyRootFilesystem(dep, true)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment has read-only root filesystem enabled")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem).To(BeTrue())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlReadOnlyRootFileSystem,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlReadOnlyRootFileSystem,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, one container without read-only root filesystem [negative]", func() {
		By("Define deployment without read-only root filesystem")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithReadOnlyRootFilesystem(dep, false)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment does not have read-only root filesystem enabled")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem).ToNot(BeNil())
		Expect(*runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem).To(BeFalse())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlReadOnlyRootFileSystem,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlReadOnlyRootFileSystem,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("two deployments, one with and one without read-only root filesystem [negative]", func() {
		By("Define deployment with read-only root filesystem")
		dep1, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithReadOnlyRootFilesystem(dep1, true)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep1, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment1 has read-only root filesystem enabled")
		runningDeployment1, err := globalhelper.GetRunningDeployment(dep1.Namespace, dep1.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment1.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
		Expect(*runningDeployment1.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem).To(BeTrue())

		By("Define deployment without read-only root filesystem")
		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithReadOnlyRootFilesystem(dep2, false)

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment2 does not have read-only root filesystem enabled")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
		Expect(*runningDeployment2.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem).To(BeFalse())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlReadOnlyRootFileSystem,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlReadOnlyRootFileSystem,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

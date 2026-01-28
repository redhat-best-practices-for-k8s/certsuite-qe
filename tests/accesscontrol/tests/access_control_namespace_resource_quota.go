package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/resourcequota"
)

var _ = Describe("Access-control namespace-resource-quota,", Label("accesscontrol3"), func() {
	var randomNamespace string
	var randomNamespace2 string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Create additional namespace for deployment2")
		randomNamespace2 = tsparams.AdditionalNamespaceForResourceQuotas + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		By("Define certsuite config file")
		err = globalhelper.DefineCertsuiteConfig(
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

		By("Delete additional namespace for deployment2")
		err := globalhelper.DeleteNamespaceAndWait(randomNamespace2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56469
	It("one deployment, one pod in a namespace with resource quota", func() {
		By("Define deployment")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Define resource quota")
		resourceQuota := resourcequota.DefineResourceQuota("quota1", randomNamespace, tsparams.CPURequest,
			tsparams.MemoryRequest, tsparams.CPULimit, tsparams.MemoryLimit)

		By("Create resource quota")
		err = globalhelper.CreateResourceQuota(resourceQuota)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56470
	It("one deployment, one pod in a namespace without resource quota [negative]", func() {
		By("Define deployment")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56471
	It("two deployments, one pod each, both in a namespace with resource quota", func() {
		By("Define deployment 1")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Define resource quota")
		resourceQuota := resourcequota.DefineResourceQuota("quota1", randomNamespace, tsparams.CPURequest,
			tsparams.MemoryRequest, tsparams.CPULimit, tsparams.MemoryLimit)

		By("Create resource quota")
		err = globalhelper.CreateResourceQuota(resourceQuota)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment 2")
		dep2, err := tshelper.DefineDeploymentWithNamespace(1, 1, "accesscontroldeployment2",
			randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Define resource quota")
		resourceQuota = resourcequota.DefineResourceQuota("quota1", randomNamespace2, tsparams.CPURequest,
			tsparams.MemoryRequest, tsparams.CPULimit, tsparams.MemoryLimit)

		By("Create resource quota")
		err = globalhelper.CreateResourceQuota(resourceQuota)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56472
	It("two deployments, one pod each, one in a namespace without resource quota [negative]", func() {
		By("Define deployment 1")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment 2")
		dep2, err := tshelper.DefineDeploymentWithNamespace(1, 1, "accesscontroldeployment2",
			randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Define resource quota")
		resourceQuota := resourcequota.DefineResourceQuota("quota1", randomNamespace2, tsparams.CPURequest,
			tsparams.MemoryRequest, tsparams.CPULimit, tsparams.MemoryLimit)

		By("Create resource quota")
		err = globalhelper.CreateResourceQuota(resourceQuota)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

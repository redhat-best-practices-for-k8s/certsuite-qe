package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/resourcequota"
)

var _ = Describe("Access-control namespace-resource-quota,", func() {
	var randomNamespace string
	var randomNamespace2 string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)

		By("Create additional namespace for deployment2")
		randomNamespace2 = tsparams.AdditionalNamespaceForResourceQuotas + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.Timeout)

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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

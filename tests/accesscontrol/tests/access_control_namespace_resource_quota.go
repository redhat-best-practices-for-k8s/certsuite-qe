package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control namespace-resource-quota,", func() {
	var randomNamespace string
	var randomNamespace2 string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)

		By("Create additional namespace for deployment2")
		randomNamespace2 = tsparams.AdditionalNamespaceForResourceQuotas + "-" + globalhelper.GenerateRandomString(5)
		err := namespaces.Create(randomNamespace2, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)

		By("Delete additional namespace for deployment2")
		err := namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, randomNamespace2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56469
	It("one deployment, one pod in a namespace with resource quota", func() {
		By("Define deployment")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Apply resource quota to namespace")
		err = tshelper.DefineAndCreateResourceQuota(randomNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56470
	It("one deployment, one pod in a namespace without resource quota [negative]", func() {
		By("Define deployment")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56471
	It("two deployments, one pod each, both in a namespace with resource quota", func() {
		By("Define deployments")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Apply resource quota to namespace")
		err = tshelper.DefineAndCreateResourceQuota(randomNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeploymentWithNamespace(1, 1, "accesscontroldeployment2",
			randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Apply resource quota to namespace")
		err = tshelper.DefineAndCreateResourceQuota(randomNamespace2, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56472
	It("two deployments, one pod each, one in a namespace without resource quota [negative]", func() {
		By("Define deployments")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeploymentWithNamespace(1, 1, "accesscontroldeployment2",
			randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Apply resource quota to namespace")
		err = tshelper.DefineAndCreateResourceQuota(randomNamespace2, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

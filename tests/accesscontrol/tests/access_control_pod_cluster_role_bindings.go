package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
)

var _ = Describe("Access-control pod cluster role binding,", func() {
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

	// 56427
	It("one deployment, one pod, does not have cluster role binding", func() {
		By("Define deployment that do not have rolebinding")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlClusterRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlClusterRoleBindings,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, does  have cluster role binding [negative]", func() {
		By("Create cluster role binding")
		err := globalhelper.CreateClusterRoleBinding(randomNamespace, "my-service-account")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment with cluster role binding ")
		dep, err := tshelper.DefineDeploymentWithClusterRoleBindingWithServiceAccount(1, 1, "accesscontroldeployment",
			randomNamespace, "my-service-account")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlClusterRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlClusterRoleBindings,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

		By("Delete service account")
		err = globalhelper.DeleteClusterRoleBinding(randomNamespace, "my-cluster-role-binding")
		Expect(err).ToNot(HaveOccurred())

		By("Delete cluster role binding")
		err = globalhelper.DeleteServiceAccount(randomNamespace, "my-service-account")
		Expect(err).ToNot(HaveOccurred())
	})
})

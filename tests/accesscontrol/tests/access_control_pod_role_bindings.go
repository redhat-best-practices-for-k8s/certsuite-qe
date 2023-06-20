package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
)

const (
	testServiceAccount  = "my-sa"
	testRoleBindingName = "my-rb"
	testRoleName        = "my-r"

	testNamespace = "my-ns"
)

func setupInitialRbacConfiguration() {
	By("Create service account")

	err := globalhelper.CreateServiceAccount(testServiceAccount, tsparams.TestAccessControlNameSpace)
	Expect(err).ToNot(HaveOccurred())

	err = globalhelper.CreateRole(testRoleName, tsparams.TestAccessControlNameSpace)
	Expect(err).ToNot(HaveOccurred())

	err = globalhelper.CreateRoleBindingWithServiceAccountSubject(testRoleBindingName, testRoleName,
		testServiceAccount, tsparams.TestAccessControlNameSpace)
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("Access-control pod-role-bindings,", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Clean(testNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		setupInitialRbacConfiguration()
	})

	It("one pod with valid role binding", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.Configuration.General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, testServiceAccount)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-role-bindings")
		err = globalhelper.LaunchTests(
			tsparams.TnfPodRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfPodRoleBindings,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with service account in different namespace", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.Configuration.General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, testServiceAccount)
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Create namespace")
		err = namespaces.Create(testNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		// Delete service account
		err = globalhelper.DeleteServiceAccount(testServiceAccount, tsparams.TestAccessControlNameSpace)
		Expect(err).ToNot(HaveOccurred())
		// Create the service account in a new namespace
		err = globalhelper.CreateServiceAccount(testServiceAccount, testNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-role-bindings")
		err = globalhelper.LaunchTests(
			tsparams.TnfPodRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfPodRoleBindings,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with role binding in different namespace", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.Configuration.General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, testServiceAccount)
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Create namespace")
		err = namespaces.Create(testNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		// Delete role binding
		err = globalhelper.DeleteRoleBinding(testRoleBindingName, tsparams.TestAccessControlNameSpace)
		Expect(err).ToNot(HaveOccurred())
		// Create role binding in a new namespace
		err = globalhelper.CreateRoleBindingWithServiceAccountSubject(testServiceAccount, testRoleName, testServiceAccount, testNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-role-bindings")
		err = globalhelper.LaunchTests(
			tsparams.TnfPodRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfPodRoleBindings,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

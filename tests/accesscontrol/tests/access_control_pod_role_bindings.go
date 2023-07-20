package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
)

func setupInitialRbacConfiguration() {
	By("Create service account")

	err := globalhelper.CreateServiceAccount(tsparams.TestServiceAccount, tsparams.TestAccessControlNameSpace)
	Expect(err).ToNot(HaveOccurred())

	err = globalhelper.CreateRole(tsparams.TestRoleName, tsparams.TestAccessControlNameSpace)
	Expect(err).ToNot(HaveOccurred())

	err = globalhelper.CreateRoleBindingWithServiceAccountSubject(tsparams.TestRoleBindingName, tsparams.TestRoleName,
		tsparams.TestServiceAccount, tsparams.TestAccessControlNameSpace, tsparams.TestAccessControlNameSpace)
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("Access-control pod-role-bindings,", func() {

	execute.BeforeAll(func() {

		By("Create additional namespace for testing")
		// these namespaces will only be used for the access-control-namespace tests
		err := namespaces.Create(tsparams.TestAnotherNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Clean(tsparams.TestAnotherNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		setupInitialRbacConfiguration()
	})

	It("one pod with valid role binding", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)

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
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with no specified service account", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

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
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		// Delete service account
		err = globalhelper.DeleteServiceAccount(tsparams.TestServiceAccount, tsparams.TestAccessControlNameSpace)
		Expect(err).ToNot(HaveOccurred())
		// Create the service account in a new namespace
		err = globalhelper.CreateServiceAccount(tsparams.TestServiceAccount, tsparams.TestAnotherNamespace)
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
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		// Delete role binding
		err = globalhelper.DeleteRoleBinding(tsparams.TestRoleBindingName, tsparams.TestAccessControlNameSpace)
		Expect(err).ToNot(HaveOccurred())
		// Create role binding in a new namespace
		err = globalhelper.CreateRoleBindingWithServiceAccountSubject(tsparams.TestRoleBindingName,
			tsparams.TestRoleName, tsparams.TestServiceAccount, tsparams.TestAccessControlNameSpace, tsparams.TestAnotherNamespace)
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

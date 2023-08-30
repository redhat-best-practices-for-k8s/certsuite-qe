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

func setupInitialRbacConfiguration(namespace string) {
	By("Create service account")

	err := globalhelper.CreateServiceAccount(tsparams.TestServiceAccount, namespace)
	Expect(err).ToNot(HaveOccurred())

	err = globalhelper.CreateRole(tsparams.TestRoleName, namespace)
	Expect(err).ToNot(HaveOccurred())

	err = globalhelper.CreateRoleBindingWithServiceAccountSubject(tsparams.TestRoleBindingName, tsparams.TestRoleName,
		tsparams.TestServiceAccount, namespace, namespace)
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("Access-control pod-role-bindings,", func() {
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

		setupInitialRbacConfiguration(randomNamespace)
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})

	It("one pod with valid role binding", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
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

	It("one pod with no specified service account (default SA) [negative]", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
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

	It("one pod with service account in different namespace", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Create 'another' namespace")
		anotherNamespace := tsparams.TestAnotherNamespace + "-" + globalhelper.GenerateRandomString(5)
		err = namespaces.Create(anotherNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		// Delete service account
		err = globalhelper.DeleteServiceAccount(tsparams.TestServiceAccount, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		// Create the service account in a new namespace
		err = globalhelper.CreateServiceAccount(tsparams.TestServiceAccount, anotherNamespace)
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

		err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, anotherNamespace, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with role binding in different namespace", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		// Delete role binding
		err = globalhelper.DeleteRoleBinding(tsparams.TestRoleBindingName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create 'another' namespace")
		anotherNamespace := tsparams.TestAnotherNamespace + "-" + globalhelper.GenerateRandomString(5)
		err = namespaces.Create(anotherNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		// Create role binding in a new namespace
		err = globalhelper.CreateRoleBindingWithServiceAccountSubject(tsparams.TestRoleBindingName,
			tsparams.TestRoleName, tsparams.TestServiceAccount, randomNamespace, anotherNamespace)
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

		err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, anotherNamespace, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())
	})
})

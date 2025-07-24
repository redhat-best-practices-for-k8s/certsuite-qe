package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
)

func setupInitialRbacConfiguration(namespace string) {
	By("Create service account")

	err := globalhelper.CreateServiceAccount(
		tsparams.TestServiceAccount, namespace)
	Expect(err).ToNot(HaveOccurred())

	role := globalhelper.DefineRole(tsparams.TestRoleName, namespace)
	err = globalhelper.CreateRole(role)
	Expect(err).ToNot(HaveOccurred())

	err = globalhelper.CreateRoleBindingWithServiceAccountSubject(
		tsparams.TestRoleBindingName, tsparams.TestRoleName,
		tsparams.TestServiceAccount, namespace, namespace)
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("Access-control pod-role-bindings,", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")

		setupInitialRbacConfiguration(randomNamespace)
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("one pod with valid role binding", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-role-bindings")
		err = globalhelper.LaunchTests(
			tsparams.CertsuitePodRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuitePodRoleBindings,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with no specified service account (default SA) [negative]", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-role-bindings")
		err = globalhelper.LaunchTests(
			tsparams.CertsuitePodRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuitePodRoleBindings,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with service account in different namespace", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Create 'another' namespace")
		anotherNamespace := tsparams.TestAnotherNamespace + "-" + globalhelper.GenerateRandomString(5)
		err = globalhelper.CreateNamespace(anotherNamespace)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(anotherNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())
		})

		// Delete service account
		err = globalhelper.DeleteServiceAccount(
			tsparams.TestServiceAccount, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		// Create the service account in a new namespace
		err = globalhelper.CreateServiceAccount(
			tsparams.TestServiceAccount, anotherNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-role-bindings")
		err = globalhelper.LaunchTests(
			tsparams.CertsuitePodRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuitePodRoleBindings,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with role binding in different namespace", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		// Delete role binding
		err = globalhelper.DeleteRoleBinding(tsparams.TestRoleBindingName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create 'another' namespace")
		anotherNamespace := tsparams.TestAnotherNamespace + "-" + globalhelper.GenerateRandomString(5)
		err = globalhelper.CreateNamespace(anotherNamespace)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(anotherNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())
		})

		// Create role binding in a new namespace
		err = globalhelper.CreateRoleBindingWithServiceAccountSubject(
			tsparams.TestRoleBindingName,
			tsparams.TestRoleName, tsparams.TestServiceAccount, randomNamespace, anotherNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-role-bindings")
		err = globalhelper.LaunchTests(
			tsparams.CertsuitePodRoleBindings,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuitePodRoleBindings,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

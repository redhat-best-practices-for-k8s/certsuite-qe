package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
)

var _ = Describe("Access-control pod-service-account,", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

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

	It("one pod with valid service account", func() {
		By("Create service account")
		err := globalhelper.CreateServiceAccount(
			tsparams.TestServiceAccount, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod with service account")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod has valid service account configured")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Spec.ServiceAccountName).To(Equal(tsparams.TestServiceAccount))

		By("Start pod-service-account")
		err = globalhelper.LaunchTests(
			tsparams.CertsuitePodServiceAccount,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuitePodServiceAccount,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with empty service account [negative]", func() {
		By("Define pod with empty service account")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, "")
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod has default service account (Kubernetes assigns 'default' when empty)")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Spec.ServiceAccountName).To(Equal("default"))

		By("Start pod-service-account")
		err = globalhelper.LaunchTests(
			tsparams.CertsuitePodServiceAccount,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuitePodServiceAccount,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with default service account [negative]", func() {
		By("Define pod with default service account")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, "default")
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod has default service account")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Spec.ServiceAccountName).To(Equal("default"))

		By("Start pod-service-account")
		err = globalhelper.LaunchTests(
			tsparams.CertsuitePodServiceAccount,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuitePodServiceAccount,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

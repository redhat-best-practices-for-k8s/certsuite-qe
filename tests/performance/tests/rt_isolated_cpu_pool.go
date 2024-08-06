package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/runtimeclass"
)

var _ = Describe("performance-isolated-cpu-pool-rt-scheduling-policy", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.PerformanceNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		// Create service account and roles and roles binding
		err = tshelper.ConfigurePrivilegedServiceAccount(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	It("One pod running in isolated cpu pool and rt cpu scheduling policy", func() {

		By("Define runtime class")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.CertsuiteRunTimeClass)
		Expect(rtc).ToNot(BeNil())

		By("Create runtime class")
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete runtime class")
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Define pod")
		testPod, err := tshelper.DefineRtPodInIsolatedCPUPool(randomNamespace, rtc)
		Expect(err).To(BeNil())

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Change to rt scheduling policy")
		command := "chrt -f -p 20 1" // To change the scheduling policy of the container start process to FIFO scheduling
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Start isolated-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteRtIsolatedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteRtIsolatedCPUPoolSchedulingPolicy,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod running in isolated cpu pool and non-rt scheduling policy", func() {
		By("Define runtime class")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.CertsuiteRunTimeClass)
		Expect(rtc).ToNot(BeNil())

		By("Create runtime class")
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete runtime class")
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Define pod")
		testPod, err := tshelper.DefineRtPodInIsolatedCPUPool(randomNamespace, rtc)
		Expect(err).To(BeNil())

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, 2*tsparams.WaitingTime)
		Expect(err).NotTo(HaveOccurred())

		By("Start isolated-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteRtIsolatedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteRtIsolatedCPUPoolSchedulingPolicy,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod running in shared cpu pool", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start isolated-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteRtIsolatedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteRtIsolatedCPUPoolSchedulingPolicy,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

package tests

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
)

var _ = Describe("performance-rt-apps-no-exec-probes", Label("performance", "ocp-required"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

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

	It("Rt app pod with no exec probes", func() {
		By("Define pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, randomNamespace,
			tsparams.RtImageName, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, 2*tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		command := "chrt -f -p 50 1" // To change the scheduling policy of the container start process to FIFO scheduling
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Start rt-apps-no-exec-probes test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteRtAppsNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteRtAppsNoExecProbes,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Rt app pod with exec probes ", func() {
		By("Define pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, randomNamespace,
			tsparams.RtImageName, tsparams.CertsuiteTargetPodLabels)

		pod.RedefineWithLivenessProbe(testPod)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, 2*tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		command := "chrt -f -p 50 1" // To change the scheduling policy of the container start process to FIFO scheduling
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Start rt-apps-no-exec-probes test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteRtAppsNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteRtAppsNoExecProbes,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One non-Rt exclusive pod with no exec probes", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "1", "1")
		pod.RedefineWithMemoryResources(testPod, "512Mi", "512Mi")

		By("Create and wait until pod is ready")

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		By("Assert pod has modified resources")
		runningPod, err := globalhelper.GetRunningPod(testPod.Namespace, testPod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal("512Mi"))
		Expect(runningPod.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal("1"))

		By("Start rt-apps-no-exec-probes test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteRtAppsNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteRtAppsNoExecProbes,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

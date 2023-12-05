package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/performance/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/performance/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
)

var _ = Describe("performance-rt-apps-no-exec-probes", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.PerformanceNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())

		// Create service account and roles and roles binding
		err = tshelper.ConfigurePrivilegedServiceAccount(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	It("Rt app pod with no exec probes", func() {

		By("Define pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, randomNamespace,
			tsparams.RtImageName, tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, 2*tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		command := "chrt -f -p 50 1" // To change the scheduling policy of the container start process to FIFO scheduling
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Start rt-apps-no-exec-probes test")
		err = globalhelper.LaunchTests(tsparams.TnfRtAppsNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfRtAppsNoExecProbes,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Rt app pod with exec probes ", func() {
		By("Define pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, randomNamespace,
			tsparams.RtImageName, tsparams.TnfTargetPodLabels)

		pod.RedefineWithLivenessProbe(testPod)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, 2*tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		command := "chrt -f -p 50 1" // To change the scheduling policy of the container start process to FIFO scheduling
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Start rt-apps-no-exec-probes test")
		err = globalhelper.LaunchTests(tsparams.TnfRtAppsNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfRtAppsNoExecProbes,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One non-Rt exclusive pod with no exec probes", func() {

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "1", "1")
		pod.RedefineWithMemoryResources(testPod, "512Mi", "512Mi")

		By("Create and wait until pod is ready")
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod has modified resources")
		runningPod, err := globalhelper.GetRunningPod(testPod.Namespace, testPod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal("512Mi"))
		Expect(runningPod.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal("1"))

		By("Start rt-apps-no-exec-probes test")
		err = globalhelper.LaunchTests(tsparams.TnfRtAppsNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfRtAppsNoExecProbes,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

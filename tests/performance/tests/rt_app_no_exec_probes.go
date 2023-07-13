package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/performance/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/performance/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
)

var _ = Describe("performance-rt-apps-no-exec-probes", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PerformanceNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.PerformanceNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Rt app pod with no exec probes", func() {

		By("Define pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, tsparams.PerformanceNamespace,
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

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfRtAppsNoExecProbes,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Rt app pod with exec probes ", func() {
		By("Define pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, tsparams.PerformanceNamespace,
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

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfRtAppsNoExecProbes,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One non-Rt exclusive pod with no exec probes", func() {

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TnfTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "1", "1")
		pod.RedefineWithMemoryResources(testPod, "512Mi", "512Mi")

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start rt-apps-no-exec-probes test")
		err = globalhelper.LaunchTests(tsparams.TnfRtAppsNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfRtAppsNoExecProbes,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

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
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/runtimeclass"
)

const disableVar = "disable"

var (
	// Each tc will save the RTC that was created in order to delete them.
	rtcNames = []string{}
)

var _ = Describe("performance-isolated-cpu-pool-rt-scheduling-policy", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PerformanceNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test in order to enable RunTimeClass deletion.")
		err := namespaces.Clean(tsparams.PerformanceNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Delete all RTC's that were created by the previous test case.")
		for _, rtc := range rtcNames {
			By("Deleting rtc " + rtc)
			err := tshelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		}

		// clear the list.
		rtcNames = []string{}
	})

	It("One pod running in isolated cpu pool with whole cpu unit and rt cpu scheduling policy", func() {

		By("Define pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			tsparams.RtImageName, tsparams.TnfTargetPodLabels)

		annotationsMap := make(map[string]string)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		testPod.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		rtcNames = append(rtcNames, tsparams.TnfRunTimeClass)

		pod.RedefineWithRunTimeClass(testPod, rtc.Name)
		pod.RedefineWithCPUResources(testPod, "1", "1")

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Change to rt scheduling policy")
		command := "chrt -f -p 20 1" // To change the scheduling policy of the container start process to FIFO scheduling
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Start isolated-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfRtIsolatedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfRtIsolatedCPUPoolSchedulingPolicy,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod running in isolated cpu pool and non-rt scheduling policy", func() {
		By("Define pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			tsparams.RtImageName, tsparams.TnfTargetPodLabels)

		annotationsMap := make(map[string]string)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		testPod.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		rtcNames = append(rtcNames, tsparams.TnfRunTimeClass)

		pod.RedefineWithRunTimeClass(testPod, rtc.Name)
		pod.RedefineWithCPUResources(testPod, "1", "1")

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, 2*tsparams.WaitingTime)
		Expect(err).To(HaveOccurred())

		By("Start isolated-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfRtIsolatedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfRtIsolatedCPUPoolSchedulingPolicy,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod running in shared cpu pool", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start isolated-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfRtIsolatedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfRtIsolatedCPUPoolSchedulingPolicy,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})

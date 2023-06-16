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

var _ = Describe("performance-exclusive-cpu-pool-rt-scheduling-policy",
	func() {

		BeforeEach(func() {
			By("Clean namespace before each test")
			err := namespaces.Clean(tsparams.PerformanceNamespace, globalhelper.APIClient)
			Expect(err).ToNot(HaveOccurred())
		})

		It("One pod running in exclusive cpu pool and shared cpu scheduling policy", func() {

			By("Define RT pod")
			testPod := tshelper.DefineRtPod(tsparams.TestPodName, tsparams.PerformanceNamespace,
				tsparams.RtImageName, tsparams.TnfTargetPodLabels)

			err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())

			By("Start exclusive-cpu-pool-rt-scheduling-policy test")
			err = globalhelper.LaunchTests(
				tsparams.TnfRtExclusiveCPUPoolSchedulingPolicy,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).NotTo(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfRtExclusiveCPUPoolSchedulingPolicy,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())
		})

		It("One pod running in exclusive cpu pool and valid rt cpu scheduling policy", func() {

			By("Define RT pod")
			testPod := tshelper.DefineRtPod(tsparams.TestPodName, tsparams.PerformanceNamespace,
				tsparams.RtImageName, tsparams.TnfTargetPodLabels)

			err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())

			By("Change to rt scheduling policy")
			command := "chrt -f -p 9 1" // To change the scheduling policy of the container start process to FIFO scheduling
			_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
			Expect(err).To(BeNil())

			By("Start exclusive-cpu-pool-rt-scheduling-policy test")
			err = globalhelper.LaunchTests(
				tsparams.TnfRtExclusiveCPUPoolSchedulingPolicy,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).NotTo(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfRtExclusiveCPUPoolSchedulingPolicy,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())
		})

		It("One pod running in exclusive cpu pool and invalid rt cpu scheduling policy", func() {

			By("Define RT pod")
			testPod := tshelper.DefineRtPod(tsparams.TestPodName, tsparams.PerformanceNamespace,
				tsparams.RtImageName, tsparams.TnfTargetPodLabels)

			err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())

			By("Change to rt scheduling policy")
			command := "chrt -f -p 20 1" // To change the scheduling policy of the container start process to FIFO scheduling
			_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
			Expect(err).To(BeNil())

			By("Start exclusive-cpu-pool-rt-scheduling-policy test")
			err = globalhelper.LaunchTests(
				tsparams.TnfRtExclusiveCPUPoolSchedulingPolicy,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).To(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfRtExclusiveCPUPoolSchedulingPolicy,
				globalparameters.TestCaseFailed)
			Expect(err).ToNot(HaveOccurred())
		})

		It("One pod running in shared cpu pool", func() {
			By("Define pod")
			testPod := pod.DefinePod(tsparams.TestPodName, tsparams.PerformanceNamespace,
				globalhelper.Configuration.General.TestImage, tsparams.TnfTargetPodLabels)

			err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())

			By("Start exclusive-cpu-pool-rt-scheduling-policy test")
			err = globalhelper.LaunchTests(
				tsparams.TnfRtExclusiveCPUPoolSchedulingPolicy,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfRtExclusiveCPUPoolSchedulingPolicy,
				globalparameters.TestCaseSkipped)
			Expect(err).ToNot(HaveOccurred())
		})
	})

package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/performance/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
)

var _ = Describe("performance-shared-cpu-pool-non-rt-scheduling-policy", func() {

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

	It("One pod with container running in shared cpu pool", func() {

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start shared-cpu-pool-non-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfSharedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfSharedCPUPoolSchedulingPolicy,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with container running in exclusive cpu pool", func() {

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TnfTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "1", "1")
		pod.RedefineWithMemoryResources(testPod, "512Mi", "512Mi")

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start shared-cpu-pool-non-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfSharedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfSharedCPUPoolSchedulingPolicy,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})

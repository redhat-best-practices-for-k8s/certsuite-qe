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

var _ = Describe("performance-exclusive-cpu-pool", Serial, func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PerformanceNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.PerformanceNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with only exclusive containers", func() {
		By("Define pod")
		testPod := tshelper.DefineExclusivePod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start exclusive-cpu-pool test")
		err = globalhelper.LaunchTests(tsparams.TnfExclusiveCPUPool,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfExclusiveCPUPool,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with one exclusive container, and one shared container", func() {

		By("Define pod")
		testPod := tshelper.DefineExclusivePod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		tshelper.RedefinePodWithSharedContainer(testPod, 0)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start exclusive-cpu-pool test")
		err = globalhelper.LaunchTests(tsparams.TnfExclusiveCPUPool,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfExclusiveCPUPool,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with only shared containers", func() {

		By("Define pod")
		testPod := tshelper.DefineExclusivePod(tsparams.TestPodName, tsparams.PerformanceNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "0.75", "0.5")

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start exclusive-cpu-pool test")
		err = globalhelper.LaunchTests(tsparams.TnfExclusiveCPUPool,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfExclusiveCPUPool,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

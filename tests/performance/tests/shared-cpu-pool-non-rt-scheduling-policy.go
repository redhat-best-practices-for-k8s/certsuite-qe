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

var _ = Describe("performance-shared-cpu-pool-non-rt-scheduling-policy", func() {
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

	It("One pod with container running in shared cpu pool", func() {

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start shared-cpu-pool-non-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfSharedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfSharedCPUPoolSchedulingPolicy,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with container running in exclusive cpu pool", func() {

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "1", "1")
		pod.RedefineWithMemoryResources(testPod, "512Mi", "512Mi")

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start shared-cpu-pool-non-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfSharedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfSharedCPUPoolSchedulingPolicy,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})

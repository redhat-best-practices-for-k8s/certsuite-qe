package tests

import (
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
)

var _ = Describe("performance-exclusive-cpu-pool", func() {
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

		if globalhelper.IsKindCluster() && runtime.NumCPU() <= 2 {
			Skip("This test requires more than 2 CPU cores")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	It("One pod with only exclusive containers", func() {
		if globalhelper.IsKindCluster() {
			// We cannot guarantee the number of available CPUs so we skip this test
			Skip("Exclusive CPU pool is not supported on Kind cluster, skipping...")
		}

		By("Define pod")
		testPod := tshelper.DefineExclusivePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Start exclusive-cpu-pool test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteExclusiveCPUPool,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteExclusiveCPUPool,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with one exclusive container, and one shared container", func() {
		By("Define pod")
		testPod := tshelper.DefineExclusivePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		tshelper.RedefinePodWithSharedContainer(testPod, 0)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Start exclusive-cpu-pool test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteExclusiveCPUPool,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteExclusiveCPUPool,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with only shared containers", func() {
		By("Define pod")
		testPod := tshelper.DefineExclusivePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "0.75", "0.5")

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Start exclusive-cpu-pool test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteExclusiveCPUPool,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteExclusiveCPUPool,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

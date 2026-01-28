package tests

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/crd"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("performance-exclusive-cpu-pool-rt-scheduling-policy", Label("performance", "ocp-required"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		By("Check if running on Kind cluster")
		if globalhelper.IsKindCluster() {
			Skip("Exclusive CPU pool tests are not supported on Kind cluster")
		}

		By("Verify cluster has worker nodes")
		if !globalhelper.HasWorkerNodes() {
			Skip("Cluster has no worker nodes - skipping exclusive CPU pool tests")
		}

		By("Verify PerformanceProfile CRD exists")
		crdExists, err := crd.EnsureCrdExists(tsparams.PerformanceProfileCrd)
		Expect(err).ToNot(HaveOccurred())
		if !crdExists {
			Skip("PerformanceProfile CRD does not exist - exclusive CPU pool tests require a performance profile")
		}

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomPrivilegedNamespace(tsparams.PerformanceNamespace)

		By("Define certsuite config file")
		err = globalhelper.DefineCertsuiteConfig(
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

	It("One pod running in exclusive cpu pool and shared cpu scheduling policy", func() {
		By("Define RT pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, randomNamespace,
			tsparams.RtImageName, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is running and has containers")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")
		Expect(len(runningPod.Status.ContainerStatuses)).To(BeNumerically(">", 0), "Pod should have containers")

		By("Assert all containers are ready")
		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Start exclusive-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteRtExclusiveCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteRtExclusiveCPUPoolSchedulingPolicy,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod running in exclusive cpu pool and valid rt cpu scheduling policy", func() {

		By("Define RT pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, randomNamespace,
			tsparams.RtImageName, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is running and has containers")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")
		Expect(len(runningPod.Status.ContainerStatuses)).To(BeNumerically(">", 0), "Pod should have containers")

		By("Assert all containers are ready")
		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Change to rt scheduling policy")
		command := "chrt -f -p 9 1" // To change the scheduling policy of the container start process to FIFO scheduling
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Start exclusive-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteRtExclusiveCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteRtExclusiveCPUPoolSchedulingPolicy,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod running in exclusive cpu pool and invalid rt cpu scheduling policy", func() {

		By("Define RT pod")
		testPod := tshelper.DefineRtPod(tsparams.TestPodName, randomNamespace,
			tsparams.RtImageName, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Change to rt scheduling policy")
		command := "chrt -f -p 20 1" // To change the scheduling policy of the container start process to FIFO scheduling
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Start exclusive-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteRtExclusiveCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteRtExclusiveCPUPoolSchedulingPolicy,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod running in shared cpu pool", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Start exclusive-cpu-pool-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteRtExclusiveCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteRtExclusiveCPUPoolSchedulingPolicy,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

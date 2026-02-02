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
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("performance-shared-cpu-pool-non-rt-scheduling-policy", Label("performance", "ocp-required"), func() {
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
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	It("One pod with container running in shared cpu pool", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is running and in shared CPU pool configuration")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")
		Expect(len(runningPod.Spec.Containers)).To(BeNumerically(">", 0), "Pod should have containers")

		// Log container resources for debugging
		GinkgoWriter.Printf("Pod has %d containers\n", len(runningPod.Spec.Containers))
		for i, container := range runningPod.Spec.Containers {
			GinkgoWriter.Printf("Container[%d] name: %s\n", i, container.Name)
			GinkgoWriter.Printf("Container[%d] CPU requests: %v\n", i, container.Resources.Requests.Cpu())
			GinkgoWriter.Printf("Container[%d] CPU limits: %v\n", i, container.Resources.Limits.Cpu())
		}

		By("Assert all containers are ready")
		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		// Note: We skip the chrt scheduling policy check here because the ubi-micro image
		// doesn't have the chrt command available. The certsuite will verify the scheduling
		// policy using its own probe pod which has the necessary tools.

		By("Start shared-cpu-pool-non-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with container running in exclusive cpu pool", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "1", "1")
		pod.RedefineWithMemoryResources(testPod, "512Mi", "512Mi")

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is running with exclusive CPU pool configuration")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")
		Expect(len(runningPod.Spec.Containers)).To(BeNumerically(">", 0), "Pod should have containers")

		// Log container resources for debugging
		GinkgoWriter.Printf("Pod has %d containers\n", len(runningPod.Spec.Containers))
		for i, container := range runningPod.Spec.Containers {
			GinkgoWriter.Printf("Container[%d] name: %s\n", i, container.Name)
			GinkgoWriter.Printf("Container[%d] CPU requests: %v\n", i, container.Resources.Requests.Cpu())
			GinkgoWriter.Printf("Container[%d] CPU limits: %v\n", i, container.Resources.Limits.Cpu())
			GinkgoWriter.Printf("Container[%d] Memory requests: %v\n", i, container.Resources.Requests.Memory())
			GinkgoWriter.Printf("Container[%d] Memory limits: %v\n", i, container.Resources.Limits.Memory())
		}

		// Verify exclusive CPU pool requirements (whole unit CPUs, limits=requests)
		cpuRequest := runningPod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
		cpuLimit := runningPod.Spec.Containers[0].Resources.Limits.Cpu().MilliValue()
		GinkgoWriter.Printf("CPU request: %dm, CPU limit: %dm\n", cpuRequest, cpuLimit)
		Expect(cpuRequest).To(Equal(cpuLimit), "CPU request should equal CPU limit for exclusive pool")
		Expect(cpuRequest%1000).To(Equal(int64(0)), "CPU should be a whole unit for exclusive pool")

		By("Assert all containers are ready")
		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Start shared-cpu-pool-non-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).NotTo(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

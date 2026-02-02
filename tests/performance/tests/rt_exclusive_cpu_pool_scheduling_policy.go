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

		// Log container resources for debugging exclusive CPU pool classification
		GinkgoWriter.Printf("Pod has %d containers\n", len(runningPod.Spec.Containers))
		for i, container := range runningPod.Spec.Containers {
			GinkgoWriter.Printf("Container[%d] name: %s\n", i, container.Name)
			GinkgoWriter.Printf("Container[%d] CPU requests: %v\n", i, container.Resources.Requests.Cpu())
			GinkgoWriter.Printf("Container[%d] CPU limits: %v\n", i, container.Resources.Limits.Cpu())
			GinkgoWriter.Printf("Container[%d] Memory requests: %v\n", i, container.Resources.Requests.Memory())
			GinkgoWriter.Printf("Container[%d] Memory limits: %v\n", i, container.Resources.Limits.Memory())
		}

		By("Verify pod meets exclusive CPU pool requirements")
		// Exclusive CPU pool requires: CPU limits=requests as whole units, memory limits=requests
		cpuRequest := runningPod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
		cpuLimit := runningPod.Spec.Containers[0].Resources.Limits.Cpu().MilliValue()
		memRequest := runningPod.Spec.Containers[0].Resources.Requests.Memory().Value()
		memLimit := runningPod.Spec.Containers[0].Resources.Limits.Memory().Value()

		GinkgoWriter.Printf("CPU request: %dm, CPU limit: %dm\n", cpuRequest, cpuLimit)
		GinkgoWriter.Printf("Memory request: %d, Memory limit: %d\n", memRequest, memLimit)

		Expect(cpuRequest).To(Equal(cpuLimit), "CPU request should equal CPU limit for exclusive pool")
		Expect(cpuRequest%1000).To(Equal(int64(0)), "CPU should be a whole unit for exclusive pool")
		Expect(memRequest).To(Equal(memLimit), "Memory request should equal memory limit for guaranteed QoS")

		By("Assert all containers are ready")
		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Verify current scheduling policy before test")
		command := "chrt -p 1"
		stdout, _, err := tshelper.ExecCommandContainer(testPod, command)
		Expect(err).ToNot(HaveOccurred(), "Failed to check scheduling policy")
		GinkgoWriter.Printf("Current scheduling policy: %s\n", stdout)
		// Default should be SCHED_OTHER with priority 0, which is valid for exclusive CPU pool
		Expect(stdout).To(ContainSubstring("priority: 0"),
			"Process should have scheduling priority 0 (SCHED_OTHER default)")

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

		// Log container resources for debugging exclusive CPU pool classification
		GinkgoWriter.Printf("Pod has %d containers\n", len(runningPod.Spec.Containers))
		for i, container := range runningPod.Spec.Containers {
			GinkgoWriter.Printf("Container[%d] name: %s\n", i, container.Name)
			GinkgoWriter.Printf("Container[%d] CPU requests: %v\n", i, container.Resources.Requests.Cpu())
			GinkgoWriter.Printf("Container[%d] CPU limits: %v\n", i, container.Resources.Limits.Cpu())
			GinkgoWriter.Printf("Container[%d] Memory requests: %v\n", i, container.Resources.Requests.Memory())
			GinkgoWriter.Printf("Container[%d] Memory limits: %v\n", i, container.Resources.Limits.Memory())
		}

		By("Verify pod meets exclusive CPU pool requirements")
		cpuRequest := runningPod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
		cpuLimit := runningPod.Spec.Containers[0].Resources.Limits.Cpu().MilliValue()
		memRequest := runningPod.Spec.Containers[0].Resources.Requests.Memory().Value()
		memLimit := runningPod.Spec.Containers[0].Resources.Limits.Memory().Value()

		GinkgoWriter.Printf("CPU request: %dm, CPU limit: %dm\n", cpuRequest, cpuLimit)
		GinkgoWriter.Printf("Memory request: %d, Memory limit: %d\n", memRequest, memLimit)

		Expect(cpuRequest).To(Equal(cpuLimit), "CPU request should equal CPU limit for exclusive pool")
		Expect(cpuRequest%1000).To(Equal(int64(0)), "CPU should be a whole unit for exclusive pool")
		Expect(memRequest).To(Equal(memLimit), "Memory request should equal memory limit for guaranteed QoS")

		By("Assert all containers are ready")
		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Change to rt scheduling policy (SCHED_FIFO with priority 9)")
		command := "chrt -f -p 9 1" // SCHED_FIFO scheduling with priority 9 (< 10, valid for exclusive CPU pool)
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Verify scheduling policy was changed successfully")
		verifyCommand := "chrt -p 1"
		stdout, _, err := tshelper.ExecCommandContainer(testPod, verifyCommand)
		Expect(err).ToNot(HaveOccurred(), "Failed to verify scheduling policy")
		GinkgoWriter.Printf("Scheduling policy after change: %s\n", stdout)
		Expect(stdout).To(ContainSubstring("SCHED_FIFO"), "Scheduling policy should be SCHED_FIFO")
		Expect(stdout).To(ContainSubstring("priority: 9"), "Scheduling priority should be 9")

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

		By("Assert pod is running and has containers")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")

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

		By("Change to invalid rt scheduling policy (SCHED_FIFO with priority 20)")
		command := "chrt -f -p 20 1" // SCHED_FIFO with priority 20 (>= 10, invalid for exclusive CPU pool)
		_, _, err = tshelper.ChangeSchedulingPolicy(testPod, command)
		Expect(err).To(BeNil())

		By("Verify scheduling policy was changed successfully")
		verifyCommand := "chrt -p 1"
		stdout, _, err := tshelper.ExecCommandContainer(testPod, verifyCommand)
		Expect(err).ToNot(HaveOccurred(), "Failed to verify scheduling policy")
		GinkgoWriter.Printf("Scheduling policy after change: %s\n", stdout)
		Expect(stdout).To(ContainSubstring("SCHED_FIFO"), "Scheduling policy should be SCHED_FIFO")
		Expect(stdout).To(ContainSubstring("priority: 20"), "Scheduling priority should be 20")

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

		By("Assert pod is running in shared CPU pool configuration")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")

		// Log container resources for debugging - should NOT have exclusive CPU pool config
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

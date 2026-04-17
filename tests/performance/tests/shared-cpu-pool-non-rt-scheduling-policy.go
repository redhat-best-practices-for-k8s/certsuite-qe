package tests

import (
	"fmt"
	"strconv"
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
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

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

		By("Assert pod matches certsuite preconditions: non-guaranteed QoS, no HostPID")
		Expect(runningPod.Status.QOSClass).ToNot(Equal(corev1.PodQOSGuaranteed),
			"Pod must not be Guaranteed QoS — certsuite skips guaranteed pods for shared CPU pool checks")
		Expect(runningPod.Spec.HostPID).To(BeFalse(),
			"Pod must not use HostPID — certsuite skips pods with HostPID")

		for i, container := range runningPod.Spec.Containers {
			GinkgoWriter.Printf("Container[%d] name=%s CPU req=%v lim=%v\n",
				i, container.Name, container.Resources.Requests.Cpu(), container.Resources.Limits.Cpu())
		}

		By("Assert all containers are ready")

		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Detect scheduling policy on the test pod to determine expected certsuite outcome")
		chrtOutput, chrtErr := globalhelper.ExecCommand(*runningPod, []string{"chrt", "-p", "1"})
		Expect(chrtErr).ToNot(HaveOccurred(), "chrt -p 1 must succeed to determine expected certsuite outcome")

		chrtStr := chrtOutput.String()
		GinkgoWriter.Printf("chrt -p 1 output: %s\n", chrtStr)

		expectedResult := globalparameters.TestCasePassed

		if strings.Contains(chrtStr, "SCHED_FIFO") || strings.Contains(chrtStr, "SCHED_RR") {
			GinkgoWriter.Printf("RT scheduling detected on PID 1 — expecting certsuite to return FAILED\n")
			expectedResult = globalparameters.TestCaseFailed
		}

		By("Start shared-cpu-pool-non-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy,
			expectedResult, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify CheckDetails report completeness")
		checkDetails, err := globalhelper.GetTestCaseCheckDetails(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		totalObjects := len(checkDetails.CompliantObjectsOut) + len(checkDetails.NonCompliantObjectsOut)
		GinkgoWriter.Printf("CheckDetails: %d compliant, %d non-compliant objects\n",
			len(checkDetails.CompliantObjectsOut), len(checkDetails.NonCompliantObjectsOut))
		Expect(totalObjects).To(BeNumerically(">=", 1),
			"At least one report object expected (no early abort)")

		for i, obj := range checkDetails.CompliantObjectsOut {
			GinkgoWriter.Printf("Compliant[%d]: type=%s reason=%s\n",
				i, obj.ObjectType, globalhelper.GetReportObjectFieldValue(obj, "Reason"))
		}

		for i, obj := range checkDetails.NonCompliantObjectsOut {
			reason := globalhelper.GetReportObjectFieldValue(obj, "Reason")
			GinkgoWriter.Printf("NonCompliant[%d]: type=%s reason=%s\n", i, obj.ObjectType, reason)
			Expect(reason).ToNot(ContainSubstring("could not determine scheduling policy"),
				"Unhandled scheduling policy error found in non-compliant objects")
		}
	})

	It("One pod with ephemeral subprocesses in shared cpu pool", func() {
		By("Define pod with ephemeral subprocesses")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.RtImageName, tsparams.CertsuiteTargetPodLabels)

		spawnCmd := []string{"/bin/bash", "-c",
			"while true; do for i in $(seq 1 5); do sleep 0.05 & done; wait; done"}
		err := pod.RedefineWithContainerExecCommand(testPod, spawnCmd, 0)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		By("Assert pod matches certsuite preconditions: non-guaranteed QoS, no HostPID")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.QOSClass).ToNot(Equal(corev1.PodQOSGuaranteed),
			"Pod must not be Guaranteed QoS — certsuite skips guaranteed pods for shared CPU pool checks")
		Expect(runningPod.Spec.HostPID).To(BeFalse(),
			"Pod must not use HostPID — certsuite skips pods with HostPID")

		By("Verify ephemeral subprocesses are running")
		psOutput, psErr := globalhelper.ExecCommand(*runningPod, []string{"pgrep", "-c", "sleep"})
		Expect(psErr).ToNot(HaveOccurred(), "pgrep must succeed to verify subprocesses are running")

		GinkgoWriter.Printf("Active sleep subprocesses: %s\n", strings.TrimSpace(psOutput.String()))

		By("Start shared-cpu-pool-non-rt-scheduling-policy test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case completed (any terminal status accepted)")
		err = globalhelper.ValidateIfReportsAreValidWithAcceptedStatuses(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy,
			[]string{globalparameters.TestCasePassed, globalparameters.TestCaseFailed,
				globalparameters.TestCaseSkipped}, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify CheckDetails for process disappearance handling")
		checkDetails, err := globalhelper.GetTestCaseCheckDetails(
			tsparams.CertsuiteSharedCPUPoolSchedulingPolicy, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		totalObjects := len(checkDetails.CompliantObjectsOut) + len(checkDetails.NonCompliantObjectsOut)
		GinkgoWriter.Printf("CheckDetails: %d compliant, %d non-compliant objects\n",
			len(checkDetails.CompliantObjectsOut), len(checkDetails.NonCompliantObjectsOut))
		Expect(totalObjects).To(BeNumerically(">=", 1),
			"At least one report object expected (no early abort)")

		// Verify that any "process disappeared" entries are in the compliant list, never non-compliant
		disappearedCount := 0

		for _, obj := range checkDetails.CompliantObjectsOut {
			reason := globalhelper.GetReportObjectFieldValue(obj, "Reason")
			if strings.Contains(reason, "process disappeared") {
				disappearedCount++
			}
		}

		for _, obj := range checkDetails.NonCompliantObjectsOut {
			reason := globalhelper.GetReportObjectFieldValue(obj, "Reason")
			Expect(reason).ToNot(ContainSubstring("process disappeared"),
				"Process disappeared entries must not appear in non-compliant list")
		}

		GinkgoWriter.Printf("Process disappeared entries in compliant list: %d\n", disappearedCount)
	})

	It("One pod with container running in exclusive cpu pool", func() {
		By("Verify cluster is configured for exclusive CPU pools")
		configured, reason, err := globalhelper.IsClusterConfiguredForExclusiveCPUs()
		Expect(err).ToNot(HaveOccurred())

		if !configured {
			Skip("Cluster not configured for exclusive CPU pools: " + reason)
		}

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		pod.RedefineWithCPUResources(testPod, "1", "1")
		pod.RedefineWithMemoryResources(testPod, "512Mi", "512Mi")

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is running with exclusive CPU pool configuration")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")
		Expect(len(runningPod.Spec.Containers)).To(BeNumerically(">", 0), "Pod should have containers")

		By("Assert pod is not using HostPID")
		Expect(runningPod.Spec.HostPID).To(BeFalse(),
			"Pod must not use HostPID — certsuite skips pods with HostPID")

		By("Assert pod has Guaranteed QoS with whole-unit exclusive CPUs")
		Expect(runningPod.Status.QOSClass).To(Equal(corev1.PodQOSGuaranteed),
			"Pod must be Guaranteed QoS for exclusive CPU pool")

		container := runningPod.Spec.Containers[0]
		cpuRequest := container.Resources.Requests.Cpu().MilliValue()
		cpuLimit := container.Resources.Limits.Cpu().MilliValue()
		memRequest := container.Resources.Requests.Memory().Value()
		memLimit := container.Resources.Limits.Memory().Value()

		GinkgoWriter.Printf("CPU request: %dm, CPU limit: %dm\n", cpuRequest, cpuLimit)
		GinkgoWriter.Printf("Memory request: %d, Memory limit: %d\n", memRequest, memLimit)
		Expect(cpuRequest).To(Equal(cpuLimit), "CPU request must equal limit for exclusive pool")
		Expect(cpuRequest%1000).To(Equal(int64(0)), "CPU must be a whole unit for exclusive pool")
		Expect(memRequest).To(Equal(memLimit), "Memory request must equal limit for Guaranteed QoS")

		By("Assert all containers are ready")

		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Verify container has pinned CPUs (not shared pool)")
		// On cgroup v2 (OCP 4.x default), the effective cpuset is at /sys/fs/cgroup/cpuset.cpus.effective.
		// On cgroup v1, it's at /sys/fs/cgroup/cpuset/cpuset.cpus.
		// A pinned container gets specific CPUs (e.g. "2-3"), while the shared pool gets the full range.
		cpusetOutput, cpusetErr := globalhelper.ExecCommand(*runningPod, []string{"cat",
			"/sys/fs/cgroup/cpuset.cpus.effective"})

		if cpusetErr != nil {
			// Fall back to cgroup v1 path
			cpusetOutput, cpusetErr = globalhelper.ExecCommand(*runningPod, []string{"cat",
				"/sys/fs/cgroup/cpuset/cpuset.cpus"})
		}

		if cpusetErr == nil {
			cpuset := strings.TrimSpace(cpusetOutput.String())
			GinkgoWriter.Printf("Container cpuset: %s\n", cpuset)

			pinnedCPUCount := countCPUsInSet(cpuset)
			GinkgoWriter.Printf("Pinned CPU count: %d, requested: %d\n", pinnedCPUCount, cpuRequest/1000)
			Expect(int64(pinnedCPUCount)).To(Equal(cpuRequest/1000),
				"Pinned CPU count should match the number of requested exclusive CPUs")
		} else {
			GinkgoWriter.Printf("Could not read cpuset (skipping pinning verification): %v\n", cpusetErr)
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

// countCPUsInSet parses a cpuset string like "0-3,5,7-9" and returns the total CPU count.
func countCPUsInSet(cpuset string) int {
	count := 0

	for _, part := range strings.Split(cpuset, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		bounds := strings.SplitN(part, "-", 2)
		if len(bounds) == 1 {
			count++
		} else {
			low, errLow := strconv.Atoi(bounds[0])
			high, errHigh := strconv.Atoi(bounds[1])

			if errLow != nil || errHigh != nil {
				continue
			}

			count += high - low + 1
		}
	}

	return count
}

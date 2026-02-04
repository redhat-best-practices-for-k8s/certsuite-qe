package tests

import (
	"context"
	"fmt"
	"strings"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("platform-alteration-boot-params", Label("platformalteration1", "ocp-required"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.PlatformAlterationNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("If Kind cluster, skip")
		if globalhelper.IsKindCluster() {
			Skip("Kind cluster does not support MCO")
		}

		By("Verify MCO is healthy and accessible")
		mcoHealthy, err := globalhelper.IsMCOHealthy()
		if err != nil || !mcoHealthy {
			Skip("MCO is not healthy or accessible on this cluster - skipping boot params tests")
		}

		By("Verify MachineConfigPools exist")
		mcpList, err := globalhelper.GetAPIClient().MachineConfigPools().List(context.TODO(), metav1.ListOptions{})
		if err != nil || len(mcpList.Items) == 0 {
			Skip("No MachineConfigPools found - skipping boot params tests")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 51302
	It("unchanged boot params", func() {
		By("Create daemonSet")
		testDaemonSet := daemonset.DefineDaemonSet(randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(testDaemonSet)
		daemonset.RedefineWithVolumeMount(testDaemonSet)

		By("Create and wait until daemonSet is ready")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(testDaemonSet, tsparams.WaitingTime)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "not schedulable") ||
				strings.Contains(errMsg, "Timed out") ||
				strings.Contains(errMsg, "not running") ||
				strings.Contains(errMsg, "not ready") {
				Skip("This test cannot run because the daemonSet is not ready: " + errMsg)
			}
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert daemonSet has ready pods on nodes")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(testDaemonSet)
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("DaemonSet status: NumberReady=%d, DesiredNumberScheduled=%d, CurrentNumberScheduled=%d\n",
			runningDaemonSet.Status.NumberReady,
			runningDaemonSet.Status.DesiredNumberScheduled,
			runningDaemonSet.Status.CurrentNumberScheduled)
		Expect(runningDaemonSet.Status.NumberReady).To(BeNumerically(">", 0), "DaemonSet should have ready pods")
		Expect(runningDaemonSet.Status.NumberReady).To(Equal(runningDaemonSet.Status.DesiredNumberScheduled),
			"All scheduled pods should be ready")

		By("Verify pods have host volume access")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">", 0), "Expected at least one pod")

		// Log pod and node details for debugging
		GinkgoWriter.Printf("Found %d pods in namespace %s\n", len(podsList.Items), randomNamespace)
		for i, pod := range podsList.Items {
			GinkgoWriter.Printf("Pod[%d] name: %s, phase: %s, node: %s\n",
				i, pod.Name, pod.Status.Phase, pod.Spec.NodeName)
			for j, container := range pod.Spec.Containers {
				GinkgoWriter.Printf("  Container[%d] name: %s, image: %s\n",
					j, container.Name, container.Image)
			}
			// Log volume mounts
			for _, vm := range pod.Spec.Containers[0].VolumeMounts {
				GinkgoWriter.Printf("  VolumeMount: %s -> %s\n", vm.Name, vm.MountPath)
			}
		}

		By("Assert pods are running with ready containers")
		for _, pod := range podsList.Items {
			Expect(pod.Status.Phase).To(Equal(corev1.PodRunning), fmt.Sprintf("Pod %s should be running", pod.Name))
			for _, cs := range pod.Status.ContainerStatuses {
				GinkgoWriter.Printf("Container %s in pod %s: ready=%v\n", cs.Name, pod.Name, cs.Ready)
				Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s in pod %s should be ready", cs.Name, pod.Name))
			}
		}

		By("Verify pod can access host filesystem")
		cmdOutput, err := globalhelper.ExecCommand(podsList.Items[0], []string{"cat", "/host/proc/cmdline"})
		if err != nil {
			GinkgoWriter.Printf("Failed to access host filesystem: %v\n", err)
			Skip("Cannot access host filesystem from pod - skipping boot params test")
		}
		kernelCmdline := cmdOutput.String()
		GinkgoWriter.Printf("Host kernel cmdline: %s\n", kernelCmdline)

		By("Start platform-alteration-boot-params test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteBootParamsName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		// This test does not alter boot params, so the certsuite test should pass.
		// The certsuite platform-alteration-boot-params test checks if kernel cmdline
		// matches what's configured in MachineConfig. If they match (as expected when
		// we don't alter anything), the test passes.
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteBootParamsName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51305
	It("change boot params using MCO", func() {
		machineConfigList, err := globalhelper.GetAPIClient().MachineConfigs().List(context.TODO(), metav1.ListOptions{})
		Expect(err).ToNot(HaveOccurred())

		machineConfigPoolList, err := globalhelper.GetAPIClient().MachineConfigPools().List(context.TODO(),
			metav1.ListOptions{})
		Expect(err).ToNot(HaveOccurred())

		// Log available MCPs for debugging
		GinkgoWriter.Printf("Found %d MachineConfigPools\n", len(machineConfigPoolList.Items))
		for i, mcp := range machineConfigPoolList.Items {
			GinkgoWriter.Printf("MCP[%d] name: %s, config: %s\n", i, mcp.Name, mcp.Spec.Configuration.Name)
		}

		foundWorkerCNF := false
		for _, machineConfig := range machineConfigList.Items {
			for _, mcp := range machineConfigPoolList.Items {
				if machineConfig.Name == mcp.Spec.Configuration.Name && mcp.Name == "worker-cnf" {
					foundWorkerCNF = true
					GinkgoWriter.Printf("Found worker-cnf MCP with machineConfig: %s\n", machineConfig.Name)
					machineConfig.Spec.KernelArguments = []string{"skew_tick=1", "nohz=off"}

					By("Update the current machineConfig")
					_, err := globalhelper.GetAPIClient().MachineConfigs().Update(context.TODO(), &machineConfig, metav1.UpdateOptions{})
					Expect(err).ToNot(HaveOccurred())

					By("Assert machineConfig has been updated")
					updatedMachineConfig, err := globalhelper.GetAPIClient().MachineConfigs().Get(context.TODO(),
						machineConfig.Name, metav1.GetOptions{})
					Expect(err).ToNot(HaveOccurred())
					Expect(updatedMachineConfig.Spec.KernelArguments).To(Equal([]string{"skew_tick=1", "nohz=off"}))
					GinkgoWriter.Printf("Updated machineConfig %s with kernel arguments: %v\n",
						machineConfig.Name, updatedMachineConfig.Spec.KernelArguments)
				}
			}
		}

		if !foundWorkerCNF {
			GinkgoWriter.Printf("No worker-cnf MachineConfigPool found - test will run without modifying boot params\n")
		}

		By("Start platform-alteration-boot-params test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteBootParamsName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteBootParamsName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

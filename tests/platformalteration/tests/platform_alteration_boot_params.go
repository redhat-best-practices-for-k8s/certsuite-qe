package tests

import (
	"context"
	"fmt"

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
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		By("Create and wait until daemonSet is ready")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert daemonSet has ready pods on nodes")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDaemonSet.Status.NumberReady).To(BeNumerically(">", 0), "DaemonSet should have ready pods")
		Expect(runningDaemonSet.Status.NumberReady).To(Equal(runningDaemonSet.Status.DesiredNumberScheduled),
			"All scheduled pods should be ready")

		By("Verify pods have host volume access")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">", 0), "Expected at least one pod")

		By("Assert pods are running with ready containers")
		for _, pod := range podsList.Items {
			Expect(pod.Status.Phase).To(Equal(corev1.PodRunning), fmt.Sprintf("Pod %s should be running", pod.Name))
			for _, cs := range pod.Status.ContainerStatuses {
				Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s in pod %s should be ready", cs.Name, pod.Name))
			}
		}

		By("Verify pod can access host filesystem")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"cat", "/host/proc/cmdline"})
		if err != nil {
			Skip("Cannot access host filesystem from pod - skipping boot params test")
		}

		By("Start platform-alteration-boot-params test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteBootParamsName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

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

		for _, machineConfig := range machineConfigList.Items {
			for _, mcp := range machineConfigPoolList.Items {
				if machineConfig.Name == mcp.Spec.Configuration.Name && mcp.Name == "worker-cnf" {
					machineConfig.Spec.KernelArguments = []string{"skew_tick=1", "nohz=off"}

					By("Update the current machineConfig")
					_, err := globalhelper.GetAPIClient().MachineConfigs().Update(context.TODO(), &machineConfig, metav1.UpdateOptions{})
					Expect(err).ToNot(HaveOccurred())

					By("Assert machineConfig has been updated")
					updatedMachineConfig, err := globalhelper.GetAPIClient().MachineConfigs().Get(context.TODO(),
						machineConfig.Name, metav1.GetOptions{})
					Expect(err).ToNot(HaveOccurred())
					Expect(updatedMachineConfig.Spec.KernelArguments).To(Equal([]string{"skew_tick=1", "nohz=off"}))
				}
			}
		}

		By("Start platform-alteration-boot-params test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteBootParamsName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteBootParamsName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

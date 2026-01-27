package tests

import (
	"context"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
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

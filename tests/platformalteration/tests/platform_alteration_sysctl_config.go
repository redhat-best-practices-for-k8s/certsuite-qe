package tests

import (
	"context"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("platform-alteration-sysctl-config", func() {
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
	It("unchanged sysctl config", func() {

		By("Create daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		By("Create and wait until daemonSet is ready")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert sysctl config is unchanged")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDaemonSet.Spec.Template.Spec.Containers[0].SecurityContext.Privileged).To(BeTrue())
		Expect(runningDaemonSet.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal("/host"))
		Expect(runningDaemonSet.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name).To(Equal("host"))

		By("Start platform-alteration-sysctl-config test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteSysctlConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteSysctlConfigName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51332
	It("change sysctl config using MCO", func() {
		Skip("This test is unstable and needs to be fixed")

		By("Create daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(podList.Items)).NotTo(BeZero())

		node, err := globalhelper.GetAPIClient().Nodes().Get(context.TODO(), podList.Items[0].Spec.NodeName, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())

		value, exists := node.Annotations["machineGetConfiguration().openshift.io/currentConfig"]
		if !exists {
			Fail("didn't get node's machine config")
		}

		mcObj, err := globalhelper.GetAPIClient().MachineConfigs().Get(context.TODO(), value, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())

		mcKernelArgs := mcObj.Spec.KernelArguments
		mcKernelArgsMap := tshelper.ArgListToMap(mcKernelArgs)

		value, exists = mcKernelArgsMap["net.ipv4.ip_forward"]
		if !exists {
			mcKernelArgs = append(mcKernelArgs, "net.ipv4.ip_forward", "0")
		} else {
			if value == "0" {
				mcKernelArgs = []string{"net.ipv4.ip_forward", "1"}
			} else {
				mcKernelArgs = []string{"net.ipv4.ip_forward", "0"}

			}
		}
		mcObj.Spec.KernelArguments = mcKernelArgs

		_, err = globalhelper.GetAPIClient().MachineConfigs().Update(context.TODO(), mcObj, metav1.UpdateOptions{})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-sysctl-config test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteSysctlConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteSysctlConfigName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

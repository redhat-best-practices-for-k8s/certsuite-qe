package tests

import (
	"context"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("platform-alteration-sysctl-config", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51302
	It("unchanged sysctl config", func() {

		By("Create daemonSet")
		daemonSet := daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-sysctl-config test")
		err = globalhelper.LaunchTests(tsparams.TnfSysctlConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfSysctlConfigName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51332
	It("change sysctl config using MCO", func() {

		By("Create daemonSet")
		daemonSet := daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(tsparams.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		node, err := globalhelper.APIClient.Nodes().Get(context.Background(), podList.Items[0].Spec.NodeName, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())

		value, exists := node.Annotations["machineconfiguration.openshift.io/currentConfig"]
		if !exists {
			Fail("didn't get node's machine config")
		}

		mcObj, err := globalhelper.APIClient.MachineConfigs().Get(context.Background(), value, metav1.GetOptions{})
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

		_, err = globalhelper.APIClient.MachineConfigs().Update(context.TODO(), mcObj, metav1.UpdateOptions{})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-sysctl-config test")
		err = globalhelper.LaunchTests(tsparams.TnfSysctlConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfSysctlConfigName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

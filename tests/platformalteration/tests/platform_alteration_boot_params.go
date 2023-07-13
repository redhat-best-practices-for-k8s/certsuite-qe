package tests

import (
	"context"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("platform-alteration-boot-params", func() {

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
	It("unchanged boot params", func() {
		By("Create daemonSet")
		daemonSet := daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-boot-params test")
		err = globalhelper.LaunchTests(tsparams.TnfBootParamsName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBootParamsName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51305
	It("change boot params using MCO", func() {
		machineConfigList, err := globalhelper.APIClient.MachineConfigs().List(context.Background(), metav1.ListOptions{})
		Expect(err).ToNot(HaveOccurred())

		machineConfigPoolList, err := globalhelper.APIClient.MachineConfigPools().List(context.Background(),
			metav1.ListOptions{})
		Expect(err).ToNot(HaveOccurred())

		for _, machineConfig := range machineConfigList.Items {
			for _, mcp := range machineConfigPoolList.Items {
				if machineConfig.Name == mcp.Spec.Configuration.Name && mcp.Name == "worker-cnf" {
					machineConfig.Spec.KernelArguments = []string{"skew_tick=1", "nohz=off"}

					By("Update the current machineConfig")
					_, err := globalhelper.APIClient.MachineConfigs().Update(context.TODO(), &machineConfig, metav1.UpdateOptions{})
					Expect(err).ToNot(HaveOccurred())
				}
			}
		}

		By("Start platform-alteration-boot-params test")
		err = globalhelper.LaunchTests(tsparams.TnfBootParamsName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBootParamsName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

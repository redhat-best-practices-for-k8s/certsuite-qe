package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("platform-alteration-tainted-node-kernel", Serial, func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 51389
	It("Untainted node", func() {

		// all nodes suppose to be untainted when the cluster is deployed.
		By("Start platform-alteration-tainted-node-kernel test")
		err := globalhelper.LaunchTests(tsparams.TnfTaintedNodeKernelName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfTaintedNodeKernelName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51390
	It("Tainted node [negative]", func() {

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(tsparams.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		if len(podList.Items) == 0 {
			Skip("no pods have been found in namespace")
		}

		// we can only set a taint flag in this way, not remove it, there is no way to untaint a running kernel,
		// the taint flag will be removed once the node is rebooted.
		By("Taint a node")
		_, err = globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", "echo 32 > /proc/sys/kernel/tainted"})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-tainted-node-kernel test")
		err = globalhelper.LaunchTests(tsparams.TnfTaintedNodeKernelName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfTaintedNodeKernelName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

		By("Reboot the node to remove the taint")
		_, err = globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", tsparams.Reboot})
		Expect(err).ToNot(HaveOccurred())

		By("Wait for the node to become not ready")
		err = tshelper.WaitForSpecificNodeCondition(globalhelper.GetAPIClient(),
			tsparams.RebootWaitingTime, tsparams.RetryInterval, podList.Items[0].Spec.NodeName, false)
		Expect(err).ToNot(HaveOccurred())

		By("Wait for the node to become ready")
		err = tshelper.WaitForSpecificNodeCondition(globalhelper.GetAPIClient(),
			tsparams.RebootWaitingTime, tsparams.RetryInterval, podList.Items[0].Spec.NodeName, true)
		Expect(err).ToNot(HaveOccurred())

	})
})

package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
)

/*
	The setup needs to be as mentioned here :
	https://developers.redhat.com/articles/2021/08/27/using-virtual-functions-dpdk-red-hat-openshift
*/

var _ = Describe("Networking dpdk-cpu-pinning-exec-probe,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNetworkingNameSpace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.IsKindCluster() {
			Skip("DPDK is not supported on Kind cluster. Skipping.")
		}

		// Check if 'openshift-sriov-network-operator' namespace exists, if not, skip
		By("Check if openshift-sriov-network-operator is installed")
		exists, err := globalhelper.NamespaceExists("openshift-sriov-network-operator")
		Expect(err).ToNot(HaveOccurred())
		if !exists {
			Skip("openshift-sriov-network-operator is not installed, skipping test")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	It("one dpdk pod with no probe", func() {
		By("Deploy dpdk pod namespace")
		dpdkPod := tshelper.DefineDpdkPod(tsparams.DpdkPodName, randomNamespace)
		err := globalhelper.CreateAndWaitUntilPodIsReady(dpdkPod, tsparams.WaitingTime)
		if err != nil {
			Skip("DPDK is not deployed. There may be some problems in setup. Hence, skipping.")
		}

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDpdkCPUPinningExecProbe,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDpdkCPUPinningExecProbe,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one dpdk pod with exec probe [negative]", func() {
		By("Deploy dpdk pod namespace")
		dpdkPod := tshelper.DefineDpdkPod(tsparams.DpdkPodName, randomNamespace)

		By("Redefine liveness probe")
		pod.RedefinePodContainerWithLivenessProbeCommand(dpdkPod, 0, []string{"cat", "/tmp/healthy"})

		err := globalhelper.CreateAndWaitUntilPodIsReady(dpdkPod, tsparams.WaitingTime)
		if err != nil {
			Skip("DPDK is not deployed. There may be some problems in setup. Hence, skipping.")
		}

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDpdkCPUPinningExecProbe,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDpdkCPUPinningExecProbe,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

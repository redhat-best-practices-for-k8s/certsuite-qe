package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
)

/*
	The setup needs to be as mentioned here:
	https://developers.redhat.com/articles/2021/08/27/using-virtual-functions-dpdk-red-hat-openshift
*/

var _ = Describe("performance-cpu-pinning-no-exec-probes", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

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
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	It("one dpdk pod with no probe", func() {
		By("Deploy dpdk pod")
		dpdkPod := tshelper.DefineDpdkPod(tsparams.DpdkPodName, randomNamespace)
		err := globalhelper.CreateAndWaitUntilPodIsReady(dpdkPod, tsparams.WaitingTime)
		if err != nil {
			Skip("DPDK is not deployed. There may be some problems in setup. Hence, skipping.")
		}

		By("Start cpu-pinning-no-exec-probes test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteCPUPinningNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteCPUPinningNoExecProbes,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one dpdk pod with exec probe [negative]", func() {
		By("Deploy dpdk pod")
		dpdkPod := tshelper.DefineDpdkPod(tsparams.DpdkPodName, randomNamespace)

		By("Redefine liveness probe")
		pod.RedefinePodContainerWithLivenessProbeCommand(dpdkPod, 0, []string{"cat", "/tmp/healthy"})

		err := globalhelper.CreateAndWaitUntilPodIsReady(dpdkPod, tsparams.WaitingTime)
		if err != nil {
			Skip("DPDK is not deployed. There may be some problems in setup. Hence, skipping.")
		}

		By("Start cpu-pinning-no-exec-probes test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteCPUPinningNoExecProbes,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteCPUPinningNoExecProbes,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

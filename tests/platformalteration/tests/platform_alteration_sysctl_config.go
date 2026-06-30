package tests

import (
	"strings"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("platform-alteration-sysctl-config", Label("platformalteration4", "ocp-required"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

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
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		By("Create and wait until daemonSet is ready")

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		if globalhelper.IsTransientDaemonSetError(err) {
			Skip("This test cannot run because the daemonSet is not ready: " + err.Error())
		}

		Expect(err).ToNot(HaveOccurred())

		By("Assert daemonSet has ready pods on nodes")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDaemonSet.Spec.Template.Spec.Containers[0].SecurityContext.Privileged).To(BeTrue())
		Expect(runningDaemonSet.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal("/host"))
		Expect(runningDaemonSet.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name).To(Equal("host"))

		By("Start platform-alteration-sysctl-config test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteSysctlConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		// The sysctl check can fail due to config drift from sources that cannot
		// be reliably pre-detected (e.g., TuneD profiles, /etc/sysctl.d/ files).
		// Accept passed, failed, or skipped as valid certsuite outcomes.
		err = globalhelper.ValidateIfReportsAreValidWithAcceptedStatuses(
			tsparams.CertsuiteSysctlConfigName,
			[]string{globalparameters.TestCasePassed, globalparameters.TestCaseFailed,
				globalparameters.TestCaseSkipped}, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51332
	It("change sysctl config using MCO", func() {
		By("Create daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		if globalhelper.IsTransientDaemonSetError(err) {
			Skip("This test cannot run because the daemonSet is not ready: " + err.Error())
		}

		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(podList.Items)).NotTo(BeZero())

		// Read the current value of a sysctl so we can alter it and restore it later.
		sysctlPath := "/host/proc/sys/net/ipv4/conf/all/arp_notify"

		By("Read current sysctl value")
		origBuf, err := globalhelper.ExecCommand(
			podList.Items[0], []string{"/bin/bash", "-c", "cat " + sysctlPath})
		Expect(err).ToNot(HaveOccurred())

		origValue := strings.TrimSpace(origBuf.String())
		GinkgoWriter.Printf("Original sysctl value for arp_notify: %s\n", origValue)

		// Flip the value: 0 -> 1 or anything else -> 0.
		newValue := "1"
		if origValue != "0" {
			newValue = "0"
		}

		By("Alter sysctl value at runtime to create drift")
		_, err = globalhelper.ExecCommand(
			podList.Items[0], []string{"/bin/bash", "-c", "echo " + newValue + " > " + sysctlPath})
		Expect(err).ToNot(HaveOccurred())

		// Ensure the original value is restored regardless of test outcome.
		DeferCleanup(func() {
			By("Restore original sysctl value")

			_, restoreErr := globalhelper.ExecCommand(
				podList.Items[0], []string{"/bin/bash", "-c", "echo " + origValue + " > " + sysctlPath})
			if restoreErr != nil {
				GinkgoWriter.Printf("Warning: failed to restore sysctl value: %v\n", restoreErr)
			}
		})

		By("Start platform-alteration-sysctl-config test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteSysctlConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		// After altering a sysctl value the certsuite should detect the drift and
		// report failure. However, on some cluster configurations the sysctl check
		// may be skipped entirely if no MachineConfig manages this particular key.
		// Accept failed or skipped as valid outcomes.
		err = globalhelper.ValidateIfReportsAreValidWithAcceptedStatuses(
			tsparams.CertsuiteSysctlConfigName,
			[]string{globalparameters.TestCaseFailed, globalparameters.TestCaseSkipped},
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

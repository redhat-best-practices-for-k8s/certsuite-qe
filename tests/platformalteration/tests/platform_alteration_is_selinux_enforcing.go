package tests

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"
)

var _ = Describe("platform-alteration-is-selinux-enforcing", Label("platformalteration3", "ocp-required"), func() {
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
			Skip("Kind cluster does not support SELinux")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 51310
	It("SELinux is enforcing on all nodes", func() {
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		By("Create and wait until daemonSet is ready")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Verify that all nodes are running with selinux on enforcing mode")
		for _, pod := range podList.Items {

			buf, err := globalhelper.ExecCommand(pod, []string{"/bin/bash", "-c", tsparams.Getenforce})
			Expect(err).ToNot(HaveOccurred())

			if !strings.Contains(buf.String(), tsparams.Enforcing) {
				_, err = globalhelper.ExecCommand(pod, []string{"/bin/bash", "-c", tsparams.SetEnforce})
				Expect(err).ToNot(HaveOccurred())
			}
		}

		By("Start platform-alteration-is-selinux-enforcing test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteIsSelinuxEnforcingName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteIsSelinuxEnforcingName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51311
	It("SELinux is permissive on one node [negative]", func() {
		if globalhelper.IsKindCluster() {
			Skip("Kind cluster does not support SELinux")
		}

		Skip("Skipping. Remove this skip when we can detect if SELinux is enabled on the node")

		daemonSet := daemonset.DefineDaemonSet(randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		By("Create and wait until daemonSet is ready")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(podList.Items)).NotTo(BeZero())

		By("Set SELinux permissive on the node")
		_, err = globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", tsparams.SetPermissive})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-is-selinux-enforcing test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteIsSelinuxEnforcingName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteIsSelinuxEnforcingName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verifying SELinux is enforcing on the node")
		_, err = globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", tsparams.SetEnforce})
		Expect(err).ToNot(HaveOccurred())

	})
})

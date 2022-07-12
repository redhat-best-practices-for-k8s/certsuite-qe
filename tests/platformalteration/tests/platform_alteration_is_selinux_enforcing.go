package tests

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("platform-alteration-is-selinux-enforcing", func() {

	execute.BeforeAll(func() {
		By("Make masters schedulable")
		err := globalhelper.EnableMasterScheduling(true)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	const (
		getenforce    = `chroot /host getenforce`
		enforcing     = "Enforcing"
		setPermissive = `chroot /host setenforce 0`
		setEnforce    = `chroot /host setenforce 1`
	)

	// 51310
	It("SELinux is enforcing on all nodes", func() {
		daemonset := daemonset.RedefineWithPriviledgedContainer(
			daemonset.RedefineWithVolumeMount(
				daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
					tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)))

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonset, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(tsparams.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Verify that all nodes are running with selinux on enforcing mode")
		for _, pod := range podList.Items {

			buf, err := globalhelper.ExecCommand(pod, []string{"/bin/bash", "-c", getenforce})
			Expect(err).ToNot(HaveOccurred())

			if !strings.Contains(buf.String(), enforcing) {
				_, err = globalhelper.ExecCommand(pod, []string{"/bin/bash", "-c", setEnforce})
				Expect(err).ToNot(HaveOccurred())
			}
		}

		By("Start platform-alteration-is-selinux-enforcing test")
		err = globalhelper.LaunchTests(tsparams.TnfIsSelinuxEnforcingName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfIsSelinuxEnforcingName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51311
	It("SELinux is permissive on one node [negative]", func() {
		daemonset := daemonset.RedefineWithPriviledgedContainer(
			daemonset.RedefineWithVolumeMount(
				daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
					tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)))

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonset, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(tsparams.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Set SELinux permissive on the node")
		_, err = globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", setPermissive})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-is-selinux-enforcing test")
		err = globalhelper.LaunchTests(tsparams.TnfIsSelinuxEnforcingName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfIsSelinuxEnforcingName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

		By("Verifying SELinux is enforcing on the node")
		_, err = globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", setEnforce})
		Expect(err).ToNot(HaveOccurred())

	})
})

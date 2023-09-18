package tests

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/crd"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
)

var _ = Describe("platform-alteration-hugepages-config", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.PlatformAlterationNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.IsKindCluster() {
			Skip("Hugepages are not supported in Kind clusters")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	// 51308
	It("unchanged configuration", func() {

		crdExists, err := crd.EnsureCrdExists(tsparams.PerformanceProfileCrd)
		Expect(err).ToNot(HaveOccurred())

		if !crdExists {
			Skip("performance profile does not exist.")
		}

		// cluster should be set with kernel hugepages = MC hugepages configuration by performance profile.
		By("Start platform-alteration-hugepages-config test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePagesConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfHugePagesConfigName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51309
	It("Change Hugepages config manually [negative]", func() {
		if globalhelper.IsKindCluster() {
			Skip("Kind cluster does not support hugepages")
		}

		By("Set rbac policy which allows authenticated users to run privileged containers")
		err := globalhelper.AllowAuthenticatedUsersRunPrivilegedContainers()
		Expect(err).ToNot(HaveOccurred())

		By("Create daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(podList.Items)).NotTo(BeZero())

		By("Get first hugepages file")
		nrHugepagesFiles, err := globalhelper.ExecCommand(
			podList.Items[0], []string{"/bin/bash", "-c", tsparams.FindHugePagesFiles})
		Expect(err).ToNot(HaveOccurred())

		hugePagesPaths := strings.Split(nrHugepagesFiles.String(), "\r\n")
		if len(hugePagesPaths) == 0 {
			Fail(fmt.Sprintf("No hugepages files have been found on node - %s ", podList.Items[0].Spec.NodeName))
		}

		By("Get hugepages config")
		currentHugepagesNumber, err := tshelper.GetHugePagesConfigNumber(hugePagesPaths[0], &podList.Items[0])
		Expect(err).ToNot(HaveOccurred())

		updatedHugePagesNumber := currentHugepagesNumber + 1

		By("Manually update hugepages config")
		err = tshelper.UpdateAndVerifyHugePagesConfig(updatedHugePagesNumber, hugePagesPaths[0], &podList.Items[0])
		Expect(err).ToNot(HaveOccurred(), "failed to update and verify hugepages file: %s, %v ", hugePagesPaths[0], err)

		By("Start platform-alteration-hugepages-config test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePagesConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePagesConfigName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

		By("Fix hugepages config")
		updatedHugePagesNumber = currentHugepagesNumber
		err = tshelper.UpdateAndVerifyHugePagesConfig(updatedHugePagesNumber, hugePagesPaths[0], &podList.Items[0])
		Expect(err).ToNot(HaveOccurred(), "failed to update and verify hugepages file: %s, %v ", hugePagesPaths[0], err)

	})

})

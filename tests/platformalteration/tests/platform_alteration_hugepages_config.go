package tests

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/crd"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"
)

var _ = Describe("platform-alteration-hugepages-config", Serial, Label("platformalteration3", "ocp-required"), func() {
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

		if globalhelper.IsKindCluster() {
			Skip("Hugepages are not supported in Kind clusters")
		}

		By("Verify MCO is healthy and accessible")

		mcoHealthy, err := globalhelper.IsMCOHealthy()
		if err != nil || !mcoHealthy {
			Skip("MCO is not healthy or accessible on this cluster - skipping hugepages config tests")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
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
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePagesConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteHugePagesConfigName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51309
	It("Change Hugepages config manually [negative]", func() {
		crdExists, err := crd.EnsureCrdExists(tsparams.PerformanceProfileCrd)
		Expect(err).ToNot(HaveOccurred())

		if !crdExists {
			Skip("performance profile does not exist.")
		}

		By("Set rbac policy which allows authenticated users to run privileged containers")
		err = globalhelper.AllowAuthenticatedUsersRunPrivilegedContainers()
		Expect(err).ToNot(HaveOccurred())

		By("Create daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)
		daemonset.RedefineWithPrivilegedContainer(daemonSet)
		daemonset.RedefineWithVolumeMount(daemonSet)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		if globalhelper.IsTransientDaemonSetError(err) {
			Skip("This test cannot run because the daemonSet is not ready: " + err.Error())
		}

		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(podList.Items)).NotTo(BeZero())

		By("Get first hugepages file")
		nrHugepagesFiles, err := globalhelper.ExecCommand(
			podList.Items[0], []string{"/bin/bash", "-c", tsparams.FindHugePagesFiles})
		Expect(err).ToNot(HaveOccurred())

		hugePagesPaths := strings.Fields(nrHugepagesFiles.String())

		if len(hugePagesPaths) == 0 {
			Skip(fmt.Sprintf("No hugepages files found on node %s - hugepages may not be configured",
				podList.Items[0].Spec.NodeName))
		}

		GinkgoWriter.Printf("Found %d hugepages files, using: %s\n", len(hugePagesPaths), hugePagesPaths[0])

		By("Get hugepages config")
		currentHugepagesNumber, err := tshelper.GetHugePagesConfigNumber(hugePagesPaths[0], &podList.Items[0])
		Expect(err).ToNot(HaveOccurred())

		GinkgoWriter.Printf("Current hugepages value: %d\n", currentHugepagesNumber)

		updatedHugePagesNumber := currentHugepagesNumber + 1

		By("Manually update hugepages config")
		err = tshelper.UpdateAndVerifyHugePagesConfig(updatedHugePagesNumber, hugePagesPaths[0], &podList.Items[0])
		Expect(err).ToNot(HaveOccurred(), "failed to update and verify hugepages file: %s, %v ", hugePagesPaths[0], err)

		// Ensure the original hugepages value is restored regardless of test outcome.
		DeferCleanup(func() {
			By("Restore original hugepages config")

			restoreErr := tshelper.UpdateAndVerifyHugePagesConfig(
				currentHugepagesNumber, hugePagesPaths[0], &podList.Items[0])
			if restoreErr != nil {
				GinkgoWriter.Printf("Warning: failed to restore hugepages config: %v\n", restoreErr)
			}
		})

		By("Start platform-alteration-hugepages-config test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePagesConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteHugePagesConfigName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

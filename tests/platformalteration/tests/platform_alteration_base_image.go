package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
)

var _ = Describe("platform-alteration-base-image", func() {
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
			// The TNF suite actually proactively skips this test if the cluster is Non-OCP.
			Skip(fmt.Sprintf("%s test is not applicable for Kind cluster", tsparams.TnfBaseImageName))
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	// 51297
	It("One deployment, one pod, running test image", func() {
		By("Define deployment")
		deployment := deployment.DefineDeployment(tsparams.TestDeploymentName,
			randomNamespace,
			globalhelper.GetConfiguration().General.TestImage,
			tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBaseImageName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51298
	It("One daemonSet, running test image", func() {
		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace,
			globalhelper.GetConfiguration().General.TestImage,
			tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBaseImageName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51299
	It("Two deployments, one pod each, change container base image by creating a file [negative]", func() {
		By("Define first deployment")
		deploymenta := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithPrivilegedContainer(deploymenta)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(podsList.Items)).NotTo(BeZero())

		By("Change container base image")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "touch /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment")
		deploymentb := deployment.DefineDeployment("platform-alteration-dpb",
			randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBaseImageName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One statefulSet, one pod, change container base image by creating a file [negative]", func() {
		By("Define statefulSet")
		statefulSet := statefulset.DefineStatefulSet(tsparams.TestStatefulSetName,
			randomNamespace,
			globalhelper.GetConfiguration().General.TestImage,
			tsparams.TnfTargetPodLabels)
		statefulset.RedefineWithPrivilegedContainer(statefulSet)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tshelper.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(podsList.Items)).NotTo(BeZero())

		By("Change container base image")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "touch /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBaseImageName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

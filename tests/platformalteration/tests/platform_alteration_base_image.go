package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/platformalterationhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/platformalterationparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("platform-alteration-base-image", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(platformalterationparameters.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51297
	It("One deployment, one pod, running test image", func() {

		By("Define deployment")
		deployment := deployment.DefineDeployment(platformalterationparameters.TestDeploymentName,
			platformalterationparameters.PlatformAlterationNamespace,
			globalhelper.Configuration.General.TestImage,
			platformalterationparameters.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, platformalterationparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			platformalterationparameters.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			platformalterationparameters.TnfBaseImageName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51298
	It("One daemonSet, running test image", func() {

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(platformalterationparameters.PlatformAlterationNamespace,
			globalhelper.Configuration.General.TestImage,
			platformalterationparameters.TnfTargetPodLabels, platformalterationparameters.TestDaemonSetName)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, platformalterationparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			platformalterationparameters.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			platformalterationparameters.TnfBaseImageName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51299
	It("Two deployments, one pod each, change container base image by creating a file [negative]", func() {

		By("Define first deployment")
		deploymenta := platformalterationhelper.DefineDeploymentWithPriviledgedContainer(
			platformalterationparameters.TestDeploymentName,
			platformalterationparameters.PlatformAlterationNamespace,
			globalhelper.Configuration.General.TestImage, platformalterationparameters.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, platformalterationparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podsList, err := globalhelper.GetListOfPodsInNamespace(platformalterationparameters.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Change container base image")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "touch /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment")
		deploymentb := deployment.DefineDeployment("platform-alteration-dpb",
			platformalterationparameters.PlatformAlterationNamespace,
			globalhelper.Configuration.General.TestImage, platformalterationparameters.TnfTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, platformalterationparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			platformalterationparameters.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			platformalterationparameters.TnfBaseImageName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One statefulSet, one pod, change container base image by creating a file [negative]", func() {

		By("Define statefulSet")
		statefulSet := platformalterationhelper.DefineStatefulSetWithPriviledgedContainer(
			platformalterationparameters.TestStatefulSetName,
			platformalterationparameters.PlatformAlterationNamespace,
			globalhelper.Configuration.General.TestImage, platformalterationparameters.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, platformalterationhelper.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podsList, err := globalhelper.GetListOfPodsInNamespace(platformalterationparameters.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Change container base image")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "touch /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			platformalterationparameters.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			platformalterationparameters.TnfBaseImageName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
)

var _ = Describe("platform-alteration-base-image", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51297
	It("One deployment, one pod, running test image", func() {

		By("Define deployment")
		deployment := deployment.DefineDeployment(tsparams.TestDeploymentName,
			tsparams.PlatformAlterationNamespace,
			globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBaseImageName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51298
	It("One daemonSet, running test image", func() {

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace,
			globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBaseImageName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51299
	It("Two deployments, one pod each, change container base image by creating a file [negative]", func() {

		By("Define first deployment")
		deploymenta := deployment.RedefineWithPriviledgedContainer(
			deployment.DefineDeployment(tsparams.TestDeploymentName,
				tsparams.PlatformAlterationNamespace,
				globalhelper.Configuration.General.TestImage,
				tsparams.TnfTargetPodLabels))

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podsList, err := globalhelper.GetListOfPodsInNamespace(tsparams.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Change container base image")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "touch /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment")
		deploymentb := deployment.DefineDeployment("platform-alteration-dpb",
			tsparams.PlatformAlterationNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TnfTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBaseImageName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One statefulSet, one pod, change container base image by creating a file [negative]", func() {

		By("Define statefulSet")
		statefulSet := statefulset.RedefineWithPriviledgedContainer(
			statefulset.DefineStatefulSet(tsparams.TestStatefulSetName,
				tsparams.PlatformAlterationNamespace,
				globalhelper.Configuration.General.TestImage,
				tsparams.TnfTargetPodLabels))

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tshelper.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podsList, err := globalhelper.GetListOfPodsInNamespace(tsparams.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Change container base image")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "touch /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.TnfBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfBaseImageName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

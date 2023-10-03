package tests

import (
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("platform-alteration-hugepages-2m-only", func() {
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

	// 55865
	It("One deployment, one pod with 2Mi hugepages", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)
		deployment.RedefineWithCPUResources(dep, "500m", "250m")
		deployment.RedefineWith2MiHugepages(dep, 4)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-2m-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages2mOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages2mOnlyName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 55866
	It("One pod with 2Mi hugepages", func() {

		By("Define pod with 2Mi hugepages")
		puta := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TnfTargetPodLabels)
		pod.RedefineWithCPUResources(puta, "500m", "250m")
		pod.RedefineWith2MiHugepages(puta, 4)

		err := globalhelper.CreateAndWaitUntilPodIsReady(puta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-2m-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages2mOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages2mOnlyName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55867
	It("One deployment, one pod, two containers, only one with 2Mi hugepages", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)
		deployment.RedefineWithCPUResources(dep, "500m", "250m")
		deployment.RedefineWith2MiHugepages(dep, 4)
		globalhelper.AppendContainersToDeployment(dep, 1, globalhelper.GetConfiguration().General.TestImage)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-2m-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages2mOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages2mOnlyName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 55868
	It("One pod, two containers, one with 2Mi hugepages, other with 1Gi [negative]", func() {

		By("Define pod")
		put := pod.DefinePod(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)
		globalhelper.AppendContainersToPod(put, 1, globalhelper.GetConfiguration().General.TestImage)
		pod.RedefineWithCPUResources(put, "500m", "250m")

		err := pod.RedefineFirstContainerWith2MiHugepages(put, 4)
		Expect(err).ToNot(HaveOccurred())

		err = pod.RedefineSecondContainerWith1GHugepages(put, 1)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-2m-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages2mOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages2mOnlyName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})
})

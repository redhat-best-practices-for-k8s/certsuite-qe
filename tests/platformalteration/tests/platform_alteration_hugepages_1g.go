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

var _ = Describe("platform-alteration-hugepages-1g-only", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.PlatformAlterationNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.IsKindCluster() {
			Skip("Hugepages are not supported in Kind clusters")
		}

		By("Check if nodes have hugepages enabled")
		if !globalhelper.NodesHaveHugePagesEnabled("1Gi") {
			Skip("Hugepages configuration is not enabled on the cluster")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.WaitingTime)
	})

	It("One deployment, one pod with 1Gi hugepages", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)
		deployment.RedefineWithCPUResources(dep, "500m", "250m")
		deployment.RedefineWith1GiHugepages(dep, 1)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages1gOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with 1Gi hugepages", func() {

		By("Define pod with 1Gi hugepages")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TnfTargetPodLabels)
		pod.RedefineWithCPUResources(put, "500m", "250m")
		pod.RedefineWith1GiHugepages(put, 1)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages1gOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment, one pod, two containers, only one with 1Gi hugepages", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)
		deployment.RedefineWithCPUResources(dep, "500m", "250m")
		deployment.RedefineWith1GiHugepages(dep, 1)
		globalhelper.AppendContainersToDeployment(dep, 1, globalhelper.GetConfiguration().General.TestImage)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages1gOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	It("One pod, two containers, both with 1Gi hugepages", func() {

		By("Define pod")
		put := pod.DefinePod(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)
		globalhelper.AppendContainersToPod(put, 1, globalhelper.GetConfiguration().General.TestImage)
		pod.RedefineWithCPUResources(put, "500m", "250m")

		err := pod.RedefineFirstContainerWith1GiHugepages(put, 1)
		Expect(err).ToNot(HaveOccurred())

		err = pod.RedefineSecondContainerWith1GHugepages(put, 1)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages1gOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod, two containers, one with 1Gi hugepages, other with 2Mi [negative]", func() {

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

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePages1gOnlyName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

package tests

import (
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("platform-alteration-hugepages-2m-only", Serial, Label("platformalteration2", "ocp-required"), func() {
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

		if globalhelper.IsKindCluster() {
			Skip("Hugepages are not supported in Kind clusters")
		}

		By("Check if nodes have hugepages enabled")
		if !globalhelper.NodesHaveHugePagesEnabled("2Mi") {
			Skip("Hugepages configuration is not enabled on the cluster")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 55865
	It("One deployment, one pod with 2Mi hugepages", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)
		deployment.RedefineWithCPUResources(dep, "500m", "250m")
		deployment.RedefineWith2MiHugepages(dep, 4)

		By("Create and wait until deployment is ready")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-2m-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages2mOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages2mOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 55866
	It("One pod with 2Mi hugepages", func() {

		By("Define pod with 2Mi hugepages")
		puta := pod.DefinePod(tsparams.TestPodName, randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels)
		pod.RedefineWithCPUResources(puta, "500m", "250m")
		pod.RedefineWith2MiHugepages(puta, 4)

		By("Create and wait until pod is ready")
		err := globalhelper.CreateAndWaitUntilPodIsReady(puta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-2m-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages2mOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages2mOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55867
	It("One deployment, one pod, two containers, only one with 2Mi hugepages", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)
		deployment.RedefineWithCPUResources(dep, "500m", "250m")
		deployment.RedefineWith2MiHugepages(dep, 4)
		globalhelper.AppendContainersToDeployment(dep, 1, tsparams.SampleWorkloadImage)

		By("Create and wait until deployment is ready")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-2m-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages2mOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages2mOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 55868
	It("One pod, two containers, one with 2Mi hugepages, other with 1Gi [negative]", func() {

		By("Define pod")
		put := pod.DefinePod(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)
		globalhelper.AppendContainersToPod(put, 1, tsparams.SampleWorkloadImage)
		pod.RedefineWithCPUResources(put, "500m", "250m")

		err := pod.RedefineFirstContainerWith2MiHugepages(put, 4)
		Expect(err).ToNot(HaveOccurred())

		err = pod.RedefineSecondContainerWith1GHugepages(put, 1)
		Expect(err).ToNot(HaveOccurred())

		By("Create and wait until pod is ready")
		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-2m-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages2mOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages2mOnlyName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

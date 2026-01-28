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

var _ = Describe("platform-alteration-hugepages-1g-only", Serial, Label("platformalteration2", "ocp-required"), func() {
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
		if !globalhelper.NodesHaveHugePagesEnabled("1Gi") {
			Skip("Hugepages configuration is not enabled on the cluster")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	It("One deployment, one pod with 1Gi hugepages", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)
		deployment.RedefineWithCPUResources(dep, "500m", "250m")
		deployment.RedefineWith1GiHugepages(dep, 1)

		By("Create and wait until deployment is ready")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages1gOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with 1Gi hugepages", func() {

		By("Define pod with 1Gi hugepages")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels)
		pod.RedefineWithCPUResources(put, "500m", "250m")
		pod.RedefineWith1GiHugepages(put, 1)

		By("Create and wait until pod is ready")
		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages1gOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment, one pod, two containers, only one with 1Gi hugepages", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)
		deployment.RedefineWithCPUResources(dep, "500m", "250m")
		deployment.RedefineWith1GiHugepages(dep, 1)
		globalhelper.AppendContainersToDeployment(dep, 1, tsparams.SampleWorkloadImage)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages1gOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	It("One pod, two containers, both with 1Gi hugepages", func() {

		By("Define pod")
		put := pod.DefinePod(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)
		globalhelper.AppendContainersToPod(put, 1, tsparams.SampleWorkloadImage)
		pod.RedefineWithCPUResources(put, "500m", "250m")

		err := pod.RedefineFirstContainerWith1GiHugepages(put, 1)
		Expect(err).ToNot(HaveOccurred())

		err = pod.RedefineSecondContainerWith1GHugepages(put, 1)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages1gOnlyName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod, two containers, one with 1Gi hugepages, other with 2Mi [negative]", func() {

		By("Define pod")
		put := pod.DefinePod(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)
		globalhelper.AppendContainersToPod(put, 1, tsparams.SampleWorkloadImage)
		pod.RedefineWithCPUResources(put, "500m", "250m")

		err := pod.RedefineFirstContainerWith2MiHugepages(put, 4)
		Expect(err).ToNot(HaveOccurred())

		err = pod.RedefineSecondContainerWith1GHugepages(put, 1)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-1g-only test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteHugePages1gOnlyName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteHugePages1gOnlyName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

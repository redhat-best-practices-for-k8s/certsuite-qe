package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/parameters"
)

var _ = Describe("Networking ocp-reserved-ports-usage,", Serial, Label("networking", "ocp-required"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNetworkingNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 59536
	It("one deployment, one pod, one container not declaring reserved ports (OCP Ports)", func() {
		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59537
	It("one deployment, one pod, one container declaring reserved ports [negative]", func() {
		By("Define and create deployment with container declaring reserved port")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(1, []corev1.ContainerPort{{ContainerPort: 22623}}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59538
	It("one deployment, one pod, two containers, neither declaring reserved ports 22222 and 22223 (OCP Ports)", func() {
		By("Define deployment with two containers")
		ports := []corev1.ContainerPort{{ContainerPort: 22222}, {ContainerPort: 22223}}
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(2, ports, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59539
	It("one deployment, one pod, two containers, one declaring reserved ports (OCP Ports) [negative]", func() {
		ports := []corev1.ContainerPort{{ContainerPort: 22222}, {ContainerPort: 22623}}

		By("Define deployment with two containers")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(2, ports, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59540
	It("one deployment, one pod not listening on reserved ports (OCP Ports)", func() {
		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(3, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59541
	It("one deployment, one pod listening on reserved ports (OCP Ports) [negative]", func() {
		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(1, []corev1.ContainerPort{{ContainerPort: 22624}}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define service and create it on cluster")
		err = tshelper.DefineAndCreateServiceOnCluster("test-service", randomNamespace, 22624, 22624, false, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59542
	It("two deployments, one pod each not listening on reserved ports (OCP Ports)", func() {
		By("Define first deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(3, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeployment(tsparams.TestDeploymentBName, randomNamespace, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59543
	It("two deployments, one pod each, one listening on reserved ports (OCP Ports) [negative]", func() {
		By("Define first deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeployment(tsparams.TestDeploymentBName, randomNamespace, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithContainerPorts(1, []corev1.ContainerPort{{ContainerPort: 22624}}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define service and create it on cluster")
		err = tshelper.DefineAndCreateServiceOnCluster("test-service", randomNamespace, 22624, 22624, false, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOcpReservedPortsUsageTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

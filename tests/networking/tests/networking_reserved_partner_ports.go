package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking reserved-partner-ports,", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNetworkingNameSpace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
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

	// 61487
	It("one deployment, one pod, one container not declaring reserved ports (Partner Ports)", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfReservedPartnerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfReservedPartnerPortsTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 61505
	It("one deployment, one pod, one container declaring reserved ports (Partner Ports) [negative]", func() {

		By("Define and create deployment with container declaring reserved port")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(1, []corev1.ContainerPort{{ContainerPort: 15443}}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment has declared reserved port")
		runningDeployment, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(15443)))

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfReservedPartnerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfReservedPartnerPortsTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 61506
	It("one deployment, one pod, two containers, neither declaring reserved ports 15002 and 15007 (Partner Ports)", func() {

		By("Define deployment with two containers")
		ports := []corev1.ContainerPort{{ContainerPort: 15002}, {ContainerPort: 15007}}
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(2, ports, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfReservedPartnerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfReservedPartnerPortsTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 61507
	It("one deployment, one pod, two containers, one declaring reserved ports (Partner Ports) [negative]", func() {
		ports := []corev1.ContainerPort{{ContainerPort: 15020}, {ContainerPort: 15019}}

		By("Define deployment with two containers")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(2, ports, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfReservedPartnerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfReservedPartnerPortsTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 61508
	It("one deployment, one pod not listening on reserved ports (Partner Ports)", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(3, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfReservedPartnerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfReservedPartnerPortsTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 61509
	It("one deployment, one pod listening on reserved ports (Partner Ports) [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(1, []corev1.ContainerPort{{ContainerPort: 15021}}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define service and create it on cluster")
		err = tshelper.DefineAndCreateServiceOnCluster("test-service", randomNamespace, 22624, 22624, false, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfReservedPartnerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfReservedPartnerPortsTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 61510
	It("two deployments, one pod each not listening on reserved ports (Partner Ports)", func() {

		By("Define first deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(3, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeployment(tsparams.TestDeploymentBName, randomNamespace, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfReservedPartnerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfReservedPartnerPortsTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 61517
	It("two deployments, one pod each, one listening on reserved ports (Partner Ports) [negative]", func() {

		By("Define first deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeployment(tsparams.TestDeploymentBName, randomNamespace, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithContainerPorts(1, []corev1.ContainerPort{{ContainerPort: 15090}}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define service and create it on cluster")
		err = tshelper.DefineAndCreateServiceOnCluster("test-service", randomNamespace, 22624, 22624, false, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfReservedPartnerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfReservedPartnerPortsTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

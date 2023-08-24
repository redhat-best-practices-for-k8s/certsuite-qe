package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking ocp-reserved-ports-usage,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		randomNamespace = tsparams.TestNetworkingNameSpace + "-" + globalhelper.GenerateRandomString(10)

		By(fmt.Sprintf("Create %s namespace", randomNamespace))
		err := namespaces.Create(randomNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Override default report directory")
		origReportDir = globalhelper.GetConfiguration().General.TnfReportDir
		reportDir := origReportDir + "/" + randomNamespace
		globalhelper.OverrideReportDir(reportDir)

		By("Override default TNF config directory")
		origTnfConfigDir = globalhelper.GetConfiguration().General.TnfConfigDir
		configDir := origTnfConfigDir + "/" + randomNamespace
		globalhelper.OverrideTnfConfigDir(configDir)

		By("Define TNF config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By(fmt.Sprintf("Remove %s namespace", randomNamespace))
		err := namespaces.DeleteAndWait(
			globalhelper.GetAPIClient().CoreV1Interface,
			randomNamespace,
			tsparams.WaitingTime,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Restore default report directory")
		globalhelper.GetConfiguration().General.TnfReportDir = origReportDir

		By("Restore default TNF config directory")
		globalhelper.GetConfiguration().General.TnfConfigDir = origTnfConfigDir
	})

	// 59536
	It("one deployment, one pod, one container not declaring reserved ports", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59537
	It("one deployment, one pod, one container declaring reserved ports [negative]", func() {

		By("Define and create deployment with container declaring reserved port")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(1, []corev1.ContainerPort{{ContainerPort: 22623}}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59538
	It("one deployment, one pod, two containers, neither declaring reserved ports 22222 and 22223", func() {

		By("Define deployment with two containers")
		ports := []corev1.ContainerPort{{ContainerPort: 22222}, {ContainerPort: 22223}}
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(2, ports, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		// By("Sleep for 5 minutes")
		// time.Sleep(5 * time.Minute)

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59539
	It("one deployment, one pod, two containers, one declaring reserved ports [negative]", func() {
		ports := []corev1.ContainerPort{{ContainerPort: 22222}, {ContainerPort: 22623}}

		By("Define deployment with two containers")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(2, ports, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59540
	It("one deployment, one pod not listening on reserved ports (OCP Ports)", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(3, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59541
	It("one deployment, one pod listening on reserved ports [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentWithContainerPorts(1, []corev1.ContainerPort{{ContainerPort: 22624}}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define service and create it on cluster")
		err = tshelper.DefineAndCreateServiceOnCluster("test-service", randomNamespace, 22624, 22624, false, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59542
	It("two deployments, one pod each not listening on reserved ports", func() {

		By("Define first deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(3, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeployment(tsparams.TestDeploymentBName, randomNamespace, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59543
	It("two deployments, one pod each, one listening on reserved ports [negative]", func() {

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
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOcpReservedPortsUsageTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

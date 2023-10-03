package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking undeclared-container-ports-usage,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNetworkingNameSpace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	It("one deployment, one pod, container declares and uses port 8080", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 1,
			[]corev1.ContainerPort{{ContainerPort: 8080}})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	It("one deployment, one pod, container declares port 8081 but does not use any", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 1,
			[]corev1.ContainerPort{{ContainerPort: 8081}})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	It("one deployment, one pod, container declares port 8081 but uses port 8080 instead [negative]", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 1,
			[]corev1.ContainerPort{{ContainerPort: 8081}})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, container uses port 8080 but does not declare any [negative]", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainers(1, 1, "networking-deployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, two containers, both containers declare used ports (8080, 8081)", func() {

		By("Define deployment and create it on cluster")
		ports := []corev1.ContainerPort{{ContainerPort: 8080}, {ContainerPort: 8081}}
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 1, ports)
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 1, []string{})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerEnvVarList(dep, 1, []corev1.EnvVar{{Name: "LIVENESS_PROBE_DEFAULT_PORT", Value: "8081"}})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	//nolint:lll
	It("one deployment, one pod, two containers, the first one uses and declares port 8080, the second one uses port 8081 but declares 8082 [negative]", func() {

		By("Define deployment and create it on cluster")
		ports := []corev1.ContainerPort{{ContainerPort: 8080}, {ContainerPort: 8082}}
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 2, ports)
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 1, []string{})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerEnvVarList(dep, 1, []corev1.EnvVar{{Name: "LIVENESS_PROBE_DEFAULT_PORT", Value: "8081"}})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, two containers, the second container uses port 8080 but does not declare any [negative]", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainers(1, 2, "networking-deployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 1, []string{})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

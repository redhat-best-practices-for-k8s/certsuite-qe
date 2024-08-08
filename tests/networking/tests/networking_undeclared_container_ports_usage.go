package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/parameters"
)

var _ = Describe("Networking undeclared-container-ports-usage,", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

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

	It("one deployment, one pod, container declares and uses port 8080", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 1,
			[]corev1.ContainerPort{{ContainerPort: 8080}})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has container port 8080")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(8080)))

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	It("one deployment, one pod, container declares port 8081 but does not use any", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 1,
			[]corev1.ContainerPort{{ContainerPort: 8081}})
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has container port 8081")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(8081)))

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		// Skipped because no listening ports are detected
		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, container declares port 8081 but uses port 8080 instead [negative]", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 1,
			[]corev1.ContainerPort{{ContainerPort: 8081}})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has container port 8081")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(8081)))

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, container uses port 8080 but does not declare any [negative]", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainers(1, 1, "networking-deployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have container port 8080")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Ports).To(BeEmpty())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, two containers, both containers declare used ports (8080, 8081)", func() {

		By("Define deployment and create it on cluster")
		ports := []corev1.ContainerPort{{ContainerPort: 8080, Protocol: "TCP"}, {ContainerPort: 8081, Protocol: "TCP"}}
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 1, ports)
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 1, []string{})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerEnvVarList(dep, 1, []corev1.EnvVar{{Name: "LIVENESS_PROBE_DEFAULT_PORT", Value: "8081"}})
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	//nolint:lll
	It("one deployment, one pod, two containers, the first one uses and declares port 8080, the second one uses port 8081 but declares 8082 [negative]", func() {

		By("Define deployment and create it on cluster")
		ports := []corev1.ContainerPort{{ContainerPort: 8080, Protocol: "TCP"}, {ContainerPort: 8082, Protocol: "TCP"}}
		dep, err := tshelper.DefineDeploymentWithContainerPorts("networking-deployment", randomNamespace, 2, ports)
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 0, []string{})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 1, []string{})
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerEnvVarList(dep, 1, []corev1.EnvVar{{Name: "LIVENESS_PROBE_DEFAULT_PORT", Value: "8081"}})
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one deployment, one pod, two containers, the second container uses port 8080 but does not declare any [negative]", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeploymentWithContainers(1, 2, "networking-deployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		err = deployment.RedefineContainerCommand(dep, 1, []string{})
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUndeclaredContainerPortsUsageTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

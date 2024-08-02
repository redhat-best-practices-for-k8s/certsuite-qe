package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-container-prestop", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

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

	// 47311
	It("One deployment, one pod with preStop field configured", func() {

		By("Define deployment with preStop field configured")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineAllContainersWithPreStopSpec(deploymenta, tsparams.PreStopCommand)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has prestop configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle.PreStop).ToNot(BeNil())

		By("Start lifecycle-container-prestop test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47315
	It("One deployment, one pod without preStop field configured [negative]", func() {

		By("Define deployment without prestop field configured")
		deployment, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have prestop configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle).To(BeNil())

		By("Start lifecycle-container-prestop test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47382
	It("One deployment, several pods, several containers that have preStop field configured", func() {

		By("Define deployment with preStop field configured")
		deploymenta, err := tshelper.DefineDeployment(3, 2, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineAllContainersWithPreStopSpec(deploymenta, tsparams.PreStopCommand)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has prestop configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle.PreStop).ToNot(BeNil())

		By("Start lifecycle-container-prestop test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47383
	It("Two deployments, several pods, several containers that have preStop field configured", func() {

		By("Define first deployment with preStop field configured")
		deploymenta, err := tshelper.DefineDeployment(3, 2, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineAllContainersWithPreStopSpec(deploymenta, tsparams.PreStopCommand)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has prestop configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle.PreStop).ToNot(BeNil())

		By("Define second deployment with preStop field configured")
		deploymentb, err := tshelper.DefineDeployment(3, 2, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineAllContainersWithPreStopSpec(deploymentb, tsparams.PreStopCommand)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has prestop configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle.PreStop).ToNot(BeNil())

		By("Start lifecycle-container-prestop test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47384
	It("One deployment, several pods, several containers one without preStop field configured [negative]", func() {

		By("Define deployment with preStop field configured")
		deploymenta, err := tshelper.DefineDeployment(3, 2, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = deployment.RedefineFirstContainerWithPreStopSpec(deploymenta, tsparams.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-prestop test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50761
	It("Two deployments, several pods, several containers that do not have preStop field configured [negative]", func() {

		By("Define and create first deployment")
		deploymenta, err := tshelper.DefineDeployment(3, 2, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have prestop configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle).To(BeNil())

		By("Define and create second deployment")
		deploymentb, err := tshelper.DefineDeployment(3, 2, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have prestop configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle).To(BeNil())

		By("Start lifecycle-container-prestop test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-pod-high-availability", Serial, Label("lifecycle5", "ocp-required"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

	BeforeEach(func() {
		if globalhelper.IsKindCluster() {
			By("Make masters schedulable")
			err := nodes.EnableMasterScheduling(globalhelper.GetAPIClient().Nodes(), true)
			Expect(err).ToNot(HaveOccurred())
		}

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.CertsuiteTargetOperatorLabels},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 48492
	It("One deployment, replicas are more than 1, podAntiAffinity is set", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define deployment")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has podAntiAffinity configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodHighAvailabilityTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48495
	It("Two deployments, replicas are more than 1, podAntiAffinity is set", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define and create first deployment")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has podAntiAffinity configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Define and create second deployment")
		deploymentb, err := tshelper.DefineDeployment(2, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymentb, tsparams.TestTargetLabels)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has podAntiAffinity configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodHighAvailabilityTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48499
	It("One deployment, replicas are more than 1, podAntiAffinity is not set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define deployment")
		deployment, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has podAntiAffinity configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity).To(BeNil())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodHighAvailabilityTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48500
	It("Two deployments, replicas are more than 1, podAntiAffinity is not set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define and create first deployment")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have podAntiAffinity configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity).To(BeNil())

		By("Define and create second deployment")
		deploymentb, err := tshelper.DefineDeployment(2, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have podAntiAffinity configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity).To(BeNil())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodHighAvailabilityTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48869
	It("One deployment, replicas equal to 1, podAntiAffinity is set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define deployment")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has podAntiAffinity configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.GetNumberOfNodes(globalhelper.GetAPIClient().K8sClient.CoreV1()) == 1 {
			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.CertsuitePodHighAvailabilityTcName,
				globalparameters.TestCaseSkipped,
				randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		} else {
			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodHighAvailabilityTcName,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		}
	})
})

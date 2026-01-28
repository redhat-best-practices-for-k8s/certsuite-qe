package tests

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-pod-recreation", Serial, Label("lifecycle7", "ocp-required"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		By("Enable intrusive tests")
		err := os.Setenv("CERTSUITE_NON_INTRUSIVE_ONLY", "false")
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.IsKindCluster() {
			By("Make masters schedulable")
			err := nodes.EnableMasterScheduling(globalhelper.GetAPIClient().Nodes(), true)
			Expect(err).ToNot(HaveOccurred())
		}

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define certsuite config file")
		err = globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.CertsuiteTargetOperatorLabels},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.GetConfiguration().General.DisableIntrusiveTests == strings.ToLower("true") {
			Skip("Intrusive tests are disabled via config")
		}
	})

	AfterEach(func() {
		By("Disable intrusive tests")
		err := os.Setenv("CERTSUITE_NON_INTRUSIVE_ONLY", "true")
		Expect(err).ToNot(HaveOccurred())

		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 47405
	It("One deployment with PodAntiAffinity, replicas are less than schedulable nodes", func() {
		if globalhelper.IsCRCCluster() {
			Skip("PodAntiAffinity test is not supported on CRC clusters")
		}

		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enough schedulable nodes.")
		}
		// at least one "clean of any resource" worker is needed.
		maxPodsPerDeployment := schedulableNodes - 1

		By("Define deployment")
		deploymenta, err := tshelper.DefineDeployment(maxPodsPerDeployment, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has pod anti-affinity")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Start lifecycle-pod-recreation test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodRecreationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodRecreationTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47406
	It("Two deployments with PodAntiAffinity, replicas are less than schedulable nodes", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enough schedulable nodes.")
		}
		// at least one "clean of any resource" worker is needed.
		maxPodsPerDeployment := (schedulableNodes / 2) - 1

		By("Define and create first deployment")
		deploymenta, err := tshelper.DefineDeployment(maxPodsPerDeployment, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deploymenta has pod anti-affinity")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Define and create second deployment")
		deploymentb, err := tshelper.DefineDeployment(maxPodsPerDeployment, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymentb, tsparams.TestTargetLabels)

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deploymentb has pod anti-affinity")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Start lifecycle-pod-recreation test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodRecreationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodRecreationTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47407
	It("One deployment with PodAntiAffinity, replicas are equal to schedulable nodes [negative]", func() {
		if globalhelper.IsCRCCluster() {
			Skip("PodAntiAffinity test is not supported on CRC clusters")
		}

		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define deployment")
		deploymenta, err := tshelper.DefineDeployment(schedulableNodes, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has pod anti-affinity")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Start lifecycle-pod-recreation test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodRecreationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodRecreationTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47408
	It("Two deployments with PodAntiAffinity, replicas are equal to schedulable nodes [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		// all nodes will be scheduled with a pod.
		if schedulableNodes < 2 {
			Skip("The cluster does not have enough schedulable nodes.")
		}
		maxPodsPerDeploymentPerFirstDeployment := (schedulableNodes / 2)
		maxPodsPerDeploymentPerSecondDeployment := schedulableNodes - maxPodsPerDeploymentPerFirstDeployment

		By("Define and create first deployment")
		deploymenta, err := tshelper.DefineDeployment(maxPodsPerDeploymentPerFirstDeployment, 1,
			tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deploymenta has pod anti-affinity")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Define and create second deployment")
		deploymentb, err := tshelper.DefineDeployment(maxPodsPerDeploymentPerSecondDeployment, 1,
			"lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymentb, tsparams.TestTargetLabels)

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deploymentb has pod anti-affinity")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity).ToNot(BeNil())

		By("Start lifecycle-pod-recreation test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodRecreationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodRecreationTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

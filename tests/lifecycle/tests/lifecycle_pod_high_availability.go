package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-pod-high-availability", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file"))
	}

	execute.BeforeAll(func() {
		By("Make masters schedulable")
		err := globalhelper.EnableMasterScheduling(true)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all nodes are labeled with 'worker-cnf' label")
		err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().CoreV1Interface, configSuite.General.CnfNodeLabel)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 48492
	It("One deployment, replicas are more than 1, podAntiAffinity is set", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define and create deployment")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodHighAvailabilityTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48495
	It("Two deployments, replicas are more than 1, podAntiAffinity is set", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define and create first deployment")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create second deployment")
		deploymentb, err := tshelper.DefineDeployment(2, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymentb, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodHighAvailabilityTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48499
	It("One deployment, replicas are more than 1, podAntiAffinity is not set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define and create deployment")
		deployment, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodHighAvailabilityTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48500
	It("Two deployments, replicas are more than 1, podAntiAffinity is not set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define and create first deployment")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create second deployment")
		deploymentb, err := tshelper.DefineDeployment(2, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodHighAvailabilityTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48869
	It("One deployment, replicas equal to 1, podAntiAffinity is set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define and create deployment")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodAntiAffinity(deploymenta, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(tsparams.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodHighAvailabilityTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

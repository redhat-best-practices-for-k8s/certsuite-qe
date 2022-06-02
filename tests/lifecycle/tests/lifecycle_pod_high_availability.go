package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
)

var _ = Describe("lifecycle-pod-high-availability", func() {

	execute.BeforeAll(func() {
		By("Make masters schedulable")
		err := helper.EnableMasterScheduling(true)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		err := helper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(parameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48492
	It("One deployment, replicas are more than 1, podAntiAffinity is set", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define & create deployment")
		deploymenta, err := helper.DefineDeployment(2, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithPodAntiAffinity(deploymenta, parameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(parameters.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(parameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48495
	It("Two deployments, replicas are more than 1, podAntiAffinity is set", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define & create first deployment")
		deploymenta, err := helper.DefineDeployment(2, 1, "lifecycleputa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithPodAntiAffinity(deploymenta, parameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		deploymentb, err := helper.DefineDeployment(2, 1, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithPodAntiAffinity(deploymentb, parameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(parameters.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(parameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48499
	It("One deployment, replicas are more than 1, podAntiAffinity is not set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define & create deployment")
		deployment, err := helper.DefineDeployment(2, 1, "lifecycleputone")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(parameters.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(parameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48500
	It("Two deployments, replicas are more than 1, podAntiAffinity is not set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define & create first deployment")
		deploymenta, err := helper.DefineDeployment(2, 1, "lifecycleputone")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		deploymentb, err := helper.DefineDeployment(2, 1, "lifecycleputtwo")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(parameters.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(parameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48869
	It("One deployment, replicas equal to 1, podAntiAffinity is set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enough schedulable nodes.")
		}

		By("Define & create deployment")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycleputone")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithPodAntiAffinity(deploymenta, parameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(parameters.TnfPodHighAvailabilityTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(parameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

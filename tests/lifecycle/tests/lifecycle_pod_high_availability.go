package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
)

var _ = Describe("lifecycle-pod-high-availability", func() {

	stringOfSkipTc := globalhelper.GetStringOfSkipTcs(lifeparameters.TnfTestCases,
		lifeparameters.TnfPodHighAvailabilityTcName)

	execute.BeforeAll(func() {
		By("Make masters schedulable")
		err := lifehelper.EnableMasterScheduling(true)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		err := lifehelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48492
	It("One deployment, replicas are more than 1, podAntiAffinity is set", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create deployment")
		deploymenta, err := lifehelper.DefineDeployment(2, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithPodAntiAffinity(deploymenta, lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48495
	It("Two deployments, replicas are more than 1, podAntiAffinity is set", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create first deployment")
		deploymenta, err := lifehelper.DefineDeployment(2, 1, "lifecycleputa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithPodAntiAffinity(deploymenta, lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		deploymentb, err := lifehelper.DefineDeployment(2, 1, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithPodAntiAffinity(deploymentb, lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48499
	It("One deployment, replicas are more than 1, podAntiAffinity is not set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create deployment")
		deployment, err := lifehelper.DefineDeployment(2, 1, "lifecycleputone")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48500
	It("Two deployments, replicas are more than 1, podAntiAffinity is not set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create first deployment")
		deploymenta, err := lifehelper.DefineDeployment(2, 1, "lifecycleputone")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		deploymentb, err := lifehelper.DefineDeployment(2, 1, "lifecycleputtwo")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48869
	It("One deployment, replicas equal to 1, podAntiAffinity is set [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create deployment")
		deploymenta, err := lifehelper.DefineDeployment(1, 1, "lifecycleputone")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithPodAntiAffinity(deploymenta, lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle pod-high-availability test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodHighAvailabilityTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

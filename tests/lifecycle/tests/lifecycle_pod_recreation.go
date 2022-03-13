package tests

import (
	"os"

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

var _ = Describe("lifecycle lifecycle-pod-recreation", func() {

	execute.BeforeAll(func() {
		By("Make masters schedulable")
		err := lifehelper.EnableMasterScheduling(true)
		Expect(err).ToNot(HaveOccurred())

		By("Enable intrusive tests")
		err = os.Setenv("TNF_NON_INTRUSIVE_ONLY", "false")
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47405
	It("One deployment with PodAntiAffinity, replicas < schedulable nodes", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enought schedulable nodes.")
		}
		// at least one "clean of any resource" worker is needed.
		maxPodsPerDeployment := schedulableNodes - 1
		By("Define & create deployment")
		deploymentStruct := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(maxPodsPerDeployment, 1, "lifecycleput"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-recreation test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			globalhelper.GetStringOfSkipTcs(lifeparameters.SkipsTcsSlice, lifeparameters.PodRecreationName))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodRecreationName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47406
	It("Two deployments with PodAntiAffinity, replicas < schedulable nodes", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enought schedulable nodes.")
		}
		// at least one "clean of any resource" worker is needed.
		maxPodsPerDeployment := (schedulableNodes / 2) - 1
		By("Define & create first deployment")
		firstDeploymentStruct := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(maxPodsPerDeployment, 1, "lifecycleputone"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(firstDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		secondDeploymentStruct := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(maxPodsPerDeployment, 1, "lifecycleputtwo"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(secondDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-recreation test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			globalhelper.GetStringOfSkipTcs(lifeparameters.SkipsTcsSlice, lifeparameters.PodRecreationName))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodRecreationName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47407
	It("One deployment with PodAntiAffinity, replicas = schedulable nodes [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create deployment")
		deploymentStruct := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(schedulableNodes, 1, "lifecycleput"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-recreation test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			globalhelper.GetStringOfSkipTcs(lifeparameters.SkipsTcsSlice, lifeparameters.PodRecreationName))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodRecreationName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47408
	It("Two deployments with PodAntiAffinity, replicas = schedulable nodes [negative]", func() {
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		// all nodes will be scheduled with a pod.
		if schedulableNodes < 2 {
			Skip("The cluster does not have enought schedulable nodes.")
		}
		maxPodsPerDeploymentPerFirstDeployment := (schedulableNodes / 2)
		maxPodsPerDeploymentPerSecondDeployment := schedulableNodes - maxPodsPerDeploymentPerFirstDeployment

		By("Define & create first deployment")
		firstDeploymentStruct := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(maxPodsPerDeploymentPerFirstDeployment, 1, "lifecycleputone"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(firstDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		secondDeploymentStruct := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(maxPodsPerDeploymentPerSecondDeployment, 1, "lifecycleputtwo"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(secondDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-recreation test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			globalhelper.GetStringOfSkipTcs(lifeparameters.SkipsTcsSlice, lifeparameters.PodRecreationName))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodRecreationName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})
})

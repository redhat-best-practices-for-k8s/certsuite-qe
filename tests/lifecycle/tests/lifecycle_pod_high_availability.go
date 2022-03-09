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

var _ = Describe("lifecycle lifecycle-pod-high-availability", func() {

	execute.BeforeAll(func() {
		By("Make masters schedulable")
		err := lifehelper.EnableMasterScheduling(true)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48492
	specName1 := "One deployment, replicas are more than 1, podAntiAffinity is set"
	It(specName1, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName1)
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create deployment")
		deploymentStruct := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(2, 1, "lifecycleput"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-high-availability test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodHighAvailabilityDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodHighAvailabilityRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodHighAvailabilityDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 48495
	specName2 := "Two deployments, replicas are more than 1, podAntiAffinity is set"
	It(specName2, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName2)
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create first deployment")
		lifecycleputone := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(2, 1, "lifecycleputone"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(lifecycleputone, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		lifecycleputtwo := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(2, 1, "lifecycleputtwo"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(lifecycleputtwo, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-high-availability test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodHighAvailabilityDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodHighAvailabilityRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodHighAvailabilityDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 48499
	specName3 := "One deployment, replicas are more than 1, podAntiAffinity is not set [negative]"
	It(specName3, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName3)
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 2 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create deployment")
		lifecycleputone := lifehelper.DefineDeployment(2, 1, "lifecycleputone")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(lifecycleputone, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-high-availability test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodHighAvailabilityDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodHighAvailabilityRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodHighAvailabilityDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 48500
	specName4 := "Two deployments, replicas are more than 1, podAntiAffinity is not set [negative]"
	It(specName4, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName4)
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes < 4 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create first deployment")
		lifecycleputone := lifehelper.DefineDeployment(2, 1, "lifecycleputone")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(lifecycleputone, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		lifecycleputtwo := lifehelper.DefineDeployment(2, 1, "lifecycleputtwo")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(lifecycleputtwo, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-high-availability test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodHighAvailabilityDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodHighAvailabilityRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodHighAvailabilityDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48869
	specName5 := "One deployment, replicas equal to 1, podAntiAffinity is set [negative]"
	It(specName5, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName5)
		schedulableNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		if schedulableNodes == 0 {
			Skip("The cluster does not have enought schedulable nodes.")
		}

		By("Define & create deployment")
		lifecycleputone := deployment.RedefineWithPodAntiAffinity(
			lifehelper.DefineDeployment(1, 1, "lifecycleputone"),
			lifeparameters.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(lifecycleputone, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-high-availability test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodHighAvailabilityDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodHighAvailabilityRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodHighAvailabilityDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})
})

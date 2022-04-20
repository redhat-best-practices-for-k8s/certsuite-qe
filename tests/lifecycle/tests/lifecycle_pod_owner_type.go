package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
)

var _ = Describe("lifecycle-pod-owner-type", func() {

	stringOfSkipTc := globalhelper.GetStringOfSkipTcs(lifeparameters.TnfTestCases, lifeparameters.TnfPodOwnerTypeTcName)

	BeforeEach(func() {
		err := lifehelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47409
	It("One ReplicaSet, several pods", func() {

		By("Define ReplicaSet with replica number")
		replicaSet := replicaset.RedefineWithReplicaNumber(lifehelper.DefineReplicaSet("lifecyclers"), 3)

		err := lifehelper.CreateAndWaitUntilReplicaSetIsReady(replicaSet, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodOwnerTypeTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47424
	It("Two deployments, several pods", func() {

		By("Define deployments")
		deploymenta, err := lifehelper.DefineDeployment(2, 1, "lifecycleputa")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentb, err := lifehelper.DefineDeployment(2, 1, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodOwnerTypeTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47426
	It("StatefulSet pod", func() {

		By("Define statefulSet")
		statefulSet := lifehelper.DefineStatefulSet("lifecyclesf")
		err := lifehelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodOwnerTypeTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47429
	It("One pod, not part of any workload resource [negative]", func() {

		By("Define pod")
		pod := pod.RedefinePodWithLabel(lifehelper.DefinePod("lifecyclepod"),
			lifeparameters.TestDeploymentLabels)
		err := lifehelper.CreateAndWaitUntilPodIsReady(pod, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodOwnerTypeTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47430
	It("Two deployments, one pod not related to any resource [negative]", func() {

		By("Define deployments")
		deploymenta, err := lifehelper.DefineDeployment(2, 1, "lifecycleputa")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentb, err := lifehelper.DefineDeployment(2, 1, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		pod := pod.RedefinePodWithLabel(lifehelper.DefinePod("lifecyclepod"),
			lifeparameters.TestDeploymentLabels)
		err = lifehelper.CreateAndWaitUntilPodIsReady(pod, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfPodOwnerTypeTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
)

var _ = Describe("lifecycle lifecycle-pod-owner-type", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47409
	It("1 ReplicaSet, several pods", func() {
		By("Define ReplicaSet with replica number")
		replicaStruct := replicaset.RedefineWithReplicaNumber(
			lifehelper.DefineReplicaSet("lifecyclers"), 3)

		err := globalhelper.CreateAndWaitUntilReplicaSetIsReady(
			replicaStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47424
	It("Multiple deployments, replica  > 1", func() {
		By("Define deployments")
		firstDeploymentStruct := lifehelper.DefineDeployment(2, 1, "lifecycleputone")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(firstDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		secondDeploymentStruct := lifehelper.DefineDeployment(2, 1, "lifecycleputtwo")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(secondDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47426
	It("StatefulSet pod", func() {
		By("Define statefulSet")
		statefulSetStruct := lifehelper.DefineStatefulSet("lifecyclesf")
		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSetStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47429
	It("1 pod, not part of any workload resource [negative]", func() {
		By("Define pod")
		podStruct := lifehelper.DefindPod("lifecyclepod")
		err := globalhelper.CreateAndWaitUntilPodIsReady(podStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47430
	It("Multiple deployments, 1 pod not related to any resource [negative]", func() {

		By("Define deployments")
		firstDeploymentStruct := lifehelper.DefineDeployment(2, 1, "lifecycleputone")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(firstDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		secondDeploymentStruct := lifehelper.DefineDeployment(2, 1, "lifecycleputtwo")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(secondDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		podStruct := lifehelper.DefindPod("lifecyclepod")
		err = globalhelper.CreateAndWaitUntilPodIsReady(podStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

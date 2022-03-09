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

var _ = Describe("lifecycle lifecycle-pod-owner-type", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47409
	specName1 := "1 ReplicaSet, several pods"
	It(specName1, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName1)
		By("Define ReplicaSet with replica number")
		replicaStruct := replicaset.RedefineWithReplicaNumber(lifehelper.DefineReplicaSet("lifecyclers"), 3)

		err := lifehelper.CreateAndWaitUntilReplicaSetIsReady(replicaStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodOwnerTypeDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47424
	specName2 := "2 deployments, several pods"
	It(specName2, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName2)
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
			lifeparameters.PodOwnerTypeDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47426
	specName3 := "StatefulSet pod"
	It(specName3, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName3)
		By("Define statefulSet")
		statefulSetStruct := lifehelper.DefineStatefulSet("lifecyclesf")
		err := lifehelper.CreateAndWaitUntilStatefulSetIsReady(statefulSetStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodOwnerTypeDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47429
	specName4 := "1 pod, not part of any workload resource [negative]"
	It(specName4, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName4)
		By("Define pod")
		podStruct := pod.RedefinePodWithLabel(lifehelper.DefindPod("lifecyclepod"),
			lifeparameters.TestDeploymentLabels)
		err := lifehelper.CreateAndWaitUntilPodIsReady(podStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodOwnerTypeDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47430
	specName5 := "Two deployments, 1 pod not related to any resource [negative]"
	It(specName5, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName5)
		By("Define deployments")
		firstDeploymentStruct := lifehelper.DefineDeployment(2, 1, "lifecycleputone")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(firstDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		secondDeploymentStruct := lifehelper.DefineDeployment(2, 1, "lifecycleputtwo")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(secondDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		podStruct := pod.RedefinePodWithLabel(lifehelper.DefindPod("lifecyclepod"),
			lifeparameters.TestDeploymentLabels)
		err = lifehelper.CreateAndWaitUntilPodIsReady(podStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.PodOwnerTypeDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButPodOwnerTypeRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.PodOwnerTypeDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

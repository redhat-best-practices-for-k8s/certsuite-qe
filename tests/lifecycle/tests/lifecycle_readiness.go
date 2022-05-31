package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
)

var _ = Describe("lifecycle-readiness", func() {

	BeforeEach(func() {
		err := helper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(parameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50145
	It("One deployment, one pod with a readiness probe", func() {
		By("Define deployment with a readiness probe")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycledp")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithReadinessProbe(deploymenta)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			parameters.TnfReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfReadinessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50146
	It("Two deployments, multiple pods each, all have a readiness probe", func() {
		By("Define first deployment with a readiness probe")
		deploymenta, err := helper.DefineDeployment(3, 1, "lifecycledpa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithReadinessProbe(deploymenta)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with a readiness probe")
		deploymentb, err := helper.DefineDeployment(3, 1, "lifecycledpb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithReadinessProbe(deploymentb)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			parameters.TnfReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfReadinessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50147
	It("One statefulSet, one pod with a readiness probe", func() {
		By("Define statefulSet with a readiness probe")
		statefulset := statefulset.RedefineWithReadinessProbe(
			helper.DefineStatefulSet("lifecycle-sf"))
		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			parameters.TnfReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfReadinessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50148
	It("One pod with a readiness probe", func() {
		By("Define pod with a readiness probe")
		pod := pod.RedefineWithReadinessProbe(pod.RedefinePodWithLabel(
			helper.DefinePod("lifecycleput"), parameters.TestDeploymentLabels))
		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			parameters.TnfReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfReadinessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50149
	It("One daemonSet without a readiness probe [negative]", func() {
		By("Define daemonSet without a readiness probe")
		daemonSet := daemonset.DefineDaemonSet(parameters.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage,
			parameters.TestDeploymentLabels, "lifecycleds")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			parameters.TnfReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfReadinessTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50150
	It("Two deployments, one pod each, one without a readiness probe [negative]", func() {
		By("Define first deployment with a readiness probe")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycledpa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithReadinessProbe(deploymenta)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment without a readiness probe")
		deploymentb, err := helper.DefineDeployment(1, 1, "lifecycledpb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(
			parameters.TnfReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfReadinessTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

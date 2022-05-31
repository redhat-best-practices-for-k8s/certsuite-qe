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

var _ = Describe("lifecycle-liveness", func() {

	BeforeEach(func() {
		err := helper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(parameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50053
	It("One deployment, one pod with a liveness probe", func() {
		By("Define deployment with a liveness probe")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycledp")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithLivenessProbe(deploymenta)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			parameters.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfLivenessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50054
	It("Two deployments, multiple pods each, all have a liveness probe", func() {
		By("Define first deployment with a liveness probe")
		deploymenta, err := helper.DefineDeployment(3, 1, "lifecycledpa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithLivenessProbe(deploymenta)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with a liveness probe")
		deploymentb, err := helper.DefineDeployment(3, 1, "lifecycledpb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithLivenessProbe(deploymentb)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			parameters.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfLivenessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50055
	It("One statefulSet, one pod with a liveness probe", func() {
		By("Define statefulSet with a liveness probe")
		statefulset := statefulset.RedefineWithLivenessProbe(

			helper.DefineStatefulSet("lifecycle-sf"))
		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			parameters.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfLivenessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50056
	It("One pod with a liveness probe", func() {
		By("Define pod with a liveness probe")
		pod := pod.RedefineWithLivenessProbe(pod.RedefinePodWithLabel(

			helper.DefinePod("lifecycleput"), parameters.TestDeploymentLabels))
		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			parameters.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfLivenessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50057
	It("One daemonSet without a liveness probe [negative]", func() {
		By("Define daemonSet without a liveness probe")
		daemonSet := daemonset.DefineDaemonSet(parameters.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage,
			parameters.TestDeploymentLabels, "lifecyclesds")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			parameters.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfLivenessTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50058
	It("Two deployments, one pod each, one without a liveness probe [negative]", func() {
		By("Define first deployment with a liveness probe")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycledpa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithLivenessProbe(deploymenta)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment without a liveness probe")
		deploymentb, err := helper.DefineDeployment(1, 1, "lifecycledpb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			parameters.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfLivenessTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

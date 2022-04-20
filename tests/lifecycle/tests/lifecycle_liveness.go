package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
)

var _ = Describe("lifecycle-liveness", func() {

	stringOfSkipTc := globalhelper.GetStringOfSkipTcs(lifeparameters.TnfTestCases,
		lifeparameters.TnfLivenessTcName)

	BeforeEach(func() {
		err := lifehelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50053
	It("One deployment, one pod with a liveness probe", func() {
		By("Define deployment with a liveness probe")
		deployment := deployment.RedefineWithLivenessProbe(
			lifehelper.DefineDeployment(1, 1, "lifecycledp"))
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfLivenessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50054
	It("Two deployments, multiple pods each, all have a liveness probe", func() {
		By("Define first deployment with a liveness probe")
		deploymenta := deployment.RedefineWithLivenessProbe(
			lifehelper.DefineDeployment(3, 1, "lifecycledpa"))
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with a liveness probe")
		deploymentb := deployment.RedefineWithLivenessProbe(
			lifehelper.DefineDeployment(3, 1, "lifecycledpb"))
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfLivenessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50055
	It("One statefulSet, one pod with a liveness probe", func() {
		By("Define statefulSet with a liveness probe")
		statefulset := statefulset.RedefineWithLivenessProbe(
			lifehelper.DefineStatefulSet("lifecycle-sf"))
		err := lifehelper.CreateAndWaitUntilStatefulSetIsReady(statefulset, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfLivenessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50056
	It("One pod with a liveness probe", func() {
		By("Define pod with a liveness probe")
		pod := pod.RedefineWithLivenessProbe(pod.RedefinePodWithLabel(
			lifehelper.DefinePod("lifecycleput"), lifeparameters.TestDeploymentLabels))
		err := lifehelper.CreateAndWaitUntilPodIsReady(pod, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfLivenessTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50057
	It("One daemonSet without a liveness probe [negative]", func() {
		By("Define daemonSet without a liveness probe")
		daemonSet := daemonset.DefineDaemonSet(lifeparameters.LifecycleNamespace,
			globalhelper.Configuration.General.TnfImage,
			lifeparameters.TestDeploymentLabels, "lifecyclesds")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfLivenessTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50058
	It("Two deployments, one pod each, one without a liveness probe [negative]", func() {
		By("Define first deployment with a liveness probe")
		deploymenta := deployment.RedefineWithLivenessProbe(
			lifehelper.DefineDeployment(1, 1, "lifecycledpa"))
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment without a liveness probe")
		deploymentb := lifehelper.DefineDeployment(1, 1, "lifecycledpb")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfLivenessTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle-container-shutdown", func() {

	stringOfSkipTc := globalhelper.GetStringOfSkipTcs(lifeparameters.TnfTestCases, lifeparameters.TnfShutdownTcName)

	BeforeEach(func() {
		err := lifehelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47311
	It("One deployment, one pod, with one container that has preStop field configured", func() {
		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := deployment.RedefineAllContainersWithPreStopSpec(
			lifehelper.DefineDeployment(1, 1, "lifecycleput"), lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47315
	It("One deployment, one pod, with one container that does not have preStop field configured [negative]", func() {
		By("Define deployment without prestop field configured")
		deploymentStructWithOutPreStop := lifehelper.DefineDeployment(1, 1, "lifecycleput")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStructWithOutPreStop, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47382
	It("One deployment, several pods, several containers that have preStop field configured", func() {
		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := deployment.RedefineAllContainersWithPreStopSpec(
			lifehelper.DefineDeployment(3, 2, "lifecycleput"), lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47383
	It("Two deployments, several pods, several containers that have preStop field configured", func() {
		By("Define first deployment with preStop field configured")
		preStopDeploymentStructA, err := deployment.RedefineAllContainersWithPreStopSpec(
			lifehelper.DefineDeployment(3, 2, "lifecycleputone"), lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStructA, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with preStop field configured")
		preStopDeploymentStructB, err := deployment.RedefineAllContainersWithPreStopSpec(
			lifehelper.DefineDeployment(3, 2, "lifecycleputtwo"), lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStructB, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47384
	It("One deployment, several pods, several containers one without preStop field configured [negative]", func() {
		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := deployment.RedefineFirstContainerWithPreStopSpec(
			lifehelper.DefineDeployment(3, 2, "lifecycleput"), lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47385
	It("Two deployments, several pods, several containers that don't have preStop field configured [negative]", func() {
		By("Define first deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
			lifehelper.DefineDeployment(3, 2, "lifecycleputone"), lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			lifehelper.DefineDeployment(3, 2, "lifecycleputtwo"), lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.LifecycleTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			stringOfSkipTc)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

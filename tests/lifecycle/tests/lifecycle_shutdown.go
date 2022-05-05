package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle-container-shutdown", func() {

	BeforeEach(func() {
		err := lifehelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47311
	It("One deployment, one pod with preStop field configured", func() {

		By("Define deployment with preStop field configured")

		deploymenta, err := lifehelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineAllContainersWithPreStopSpec(
			deploymenta, lifeparameters.PreStopCommand)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47315
	It("One deployment, one pod without preStop field configured [negative]", func() {

		By("Define deployment without prestop field configured")
		deployment, err := lifehelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
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
		deploymenta, err := lifehelper.DefineDeployment(3, 2, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineAllContainersWithPreStopSpec(
			deploymenta, lifeparameters.PreStopCommand)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
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
		deploymenta, err := lifehelper.DefineDeployment(3, 2, "lifecycleputa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineAllContainersWithPreStopSpec(
			deploymenta, lifeparameters.PreStopCommand)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with preStop field configured")
		deploymentb, err := lifehelper.DefineDeployment(3, 2, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineAllContainersWithPreStopSpec(
			deploymentb, lifeparameters.PreStopCommand)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
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
		deploymenta, err := lifehelper.DefineDeployment(3, 2, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta, err = deployment.RedefineFirstContainerWithPreStopSpec(
			deploymenta, lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47385
	It("Two deployments, several pods, several containers that don't have preStop field configured [negative]", func() {

		By("Define & create first deployment")
		deploymenta, err := lifehelper.DefineDeployment(3, 2, "lifecycleputa")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define & create second deployment")
		deploymentb, err := lifehelper.DefineDeployment(3, 2, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfShutdownTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

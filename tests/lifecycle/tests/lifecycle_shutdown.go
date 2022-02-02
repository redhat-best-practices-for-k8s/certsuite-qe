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

var _ = Describe("lifecycle lifecycle-container-shutdown", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47311
	It("One deployment, one pod with one container that "+
		"has preStop field configured", func() {

		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := lifehelper.DefineDeployment(true, 1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButShutdownRegex,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47315
	It("One deployment, one pod with one container that does not "+
		"have preStop field configured [negative]", func() {

		By("Define deployment without prestop field configured")
		deploymentStructWithOutPreStop, err := lifehelper.DefineDeployment(false, 1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStructWithOutPreStop, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47382
	It("One deployment, several pods, several containers that "+
		"have preStop field configured", func() {

		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := lifehelper.DefineDeployment(true, 3, 2, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47383
	It("Several deployments, several pods, several containers "+
		"that have preStop field configured", func() {

		By("Define first deployment with preStop field configured")
		preStopDeploymentStructA, err := lifehelper.DefineDeployment(true, 3, 2, "lifecycleputone")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStructA, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with preStop field configured")
		preStopDeploymentStructB, err := lifehelper.DefineDeployment(true, 3, 2, "lifecycleputtwo")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStructB, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47384
	It("One deployment, several pods, several containers, one without preStop "+
		"field configured [negative]", func() {

		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := lifehelper.DefineDeployment(false, 3, 2, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		preStopDeploymentStruct, err = deployment.RedefineFirstContainerWithPreStopSpec(
			preStopDeploymentStruct, lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47385
	It("Several deployments, several pods, several containers "+
		"that don't have preStop field configured [negative]", func() {

		By("Define first deployment")
		replicaDefinedDeploymentA, err := lifehelper.DefineDeployment(false, 3, 2, "lifecycleputone")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			replicaDefinedDeploymentA, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment")
		replicaDefinedDeploymentB, err := lifehelper.DefineDeployment(false, 3, 2, "lifecycleputtwo")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			replicaDefinedDeploymentB, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

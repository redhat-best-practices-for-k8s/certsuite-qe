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
	It("One deployment, one pod with one container that has preStop field configured", func() {

		By("Define deployment with preStop field configured")
		preStopDeploymentStruct := deployment.RedefineWithPreStopSpec(
			lifehelper.DefineLifecycleDeployment(false, "lifecycleput"),
			lifeparameters.PreStopCommand)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(preStopDeploymentStruct, lifeparameters.WaitingTime)
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
	It("One deployment, one pod with one container that does not have preStop field configured [negative]", func() {

		By("Define deployment without prestop field configured")
		deploymentStructWithOutPreStop := lifehelper.DefineLifecycleDeployment(false, "lifecycleput")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStructWithOutPreStop, lifeparameters.WaitingTime)
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
	It("One deployment, several pods, several containers that has preStop field configured", func() {

		By("Define deployment with preStop field configured")
		replicaDefinedDeployment := lifehelper.DefineLifecycleDeploymentSeveralPodsWithTwoContainers(
			"lifecycleputone", 3)
		preStopDeploymentStruct := deployment.RedefineAllContainersWithPreStopSpec(
			replicaDefinedDeployment, lifeparameters.PreStopCommand)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
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
	It("Several deployments, several pods, several containers that has preStop field configured", func() {

		By("Define first deployment with preStop field configured")
		replicaDefinedDeploymentA := lifehelper.DefineLifecycleDeploymentSeveralPodsWithTwoContainers(
			"lifecycleputone", 3)
		preStopDeploymentStructA := deployment.RedefineAllContainersWithPreStopSpec(
			replicaDefinedDeploymentA, lifeparameters.PreStopCommand)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStructA, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with preStop field configured", func() {
			replicaDefinedDeploymentB := lifehelper.DefineLifecycleDeploymentSeveralPodsWithTwoContainers(
				"lifecycleputtwo", 3)
			preStopDeploymentStructB := deployment.RedefineAllContainersWithPreStopSpec(
				replicaDefinedDeploymentB, lifeparameters.PreStopCommand)

			err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
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
	})

	// 47384
	It("One deployment, several pods, several containers, one without preStop field configured [negative]", func() {

		By("Define deployment with preStop field configured")
		replicaDefinedDeployment := lifehelper.DefineLifecycleDeploymentSeveralPodsWithTwoContainers(
			"lifecycleputone", 3)
		preStopDeploymentStruct := deployment.RedefineFirstContainerWithPreStopSpec(
			replicaDefinedDeployment, lifeparameters.PreStopCommand)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
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
	It("Several deployments, several pods, several containers that does not have preStop field configured [negative]", func() {

		By("Define first deployment")
		replicaDefinedDeploymentA := lifehelper.DefineLifecycleDeploymentSeveralPodsWithTwoContainers(
			"lifecycleputone", 3)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
			replicaDefinedDeploymentA, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment")
		replicaDefinedDeploymentB := lifehelper.DefineLifecycleDeploymentSeveralPodsWithTwoContainers(
			"lifecycleputtwo", 3)

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

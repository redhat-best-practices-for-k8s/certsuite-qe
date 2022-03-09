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
	specName1 := "One deployment, one pod, with one container that has preStop field configured"
	It(specName1, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName1)
		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := deployment.RedefineAllContainersWithPreStopSpec(
			lifehelper.DefineDeployment(1, 1, "lifecycleput"), lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.ShutdownDefaultName,
			tcNameForReport,
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
	specName2 := "One deployment, one pod, with one container that does not have preStop field configured [negative]"
	It(specName2, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName2)
		By("Define deployment without prestop field configured")
		deploymentStructWithOutPreStop := lifehelper.DefineDeployment(1, 1, "lifecycleput")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStructWithOutPreStop, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.ShutdownDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47382
	specName3 := "One deployment, several pods, several containers that have preStop field configured"
	It(specName3, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName3)

		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := deployment.RedefineAllContainersWithPreStopSpec(
			lifehelper.DefineDeployment(3, 2, "lifecycleput"), lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.ShutdownDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47383
	specName4 := "Two deployments, several pods, several containers that have preStop field configured"
	It(specName4, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName4)

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
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.ShutdownDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47384
	specName5 := "One deployment, several pods, several containers one without preStop field configured [negative]"
	It(specName5, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName5)

		By("Define deployment with preStop field configured")
		preStopDeploymentStruct, err := deployment.RedefineFirstContainerWithPreStopSpec(
			lifehelper.DefineDeployment(3, 2, "lifecycleput"), lifeparameters.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.ShutdownDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47385
	specName6 := "Two deployments, several pods, several containers that don't have preStop field configured [negative]"
	It(specName6, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName6)

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
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.ShutdownDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

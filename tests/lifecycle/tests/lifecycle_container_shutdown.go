package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-container-shutdown", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all nodes are labeled with 'worker-cnf' label")
		err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().CoreV1Interface, configSuite.General.CnfNodeLabel)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 47311
	It("One deployment, one pod with preStop field configured", func() {

		By("Define deployment with preStop field configured")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineAllContainersWithPreStopSpec(deploymenta, tsparams.PreStopCommand)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47315
	It("One deployment, one pod without preStop field configured [negative]", func() {

		By("Define deployment without prestop field configured")
		deployment, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47382
	It("One deployment, several pods, several containers that have preStop field configured", func() {

		By("Define deployment with preStop field configured")
		deploymenta, err := tshelper.DefineDeployment(3, 2, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineAllContainersWithPreStopSpec(deploymenta, tsparams.PreStopCommand)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47383
	It("Two deployments, several pods, several containers that have preStop field configured", func() {

		By("Define first deployment with preStop field configured")
		deploymenta, err := tshelper.DefineDeployment(3, 2, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineAllContainersWithPreStopSpec(deploymenta, tsparams.PreStopCommand)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with preStop field configured")
		deploymentb, err := tshelper.DefineDeployment(3, 2, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineAllContainersWithPreStopSpec(deploymentb, tsparams.PreStopCommand)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47384
	It("One deployment, several pods, several containers one without preStop field configured [negative]", func() {

		By("Define deployment with preStop field configured")
		deploymenta, err := tshelper.DefineDeployment(3, 2, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = deployment.RedefineFirstContainerWithPreStopSpec(deploymenta, tsparams.PreStopCommand)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50761
	It("Two deployments, several pods, several containers that do not have preStop field configured [negative]", func() {

		By("Define and create first deployment")
		deploymenta, err := tshelper.DefineDeployment(3, 2, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create second deployment")
		deploymentb, err := tshelper.DefineDeployment(3, 2, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			tsparams.TnfShutdownTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfShutdownTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

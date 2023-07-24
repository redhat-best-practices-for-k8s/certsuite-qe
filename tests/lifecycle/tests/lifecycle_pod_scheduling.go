package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-pod-scheduling", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file"))
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
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 48120
	It("One deployment, no nodeSelector nor nodeAffinity", func() {

		By("Define Deployment")
		deployment, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(tsparams.TnfPodSchedulingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodSchedulingTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48458
	It("One deployment with nodeSelector [negative]", func() {

		By("Define Deployment with nodeSelector")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithNodeSelector(deploymenta, map[string]string{configSuite.General.CnfNodeLabel: ""})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(tsparams.TnfPodSchedulingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodSchedulingTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48470
	It("One deployment with nodeAffinity [negative]", func() {

		By("Define Deployment with nodeAffinity")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithNodeAffinity(deploymenta, configSuite.General.CnfNodeLabel)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(tsparams.TnfPodSchedulingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodSchedulingTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48471
	It("Two deployments, one pod each, one pod with nodeAffinity [negative]", func() {

		By("Define Deployment without nodeAffinity")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define Deployment with nodeAffinity")
		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithNodeAffinity(deploymentb, configSuite.General.CnfNodeLabel)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(tsparams.TnfPodSchedulingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodSchedulingTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48472
	It("One deployment, one daemonSet [negative]", func() {

		By("Define Deployment without nodeAffinity/ nodeSelector")
		deployment, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(tsparams.LifecycleNamespace,
			globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(tsparams.TnfPodSchedulingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodSchedulingTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

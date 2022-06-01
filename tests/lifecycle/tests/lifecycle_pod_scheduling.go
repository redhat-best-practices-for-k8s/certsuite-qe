package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle-pod-scheduling", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file"))
	}

	BeforeEach(func() {
		err := helper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(parameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48120
	It("One deployment, no nodeSelector nor nodeAffinity", func() {

		By("Define Deployment")
		deployment, err := helper.DefineDeployment(1, 1, "lifecycledp")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(
  			parameters.TnfPodSchedulingTcName,
	  		globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfPodSchedulingTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48458
	It("One deployment with nodeSelector [negative]", func() {

		By("Define Deployment with nodeSelector")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycledp")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithNodeSelector(deploymenta,
			map[string]string{configSuite.General.CnfNodeLabel: ""})
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(
  			parameters.TnfPodSchedulingTcName,
	  		globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfPodSchedulingTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48470
	It("One deployment with nodeAffinity [negative]", func() {

		By("Define Deployment with nodeAffinity")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycledp")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithNodeAffinity(deploymenta, configSuite.General.CnfNodeLabel)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(
  			parameters.TnfPodSchedulingTcName,
	  		globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfPodSchedulingTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48471
	It("Two deployments, one pod each, one pod with nodeAffinity [negative]", func() {

		By("Define Deployment without nodeAffinity")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycledpa")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define Deployment with nodeAffinity")
		deploymentb, err := helper.DefineDeployment(1, 1, "lifecycledpb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithNodeAffinity(deploymentb,
			configSuite.General.CnfNodeLabel)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(
  			parameters.TnfPodSchedulingTcName,
	  		globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfPodSchedulingTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48472
	It("One deployment, one daemonSet [negative]", func() {

		By("Define Deployment without nodeAffinity/ nodeSelector")
		deployment, err := helper.DefineDeployment(1, 1, "lifecycledp")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(parameters.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage,
			parameters.TestDeploymentLabels, "lifecycleds")

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-scheduling test")
		err = globalhelper.LaunchTests(
			parameters.TnfPodSchedulingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfPodSchedulingTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

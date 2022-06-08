package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("lifecycle-image-pull-policy", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48473
	It("One deployment with ifNotPresent as ImagePullPolicy", func() {

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48474
	It("Several deployments with ifNotPresent as ImagePullPolicy", func() {

		By("Define deployments with ifNotPresent as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, "lifecycleputa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithImagePullPolicy(deploymentb, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentc, err := tshelper.DefineDeployment(1, 1, "lifecycleputc")
		Expect(err).ToNot(HaveOccurred())

		deploymentc = deployment.RedefineWithImagePullPolicy(deploymentc, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentc, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48478
	It("One DaemonSet with ifNotPresent as ImagePullPolicy", func() {

		By("Define DaemonSet with ifNotPresent as ImagePullPolicy")
		daemonSet := tshelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycleds", globalhelper.Configuration.General.TestImage, v1.PullIfNotPresent)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48479
	It("Several DaemonSets with ifNotPresent as ImagePullPolicy", func() {

		By("Define DaemonSets with ifNotPresent as ImagePullPolicy")
		daemonSeta := tshelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycledsa", globalhelper.Configuration.General.TestImage, v1.PullIfNotPresent)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSeta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		daemonSetb := tshelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycledsb", globalhelper.Configuration.General.TestImage, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSetb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		daemonSetc := tshelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycledsc", globalhelper.Configuration.General.TestImage, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSetc, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48480
	It("One DaemonSet without ImagePullPolicy, image tag is not specified [negative]", func() {
		// if you omit the imagePullPolicy field,
		// and you don't specify the tag for the container image,
		// imagePullPolicy is automatically set to Always;
		By("Define DaemonSet without ImagePullPolicy")
		daemonSet := daemonset.DefineDaemonSet(
			tsparams.LifecycleNamespace,
			"registry.access.redhat.com/ubi8/ubi",
			tsparams.TestDeploymentLabels, "lifecycleds")

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48481
	It("One deployment without ImagePullPolicy, image tag is latest [negative]", func() {
		// if you omit the imagePullPolicy field,
		// and the tag for the container image is :latest,
		// imagePullPolicy is automatically set to Always;
		By("Define deployment without ImagePullPolicy")
		deployment := deployment.DefineDeployment("lifecycleput",
			tsparams.LifecycleNamespace,
			"registry.access.redhat.com/ubi8/ubi:latest",
			tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48482
	It("One deployment with Always as ImagePullPolicy [negative]", func() {

		By("Define deployment with 'Always' as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullAlways)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48484
	It("Two deployments one with Never other with ifNotPresent as ImagePullPolicy [negative]", func() {

		By("Define deployment with Never as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullNever)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithImagePullPolicy(deploymentb, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48485
	It("One DaemonSet with Never one deployment with ifNotPresent as ImagePullPolicy [negative]", func() {

		By("Define DaemonSet with Never as ImagePullPolicy")
		daemonSet := tshelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycleds", globalhelper.Configuration.General.TestImage, v1.PullNever)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

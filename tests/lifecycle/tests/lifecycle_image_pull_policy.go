package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("lifecycle-image-pull-policy", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48473
	It("One deployment with ifNotPresent as ImagePullPolicy", func() {

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymenta, err := lifehelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48474
	It("Several deployments with ifNotPresent as ImagePullPolicy", func() {

		By("Define deployments with ifNotPresent as ImagePullPolicy")
		deploymenta, err := lifehelper.DefineDeployment(1, 1, "lifecycleputa")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentb, err := lifehelper.DefineDeployment(1, 1, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithImagePullPolicy(deploymentb, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentc, err := lifehelper.DefineDeployment(1, 1, "lifecycleputc")
		Expect(err).ToNot(HaveOccurred())

		deploymentc = deployment.RedefineWithImagePullPolicy(deploymentc, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentc, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48478
	It("One DaemonSet with ifNotPresent as ImagePullPolicy", func() {

		By("Define DaemonSet with ifNotPresent as ImagePullPolicy")
		daemonSet := lifehelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycleds", globalhelper.Configuration.General.TestImage, v1.PullIfNotPresent)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48479
	It("Several DaemonSets with ifNotPresent as ImagePullPolicy", func() {

		By("Define DaemonSets with ifNotPresent as ImagePullPolicy")
		daemonSeta := lifehelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycledsa", globalhelper.Configuration.General.TestImage, v1.PullIfNotPresent)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSeta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		daemonSetb := lifehelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycledsb", globalhelper.Configuration.General.TestImage, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSetb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		daemonSetc := lifehelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycledsc", globalhelper.Configuration.General.TestImage, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSetc, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
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
			lifeparameters.LifecycleNamespace,
			"registry.access.redhat.com/ubi8/ubi",
			lifeparameters.TestDeploymentLabels, "lifecycleds")

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
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
			lifeparameters.LifecycleNamespace,
			"registry.access.redhat.com/ubi8/ubi:latest",
			lifeparameters.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48482
	It("One deployment with Always as ImagePullPolicy [negative]", func() {

		By("Define deployment with 'Always' as ImagePullPolicy")
		deploymenta, err := lifehelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullAlways)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48484
	It("Two deployments one with Never other with ifNotPresent as ImagePullPolicy [negative]", func() {

		By("Define deployment with Never as ImagePullPolicy")
		deploymenta, err := lifehelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullNever)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymentb, err := lifehelper.DefineDeployment(1, 1, "lifecycleputb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithImagePullPolicy(deploymentb, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48485
	It("One DaemonSet with Never one deployment with ifNotPresent as ImagePullPolicy [negative]", func() {

		By("Define DaemonSet with Never as ImagePullPolicy")
		daemonSet := lifehelper.DefineDaemonSetWithImagePullPolicy(
			"lifecycleds", globalhelper.Configuration.General.TestImage, v1.PullNever)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymenta, err := lifehelper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithImagePullPolicy(deploymenta, v1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(
			lifeparameters.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TnfImagePullPolicyTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

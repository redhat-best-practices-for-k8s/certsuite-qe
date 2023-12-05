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
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("lifecycle-image-pull-policy", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	// 48473
	It("One deployment with ifNotPresent as ImagePullPolicy", func() {

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithImagePullPolicy(deploymenta, corev1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has ifNotPresent configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullIfNotPresent))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48474
	It("Several deployments with ifNotPresent as ImagePullPolicy", func() {

		By("Define deployments with ifNotPresent as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithImagePullPolicy(deploymenta, corev1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has ifNotPresent configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullIfNotPresent))

		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithImagePullPolicy(deploymentb, corev1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has ifNotPresent configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullIfNotPresent))

		deploymentc, err := tshelper.DefineDeployment(1, 1, "lifecycle-dpc", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithImagePullPolicy(deploymentc, corev1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentc, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has ifNotPresent configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentc.Namespace, deploymentc.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullIfNotPresent))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48478
	It("One DaemonSet with ifNotPresent as ImagePullPolicy", func() {

		By("Define DaemonSet with ifNotPresent as ImagePullPolicy")
		daemonSet := tshelper.DefineDaemonSetWithImagePullPolicy(tsparams.TestDaemonSetName,
			randomNamespace, globalhelper.GetConfiguration().General.TestImage, corev1.PullIfNotPresent)

		By("Create DaemonSet")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that DaemonSet has ifNotPresent as ImagePullPolicy")
		pullPolicy, err := globalhelper.GetDaemonSetPullPolicy(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(pullPolicy).To(Equal(corev1.PullIfNotPresent))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48479
	It("Several DaemonSets with ifNotPresent as ImagePullPolicy", func() {

		By("Define DaemonSets with ifNotPresent as ImagePullPolicy")
		daemonSeta := tshelper.DefineDaemonSetWithImagePullPolicy(tsparams.TestDaemonSetName,
			randomNamespace, globalhelper.GetConfiguration().General.TestImage, corev1.PullIfNotPresent)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSeta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that DaemonSetA has ifNotPresent as ImagePullPolicy")
		pullPolicy, err := globalhelper.GetDaemonSetPullPolicy(daemonSeta)
		Expect(err).ToNot(HaveOccurred())
		Expect(pullPolicy).To(Equal(corev1.PullIfNotPresent))

		daemonSetb := tshelper.DefineDaemonSetWithImagePullPolicy("lifecycle-dsb",
			randomNamespace, globalhelper.GetConfiguration().General.TestImage, corev1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSetb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that DaemonSetB has ifNotPresent as ImagePullPolicy")
		pullPolicy, err = globalhelper.GetDaemonSetPullPolicy(daemonSetb)
		Expect(err).ToNot(HaveOccurred())
		Expect(pullPolicy).To(Equal(corev1.PullIfNotPresent))

		daemonSetc := tshelper.DefineDaemonSetWithImagePullPolicy("lifecycle-dsc",
			randomNamespace, globalhelper.GetConfiguration().General.TestImage, corev1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSetc, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that DaemonSetC has ifNotPresent as ImagePullPolicy")
		pullPolicy, err = globalhelper.GetDaemonSetPullPolicy(daemonSetc)
		Expect(err).ToNot(HaveOccurred())
		Expect(pullPolicy).To(Equal(corev1.PullIfNotPresent))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48480
	It("One DaemonSet without ImagePullPolicy, image tag is not specified [negative]", func() {
		// if you omit the imagePullPolicy field,
		// and you do not specify the tag for the container image,
		// imagePullPolicy is automatically set to Always;
		By("Define DaemonSet without ImagePullPolicy")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, "registry.access.redhat.com/ubi8/ubi",
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Create DaemonSet")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that DaemonSet has Always as ImagePullPolicy")
		pullPolicy, err := globalhelper.GetDaemonSetPullPolicy(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(pullPolicy).To(Equal(corev1.PullAlways))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48481
	It("One deployment without ImagePullPolicy, image tag is latest [negative]", func() {
		// if you omit the imagePullPolicy field,
		// and the tag for the container image is :latest,
		// imagePullPolicy is automatically set to Always;
		By("Define deployment without ImagePullPolicy")
		deployment := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			"registry.access.redhat.com/ubi8/ubi:latest", tsparams.TestTargetLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has Always pull policy configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullAlways))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48482
	It("One deployment with Always as ImagePullPolicy [negative]", func() {

		By("Define deployment with 'Always' as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithImagePullPolicy(deploymenta, corev1.PullAlways)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has Always pull policy configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullAlways))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48484
	It("Two deployments one with Never other with ifNotPresent as ImagePullPolicy [negative]", func() {

		By("Define deployment with Never as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithImagePullPolicy(deploymenta, corev1.PullNever)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has Never pull policy configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullNever))

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithImagePullPolicy(deploymentb, corev1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has ifNotPresent pull policy configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullIfNotPresent))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48485
	It("One DaemonSet with Never one deployment with ifNotPresent as ImagePullPolicy [negative]", func() {

		By("Define DaemonSet with Never as ImagePullPolicy")
		daemonSet := tshelper.DefineDaemonSetWithImagePullPolicy(tsparams.TestDaemonSetName,
			randomNamespace, globalhelper.GetConfiguration().General.TestImage, corev1.PullNever)

		By("Create DaemonSet")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that DaemonSet has Never as ImagePullPolicy")
		pullPolicy, err := globalhelper.GetDaemonSetPullPolicy(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(pullPolicy).To(Equal(corev1.PullNever))

		By("Define deployment with ifNotPresent as ImagePullPolicy")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithImagePullPolicy(deploymenta, corev1.PullIfNotPresent)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has ifNotPresent pull policy configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullIfNotPresent))

		By("Start lifecycle-image-pull-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfImagePullPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfImagePullPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/statefulset"
)

var _ = Describe("lifecycle-container-poststart", Label("lifecycle2"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 55910
	It("One deployment, one pod with postStart spec", func() {
		By("Define deployment with postStart spec")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPostStart(deploymenta)

		By("Deploy deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has postStart configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle.PostStart).ToNot(BeNil())

		By("Start lifecycle-container-poststart test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerStartUpTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55911
	It("Two deployments, two containers each, all have postStart spec", func() {
		By("Define first deployment with postStart spec")
		deploymenta, err := tshelper.DefineDeployment(1, 2, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPostStart(deploymenta)

		By("Deploy deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has postStart configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle.PostStart).ToNot(BeNil())

		By("Define second deployment with postStart spec")
		deploymentb, err := tshelper.DefineDeployment(1, 2, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPostStart(deploymentb)

		By("Deploy deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has postStart configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle.PostStart).ToNot(BeNil())

		By("Start lifecycle-container-poststart test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerStartUpTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55913
	It("One statefulSet, one pod with postStart spec", func() {
		By("Define statefulSet with postStart spec")
		statefulSet := tshelper.DefineStatefulSet(tsparams.TestStatefulSetName, randomNamespace)
		statefulset.RedefineWithPostStart(statefulSet)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulSet has postStart configured")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(statefulSet.Namespace, statefulSet.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningStatefulSet.Spec.Template.Spec.Containers[0].Lifecycle).ToNot(BeNil())

		By("Start lifecycle-container-poststart test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerStartUpTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55915
	It("One pod with postStart spec", func() {
		By("Define pod with postStart spec")
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		pod.RedefineWithPostStart(put)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod has postStart configured")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, put.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Spec.Containers[0].Lifecycle).ToNot(BeNil())
		Expect(runningPod.Spec.Containers[0].Lifecycle.PostStart).ToNot(BeNil())

		By("Start lifecycle-container-poststart test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerStartUpTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55916
	It("One daemonSet without postStart spec [negative]", func() {
		By("Define daemonSet without postStart spec")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace,
			tsparams.SampleWorkloadImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Deploy daemonSet")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert DaemonSet is without postStart spec")
		runningDaemonset, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDaemonset.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(runningDaemonset.Spec.Template.Spec.Containers[0].Lifecycle).To(BeNil())

		By("Start lifecycle-container-poststart test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerStartUpTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55914
	It("Two deployments, one pod each, one without postStart spec [negative]", func() {
		By("Define first deployment with postStart spec")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPostStart(deploymenta)

		By("Deploy deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has postStart configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Lifecycle.PostStart).ToNot(BeNil())

		By("Define second deployment without postStart spec")
		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Deploy deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have postStart configured")
		runningDeployment2, err := globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].Lifecycle).To(BeNil())

		By("Start lifecycle-container-poststart test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerStartUpTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

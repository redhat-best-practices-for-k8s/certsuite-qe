package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/statefulset"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-readiness", Label("lifecycle8"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.CertsuiteTargetOperatorLabels},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 50145
	It("One deployment, one pod with a readiness probe", func() {
		By("Define deployment with a readiness probe")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithReadinessProbe(deploymenta)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has readiness probe")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteReadinessTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50146
	It("Two deployments, multiple pods each, all have a readiness probe", func() {
		By("Define first deployment with a readiness probe")
		deploymenta, err := tshelper.DefineDeployment(3, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithReadinessProbe(deploymenta)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deploymenta has readiness probe")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())

		By("Define second deployment with a readiness probe")
		deploymentb, err := tshelper.DefineDeployment(3, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithReadinessProbe(deploymentb)

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deploymentb has readiness probe")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteReadinessTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50147
	It("One statefulSet, one pod with a readiness probe", func() {
		By("Define statefulSet with a readiness probe")
		statefulSet := tshelper.DefineStatefulSet(tsparams.TestStatefulSetName, randomNamespace)
		statefulset.RedefineWithReadinessProbe(statefulSet)
		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulSet has readiness probe")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(statefulSet.Namespace, statefulSet.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningStatefulSet.Spec.Template.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteReadinessTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50148
	It("One pod with a readiness probe", func() {
		By("Define pod with a readiness probe")
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		pod.RedefineWithReadinessProbe(put)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod has readiness probe configured")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, put.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteReadinessTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50149
	It("One daemonSet without a readiness probe [negative]", func() {
		By("Define daemonSet without a readiness probe")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Create daemonSet")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert daemonSet has no readiness probe")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDaemonSet.Spec.Template.Spec.Containers[0].ReadinessProbe).To(BeNil())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteReadinessTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50150
	It("Two deployments, one pod each, one without a readiness probe [negative]", func() {
		By("Define first deployment with a readiness probe")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithReadinessProbe(deploymenta)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deploymenta has readiness probe")
		runningDeployment, err := globalhelper.GetRunningDeployment(deploymenta.Namespace, deploymenta.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())

		By("Define second deployment without a readiness probe")
		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deploymentb has no readiness probe")
		runningDeployment, err = globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].ReadinessProbe).To(BeNil())

		By("Start lifecycle-readiness test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteReadinessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteReadinessTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

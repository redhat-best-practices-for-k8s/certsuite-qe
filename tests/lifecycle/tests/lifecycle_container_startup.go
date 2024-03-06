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
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
)

var _ = Describe("lifecycle-container-poststart", func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.WaitingTime)
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
		err = globalhelper.LaunchTests(tsparams.TnfContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerStartUpTcName, globalparameters.TestCasePassed, randomReportDir)
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
		err = globalhelper.LaunchTests(tsparams.TnfContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerStartUpTcName, globalparameters.TestCasePassed, randomReportDir)
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
		err = globalhelper.LaunchTests(tsparams.TnfContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerStartUpTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55915
	It("One pod with postStart spec", func() {
		By("Define pod with postStart spec")
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		pod.RedefineWithPostStart(put)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-container-poststart test")
		err = globalhelper.LaunchTests(tsparams.TnfContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerStartUpTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55916
	It("One daemonSet without postStart spec [negative]", func() {
		By("Define daemonSet without postStart spec")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace,
			globalhelper.GetConfiguration().General.TestImage,
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
		err = globalhelper.LaunchTests(tsparams.TnfContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerStartUpTcName, globalparameters.TestCaseFailed, randomReportDir)
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
		err = globalhelper.LaunchTests(tsparams.TnfContainerStartUpTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerStartUpTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

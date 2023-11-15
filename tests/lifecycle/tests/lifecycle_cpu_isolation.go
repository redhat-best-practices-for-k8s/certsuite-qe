package tests

import (
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/runtimeclass"
)

const (
	DeletingRTC = "Deleting RTC: "
)

var _ = Describe("lifecycle-cpu-isolation", Serial, func() {
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

		if globalhelper.IsKindCluster() && runtime.NumCPU() <= 2 {
			Skip("This test requires more than 2 CPU cores")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	const disableVar = "disable"

	// 54723
	It("One pod with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources and runTimeClass")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		put.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		pod.RedefineWithRunTimeClass(put, rtc.Name)
		pod.RedefineWithCPUResources(put, "1", "1")

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54728
	It("One pod two containers with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources and runTimeClass")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		globalhelper.AppendContainersToPod(put, 1, globalhelper.GetConfiguration().General.TestImage)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		put.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		pod.RedefineWithRunTimeClass(put, rtc.Name)
		pod.RedefineWithCPUResources(put, "1", "1")

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54732
	It("One deployment one pod with conditions met", func() {

		annotationsMap := make(map[string]string)

		By("Define deployment with resources and runTimeClass")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		dep.Spec.Template.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		deployment.RedefineWithRunTimeClass(dep, rtc.Name)
		deployment.RedefineWithCPUResources(dep, "1", "1")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has runTimeClass configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Template.Spec.RuntimeClassName).To(Equal(rtc.Name))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal("1"))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal("1"))

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54733
	It("One daemonSet with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define daemonSet with resources and runTimeClass")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		daemonSet.Spec.Template.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		daemonset.RedefineWithRunTimeClass(daemonSet, rtc.Name)
		daemonset.RedefineWithCPUResources(daemonSet, "1", "1")

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert DaemonSet is with runTimeClass and CPU resources")
		runningDaemonset, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDaemonset.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(runningDaemonset.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal("1"))
		Expect(runningDaemonset.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal("1"))
		Expect(*runningDaemonset.Spec.Template.Spec.RuntimeClassName).To(Equal(rtc.Name))

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54734
	It("One daemonSet no annotations [negative]", func() {
		By("Define daemonSet with resources and runTimeClass")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		daemonset.RedefineWithRunTimeClass(daemonSet, rtc.Name)
		daemonset.RedefineWithCPUResources(daemonSet, "1", "1")

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54735
	It("One deployment no runTimeClass [negative]", func() {

		annotationsMap := make(map[string]string)

		By("Define deployment with resources and no runTimeClass")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		dep.Spec.Template.SetAnnotations(annotationsMap)
		deployment.RedefineWithCPUResources(dep, "1", "1")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has runTimeClass configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.RuntimeClassName).To(BeNil())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54737
	It("Two pods one with conditions met, other lacks runTimeClass [negative]", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources and runTimeClass")
		puta := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		putb := pod.DefinePod("lifecycle-podb", randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		puta.SetAnnotations(annotationsMap)
		putb.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Define runTimeClass for the first pod")
		pod.RedefineWithRunTimeClass(puta, rtc.Name)
		pod.RedefineWithCPUResources(puta, "1", "1")

		By("Redfine the second pod with CPU resources")
		pod.RedefineWithCPUResources(putb, "1", "1")

		By("Create the first pod")
		err = globalhelper.CreateAndWaitUntilPodIsReady(puta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Create the second pod")
		err = globalhelper.CreateAndWaitUntilPodIsReady(putb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

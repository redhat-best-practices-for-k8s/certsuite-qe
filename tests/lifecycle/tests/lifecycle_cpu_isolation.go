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
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/runtimeclass"
)

var _ = Describe("lifecycle-cpu-isolation", func() {

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DeleteRunTimeClasses()
		Expect(err).ToNot(HaveOccurred())
	})

	const disableVar = "disable"

	// 54723
	It("One pod with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources & runTimeClass")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		put.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass("test")
		err := runtimeclass.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		pod.RedefineWithRunTimeClass(put, rtc)
		pod.RedefineWithResources(put, "1", "1")
		pod.RedefinePodWithLabel(put, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54728
	It("One pod two containers with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources & runTimeClass")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage)
		globalhelper.AppendContainersToPod(put, 1, globalhelper.Configuration.General.TestImage)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		put.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass("test")
		err := runtimeclass.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		pod.RedefineWithRunTimeClass(put, rtc)
		pod.RedefineWithResources(put, "1", "1")
		pod.RedefinePodWithLabel(put, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54732
	It("One deployment one pod with conditions met", func() {

		annotationsMap := make(map[string]string)

		By("Define deployment with resources & runTimeClass")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		dep.Spec.Template.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass("test")
		err := runtimeclass.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithRunTimeClass(dep, rtc)
		deployment.RedefineWithResources(dep, "1", "1")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54733
	It("One daemonSet with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define daemonSet with resources & runTimeClass")
		daemonSet := daemonset.DefineDaemonSet(tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		daemonSet.Spec.Template.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass("test")
		err := runtimeclass.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		daemonset.RedefineWithRunTimeClass(daemonSet, rtc)
		daemonset.RedefineWithResources(daemonSet, "1", "1")

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54734
	It("One daemonSet no annotations [negative]", func() {
		By("Define daemonSet with resources & runTimeClass")
		daemonSet := daemonset.DefineDaemonSet(tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass("test")
		err := runtimeclass.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		daemonset.RedefineWithRunTimeClass(daemonSet, rtc)
		daemonset.RedefineWithResources(daemonSet, "1", "1")

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54735
	It("One deployment no runTimeClass [negative]", func() {

		annotationsMap := make(map[string]string)

		By("Define deployment with resources & runTimeClass")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		dep.Spec.Template.SetAnnotations(annotationsMap)
		deployment.RedefineWithResources(dep, "1", "1")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54737
	It("Two pods one with conditions met, other lacks runTimeClass [negative]", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources & runTimeClass")
		puta := pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage)
		putb := pod.DefinePod("lifecycle-podb", tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		puta.SetAnnotations(annotationsMap)
		putb.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass("test")
		err := runtimeclass.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		By("Define runTimeClass for the first pod")
		pod.RedefineWithRunTimeClass(puta, rtc)
		pod.RedefineWithResources(puta, "1", "1")
		pod.RedefinePodWithLabel(puta, tsparams.TestTargetLabels)

		pod.RedefineWithResources(putb, "1", "1")
		pod.RedefinePodWithLabel(putb, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(puta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(putb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

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

var (
	// Each tc will save the RTC that was created in order to delete them.
	rtcNames = []string{}
)

var _ = Describe("lifecycle-cpu-isolation", func() {

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

	})

	AfterEach(func() {
		By("Clean namespace after each test in order to enable RunTimeClass deletion.")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Delete all RTC's that were created by the previous test case.")
		for _, rtc := range rtcNames {
			By("Deleting rtc " + rtc)
			err := tshelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		}

		// clear the list.
		rtcNames = []string{}
	})

	const disableVar = "disable"

	// 54723
	It("One pod with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources and runTimeClass")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		put.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		rtcNames = append(rtcNames, tsparams.TnfRunTimeClass)

		pod.RedefineWithRunTimeClass(put, rtc.Name)
		pod.RedefineWithCPUResources(put, "1", "1")

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54728
	It("One pod two containers with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources and runTimeClass")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels)
		globalhelper.AppendContainersToPod(put, 1, globalhelper.Configuration.General.TestImage)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		put.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		rtcNames = append(rtcNames, tsparams.TnfRunTimeClass)

		pod.RedefineWithRunTimeClass(put, rtc.Name)
		pod.RedefineWithCPUResources(put, "1", "1")

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54732
	It("One deployment one pod with conditions met", func() {

		annotationsMap := make(map[string]string)

		By("Define deployment with resources and runTimeClass")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		dep.Spec.Template.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		rtcNames = append(rtcNames, tsparams.TnfRunTimeClass)

		deployment.RedefineWithRunTimeClass(dep, rtc.Name)
		deployment.RedefineWithCPUResources(dep, "1", "1")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54733
	It("One daemonSet with conditions met", func() {
		annotationsMap := make(map[string]string)

		By("Define daemonSet with resources and runTimeClass")
		daemonSet := daemonset.DefineDaemonSet(tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		daemonSet.Spec.Template.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		rtcNames = append(rtcNames, tsparams.TnfRunTimeClass)

		daemonset.RedefineWithRunTimeClass(daemonSet, rtc.Name)
		daemonset.RedefineWithCPUResources(daemonSet, "1", "1")

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54734
	It("One daemonSet no annotations [negative]", func() {
		By("Define daemonSet with resources and runTimeClass")
		daemonSet := daemonset.DefineDaemonSet(tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.TnfRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		rtcNames = append(rtcNames, tsparams.TnfRunTimeClass)

		daemonset.RedefineWithRunTimeClass(daemonSet, rtc.Name)
		daemonset.RedefineWithCPUResources(daemonSet, "1", "1")

		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54735
	It("One deployment no runTimeClass [negative]", func() {

		annotationsMap := make(map[string]string)

		By("Define deployment with resources and runTimeClass")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TestTargetLabels)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar

		dep.Spec.Template.SetAnnotations(annotationsMap)
		deployment.RedefineWithCPUResources(dep, "1", "1")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54737
	It("Two pods one with conditions met, other lacks runTimeClass [negative]", func() {
		annotationsMap := make(map[string]string)

		By("Define pod with resources and runTimeClass")
		puta := pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels)
		putb := pod.DefinePod("lifecycle-podb", tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
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

		rtcNames = append(rtcNames, tsparams.TnfRunTimeClass)

		By("Define runTimeClass for the first pod")
		pod.RedefineWithRunTimeClass(puta, rtc.Name)
		pod.RedefineWithCPUResources(puta, "1", "1")

		pod.RedefineWithCPUResources(putb, "1", "1")

		err = globalhelper.CreateAndWaitUntilPodIsReady(puta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(putb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.TnfCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCPUIsolationTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

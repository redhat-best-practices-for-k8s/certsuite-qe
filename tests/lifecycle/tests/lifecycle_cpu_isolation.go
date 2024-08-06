package tests

import (
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/runtimeclass"
)

const (
	DeletingRTC = "Deleting RTC: "
)

var _ = Describe("lifecycle-cpu-isolation", Serial, func() {
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
			[]string{tsparams.CertsuiteTargetOperatorLabels},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.IsKindCluster() && runtime.NumCPU() <= 2 {
			Skip("This test requires more than 2 CPU cores")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
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
		rtc := runtimeclass.DefineRunTimeClass(tsparams.CertsuiteRunTimeClass)
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
		err = globalhelper.LaunchTests(tsparams.CertsuiteCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCPUIsolationTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54728
	It("One pod two containers with conditions met", func() {
		annotationsMap := make(map[string]string)

		// Check if nodes in the cluster are overcommitted as far as resources
		// are concerned. If so, skip the test.
		overcommitted, err := globalhelper.IsClusterOvercommitted()
		Expect(err).ToNot(HaveOccurred())
		if overcommitted {
			Skip("This test requires nodes to be undercommitted")
		}

		By("Define pod with resources and runTimeClass")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		globalhelper.AppendContainersToPod(put, 1, globalhelper.GetConfiguration().General.TestImage)

		By("Add annotations to the pod")
		annotationsMap["cpu-load-balancing.crio.io"] = disableVar
		annotationsMap["irq-load-balancing.crio.io"] = disableVar
		put.SetAnnotations(annotationsMap)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.CertsuiteRunTimeClass)
		err = globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Redefine the first container with CPU resources")
		pod.RedefineWithRunTimeClass(put, rtc.Name)
		pod.RedefineWithCPUResources(put, "1", "1")

		By("Create pod and wait until it is ready")
		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCPUIsolationTcName, globalparameters.TestCasePassed, randomReportDir)
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
		rtc := runtimeclass.DefineRunTimeClass(tsparams.CertsuiteRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		deployment.RedefineWithRunTimeClass(dep, rtc.Name)
		deployment.RedefineWithCPUResources(dep, "1", "1")

		By("Deploy deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has runTimeClass configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Template.Spec.RuntimeClassName).To(Equal(rtc.Name))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal("1"))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal("1"))

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCPUIsolationTcName, globalparameters.TestCasePassed, randomReportDir)
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
		rtc := runtimeclass.DefineRunTimeClass(tsparams.CertsuiteRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		daemonset.RedefineWithRunTimeClass(daemonSet, rtc.Name)
		daemonset.RedefineWithCPUResources(daemonSet, "1", "1")

		By("Deploy daemonSet")
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
		err = globalhelper.LaunchTests(tsparams.CertsuiteCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCPUIsolationTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54734
	It("One daemonSet no annotations [negative]", func() {
		By("Define daemonSet with resources and runTimeClass")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		By("Define runTimeClass")
		rtc := runtimeclass.DefineRunTimeClass(tsparams.CertsuiteRunTimeClass)
		err := globalhelper.CreateRunTimeClass(rtc)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By(DeletingRTC + rtc.Name)
			err := globalhelper.DeleteRunTimeClass(rtc)
			Expect(err).ToNot(HaveOccurred())
		})

		daemonset.RedefineWithRunTimeClass(daemonSet, rtc.Name)
		daemonset.RedefineWithCPUResources(daemonSet, "1", "1")

		By("Deploy daemonSet")
		err = globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCPUIsolationTcName, globalparameters.TestCaseFailed, randomReportDir)
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

		By("Deploy deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has runTimeClass configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.RuntimeClassName).To(BeNil())

		By("Start lifecycle-cpu-isolation test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCPUIsolationTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54737
	It("Two pods one with conditions met, other lacks runTimeClass [negative]", func() {
		annotationsMap := make(map[string]string)

		// Check if nodes in the cluster are overcommitted as far as resources
		// are concerned. If so, skip the test.
		overcommitted, err := globalhelper.IsClusterOvercommitted()
		Expect(err).ToNot(HaveOccurred())
		if overcommitted {
			Skip("This test requires nodes to be undercommitted")
		}

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
		rtc := runtimeclass.DefineRunTimeClass(tsparams.CertsuiteRunTimeClass)
		err = globalhelper.CreateRunTimeClass(rtc)
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
		err = globalhelper.LaunchTests(tsparams.CertsuiteCPUIsolationTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCPUIsolationTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

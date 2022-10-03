package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-startup-probe", func() {

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54808
	It("One deployment, one pod with a startup probe", func() {
		By("Define deployment with a startup probe")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithStartUpProbe(deploymenta)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-startup-probe test")
		err = globalhelper.LaunchTests(tsparams.TnfStartUpProbeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStartUpProbeTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54809
	It("Two deployments, multiple pods each, all have a startup probe", func() {
		By("Define first deployment with a startup probe")
		deploymenta, err := tshelper.DefineDeployment(3, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithStartUpProbe(deploymenta)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with a startup probe")
		deploymentb, err := tshelper.DefineDeployment(3, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		deploymentb = deployment.RedefineWithStartUpProbe(deploymentb)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-startup-probe test")
		err = globalhelper.LaunchTests(tsparams.TnfStartUpProbeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStartUpProbeTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54810
	It("One statefulSet, one pod with a startup probe", func() {
		By("Define statefulSet with a startup probe")
		statefulset := statefulset.RedefineWithStartUpProbe(tshelper.DefineStatefulSet(tsparams.TestStatefulSetName))

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-startup-probe test")
		err = globalhelper.LaunchTests(tsparams.TnfStartUpProbeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStartUpProbeTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54811
	It("One pod with a startup probe", func() {
		By("Define pod with a startup probe")
		pod := pod.RedefineWithStartUpProbe(pod.RedefinePodWithLabel(
			tshelper.DefinePod(tsparams.TestPodName), tsparams.TestTargetLabels))

		err := globalhelper.CreateAndWaitUntilPodIsReady(pod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-startup-probe test")
		err = globalhelper.LaunchTests(tsparams.TnfStartUpProbeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStartUpProbeTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54814
	It("One deployment two containers with a startup probe", func() {
		By("Define deployment with a startup probe")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deploymenta = globalhelper.AppendContainersToDeployment(deploymenta, 1, globalhelper.Configuration.General.TestImage)
		deploymenta = deployment.RedefineWithStartUpProbe(deploymenta)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-startup-probe test")
		err = globalhelper.LaunchTests(tsparams.TnfStartUpProbeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStartUpProbeTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54812
	It("One daemonSet without a startup probe [negative]", func() {
		By("Define daemonSet without a startup probe")
		daemonSet := daemonset.DefineDaemonSet(tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-startup-probe test")
		err = globalhelper.LaunchTests(tsparams.TnfStartUpProbeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStartUpProbeTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54813
	It("Two deployments, one pod each, one without a startup probe [negative]", func() {
		By("Define first deployment with a startup probe")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithStartUpProbe(deploymenta)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment without a startup probe")
		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-startup-probe test")
		err = globalhelper.LaunchTests(tsparams.TnfStartUpProbeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStartUpProbeTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54815
	It("One deployment two containers one has a startup probe, other does not [negative]", func() {
		By("Define deployment with a startup probe")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deploymenta = deployment.RedefineWithStartUpProbe(deploymenta)
		deploymenta = globalhelper.AppendContainersToDeployment(deploymenta, 1, globalhelper.Configuration.General.TestImage)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-startup-probe test")
		err = globalhelper.LaunchTests(tsparams.TnfStartUpProbeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStartUpProbeTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

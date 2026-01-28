package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/replicaset"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-pod-owner-type", Label("lifecycle6"), func() {
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
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 47409
	It("One ReplicaSet, several pods", func() {

		By("Define ReplicaSet with replica number")
		replicaSet := tshelper.DefineReplicaSet(tsparams.TestReplicaSetName, randomNamespace)
		replicaset.RedefineWithReplicaNumber(replicaSet, 3)

		err := globalhelper.CreateAndWaitUntilReplicaSetIsReady(replicaSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert replicaSet has pods")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">=", 3))

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodOwnerTypeTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47424
	It("Two deployments, several pods", func() {

		By("Define deployments")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentb, err := tshelper.DefineDeployment(2, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert both deployments have pods")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">=", 4))

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodOwnerTypeTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47426
	It("StatefulSet pod", func() {
		By("Define statefulSet")
		statefulSet := tshelper.DefineStatefulSet(tsparams.TestStatefulSetName, randomNamespace)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulSet has pod")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(statefulSet.Namespace, statefulSet.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningStatefulSet).ToNot(BeNil())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodOwnerTypeTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47429
	It("One pod, not part of any workload resource [negative]", func() {

		By("Define pod")
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is ready and not owned by any workload")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, put.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())
		Expect(len(runningPod.OwnerReferences)).To(Equal(0))

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodOwnerTypeTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47430
	It("Two deployments, one pod not related to any resource [negative]", func() {
		By("Define deployments")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentb, err := tshelper.DefineDeployment(2, 1, "lifecycle-dpb", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodOwnerTypeTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

})

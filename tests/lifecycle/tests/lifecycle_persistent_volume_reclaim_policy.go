package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolume"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
)

var _ = Describe("lifecycle-persistent-volume-reclaim-policy", func() {

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54201
	It("One deployment, one pod with a volume that uses a reclaim policy of delete", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TestTargetLabels)

		dep = deployment.RedefineWithVolume(dep, "test")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pv := persistentvolume.RedefineWithPVReclaimPolicy(
			persistentvolume.DefinePersistentVolume("test", tsparams.LifecycleNamespace), corev1.PersistentVolumeReclaimDelete)

		err = tshelper.CreatePersistentVolume(pv, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicy, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54202
	It("One pod with a volume that uses a reclaim policy of delete", func() {

		By("Define pod")
		put := pod.RedefineWithVolume(pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage), "test")

		put = pod.RedefinePodWithLabel(put, tsparams.TestTargetLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pv := persistentvolume.RedefineWithPVReclaimPolicy(
			persistentvolume.DefinePersistentVolume("test", tsparams.LifecycleNamespace), corev1.PersistentVolumeReclaimDelete)

		err = tshelper.CreatePersistentVolume(pv, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicy, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54203
	It("One replicaSet with a volume that uses a reclaim policy of delete", func() {

		By("Define replicaSet")
		rs := replicaset.RedefineWithVolume(replicaset.DefineReplicaSet(tsparams.TestReplicaSetName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TestTargetLabels), "test")

		err := tshelper.CreateAndWaitUntilReplicaSetIsReady(rs, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pv := persistentvolume.RedefineWithPVReclaimPolicy(
			persistentvolume.DefinePersistentVolume("test", tsparams.LifecycleNamespace), corev1.PersistentVolumeReclaimDelete)

		err = tshelper.CreatePersistentVolume(pv, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicy, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54204
	It("One deployment, one pod with a volume that uses a reclaim policy of retain [negative]", func() {

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TestTargetLabels)

		dep = deployment.RedefineWithVolume(dep, "test")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pv := persistentvolume.RedefineWithPVReclaimPolicy(
			persistentvolume.DefinePersistentVolume("test", tsparams.LifecycleNamespace), corev1.PersistentVolumeReclaimRetain)

		err = tshelper.CreatePersistentVolume(pv, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicy, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54206
	It("One pod with a volume that uses a reclaim policy of recycle [negative]", func() {

		By("Define pod")
		put := pod.RedefineWithVolume(pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage), "test")

		put = pod.RedefinePodWithLabel(put, tsparams.TestTargetLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pv := persistentvolume.RedefineWithPVReclaimPolicy(
			persistentvolume.DefinePersistentVolume("test", tsparams.LifecycleNamespace), corev1.PersistentVolumeReclaimRecycle)

		err = tshelper.CreatePersistentVolume(pv, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicy, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54207
	It("Two deployments, one with reclaim policy of delete, other with recycle [negative]", func() {

		By("Define deployments")
		depa := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage, tsparams.TestTargetLabels)

		depa = deployment.RedefineWithVolume(depa, "test")

		depb := deployment.DefineDeployment("lifecycle-dpb", tsparams.LifecycleNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TestTargetLabels)

		depb = deployment.RedefineWithVolume(depb, "testb")

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(depa, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(depb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pv := persistentvolume.RedefineWithPVReclaimPolicy(
			persistentvolume.DefinePersistentVolume("test", tsparams.LifecycleNamespace), corev1.PersistentVolumeReclaimDelete)

		err = tshelper.CreatePersistentVolume(pv, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvb := persistentvolume.RedefineWithPVReclaimPolicy(
			persistentvolume.DefinePersistentVolume("testb", tsparams.LifecycleNamespace), corev1.PersistentVolumeReclaimRecycle)

		err = tshelper.CreatePersistentVolume(pvb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicy,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicy, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

})

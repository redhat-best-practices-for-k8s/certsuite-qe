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
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolume"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolumeclaim"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
)

var (
	// Each tc will save the PV that was created in order to delete them.
	pvNames = []string{}
)

var _ = Describe("lifecycle-persistent-volume-reclaim-policy", Serial, func() {
	var randomNamespace string
	var randomPV string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		randomPV = tsparams.TestPVName + "-" + globalhelper.GenerateRandomString(10)
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
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)

		By("Delete all PVs that were created by the previous test case.")
		for _, pv := range pvNames {
			By("Deleting pv " + pv)
			err := tshelper.DeletePV(pv, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())
		}

		// clear the list.
		pvNames = []string{}
	})

	// 54201
	It("One deployment, one pod with a volume that uses a reclaim policy of delete", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, persistentVolume.Name)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, randomNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(dep, persistentVolume.Name, pvc.Name)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54202
	It("One pod with a volume that uses a reclaim policy of delete", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, persistentVolume.Name)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, randomNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		pod.RedefineWithPVC(put, persistentVolume.Name, pvc.Name)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54203
	It("One replicaSet with a volume that uses a reclaim policy of delete", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, persistentVolume.Name)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, randomNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define replicaSet")
		rs := replicaset.DefineReplicaSet(tsparams.TestReplicaSetName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)
		replicaset.RedefineWithPVC(rs, persistentVolume.Name, pvc.Name)

		err = globalhelper.CreateAndWaitUntilReplicaSetIsReady(rs, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54204
	It("One deployment, one pod with a volume that uses a reclaim policy of retain [negative]", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimRetain)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, persistentVolume.Name)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, randomNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(dep, persistentVolume.Name, pvc.Name)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54206
	It("One pod with a volume that uses a reclaim policy of recycle [negative]", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimRecycle)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, persistentVolume.Name)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, randomNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		pod.RedefineWithPVC(put, persistentVolume.Name, pvc.Name)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54207
	It("Two deployments, one with reclaim policy of delete, other with recycle [negative]", func() {

		By("Define and create first pv")
		pvca := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolumea := persistentvolume.DefinePersistentVolume(randomPV, pvca.Name, pvca.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolumea, corev1.PersistentVolumeReclaimDelete)
		err := tshelper.CreatePersistentVolume(persistentVolumea, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, persistentVolumea.Name)

		By("Define and create second pv")
		pvcb := persistentvolumeclaim.DefinePersistentVolumeClaim("lifecycle-pvcb", randomNamespace)
		persistentVolumeb := persistentvolume.DefinePersistentVolume("lifecycle-pvb-"+globalhelper.GenerateRandomString(10),
			pvcb.Name, pvcb.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolumeb, corev1.PersistentVolumeReclaimRecycle)
		err = tshelper.CreatePersistentVolume(persistentVolumeb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, persistentVolumeb.Name)

		By("create first pvc")

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvca, randomNamespace, tsparams.WaitingTime, persistentVolumea.Name)
		Expect(err).ToNot(HaveOccurred())

		By("create second pvc")

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvcb, randomNamespace, tsparams.WaitingTime, persistentVolumeb.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployments")
		depa := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(depa, persistentVolumea.Name, pvca.Name)

		depb := deployment.DefineDeployment("lifecycle-dpb", randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(depb, persistentVolumeb.Name, pvcb.Name)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), depa, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(), depb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

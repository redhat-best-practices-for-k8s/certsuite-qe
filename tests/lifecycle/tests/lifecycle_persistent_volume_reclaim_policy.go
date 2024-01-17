package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolume"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolumeclaim"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
)

var _ = Describe("lifecycle-persistent-volume-reclaim-policy", Serial, func() {
	var randomNamespace string
	var randomPV string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		randomPV = tsparams.TestPVName + "-" + globalhelper.GenerateRandomString(10)
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.WaitingTime)
	})

	// 54201
	It("One deployment, one pod with a volume that uses a reclaim policy of delete", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		By("Create persistent volume")
		err := globalhelper.CreatePersistentVolume(persistentVolume)
		Expect(err).ToNot(HaveOccurred())

		By("Create persistent volume claim")
		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete persistent volume claim")
			err = globalhelper.DeletePersistentVolumeClaim(pvc)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume")
			err = globalhelper.DeletePersistentVolume(persistentVolume.Name, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(dep, persistentVolume.Name, pvc.Name)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has persistent volume claim configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName).To(Equal(pvc.Name))

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 54202
	It("One pod with a volume that uses a reclaim policy of delete", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := globalhelper.CreatePersistentVolume(persistentVolume)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete persistent volume claim")
			err = globalhelper.DeletePersistentVolumeClaim(pvc)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume")
			err = globalhelper.DeletePersistentVolume(persistentVolume.Name, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Define pod")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		pod.RedefineWithPVC(put, persistentVolume.Name, pvc.Name)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54203
	It("One replicaSet with a volume that uses a reclaim policy of delete", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := globalhelper.CreatePersistentVolume(persistentVolume)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete persistent volume claim")
			err = globalhelper.DeletePersistentVolumeClaim(pvc)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume")
			err = globalhelper.DeletePersistentVolume(persistentVolume.Name, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())
		})

		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define replicaSet")
		rs := replicaset.DefineReplicaSet(tsparams.TestReplicaSetName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)
		replicaset.RedefineWithPVC(rs, persistentVolume.Name, pvc.Name)

		err = globalhelper.CreateAndWaitUntilReplicaSetIsReady(rs, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54204
	It("One deployment, one pod with a volume that uses a reclaim policy of retain [negative]", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimRetain)

		err := globalhelper.CreatePersistentVolume(persistentVolume)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete persistent volume claim")
			err = globalhelper.DeletePersistentVolumeClaim(pvc)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume")
			err = globalhelper.DeletePersistentVolume(persistentVolume.Name, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())
		})

		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(dep, persistentVolume.Name, pvc.Name)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has persistent volume claim configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName).To(Equal(pvc.Name))

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54206
	It("One pod with a volume that uses a reclaim policy of recycle [negative]", func() {
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimRecycle)

		err := globalhelper.CreatePersistentVolume(persistentVolume)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete persistent volume claim")
			err = globalhelper.DeletePersistentVolumeClaim(pvc)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume")
			err = globalhelper.DeletePersistentVolume(persistentVolume.Name, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())
		})

		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		pod.RedefineWithPVC(put, persistentVolume.Name, pvc.Name)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 54207
	It("Two deployments, one with reclaim policy of delete, other with recycle [negative]", func() {

		By("Define and create first pv")
		pvca := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentVolumea := persistentvolume.DefinePersistentVolume(randomPV, pvca.Name, pvca.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolumea, corev1.PersistentVolumeReclaimDelete)
		err := globalhelper.CreatePersistentVolume(persistentVolumea)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create second pv")
		pvcb := persistentvolumeclaim.DefinePersistentVolumeClaim("lifecycle-pvcb", randomNamespace)
		persistentVolumeb := persistentvolume.DefinePersistentVolume("lifecycle-pvb-"+globalhelper.GenerateRandomString(10),
			pvcb.Name, pvcb.Namespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolumeb, corev1.PersistentVolumeReclaimRecycle)
		err = globalhelper.CreatePersistentVolume(persistentVolumeb)
		Expect(err).ToNot(HaveOccurred())

		By("create first pvc")
		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvca, tsparams.WaitingTime, persistentVolumea.Name)
		Expect(err).ToNot(HaveOccurred())

		By("create second pvc")
		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvcb, tsparams.WaitingTime, persistentVolumeb.Name)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete persistent volume claim")
			err = globalhelper.DeletePersistentVolumeClaim(pvca)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume")
			err = globalhelper.DeletePersistentVolume(persistentVolumea.Name, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume claim")
			err = globalhelper.DeletePersistentVolumeClaim(pvcb)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume")
			err = globalhelper.DeletePersistentVolume(persistentVolumeb.Name, tsparams.WaitingTime)
		})

		By("Define deployments")
		depa := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(depa, persistentVolumea.Name, pvca.Name)

		depb := deployment.DefineDeployment("lifecycle-dpb", randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(depb, persistentVolumeb.Name, pvcb.Name)

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(depa, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has persistent volume claim configured")
		runningDeployment, err := globalhelper.GetRunningDeployment(depa.Namespace, depa.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName).To(Equal(pvca.Name))

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(depb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has persistent volume claim configured")
		runningDeployment, err = globalhelper.GetRunningDeployment(depb.Namespace, depb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName).To(Equal(pvcb.Name))

		By("Start lifecycle-persistent-volume-reclaim-policy test")
		err = globalhelper.LaunchTests(tsparams.TnfPersistentVolumeReclaimPolicyTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPersistentVolumeReclaimPolicyTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

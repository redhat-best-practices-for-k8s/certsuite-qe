package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolume"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolumeclaim"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
)

var (
	// Each tc will save the PV that was created in order to delete them.
	pvNames = []string{}
)

var _ = Describe("lifecycle-persistent-volume-reclaim-policy", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file"))
	}

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all nodes are labeled with 'worker-cnf' label")
		err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().CoreV1Interface, configSuite.General.CnfNodeLabel)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test in order to enable PVs deletion.")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

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

		persistentVolume := persistentvolume.DefinePersistentVolume(tsparams.TestPVName, tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, tsparams.TestPVName)

		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, tsparams.LifecycleNamespace)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.LifecycleNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(dep, tsparams.TestVolumeName, tsparams.TestPVCName)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
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

		persistentVolume := persistentvolume.DefinePersistentVolume(tsparams.TestPVName, tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, tsparams.TestPVName)

		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, tsparams.LifecycleNamespace)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.LifecycleNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		pod.RedefineWithPVC(put, tsparams.TestVolumeName, tsparams.TestPVCName)

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

		persistentVolume := persistentvolume.DefinePersistentVolume(tsparams.TestPVName, tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, tsparams.TestPVName)

		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, tsparams.LifecycleNamespace)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.LifecycleNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define replicaSet")
		rs := replicaset.DefineReplicaSet(tsparams.TestReplicaSetName, tsparams.LifecycleNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)
		replicaset.RedefineWithPVC(rs, tsparams.TestVolumeName, tsparams.TestPVCName)

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

		persistentVolume := persistentvolume.DefinePersistentVolume(tsparams.TestPVName, tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimRetain)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, tsparams.TestPVName)

		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, tsparams.LifecycleNamespace)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.LifecycleNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(dep, tsparams.TestVolumeName, tsparams.TestPVCName)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
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

		persistentVolume := persistentvolume.DefinePersistentVolume(tsparams.TestPVName, tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimRecycle)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, tsparams.TestPVName)

		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, tsparams.LifecycleNamespace)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.LifecycleNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.LifecycleNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)
		pod.RedefineWithPVC(put, tsparams.TestVolumeName, tsparams.TestPVCName)

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
		persistentVolumea := persistentvolume.DefinePersistentVolume(tsparams.TestPVName, tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolumea, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolumea, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, tsparams.TestPVName)

		By("Define and create second pv")
		persistentVolumeb := persistentvolume.DefinePersistentVolume("lifecycle-pvb", tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolumeb, corev1.PersistentVolumeReclaimRecycle)

		err = tshelper.CreatePersistentVolume(persistentVolumeb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, "lifecycle-pvb")

		By("Define and create first pvc")
		pvca := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, tsparams.LifecycleNamespace)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvca, tsparams.LifecycleNamespace, tsparams.WaitingTime, persistentVolumea.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create second pvc")
		pvcb := persistentvolumeclaim.DefinePersistentVolumeClaim("lifecycle-pvcb", tsparams.LifecycleNamespace)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvcb, tsparams.LifecycleNamespace, tsparams.WaitingTime, persistentVolumeb.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployments")
		depa := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.LifecycleNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(depa, tsparams.TestVolumeName, tsparams.TestPVCName)

		depb := deployment.DefineDeployment("lifecycle-dpb", tsparams.LifecycleNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels)

		deployment.RedefineWithPVC(depb, tsparams.TestVolumeName, "lifecycle-pvcb")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(depa, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(depb, tsparams.WaitingTime)
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

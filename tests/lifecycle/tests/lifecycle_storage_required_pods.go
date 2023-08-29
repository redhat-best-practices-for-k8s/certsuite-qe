package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolume"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolumeclaim"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("lifecycle-storage-required-pods", func() {
	var randomNamespace string
	var randomStorageClassName string
	var randomPV string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		randomStorageClassName = tsparams.TestLocalStorageClassName + "-" + globalhelper.GenerateRandomString(10)
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

		By(fmt.Sprintf("Create %s storageclass", randomStorageClassName))
		err = tshelper.CreateStorageClass(randomStorageClassName, false)
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

		By(fmt.Sprintf("Remove %s storageclass", randomStorageClassName))
		err := tshelper.DeleteStorageClass(randomStorageClassName)
		Expect(err).ToNot(HaveOccurred())

		// clear the list.
		pvNames = []string{}
	})

	It("One pod with a storage, PVC with no storageclass defined", func() {
		By("Define PV")
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, persistentVolume.Name)

		By("Define PVC")
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, randomNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod with a pvc")
		testPod := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)

		pod.RedefineWithPVC(testPod, persistentVolume.Name, pvc.Name)
		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start storage-required-pods test")
		err = globalhelper.LaunchTests(tsparams.TnfStorageRequiredPods,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStorageRequiredPods, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with local storage, PVC with storageclass defined", func() {
		By("Define PV")
		testPv := persistentvolume.DefinePersistentVolume(randomPV)
		persistentvolume.RedefineWithPVReclaimPolicy(testPv, corev1.PersistentVolumeReclaimDelete)
		persistentvolume.RedefineWithStorageClass(testPv, randomStorageClassName)

		err := tshelper.CreatePersistentVolume(testPv, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, testPv.Name)

		By("Define PVC")
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)
		persistentvolumeclaim.RedefineWithStorageClass(pvc, randomStorageClassName)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, randomNamespace, tsparams.WaitingTime, testPv.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod with a pvc")
		testPod := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)

		pod.RedefineWithPVC(testPod, testPv.Name, pvc.Name)
		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start storage-required-pods test")
		err = globalhelper.LaunchTests(tsparams.TnfStorageRequiredPods,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStorageRequiredPods, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

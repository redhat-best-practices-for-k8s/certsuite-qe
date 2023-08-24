package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolume"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/persistentvolumeclaim"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-storage-required-pods", Serial, func() {
	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Create local-storage storageclass")
		err = tshelper.CreateStorageClass("local-storage")
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

		By("Delete local-storage storageclass")
		err = tshelper.DeleteStorageClass("local-storage")
		Expect(err).ToNot(HaveOccurred())

		// clear the list.
		pvNames = []string{}
	})

	It("One pod with a storage, PVC with no storageclass defined", func() {
		By("Define PV")
		persistentVolume := persistentvolume.DefinePersistentVolume(tsparams.TestPVName, tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(persistentVolume, corev1.PersistentVolumeReclaimDelete)

		err := tshelper.CreatePersistentVolume(persistentVolume, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, tsparams.TestPVName)

		By("Define PVC")
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, tsparams.LifecycleNamespace)
		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.LifecycleNamespace, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod with a pvc")
		testPod := tshelper.DefinePod(tsparams.TestPodName)

		pod.RedefineWithPVC(testPod, tsparams.TestPVCName, tsparams.TestPVCName)
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
		testPv := persistentvolume.DefinePersistentVolume(tsparams.TestPVName, tsparams.LifecycleNamespace)
		persistentvolume.RedefineWithPVReclaimPolicy(testPv, corev1.PersistentVolumeReclaimDelete)
		persistentvolume.RedefineWithStorageClass(testPv, tsparams.TestLocalStorageClassName)

		err := tshelper.CreatePersistentVolume(testPv, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		pvNames = append(pvNames, tsparams.TestPVName)

		By("Define PVC")
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, tsparams.LifecycleNamespace)
		persistentvolumeclaim.RedefineWithStorageClass(pvc, tsparams.TestLocalStorageClassName)

		err = tshelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.LifecycleNamespace, tsparams.WaitingTime, testPv.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod with a pvc")
		testPod := tshelper.DefinePod(tsparams.TestPodName)

		pod.RedefineWithPVC(testPod, tsparams.TestPVCName, tsparams.TestPVCName)
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

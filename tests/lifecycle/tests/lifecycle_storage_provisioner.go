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

var _ = Describe("lifecycle-storage-provisioner", func() {
	var randomNamespace string
	var randomStorageClassName string
	var randomPV string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		randomStorageClassName = tsparams.TestLocalStorageClassName + "-" + globalhelper.GenerateRandomString(10)
		randomPV = tsparams.TestPVName + "-" + globalhelper.GenerateRandomString(10)

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

		By(fmt.Sprintf("Create %s storageclass", randomStorageClassName))
		err = globalhelper.CreateStorageClass(randomStorageClassName, false)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)

		By(fmt.Sprintf("Remove %s storageclass", randomStorageClassName))
		err := globalhelper.DeleteStorageClass(randomStorageClassName)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with a storage, PVC with no storageclass defined", func() {
		By("Define PVC")
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)

		By("Define PV")
		persistentVolume := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)

		By("Redefine PV with reclaim policy")
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

		By("Create PVC and wait until it is bound")
		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.WaitingTime, persistentVolume.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod with a pvc")
		testPod := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)

		pod.RedefineWithPVC(testPod, persistentVolume.Name, pvc.Name)

		By("Create pod and wait until it is ready")
		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start storage-provisioner test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteStorageProvisioner,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.GetNumberOfNodes(globalhelper.GetAPIClient().K8sClient.CoreV1()) == 1 {
			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteStorageProvisioner, globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		} else {
			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteStorageProvisioner, globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	It("One pod with local storage, PVC with storageclass defined", func() {
		By("Define PVC")
		pvc := persistentvolumeclaim.DefinePersistentVolumeClaim(tsparams.TestPVCName, randomNamespace)

		By("Redefine PVC with storageclass")
		persistentvolumeclaim.RedefineWithStorageClass(pvc, randomStorageClassName)

		By("Define PV")
		testPv := persistentvolume.DefinePersistentVolume(randomPV, pvc.Name, pvc.Namespace)

		By("Redefine PV with reclaim policy and storageclass")
		persistentvolume.RedefineWithPVReclaimPolicy(testPv, corev1.PersistentVolumeReclaimDelete)
		persistentvolume.RedefineWithStorageClass(testPv, randomStorageClassName)

		By("Create PV")
		err := globalhelper.CreatePersistentVolume(testPv)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete persistent volume claim")
			err = globalhelper.DeletePersistentVolumeClaim(pvc)
			Expect(err).ToNot(HaveOccurred())

			By("Delete persistent volume")
			err = globalhelper.DeletePersistentVolume(testPv.Name, tsparams.WaitingTime)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Create PVC and wait until it is bound")
		err = globalhelper.CreateAndWaitUntilPVCIsBound(pvc, tsparams.WaitingTime, testPv.Name)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod with a pvc")
		testPod := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)

		pod.RedefineWithPVC(testPod, testPv.Name, pvc.Name)
		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start storage-provisioner test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteStorageProvisioner,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.GetNumberOfNodes(globalhelper.GetAPIClient().K8sClient.CoreV1()) == 1 {
			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteStorageProvisioner, globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		} else {
			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteStorageProvisioner, globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		}
	})
})

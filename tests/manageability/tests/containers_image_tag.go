package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/manageability/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	corev1 "k8s.io/api/core/v1"
)

const (
	sampleWorkloadImage = "quay.io/redhat-best-practices-for-k8s/certsuite-sample-workload"
)

var _ = Describe("manageability-containers-image-tag", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.ManageabilityNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	It("One pod with valid image tag", func() {
		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.TestImageWithValidTag, tsparams.CertsuiteTargetPodLabels)

		By("Create pod")
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is running and has containers")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, testPod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")
		Expect(runningPod.Spec.Containers[0].Image).To(ContainSubstring(":"))

		By("Assert all containers are ready")
		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Start containers-image-tag test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerImageTag,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerImageTag,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with invalid image tag", func() {

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			sampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		By("Create and wait until pod is ready")
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is running and has containers")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, testPod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")

		By("Assert all containers are ready")
		for _, cs := range runningPod.Status.ContainerStatuses {
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Start containers-image-tag test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerImageTag,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerImageTag,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

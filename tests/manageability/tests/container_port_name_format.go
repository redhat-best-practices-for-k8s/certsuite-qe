package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/manageability/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/manageability/parameters"
)

var _ = Describe("manageability-container-port-name", func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

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

	It("One pod with valid port name", func() {
		By("Define pod")
		testPod := tshelper.DefineManageabilityPod(tsparams.TestPodName, randomNamespace,
			tsparams.TestImageWithValidTag, tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is ready with container port configured")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, testPod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())
		Expect(len(runningPod.Spec.Containers[0].Ports)).To(BeNumerically(">", 0))

		By("Start container-port-name test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerPortName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerPortName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with invalid port name", func() {
		By("Define pod")
		testPod := tshelper.DefineManageabilityPod(tsparams.TestPodName, randomNamespace,
			tsparams.TestImageWithValidTag, tsparams.CertsuiteTargetPodLabels)

		tshelper.RedefinePodWithContainerPort(testPod, 0, tsparams.InvalidPortName)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is ready with container port configured")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, testPod.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())
		Expect(len(runningPod.Spec.Containers[0].Ports)).To(BeNumerically(">", 0))

		By("Start container-port-name test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteContainerPortName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteContainerPortName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

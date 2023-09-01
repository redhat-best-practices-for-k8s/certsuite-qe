package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
)

var _ = Describe("Access-control ssh-daemons,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})

	It("one pod with no ssh running", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start ssh-daemons")
		err = globalhelper.LaunchTests(
			tsparams.TnfNoSSHDaemonsAllowed,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNoSSHDaemonsAllowed,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with ssh daemon running", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		err := pod.RedefineWithContainerExecCommand(testPod, tsparams.SSHDaemonStartContainerCommand, 0)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start ssh-daemons")
		err = globalhelper.LaunchTests(
			tsparams.TnfNoSSHDaemonsAllowed,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))

		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNoSSHDaemonsAllowed,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

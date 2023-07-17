package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
)

var _ = Describe("Access-control ssh-daemons,", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with no ssh running", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.Configuration.General.TestImage, tsparams.TestDeploymentLabels)

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

	FIt("one pod with ssh daemon running", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.Configuration.General.TestImage, tsparams.TestDeploymentLabels)

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

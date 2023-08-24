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

var _ = Describe("Access-control pod-service-account,", Serial, func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with valid service account", func() {

		By("Create service account")
		err := globalhelper.CreateServiceAccount(tsparams.TestServiceAccount, tsparams.TestAccessControlNameSpace)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod with service account")
		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, tsparams.TestServiceAccount)

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-service-account")
		err = globalhelper.LaunchTests(
			tsparams.TnfPodServiceAccount,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfPodServiceAccount,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with empty service account", func() {

		By("Define pod with empty service account")

		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, "")
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-service-account")
		err = globalhelper.LaunchTests(
			tsparams.TnfPodServiceAccount,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfPodServiceAccount,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with default service account", func() {

		By("Define pod with empty service account")

		testPod := pod.DefinePod(tsparams.TestPodName, tsparams.TestAccessControlNameSpace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

		pod.RedefineWithServiceAccount(testPod, "default")
		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start pod-service-account")
		err = globalhelper.LaunchTests(
			tsparams.TnfPodServiceAccount,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfPodServiceAccount,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

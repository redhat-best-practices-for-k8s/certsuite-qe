package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/manageability/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/manageability/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("manageability-container-port-name", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.ManageabilityNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.ManageabilityNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with valid port name", func() {

		By("Define pod")
		testPod := tshelper.DefineManageabilityPod(tsparams.TestPodName, tsparams.ManageabilityNamespace,
			tsparams.TestImageWithValidTag, tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start container-port-name test")
		err = globalhelper.LaunchTests(tsparams.TnfContainerPortName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerPortName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One pod with invalid port name", func() {

		By("Define pod")
		testPod := tshelper.DefineManageabilityPod(tsparams.TestPodName, tsparams.ManageabilityNamespace,
			tsparams.TestImageWithValidTag, tsparams.TnfTargetPodLabels)

		tshelper.RedefinePodWithContainerPort(testPod, 0, tsparams.InvalidPortName)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start container-port-name test")
		err = globalhelper.LaunchTests(tsparams.TnfContainerPortName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfContainerPortName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

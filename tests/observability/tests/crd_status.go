package tests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/observability/observabilityhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/observability/observabilityparameters"
)

var (
	// Each TC save the CRDs that has created so they can be automatically
	// removed before the next TC starts.
	crdNames = []string{}
)

var _ = Describe(observabilityparameters.TnfCrdStatusTcName, func() {
	const tnfTestCaseName = observabilityparameters.TnfCrdStatusTcName
	qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

	AfterEach(func() {
		By("Removing all CRDs created by previous test case.")
		for _, crd := range crdNames {
			By("Removing CRD " + crd)
			globalhelper.DeleteCrdAndWaitUntilIsRemoved(crd, 10*time.Second)
		}
		// Clear list.
		crdNames = []string{}
	})

	// Positive #1
	It("One CRD created with status subresource", func() {

		By("Create CRD in the cluster with suffix " + observabilityparameters.CrdSuffix1)
		crd1 := observabilityhelper.DefineCrdWithStatusSubresource("TestCrd", observabilityparameters.CrdSuffix1)

		err := globalhelper.CreateAndWaitUntilCrdIsReady(crd1, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Positive #2
	It("Two CRDs created, both with status subresource", func() {

		By("Create CRD in the cluster with suffix " + observabilityparameters.CrdSuffix1)
		crd1 := observabilityhelper.DefineCrdWithStatusSubresource("TestCrdOne", observabilityparameters.CrdSuffix1)

		err := globalhelper.CreateAndWaitUntilCrdIsReady(crd1, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create CRD in the cluster with suffix " + observabilityparameters.CrdSuffix2)
		crd2 := observabilityhelper.DefineCrdWithStatusSubresource("TestCrdTwo", observabilityparameters.CrdSuffix2)

		err = globalhelper.CreateAndWaitUntilCrdIsReady(crd2, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #1
	It("One CRD created without status subresource [negative]", func() {

		By("Create CRD in the cluster with suffix " + observabilityparameters.CrdSuffix1)
		crd1 := observabilityhelper.DefineCrdWithoutStatusSubresource("TestCrd", observabilityparameters.CrdSuffix1)

		err := globalhelper.CreateAndWaitUntilCrdIsReady(crd1, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// Negative #2
	It("Two CRDs created, one with and the other without status subresource [negative]", func() {

		By("Create CRD in the cluster with suffix " + observabilityparameters.CrdSuffix1)
		crd1 := observabilityhelper.DefineCrdWithStatusSubresource("TestCrdOne", observabilityparameters.CrdSuffix1)

		err := globalhelper.CreateAndWaitUntilCrdIsReady(crd1, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create CRD in the cluster with suffix " + observabilityparameters.CrdSuffix2)
		crd2 := observabilityhelper.DefineCrdWithoutStatusSubresource("TestCrdTwo", observabilityparameters.CrdSuffix2)

		err = globalhelper.CreateAndWaitUntilCrdIsReady(crd2, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Two CRDs created, both without status subresource [negative]", func() {

		By("Create CRD in the cluster with suffix " + observabilityparameters.CrdSuffix1)
		crd1 := observabilityhelper.DefineCrdWithoutStatusSubresource("TestCrdOne", observabilityparameters.CrdSuffix1)

		err := globalhelper.CreateAndWaitUntilCrdIsReady(crd1, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create CRD in the cluster with suffix " + observabilityparameters.CrdSuffix2)
		crd2 := observabilityhelper.DefineCrdWithoutStatusSubresource("TestCrdTwo", observabilityparameters.CrdSuffix2)

		err = globalhelper.CreateAndWaitUntilCrdIsReady(crd2, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One CRD deployed not having any of the configured suffixes [skip]", func() {

		By("Create CRD in the cluster with suffix " + observabilityparameters.NotConfiguredCrdSuffix)
		crd1 := observabilityhelper.DefineCrdWithoutStatusSubresource("TestCrdOne",
			observabilityparameters.NotConfiguredCrdSuffix)

		err := globalhelper.CreateAndWaitUntilCrdIsReady(crd1, observabilityparameters.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})

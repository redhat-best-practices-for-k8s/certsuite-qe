package tests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/observability/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"
)

var (
	// Each TC save the CRDs that has created so they can be automatically
	// removed before the next TC starts.
	crdNames = []string{}
)

var _ = Describe(tsparams.TnfCrdStatusTcName, Serial, func() {
	AfterEach(func() {
		By("Removing all CRDs created by previous test case.")
		for _, crd := range crdNames {
			By("Removing CRD " + crd)
			tshelper.DeleteCrdAndWaitUntilIsRemoved(crd, 10*time.Second)
		}
		// Clear list.
		crdNames = []string{}
	})

	// 52444
	It("One CRD created with status subresource", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create CRD in the cluster with suffix " + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithStatusSubresource("TestCrd", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start TNF " + tsparams.TnfCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfCrdStatusTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdStatusTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52445
	It("Two CRDs created, both with status subresource", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create CRD in the cluster with suffix " + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithStatusSubresource("TestCrdOne", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create CRD in the cluster with suffix " + tsparams.CrdSuffix2)
		crd2 := tshelper.DefineCrdWithStatusSubresource("TestCrdTwo", tsparams.CrdSuffix2)

		err = tshelper.CreateAndWaitUntilCrdIsReady(crd2, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start TNF " + tsparams.TnfCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfCrdStatusTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdStatusTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52446
	It("One CRD created without status subresource [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create CRD in the cluster with suffix " + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithoutStatusSubresource("TestCrd", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start TNF " + tsparams.TnfCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfCrdStatusTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdStatusTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52447
	It("Two CRDs created, one with and the other without status subresource [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create CRD in the cluster with suffix " + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithStatusSubresource("TestCrdOne", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create CRD in the cluster with suffix " + tsparams.CrdSuffix2)
		crd2 := tshelper.DefineCrdWithoutStatusSubresource("TestCrdTwo", tsparams.CrdSuffix2)

		err = tshelper.CreateAndWaitUntilCrdIsReady(crd2, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start TNF " + tsparams.TnfCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfCrdStatusTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdStatusTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52448
	It("Two CRDs created, both without status subresource [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create CRD in the cluster with suffix " + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithoutStatusSubresource("TestCrdOne", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create CRD in the cluster with suffix " + tsparams.CrdSuffix2)
		crd2 := tshelper.DefineCrdWithoutStatusSubresource("TestCrdTwo", tsparams.CrdSuffix2)

		err = tshelper.CreateAndWaitUntilCrdIsReady(crd2, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start TNF " + tsparams.TnfCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfCrdStatusTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdStatusTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52449
	It("One CRD deployed not having any of the configured suffixes [skip]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create CRD in the cluster with suffix " + tsparams.NotConfiguredCrdSuffix)
		crd1 := tshelper.DefineCrdWithoutStatusSubresource("TestCrdOne",
			tsparams.NotConfiguredCrdSuffix)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start TNF " + tsparams.TnfCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfCrdStatusTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdStatusTcName, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})

package tests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/observability/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/observability/parameters"
)

var (
	// Each TC save the CRDs that has created so they can be automatically
	// removed before the next TC starts.
	crdNames = []string{}
)

const (
	CreateCRDInClusterStr = "Create CRD in the cluster with suffix "
)

var _ = Describe(tsparams.CertsuiteCrdStatusTcName, Serial, func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			tshelper.GetCertsuiteTargetPodLabelsSlice(),
			[]string{},
			[]string{},
			[]string{tsparams.CrdSuffix1, tsparams.CrdSuffix2}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Removing all CRDs created by previous test case.")

		for _, crd := range crdNames {
			By("Removing CRD " + crd)
			tshelper.DeleteCrdAndWaitUntilIsRemoved(crd, 10*time.Second)
		}
		// Clear list.
		crdNames = []string{}

		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.CrdDeployTimeoutMins)
	})

	// 52444
	It("One CRD created with status subresource", func() {
		By(CreateCRDInClusterStr + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithStatusSubresource("TestCrd", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start Certsuite " + tsparams.CertsuiteCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdStatusTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdStatusTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52445
	It("Two CRDs created, both with status subresource", func() {
		By(CreateCRDInClusterStr + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithStatusSubresource("TestCrdOne", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By(CreateCRDInClusterStr + tsparams.CrdSuffix2)
		crd2 := tshelper.DefineCrdWithStatusSubresource("TestCrdTwo", tsparams.CrdSuffix2)

		err = tshelper.CreateAndWaitUntilCrdIsReady(crd2, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start Certsuite " + tsparams.CertsuiteCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdStatusTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdStatusTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52446
	It("One CRD created without status subresource [negative]", func() {
		By(CreateCRDInClusterStr + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithoutStatusSubresource("TestCrd", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start Certsuite " + tsparams.CertsuiteCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdStatusTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdStatusTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52447
	It("Two CRDs created, one with and the other without status subresource [negative]", func() {
		By(CreateCRDInClusterStr + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithStatusSubresource("TestCrdOne", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By(CreateCRDInClusterStr + tsparams.CrdSuffix2)
		crd2 := tshelper.DefineCrdWithoutStatusSubresource("TestCrdTwo", tsparams.CrdSuffix2)

		err = tshelper.CreateAndWaitUntilCrdIsReady(crd2, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start Certsuite " + tsparams.CertsuiteCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdStatusTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdStatusTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52448
	It("Two CRDs created, both without status subresource [negative]", func() {
		By(CreateCRDInClusterStr + tsparams.CrdSuffix1)
		crd1 := tshelper.DefineCrdWithoutStatusSubresource("TestCrdOne", tsparams.CrdSuffix1)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By(CreateCRDInClusterStr + tsparams.CrdSuffix2)
		crd2 := tshelper.DefineCrdWithoutStatusSubresource("TestCrdTwo", tsparams.CrdSuffix2)

		err = tshelper.CreateAndWaitUntilCrdIsReady(crd2, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRDs to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)
		crdNames = append(crdNames, crd2.Name)

		By("Start Certsuite " + tsparams.CertsuiteCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdStatusTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdStatusTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 52449
	It("One CRD deployed not having any of the configured suffixes [skip]", func() {
		By(CreateCRDInClusterStr + tsparams.NotConfiguredCrdSuffix)
		crd1 := tshelper.DefineCrdWithoutStatusSubresource("TestCrdOne",
			tsparams.NotConfiguredCrdSuffix)

		err := tshelper.CreateAndWaitUntilCrdIsReady(crd1, tsparams.CrdDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		// Save CRD to be removed after the TC has finished.
		crdNames = append(crdNames, crd1.Name)

		By("Start Certsuite " + tsparams.CertsuiteCrdStatusTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdStatusTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdStatusTcName, globalparameters.TestCaseSkipped,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

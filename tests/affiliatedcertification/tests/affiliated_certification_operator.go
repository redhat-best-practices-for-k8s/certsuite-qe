package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
)

const (
	ErrorLabelingOperatorStr = "Error labeling operator "
)

// Shared operator namespace - matches the suite-level deployment.
var (
	sharedOperatorNamespace = tsparams.TestCertificationNameSpace + "-operators"
	grafanaOperatorName     = "grafana-operator" // Will be set by operator discovery
)

var _ = Describe("Affiliated-certification operator certification,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestCertificationNameSpace)

		// If Kind cluster, skip.
		if globalhelper.IsKindCluster() {
			Skip("This test is not supported on Kind cluster")
		}

		// 🚀 PERFORMANCE OPTIMIZATION: Use pre-deployed shared operators
		// The operators are now deployed once at suite level in SynchronizedBeforeSuite()
		// instead of being deployed fresh for each test in BeforeEach()

		By("Configure environment for operator tests (using shared operators)")
		preConfigureAffiliatedCertificationEnvironment(randomNamespace, randomCertsuiteConfigDir)

		By("✅ Using suite-level shared operators - no individual operator deployment needed!")
		// Note: Operators are already deployed and ready in the shared namespace:
		// - cockroachdb operator (uncertified)
		// - cockroachdb-certified operator (certified)
		// - grafana operator (community)
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 46699
	It("one operator to test, operator is not in certified-operators organization [negative]",
		func() {
			By("Label shared cockroachdb operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					tsparams.UncertifiedOperatorPrefixCockroach,
					sharedOperatorNamespace, // Use shared namespace
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixCockroach)

			// Update config to point to shared operator namespace for this test
			By("Update certsuite config to test shared operators")
			err := globalhelper.DefineCertsuiteConfig(
				[]string{sharedOperatorNamespace}, // Point to shared operators
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error updating certsuite config for shared operators")

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				tsparams.TestCaseOperatorAffiliatedCertName)

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

	// 46700
	It("one operator to test, operator is in certified-operators organization",
		func() {
			By("Label shared cockroachdb-certified operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					tsparams.CertifiedOperatorPrefixCockroachCertified,
					sharedOperatorNamespace, // Use shared namespace
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixCockroachCertified)

			// Update config to point to shared operator namespace for this test
			By("Update certsuite config to test shared operators")
			err := globalhelper.DefineCertsuiteConfig(
				[]string{sharedOperatorNamespace}, // Point to shared operators
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error updating certsuite config for shared operators")

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				tsparams.TestCaseOperatorAffiliatedCertName)

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

	// 46701
	It("two operators to test, both are in certified-operators organization",
		func() {
			By("Label shared cockroachdb-certified operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					tsparams.CertifiedOperatorPrefixCockroachCertified,
					sharedOperatorNamespace, // Use shared namespace
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixCockroachCertified)

			By("Label shared grafana operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					grafanaOperatorName,     // Use suite-level variable
					sharedOperatorNamespace, // Use shared namespace
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+grafanaOperatorName)

			// Update config to point to shared operator namespace for this test
			By("Update certsuite config to test shared operators")
			err := globalhelper.DefineCertsuiteConfig(
				[]string{sharedOperatorNamespace}, // Point to shared operators
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error updating certsuite config for shared operators")

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				tsparams.TestCaseOperatorAffiliatedCertName)

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

	// 46703
	It("three operators to test, two are in certified-operators and one is not [negative]",
		func() {
			By("Label shared cockroachdb operator (uncertified) to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					tsparams.UncertifiedOperatorPrefixCockroach,
					sharedOperatorNamespace, // Use shared namespace
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixCockroach)

			By("Label shared cockroachdb-certified operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					tsparams.CertifiedOperatorPrefixCockroachCertified,
					sharedOperatorNamespace, // Use shared namespace
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixCockroachCertified)

			By("Label shared grafana operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					grafanaOperatorName,     // Use suite-level variable
					sharedOperatorNamespace, // Use shared namespace
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+grafanaOperatorName)

			// Update config to point to shared operator namespace for this test
			By("Update certsuite config to test shared operators")
			err := globalhelper.DefineCertsuiteConfig(
				[]string{sharedOperatorNamespace}, // Point to shared operators
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error updating certsuite config for shared operators")

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				tsparams.TestCaseOperatorAffiliatedCertName)

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})
})

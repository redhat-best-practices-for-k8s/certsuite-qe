package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
)

var _ = Describe("Affiliated-certification operator certification,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string
	var grafanaOperatorName string
	var certifiedOperatorName string
	var certifiedOperatorCSVPrefix string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestCertificationNameSpace)

		// If Kind cluster, skip.
		if globalhelper.IsKindCluster() {
			Skip("This test is not supported on Kind cluster")
		}

		preConfigureAffiliatedCertificationEnvironment(randomNamespace, randomCertsuiteConfigDir)

		By("Deploy cockroachdb for testing")
		// cockroachdb: not in certified-operators group in catalog, for negative test cases
		err := tshelper.DeployOperatorSubscription(
			"cockroachdb",
			"cockroachdb",
			"stable-v6.x",
			randomNamespace,
			tsparams.CommunityOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.UncertifiedOperatorPrefixCockroach)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixCockroach,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.UncertifiedOperatorPrefixCockroach+
			" is not ready")

		By("Query the packagemanifest for a certified operator - trying cockroachdb-certified first")
		// Try to find cockroachdb-certified first (for OCP 4.19 and earlier)
		// If not found, look for other certified operators (for OCP 4.20+)
		var catalogSource string
		certifiedOperatorName, catalogSource, err = globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			"cockroachdb-certified",
			randomNamespace,
		)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for cockroachdb-certified")

		// If cockroachdb-certified is not available, try alternative certified operators
		if certifiedOperatorName == "not found" {
			By("cockroachdb-certified not found, searching for alternative certified operators")
			// Try common certified operators that should be available across versions
			alternativeCertifiedOperators := []string{
				"mongodb-enterprise",
				"redis-enterprise",
				"crunchy-postgres",
			}

			for _, altOperator := range alternativeCertifiedOperators {
				certifiedOperatorName, catalogSource, err = globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
					altOperator,
					randomNamespace,
				)
				if err == nil && certifiedOperatorName != "not found" {
					By(fmt.Sprintf("Found alternative certified operator: %s", certifiedOperatorName))

					break
				}
			}
		}

		// If still not found, fail with helpful message
		if certifiedOperatorName == "not found" {
			Skip("No suitable certified operators found in this OCP version - skipping affiliated certification tests")
		}

		By(fmt.Sprintf("Using certified operator: %s from catalog: %s", certifiedOperatorName, catalogSource))

		By("Query the packagemanifest for available channel, version and CSV for " + certifiedOperatorName)
		var channel, version, csvName string
		channel, version, csvName = globalhelper.CheckOperatorChannelAndVersionOrFail(certifiedOperatorName, randomNamespace)

		// Extract the CSV prefix (without version) for use in labels and test assertions
		// CSV names are typically in format "operator-name.vX.Y.Z"
		certifiedOperatorCSVPrefix = certifiedOperatorName

		By(fmt.Sprintf("Deploy %s operator (channel %s, version %s) for testing", certifiedOperatorName, channel, version))
		// certified operator: in certified-operators group and version is certified
		err = tshelper.DeployOperatorSubscription(
			certifiedOperatorName,
			certifiedOperatorName,
			channel,
			randomNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+certifiedOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(certifiedOperatorCSVPrefix,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+" is not ready")

		By("Query the packagemanifest for Grafana operator package name and catalog source")
		grafanaOperatorName, catalogSource = globalhelper.CheckOperatorExistsOrFail("grafana", randomNamespace)

		By("Query the packagemanifest for available channel, version and CSV for " + grafanaOperatorName)
		channel, version, csvName = globalhelper.CheckOperatorChannelAndVersionOrFail(grafanaOperatorName, randomNamespace)

		By(fmt.Sprintf("Deploy Grafana operator (channel %s, version %s) for testing", channel, version))
		// grafana-operator: in community-operators group
		err = tshelper.DeployOperatorSubscription(
			grafanaOperatorName,
			grafanaOperatorName,
			channel,
			randomNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			grafanaOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(grafanaOperatorName,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 46699
	It("one operator to test, operator is not in certified-operators organization [negative]",
		func() {
			By("Label operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					tsparams.UncertifiedOperatorPrefixCockroach,
					randomNamespace,
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixCockroach)

			By("Assert operator CSV is ready")
			csv, err := tshelper.GetCsvByPrefix(tsparams.UncertifiedOperatorPrefixCockroach, randomNamespace)
			Expect(err).ToNot(HaveOccurred())
			Expect(csv).ToNot(BeNil())

			// Assert that the random report dir exists
			Expect(randomReportDir).To(BeADirectory(), "Random report dir does not exist")

			// Assert that the random certsuite config dir exists
			Expect(randomCertsuiteConfigDir).To(BeADirectory(), "Random certsuite config dir does not exist")

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				tsparams.TestCaseOperatorAffiliatedCertName+" test")

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseOperatorAffiliatedCertName,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
		})

	// 46697
	It("two operators to test, one is in certified-operators organization and its version is certified,"+
		" one is not in certified-operators organization [negative]", func() {
		By("Label operators to be certified")

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				certifiedOperatorCSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+certifiedOperatorCSVPrefix)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixCockroach,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixCockroach)

		By("Assert both operator CSVs are ready")
		certifiedCSV, err := tshelper.GetCsvByPrefix(certifiedOperatorCSVPrefix, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(certifiedCSV).ToNot(BeNil())
		uncertifiedCSV, err := tshelper.GetCsvByPrefix(tsparams.UncertifiedOperatorPrefixCockroach, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(uncertifiedCSV).ToNot(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46582
	It("one operator to test, operator is in certified-operators organization"+
		" and its version is certified", func() {
		By("Label operator to be certified")

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				grafanaOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+grafanaOperatorName)

		By("Assert operator CSV is ready")
		csv, err := tshelper.GetCsvByPrefix(grafanaOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(csv).ToNot(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46696
	It("two operators to test, both are in certified-operators organization and their"+
		" versions are certified", func() {
		By("Label operators to be certified")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				grafanaOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+grafanaOperatorName)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				certifiedOperatorCSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+certifiedOperatorCSVPrefix)

		By("Assert both operator CSVs are ready")
		grafanaCSV, err := tshelper.GetCsvByPrefix(grafanaOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(grafanaCSV).ToNot(BeNil())
		certifiedCSV, err := tshelper.GetCsvByPrefix(certifiedOperatorCSVPrefix, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(certifiedCSV).ToNot(BeNil())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46698
	It("no operators are labeled for testing [negative]", func() {
		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

})

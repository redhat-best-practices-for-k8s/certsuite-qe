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

		By("Query the packagemanifest for the cockroachdb-certified operator default channel")
		channel, err := globalhelper.QueryPackageManifestForDefaultChannel(
			"cockroachdb-certified",
			randomNamespace,
		)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for cockroachdb-certified")
		Expect(channel).ToNot(Equal("not found"), "Channel not found")

		By("Query the packagemanifest for the cockroachdb-certified operator")
		version, err := globalhelper.QueryPackageManifestForVersion("cockroachdb-certified", randomNamespace, channel)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for cockroachdb-certified")
		Expect(version).ToNot(Equal("not found"), "Version not found")

		By(fmt.Sprintf("Deploy cockroachdb-certified operator %s for testing", "v"+version))
		// cockroachdb-certified operator: in certified-operators group and version is certified
		err = tshelper.DeployOperatorSubscription(
			"cockroachdb-certified",
			"cockroachdb-certified",
			channel,
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorPrefixCockroachCertified+".v"+version,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.CertifiedOperatorPrefixCockroachCertified)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixCockroachCertified,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixCockroachCertified+".v"+version+
			" is not ready")

		By("Query the packagemanifest for Grafana operator package name and catalog source")
		var catalogSource string
		grafanaOperatorName, catalogSource = globalhelper.CheckOperatorExistsOrFail("grafana", randomNamespace)

		By("Query the packagemanifest for available channel, version and CSV for " + grafanaOperatorName)
		var csvName string
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
				tsparams.CertifiedOperatorPrefixCockroachCertified,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixCockroachCertified)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixCockroach,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixCockroach)

		By("Assert both operator CSVs are ready")
		certifiedCSV, err := tshelper.GetCsvByPrefix(tsparams.CertifiedOperatorPrefixCockroachCertified, randomNamespace)
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
				tsparams.CertifiedOperatorPrefixCockroachCertified,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixCockroachCertified)

		By("Assert both operator CSVs are ready")
		grafanaCSV, err := tshelper.GetCsvByPrefix(grafanaOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(grafanaCSV).ToNot(BeNil())
		cockroachCSV, err := tshelper.GetCsvByPrefix(tsparams.CertifiedOperatorPrefixCockroachCertified, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(cockroachCSV).ToNot(BeNil())

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

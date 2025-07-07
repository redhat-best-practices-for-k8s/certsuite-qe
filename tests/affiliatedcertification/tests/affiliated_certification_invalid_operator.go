package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
)

var _ = Describe("Affiliated-certification invalid operator certification,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string
	var grafanaOperatorName string

	BeforeEach(func() {
		if globalhelper.IsKindCluster() {
			Skip("Skipping operator certification tests on kind cluster")
		}

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestCertificationNameSpace)

		preConfigureAffiliatedCertificationEnvironment(
			randomNamespace,
			randomCertsuiteConfigDir,
		)

		By("Query the packagemanifest for Grafana operator package name and catalog source")
		var catalogSource string
		var err error
		grafanaOperatorName, catalogSource, err = globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			"grafana", randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for Grafana operator")
		Expect(grafanaOperatorName).ToNot(Equal("not found"), "Grafana operator package not found")
		Expect(catalogSource).ToNot(Equal("not found"), "Grafana operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + grafanaOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			grafanaOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+grafanaOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

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

		err = waitUntilOperatorIsReady(grafanaOperatorName,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")

		// sriov-fec.v1.1.0 operator : in certified-operators group, version is not certified
		By("Deploy alternate operator catalog source")
		err = globalhelper.DisableCatalogSource(tsparams.CertifiedOperatorGroup)
		Expect(err).ToNot(HaveOccurred(), "Error disabling "+
			tsparams.CertifiedOperatorGroup+" catalog source")
		Eventually(func() bool {
			stillEnabled, err := globalhelper.IsCatalogSourceEnabled(
				tsparams.CertifiedOperatorGroup,
				tsparams.OperatorSourceNamespace,
				tsparams.CertifiedOperatorDisplayName)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

			return !stillEnabled
		}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
			"Default catalog source is still enabled")

		// Deploying certified operator with invalid catalog version is necessary in order to cover negative scenarios
		err = globalhelper.DeployRHCertifiedOperatorSource("4.7")
		Expect(err).ToNot(HaveOccurred(), "Error deploying catalog source")

		By("Deploy sriov-fec operator with uncertified version")
		err = tshelper.DeployOperatorSubscription(
			tsparams.UncertifiedOperatorPrefixSriov,
			tsparams.UncertifiedOperatorPrefixSriov,
			"stable",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.UncertifiedOperatorFullSriov,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.UncertifiedOperatorPrefixSriov)

		approveInstallPlanWhenReady(tsparams.UncertifiedOperatorFullSriov,
			randomNamespace)

		err = waitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixSriov,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.UncertifiedOperatorPrefixSriov+
			" is not ready")

		By("Re-enable default catalog source")
		err = globalhelper.DeleteCatalogSource(tsparams.CertifiedOperatorGroup,
			randomNamespace,
			"redhat-certified")
		Expect(err).ToNot(HaveOccurred(), "Error removing alternate catalog source")

		err = globalhelper.EnableCatalogSource(tsparams.CertifiedOperatorGroup)
		Expect(err).ToNot(HaveOccurred(), "Error enabling default catalog source")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 46695
	It("one operator to test, operator is in certified-operators organization but its version"+
		" is not certified [negative]", func() {

		if globalhelper.IsKindCluster() {
			Skip("Skip on kind cluster")
		}

		By("Label operator to be certified")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixSriov,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixSriov)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46700
	It("two operators to test, both are in certified-operators organization,"+
		" one’s version is certified, the other’s is not [negative]", func() {

		if globalhelper.IsKindCluster() {
			Skip("Skip on kind cluster")
		}

		By("Label operators to be certified")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixSriov,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixSriov)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				grafanaOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+grafanaOperatorName)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})
})

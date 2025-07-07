package operator

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
)

var _ = Describe("Operator crd-versioning,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.OperatorNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.CertsuiteTargetOperatorLabels},
			[]string{},
			tsparams.CertsuiteTargetCrdFilters, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("operator crd has valid versioning", func() {
		By("Query the packagemanifest for Grafana operator package name and catalog source")
		grafanaOperatorName, catalogSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
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

		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				grafanaOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+grafanaOperatorName)

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorCrdVersioning,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorCrdVersioning,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

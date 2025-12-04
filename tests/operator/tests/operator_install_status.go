package operator

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
)

var _ = Describe("Operator install-source,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string
	var operatorName string
	var catalogSource string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.OperatorNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			tsparams.CertsuiteTargetCrdFilters, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Query the packagemanifest for grafana operator package name and catalog source")
		operatorName, catalogSource = globalhelper.CheckOperatorExistsOrFail("grafana", randomNamespace)

		By("Query the packagemanifest for available channel, version and CSV for " + operatorName)
		channel, _, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(operatorName, randomNamespace)

		By("Deploy grafana operator for testing")
		err = tshelper.DeployOperatorSubscription(
			operatorName,
			operatorName,
			channel,
			randomNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			operatorName)

		err = tshelper.WaitUntilOperatorIsReady(operatorName,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("one operator that reports Succeeded as its installation status", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		By("Assert operator CSV is in Succeeded phase")
		csv, err := tshelper.GetCsvByPrefix(operatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(csv.Status.Phase).To(Equal(v1alpha1.CSVPhaseSucceeded))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatus,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("two operators, one does not reports Succeeded as its installation status (quick failure) [negative]", func() {
		// TODO: Known issue with OCP 4.20 certified-operators. Fix later.
		if !globalhelper.IsKindCluster() {
			if ocpVersion, err := globalhelper.GetClusterVersion(); err == nil && strings.HasPrefix(ocpVersion, "4.20") {
				Skip("TODO: Known issue with OCP 4.20 certified-operators. Fix later.")
			}
		}

		By("Query the packagemanifest for postgresql operator package name and catalog source")
		postgresOperatorName, catalogSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			tsparams.OperatorPackageNamePrefixLightweight, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for postgresql operator")
		Expect(postgresOperatorName).ToNot(Equal("not found"), "PostgreSQL operator package not found")
		Expect(catalogSource).ToNot(Equal("not found"), "PostgreSQL operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + postgresOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			postgresOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+postgresOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		By("Deploy postgresql operator for testing")
		// Deploy PostgreSQL operator with nodeSelector that will cause quick failure
		nodeSelector := map[string]string{"target": "nonexistent-node"}
		err = tshelper.DeployOperatorSubscriptionWithNodeSelector(
			postgresOperatorName,
			channel,
			randomNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
			nodeSelector,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			postgresOperatorName)

		// Do not wait for the PostgreSQL operator to be ready - it should fail due to nodeSelector

		By("Verify that PostgreSQL operator CSV is not in Succeeded phase")

		Eventually(func() bool {
			isNotSucceeded, err := tshelper.IsCSVNotSucceeded(tsparams.OperatorPrefixLightweight, randomNamespace)
			if err != nil {
				fmt.Printf("Error checking CSV status for %s: %v\n", tsparams.OperatorPrefixLightweight, err)

				return false
			}
			fmt.Printf("PostgreSQL operator %s CSV status is not Succeeded: %t\n", tsparams.OperatorPrefixLightweight, isNotSucceeded)

			return isNotSucceeded
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Equal(true),
			"PostgreSQL operator CSV should not be in Succeeded phase for this negative test")

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixLightweight,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+postgresOperatorName)

		By("Update certsuite config to include both operators")
		err = globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			tsparams.CertsuiteTargetCrdFilters, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Assert grafana operator CSV is in Succeeded phase")
		csv, err := tshelper.GetCsvByPrefix(operatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(csv.Status.Phase).To(Equal(v1alpha1.CSVPhaseSucceeded))

		By("Assert PostgreSQL operator CSV is not in Succeeded phase")
		postgresCSV, err := tshelper.GetCsvByPrefix(tsparams.OperatorPrefixLightweight, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(postgresCSV.Status.Phase).ToNot(Equal(v1alpha1.CSVPhaseSucceeded))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatus,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("two operators, one does not reports Succeeded as its installation status (delayed failure) [negative]", Serial, func() {
		// TODO: Known issue with OCP 4.20 certified-operators. Fix later.
		if !globalhelper.IsKindCluster() {
			if ocpVersion, err := globalhelper.GetClusterVersion(); err == nil && strings.HasPrefix(ocpVersion, "4.20") {
				Skip("TODO: Known issue with OCP 4.20 certified-operators. Fix later.")
			}
		}

		By("Query the packagemanifest for postgresql operator package name and catalog source")
		postgresqlOperatorName, catalogSource2, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			tsparams.OperatorPackageNamePrefixLightweight, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for postgresql operator")
		Expect(postgresqlOperatorName).ToNot(Equal("not found"), "postgresql operator package not found")
		Expect(catalogSource2).ToNot(Equal("not found"), "postgresql operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + postgresqlOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			postgresqlOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+postgresqlOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		By("Deploy postgresql operator for testing")
		// The postgresql operator fails to deploy, which creates a delayed failure scenario
		// This allows testing of the CNF Certification Suite timeout mechanism
		// for operator readiness.
		nodeSelector := map[string]string{"target": "none"}
		err = tshelper.DeployOperatorSubscriptionWithNodeSelector(
			postgresqlOperatorName,
			channel,
			randomNamespace,
			catalogSource2,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
			nodeSelector,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			postgresqlOperatorName)

		// Do not wait until the operator is ready. This time the CNF Certification suite must handle the situation.

		By("Verify that postgresql operator CSV is not in Succeeded phase")
		Eventually(func() bool {
			isNotSucceeded, err := tshelper.IsCSVNotSucceeded(tsparams.OperatorPrefixLightweight, randomNamespace)
			if err != nil {
				fmt.Printf("Error checking CSV status for %s: %v\n", tsparams.OperatorPrefixLightweight, err)

				return false
			}
			fmt.Printf("postgresql operator %s CSV status is not Succeeded: %t\n", tsparams.OperatorPrefixLightweight, isNotSucceeded)

			return isNotSucceeded
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Equal(true),
			"postgresql operator CSV should not be in Succeeded phase for this negative test")

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixLightweight,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixLightweight,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+postgresqlOperatorName)

		By("Assert PostgreSQL operator CSV is not in Succeeded phase")
		postgresCSV, err := tshelper.GetCsvByPrefix(tsparams.OperatorPrefixLightweight, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(postgresCSV.Status.Phase).ToNot(Equal(v1alpha1.CSVPhaseSucceeded))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatus,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

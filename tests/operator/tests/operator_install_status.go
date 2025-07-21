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

		By("Start test")
		err := globalhelper.LaunchTests(
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
		By("Query the packagemanifest for postgresql operator package name and catalog source")
		postgresOperatorName, catalogSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			"cloud-native-postgresql", randomNamespace)
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
			isNotSucceeded, err := tshelper.IsCSVNotSucceeded(postgresOperatorName, randomNamespace)
			if err != nil {
				fmt.Printf("Error checking CSV status for %s: %v\n", postgresOperatorName, err)

				return false
			}
			fmt.Printf("PostgreSQL operator %s CSV status is not Succeeded: %t\n", postgresOperatorName, isNotSucceeded)

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
				postgresOperatorName,
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
		By("Query the packagemanifest for Jaeger operator package name and catalog source")
		jaegerOperatorName, catalogSource2, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			"jaeger", randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for Jaeger operator")
		Expect(jaegerOperatorName).ToNot(Equal("not found"), "Jaeger operator package not found")
		Expect(catalogSource2).ToNot(Equal("not found"), "Jaeger operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + jaegerOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			jaegerOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+jaegerOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		By("Deploy Jaeger operator for testing")
		// The jaeger operator fails to deploy, which creates a delayed failure scenario
		// This allows testing of the CNF Certification Suite timeout mechanism
		// for operator readiness.
		nodeSelector := map[string]string{"target": "none"}
		err = tshelper.DeployOperatorSubscriptionWithNodeSelector(
			jaegerOperatorName,
			channel,
			randomNamespace,
			catalogSource2,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
			nodeSelector,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			jaegerOperatorName)

		// Do not wait until the operator is ready. This time the CNF Certification suite must handle the situation.

		By("Verify that Jaeger operator CSV is not in Succeeded phase")
		Eventually(func() bool {
			isNotSucceeded, err := tshelper.IsCSVNotSucceeded(jaegerOperatorName, randomNamespace)
			if err != nil {
				fmt.Printf("Error checking CSV status for %s: %v\n", jaegerOperatorName, err)

				return false
			}
			fmt.Printf("Jaeger operator %s CSV status is not Succeeded: %t\n", jaegerOperatorName, isNotSucceeded)

			return isNotSucceeded
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Equal(true),
			"Jaeger operator CSV should not be in Succeeded phase for this negative test")

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
				jaegerOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+jaegerOperatorName)

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

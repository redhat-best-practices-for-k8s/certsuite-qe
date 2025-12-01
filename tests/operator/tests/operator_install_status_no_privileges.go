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

var _ = Describe("Operator install-status-no-privileges,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string
	var operatorName string
	var catalogSource string
	var lightweightOperatorName string
	var lightweightOperatorCSVPrefix string

	// findAvailableLightweightOperator searches for an operator with NO clusterPermissions
	// This is specifically for testing operators without cluster-level privileges
	findAvailableLightweightOperator := func(namespace string) (packageName, catalogSource, csvPrefix string) {
		// List of operators known to have NO clusterPermissions
		// These are namespace-scoped operators
		operatorsWithoutClusterPermissions := []string{
			"postgresql",       // Original choice, good for OCP 4.19 and earlier - has no clusterPermissions
			"mongodb-operator", // Namespace-scoped MongoDB operator
			"mariadb-operator", // Namespace-scoped MariaDB operator
		}

		for _, operatorPrefix := range operatorsWithoutClusterPermissions {
			opName, catSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
				operatorPrefix,
				namespace,
			)
			if err == nil && opName != "not found" && catSource != "not found" {
				By(fmt.Sprintf("Found available operator without clusterPermissions: %s from catalog: %s", opName, catSource))

				return opName, catSource, operatorPrefix
			}
		}

		Skip("No suitable operators without clusterPermissions found in this OCP version - " +
			"skipping operator install status no privileges tests")

		return "", "", ""
	}

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

		// grafana operator has clusterPermissions but no resourceNames
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

		// Find a lightweight operator that typically has no/minimal clusterPermissions
		By("Find an available lightweight operator for testing")
		var lightweightCatalogSource string
		lightweightOperatorName, lightweightCatalogSource, lightweightOperatorCSVPrefix = findAvailableLightweightOperator(randomNamespace)

		By(fmt.Sprintf("Using lightweight operator: %s from catalog: %s", lightweightOperatorName, lightweightCatalogSource))

		By("Query the packagemanifest for available channel, version and CSV for " + lightweightOperatorName)
		channel2, version2, csvName2, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			lightweightOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+lightweightOperatorName)
		Expect(channel2).ToNot(Equal("not found"), "Channel not found")
		Expect(version2).ToNot(Equal("not found"), "Version not found")
		Expect(csvName2).ToNot(Equal("not found"), "CSV name not found")

		By(fmt.Sprintf("Deploy %s operator (channel %s, version %s) for testing", lightweightOperatorName, channel2, version2))
		err = tshelper.DeployOperatorSubscription(
			lightweightOperatorName,
			lightweightOperatorName,
			channel2,
			randomNamespace,
			lightweightCatalogSource,
			tsparams.OperatorSourceNamespace,
			csvName2,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			lightweightOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(lightweightOperatorCSVPrefix,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName2+
			" is not ready")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 66381
	It("one operator with no clusterPermissions", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				lightweightOperatorCSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+lightweightOperatorCSVPrefix)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66383
	It("one operator with clusterPermissions [negative]", func() {
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66384
	It("two operators, one with no clusterPermissions and one with clusterPermissions", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				lightweightOperatorCSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+lightweightOperatorCSVPrefix)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

})

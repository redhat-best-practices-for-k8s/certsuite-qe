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

	// findAvailableLightweightOperator searches for a lightweight operator available in the cluster
	// It tries a list of commonly available operators and returns the first one found
	findAvailableLightweightOperator := func(namespace string) (packageName, catalogSource, csvPrefix string) {
		// List of lightweight operators to try, in order of preference
		lightweightOperators := []string{
			"postgresql",                  // Original choice, good for OCP 4.19 and earlier
			"redis-operator",              // Alternative for newer versions
			"etcd",                        // Usually available across versions
			"prometheus",                  // Common monitoring operator
			"node-observability-operator", // Lightweight observability operator
		}

		for _, operatorPrefix := range lightweightOperators {
			opName, catSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
				operatorPrefix,
				namespace,
			)
			if err == nil && opName != "not found" && catSource != "not found" {
				By(fmt.Sprintf("Found available lightweight operator: %s from catalog: %s", opName, catSource))

				return opName, catSource, operatorPrefix
			}
		}

		Skip("No suitable lightweight operators found in this OCP version - skipping operator install status tests")

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
		By("Find an available lightweight operator for testing")
		lightweightOperatorName, lightweightCatalogSource, operatorCSVPrefix := findAvailableLightweightOperator(randomNamespace)

		By(fmt.Sprintf("Using lightweight operator: %s from catalog: %s", lightweightOperatorName, lightweightCatalogSource))

		By("Query the packagemanifest for available channel, version and CSV for " + lightweightOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			lightweightOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+lightweightOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		By(fmt.Sprintf("Deploy %s operator for testing", lightweightOperatorName))
		// Deploy lightweight operator with nodeSelector that will cause quick failure
		nodeSelector := map[string]string{"target": "nonexistent-node"}
		err = tshelper.DeployOperatorSubscriptionWithNodeSelector(
			lightweightOperatorName,
			channel,
			randomNamespace,
			lightweightCatalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
			nodeSelector,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			lightweightOperatorName)

		// Do not wait for the lightweight operator to be ready - it should fail due to nodeSelector

		By(fmt.Sprintf("Verify that %s operator CSV is not in Succeeded phase", lightweightOperatorName))

		Eventually(func() bool {
			isNotSucceeded, err := tshelper.IsCSVNotSucceeded(operatorCSVPrefix, randomNamespace)
			if err != nil {
				fmt.Printf("Error checking CSV status for %s: %v\n", operatorCSVPrefix, err)

				return false
			}
			fmt.Printf("%s operator %s CSV status is not Succeeded: %t\n", lightweightOperatorName, operatorCSVPrefix, isNotSucceeded)

			return isNotSucceeded
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Equal(true),
			fmt.Sprintf("%s operator CSV should not be in Succeeded phase for this negative test", lightweightOperatorName))

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
				operatorCSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+lightweightOperatorName)

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

		By(fmt.Sprintf("Assert %s operator CSV is not in Succeeded phase", lightweightOperatorName))
		lightweightCSV, err := tshelper.GetCsvByPrefix(operatorCSVPrefix, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(lightweightCSV.Status.Phase).ToNot(Equal(v1alpha1.CSVPhaseSucceeded))

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
		By("Find an available lightweight operator for testing")
		lightweightOperatorName, lightweightCatalogSource, operatorCSVPrefix := findAvailableLightweightOperator(randomNamespace)

		By(fmt.Sprintf("Using lightweight operator: %s from catalog: %s", lightweightOperatorName, lightweightCatalogSource))

		By("Query the packagemanifest for available channel, version and CSV for " + lightweightOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			lightweightOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+lightweightOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		By(fmt.Sprintf("Deploy %s operator for testing", lightweightOperatorName))
		// The lightweight operator fails to deploy, which creates a delayed failure scenario
		// This allows testing of the CNF Certification Suite timeout mechanism
		// for operator readiness.
		nodeSelector := map[string]string{"target": "none"}
		err = tshelper.DeployOperatorSubscriptionWithNodeSelector(
			lightweightOperatorName,
			channel,
			randomNamespace,
			lightweightCatalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
			nodeSelector,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			lightweightOperatorName)

		// Do not wait until the operator is ready. This time the CNF Certification suite must handle the situation.

		By(fmt.Sprintf("Verify that %s operator CSV is not in Succeeded phase", lightweightOperatorName))
		Eventually(func() bool {
			isNotSucceeded, err := tshelper.IsCSVNotSucceeded(operatorCSVPrefix, randomNamespace)
			if err != nil {
				fmt.Printf("Error checking CSV status for %s: %v\n", operatorCSVPrefix, err)

				return false
			}
			fmt.Printf("%s operator %s CSV status is not Succeeded: %t\n", lightweightOperatorName, operatorCSVPrefix, isNotSucceeded)

			return isNotSucceeded
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Equal(true),
			fmt.Sprintf("%s operator CSV should not be in Succeeded phase for this negative test", lightweightOperatorName))

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorCSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorCSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+lightweightOperatorName)

		By(fmt.Sprintf("Assert %s operator CSV is not in Succeeded phase", lightweightOperatorName))
		lightweightCSV, err := tshelper.GetCsvByPrefix(operatorCSVPrefix, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(lightweightCSV.Status.Phase).ToNot(Equal(v1alpha1.CSVPhaseSucceeded))

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

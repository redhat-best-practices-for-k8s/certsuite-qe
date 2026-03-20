package tests

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/operatorversions"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
)

// waitForOperatorReadyOrSkip waits for an operator to be ready and skips the test on timeout.
func waitForOperatorReadyOrSkip(csvPrefix, namespace, displayName string) {
	err := tshelper.WaitUntilOperatorIsReady(csvPrefix, namespace)
	if err != nil {
		if strings.Contains(err.Error(), "timed out") {
			Skip(fmt.Sprintf("Operator %s failed to become ready: %v", displayName, err))
		}

		Expect(err).ToNot(HaveOccurred(), "Operator "+displayName+" is not ready")
	}
}

// deployUncertifiedOperator deploys the uncertified operator and waits for it to be ready.
// Skips the test if the operator times out.
func deployUncertifiedOperator(operatorInfo operatorversions.OperatorInfo, namespace string) {
	By("Deploy uncertified operator for testing: " + operatorInfo.PackageName)

	err := tshelper.DeployOperatorSubscription(
		operatorInfo.PackageName,
		operatorInfo.PackageName,
		operatorInfo.Channel,
		namespace,
		operatorInfo.CatalogSource,
		tsparams.OperatorSourceNamespace,
		"",
		v1alpha1.ApprovalAutomatic,
	)
	Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+operatorInfo.PackageName)

	waitForOperatorReadyOrSkip(operatorInfo.CSVPrefix, namespace, operatorInfo.PackageName)
}

// deployCertifiedOperator queries the package manifest and deploys the certified operator.
// Skips the test if the operator times out.
func deployCertifiedOperator(operatorInfo operatorversions.OperatorInfo, namespace string) {
	By("Query the packagemanifest for the certified operator: " + operatorInfo.PackageName)

	channel, err := globalhelper.QueryPackageManifestForDefaultChannel(
		operatorInfo.PackageName,
		namespace,
	)
	Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+operatorInfo.PackageName)
	Expect(channel).ToNot(Equal("not found"), "Channel not found")

	By("Query the packagemanifest for the certified operator version")

	version, err := globalhelper.QueryPackageManifestForVersion(operatorInfo.PackageName, namespace, channel)
	Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+operatorInfo.PackageName)
	Expect(version).ToNot(Equal("not found"), "Version not found")

	By(fmt.Sprintf("Deploy certified operator %s v%s for testing", operatorInfo.PackageName, version))

	err = tshelper.DeployOperatorSubscription(
		operatorInfo.PackageName,
		operatorInfo.PackageName,
		channel,
		namespace,
		operatorInfo.CatalogSource,
		tsparams.OperatorSourceNamespace,
		operatorInfo.CSVPrefix+".v"+version,
		v1alpha1.ApprovalAutomatic,
	)
	Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+operatorInfo.PackageName)

	waitForOperatorReadyOrSkip(operatorInfo.CSVPrefix, namespace, operatorInfo.PackageName)
}

// deployGrafanaOperator queries the package manifest and deploys the grafana operator.
// Skips the test if the operator times out. Returns the operator name.
func deployGrafanaOperator(namespace string) string {
	By("Query the packagemanifest for Grafana operator package name and catalog source")

	operatorName, catalogSource := globalhelper.CheckOperatorExistsOrFail("grafana", namespace)

	By("Query the packagemanifest for available channel, version and CSV for " + operatorName)

	channel, version, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(operatorName, namespace)

	By(fmt.Sprintf("Deploy Grafana operator (channel %s, version %s) for testing", channel, version))

	err := tshelper.DeployOperatorSubscription(
		operatorName,
		operatorName,
		channel,
		namespace,
		catalogSource,
		tsparams.OperatorSourceNamespace,
		csvName,
		v1alpha1.ApprovalAutomatic,
	)
	Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+operatorName)

	waitForOperatorReadyOrSkip(operatorName, namespace, operatorName)

	return operatorName
}

var _ = Describe("Affiliated-certification operator certification,", Serial, Label("affiliatedcertification", "ocp-required"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
		certifiedOperator        operatorversions.OperatorInfo
		uncertifiedOperator      operatorversions.OperatorInfo
	)

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestCertificationNameSpace)

		// If Kind cluster, skip.
		if globalhelper.IsKindCluster() {
			Skip("This test is not supported on Kind cluster")
		}

		// Get the OCP version and select the appropriate certified operator
		ocpVersion, err := globalhelper.GetClusterVersion()
		Expect(err).ToNot(HaveOccurred(), "Error getting cluster version")
		certifiedOperator = operatorversions.GetCertifiedOperator(ocpVersion)
		uncertifiedOperator = operatorversions.GetUncertifiedOperator(ocpVersion)
		By(fmt.Sprintf("Using certified operator for OCP %s: %s", ocpVersion, certifiedOperator.String()))
		By(fmt.Sprintf("Using uncertified operator for OCP %s: %s", ocpVersion, uncertifiedOperator.String()))

		preConfigureAffiliatedCertificationEnvironment(randomNamespace, randomCertsuiteConfigDir)
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 46699
	It("one operator to test, operator is not in certified-operators organization [negative]",
		func() {
			deployUncertifiedOperator(uncertifiedOperator, randomNamespace)

			By("Label operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					uncertifiedOperator.CSVPrefix,
					randomNamespace,
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				ErrorLabelingOperatorStr+uncertifiedOperator.CSVPrefix)

			By("Assert operator CSV is ready")
			csv, err := tshelper.GetCsvByPrefix(uncertifiedOperator.CSVPrefix, randomNamespace)
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
		deployCertifiedOperator(certifiedOperator, randomNamespace)
		deployUncertifiedOperator(uncertifiedOperator, randomNamespace)

		By("Label operators to be certified")

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				certifiedOperator.CSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+certifiedOperator.CSVPrefix)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				uncertifiedOperator.CSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+uncertifiedOperator.CSVPrefix)

		By("Assert both operator CSVs are ready")
		certifiedCSV, err := tshelper.GetCsvByPrefix(certifiedOperator.CSVPrefix, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(certifiedCSV).ToNot(BeNil())
		uncertifiedCSV, err := tshelper.GetCsvByPrefix(uncertifiedOperator.CSVPrefix, randomNamespace)
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
		grafanaOperatorName := deployGrafanaOperator(randomNamespace)

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
		grafanaOperatorName := deployGrafanaOperator(randomNamespace)
		deployCertifiedOperator(certifiedOperator, randomNamespace)

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
				certifiedOperator.CSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+certifiedOperator.CSVPrefix)

		By("Assert both operator CSVs are ready")
		grafanaCSV, err := tshelper.GetCsvByPrefix(grafanaOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(grafanaCSV).ToNot(BeNil())
		certifiedCSV, err := tshelper.GetCsvByPrefix(certifiedOperator.CSVPrefix, randomNamespace)
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
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})
})

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
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/operatorversions"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
	ErrorRemovingLabelStr    = "Error removing label from operator "
)

var _ = Describe("Operator install-source,", Serial, Label("operator", "ocp-required"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

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
			[]string{"nginxingresses.charts.nginx.org"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		// Install 3 separate operators for testing
		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("deploy cluster-wide cluster-logging operator", func() {
		const (
			clusterLoggingOperatorName = "cluster-logging"
		)
		openshiftLoggingNamespace := randomNamespace

		By("Create openshift-logging namespace")
		err := globalhelper.CreateNamespace(openshiftLoggingNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create fake operator group for cluster-logging operator")
		err = tshelper.DeployTestOperatorGroup(openshiftLoggingNamespace, true)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Query the packagemanifest for defaultChannel for " + clusterLoggingOperatorName)
		channel, err := globalhelper.QueryPackageManifestForDefaultChannel(clusterLoggingOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+clusterLoggingOperatorName)

		fmt.Printf("CHANNEL FOUND: %s\n", channel)

		By("Query the packagemanifest for the " + clusterLoggingOperatorName)
		version, err := globalhelper.QueryPackageManifestForVersion(clusterLoggingOperatorName, randomNamespace, channel)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+clusterLoggingOperatorName)

		fmt.Printf("VERSION FOUND: %s\n", version)

		By("Deploy cluster-logging operator for testing")
		err = tshelper.DeployOperatorSubscription(
			clusterLoggingOperatorName,
			clusterLoggingOperatorName,
			channel,
			openshiftLoggingNamespace,
			tsparams.RedhatOperatorGroup,
			tsparams.OperatorSourceNamespace,
			clusterLoggingOperatorName+".v"+version,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+clusterLoggingOperatorName)

		By("Wait until operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(clusterLoggingOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+clusterLoggingOperatorName+" is not ready")

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				clusterLoggingOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+clusterLoggingOperatorName)

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66142
	It("one operator installed with OLM", func() {
		By("Query the packagemanifest for Grafana operator package name and catalog source")
		grafanaOperatorName, catalogSource := globalhelper.CheckOperatorExistsOrFail("grafana", randomNamespace)

		By("Query the packagemanifest for available channel, version and CSV for " + grafanaOperatorName)
		channel, version, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(grafanaOperatorName, randomNamespace)

		By(fmt.Sprintf("Deploy Grafana operator (channel %s, version %s) for testing", channel, version))
		// grafana-operator: in community-operators group
		err := tshelper.DeployOperatorSubscription(
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
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66143
	It("one operator not installed with OLM [negative]", func() {
		// Note: This test uses a lightweight operator that varies by OCP version
		// See issue #1283 and operatorversions package for operator catalog availability
		ocpVersion := "4.19"

		if !globalhelper.IsKindCluster() {
			var verr error
			ocpVersion, verr = globalhelper.GetClusterVersion()
			Expect(verr).ToNot(HaveOccurred(), "Error getting cluster version")
		}
		lightweightOp := operatorversions.GetLightweightOperator(ocpVersion)

		By(fmt.Sprintf("Query the packagemanifest for %s operator package name and catalog source", lightweightOp.PackageName))
		lightweightOperatorName, lwCatalogSource := globalhelper.CheckOperatorExistsOrFail(lightweightOp.PackageName,
			randomNamespace)

		// Lightweight operator can be deployed in the same namespace (OwnNamespace install mode)
		By(fmt.Sprintf("Deploy %s operator group", lightweightOp.PackageName))
		err := tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group for "+lightweightOp.PackageName+" operator")

		By("Query the packagemanifest for available channel, version and CSV for " + lightweightOperatorName)
		channel, _, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(lightweightOperatorName, randomNamespace)

		By(fmt.Sprintf("Deploy %s operator for testing", lightweightOp.PackageName))
		err = tshelper.DeployOperatorSubscription(
			lightweightOperatorName,
			lightweightOperatorName,
			channel,
			randomNamespace,
			lwCatalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			lightweightOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(lightweightOp.CSVPrefix,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")

		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				lightweightOp.CSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+lightweightOperatorName)

		By("Delete operator's subscription")
		err = globalhelper.DeleteSubscription(randomNamespace,
			lightweightOperatorName+"-subscription")
		Expect(err).ToNot(HaveOccurred())

		By(fmt.Sprintf("Update certsuite config to include %s operator", lightweightOp.PackageName))
		err = globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"nginxingresses.charts.nginx.org"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66144
	It("two operators, both installed with OLM", func() {
		// Note: This test uses grafana-operator and a lightweight operator that varies by OCP version
		// See issue #1283 and operatorversions package for operator catalog availability
		ocpVersion := "4.19"

		if !globalhelper.IsKindCluster() {
			var verr error
			ocpVersion, verr = globalhelper.GetClusterVersion()
			Expect(verr).ToNot(HaveOccurred(), "Error getting cluster version")
		}
		lightweightOp := operatorversions.GetLightweightOperator(ocpVersion)

		By("Query the packagemanifest for Grafana operator package name and catalog source")
		grafanaOperatorName, gCatalogSource := globalhelper.CheckOperatorExistsOrFail("grafana", randomNamespace)

		By("Query the packagemanifest for available channel, version and CSV for " + grafanaOperatorName)
		channel, version, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(grafanaOperatorName, randomNamespace)

		By(fmt.Sprintf("Deploy Grafana operator (channel %s, version %s) for testing", channel, version))
		// grafana-operator: in community-operators group
		err := tshelper.DeployOperatorSubscription(
			grafanaOperatorName,
			grafanaOperatorName,
			channel,
			randomNamespace,
			gCatalogSource,
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

		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				grafanaOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+grafanaOperatorName)

		By(fmt.Sprintf("Query packagemanifest for %s operator (second)", lightweightOp.PackageName))
		lightweightOperatorName, catalogSource2 := globalhelper.CheckOperatorExistsOrFail(lightweightOp.PackageName,
			randomNamespace)

		By("Query the packagemanifest for available channel, version and CSV for " + lightweightOperatorName + " for second operator")
		channel2, version2, csvName2 := globalhelper.CheckOperatorChannelAndVersionOrFail(lightweightOperatorName, randomNamespace)

		By(fmt.Sprintf("Deploy %s operator (channel %s, version %s) for testing", lightweightOp.PackageName, channel2, version2))
		err = tshelper.DeployOperatorSubscription(
			lightweightOperatorName+"-2",
			lightweightOperatorName,
			channel2,
			randomNamespace,
			catalogSource2,
			tsparams.OperatorSourceNamespace,
			csvName2,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			lightweightOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(lightweightOp.CSVPrefix,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName2+
			" is not ready")

		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				lightweightOp.CSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+lightweightOperatorName)

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66145
	It("two operators, one not installed with OLM [negative]", func() {
		// Note: This test uses a lightweight operator that varies by OCP version
		// See issue #1283 and operatorversions package for operator catalog availability
		ocpVersion := "4.19"

		if !globalhelper.IsKindCluster() {
			var verr error
			ocpVersion, verr = globalhelper.GetClusterVersion()
			Expect(verr).ToNot(HaveOccurred(), "Error getting cluster version")
		}
		lightweightOp := operatorversions.GetLightweightOperator(ocpVersion)

		By(fmt.Sprintf("Query the packagemanifest for %s operator package name and catalog source", lightweightOp.PackageName))
		lightweightOperatorName, lwCatalogSource := globalhelper.CheckOperatorExistsOrFail(lightweightOp.PackageName,
			randomNamespace)

		// Lightweight operator can be deployed in the same namespace (OwnNamespace install mode)
		By(fmt.Sprintf("Deploy %s operator group", lightweightOp.PackageName))
		err := tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group for "+lightweightOp.PackageName+" operator")

		By("Query the packagemanifest for available channel, version and CSV for " + lightweightOperatorName)
		channel, _, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(lightweightOperatorName, randomNamespace)

		By(fmt.Sprintf("Deploy %s operator for testing", lightweightOp.PackageName))
		err = tshelper.DeployOperatorSubscription(
			lightweightOperatorName,
			lightweightOperatorName,
			channel,
			randomNamespace,
			lwCatalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			lightweightOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(lightweightOp.CSVPrefix,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")

		By("Query the packagemanifest for Grafana operator package name and catalog source for second operator")
		grafanaOperatorName, catalogSource2, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			"grafana", randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for Grafana operator")
		Expect(grafanaOperatorName).ToNot(Equal("not found"), "Grafana operator package not found")
		Expect(catalogSource2).ToNot(Equal("not found"), "Grafana operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + grafanaOperatorName + " for second operator")
		channel2, version2, csvName2, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			grafanaOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+grafanaOperatorName)
		Expect(channel2).ToNot(Equal("not found"), "Channel not found")
		Expect(version2).ToNot(Equal("not found"), "Version not found")
		Expect(csvName2).ToNot(Equal("not found"), "CSV name not found")

		By(fmt.Sprintf("Deploy Grafana operator (channel %s, version %s) for testing", channel2, version2))
		err = tshelper.DeployOperatorSubscription(
			"grafana-operator-2",
			grafanaOperatorName,
			channel2,
			randomNamespace,
			catalogSource2,
			tsparams.OperatorSourceNamespace,
			csvName2,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			grafanaOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(grafanaOperatorName,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName2+
			" is not ready")

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				grafanaOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+grafanaOperatorName)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				lightweightOp.CSVPrefix,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+lightweightOperatorName)

		By("Delete operator's subscription")
		err = globalhelper.DeleteSubscription(randomNamespace,
			lightweightOperatorName+"-subscription")
		Expect(err).ToNot(HaveOccurred())

		By(fmt.Sprintf("Update certsuite config to include %s operator", lightweightOp.PackageName))
		err = globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"nginxingresses.charts.nginx.org"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

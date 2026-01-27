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
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
)

var _ = Describe("Operator single-or-multi-namespaced-allowed-in-tenant-namespaces", Serial, Label("operator", "ocp-required"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string
	var randomTargetingNamespace string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.OperatorNamespace)

		randomTargetingNamespace = randomNamespace + "-targeting"

		createTestOperatorGroup(randomTargetingNamespace, tsparams.SingleOrMultiNamespacedOperatorGroup, []string{randomNamespace})

		DeferCleanup(func() {
			err := globalhelper.DeleteNamespaceAndWait(randomTargetingNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace "+randomTargetingNamespace)
		})

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace, randomTargetingNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.CertsuiteTargetOperatorLabels},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)

		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("operator namespace contains only single/multi namespace operator", func() {
		createTestOperatorGroup(randomNamespace, tsparams.SingleOrMultiNamespacedOperatorGroup,
			[]string{randomNamespace + "-one"})
		installAndLabelOperator(randomNamespace)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)

		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalparameters.TestCasePassed, randomReportDir)

		Expect(err).ToNot(HaveOccurred())
	})

	// negative
	It("operator namespace contains own-namespaced namespace operator", func() {
		By("Deploy operator group")
		err := tshelper.DeployTestOperatorGroup(randomNamespace, false)

		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		installAndLabelOperator(randomNamespace)

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)

		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalparameters.TestCaseFailed, randomReportDir)

		Expect(err).ToNot(HaveOccurred())
	})

	// positive
	It("operator namespace contains single namespaced operator with cluster-wide operator installed in a different namespace", func() {
		installClusterWideOperator()

		createTestOperatorGroup(randomNamespace, tsparams.SingleOrMultiNamespacedOperatorGroup,
			[]string{randomNamespace + "-one", randomNamespace + "-two"})
		installAndLabelOperator(randomNamespace)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)

		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalparameters.TestCasePassed, randomReportDir)

		Expect(err).ToNot(HaveOccurred())
	})

	// negative - not possible for InterOperatorGroupOwnerConflict
	/*It("operator namespace contains single namespaced operator with cluster-wide operator installed in the same namespace", func() {
		installClusterWideOperator("cluster-logging", randomNamespace)
		installAndLabelOperator(randomNamespace)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})*/

	// negative
	It("operator namespace contains single namespaced operator with non-operator pods", func() {
		createTestOperatorGroup(randomNamespace, tsparams.SingleOrMultiNamespacedOperatorGroup,
			[]string{randomNamespace + "-one", randomNamespace + "-two"})
		installAndLabelOperator(randomNamespace)

		By("Define pod")
		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)

		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalparameters.TestCaseFailed, randomReportDir)

		Expect(err).ToNot(HaveOccurred())
	})

	// negative - test is failing due to missing implementation in certsuite?
	XIt("operator namespace contains single namespaced operator with operators targeting this namespace", func() {
		createTestOperatorGroup(randomNamespace, tsparams.SingleOrMultiNamespacedOperatorGroup, []string{randomNamespace + "-one"})
		installAndLabelOperator(randomNamespace)

		// Get the lightweight operator based on OCP version
		ocpVersion := "4.19"
		if !globalhelper.IsKindCluster() {
			var verr error
			ocpVersion, verr = globalhelper.GetClusterVersion()
			Expect(verr).ToNot(HaveOccurred(), "Error getting cluster version")
		}
		lightweightOp := operatorversions.GetLightweightOperator(ocpVersion)

		By(fmt.Sprintf("Query the packagemanifest for %s operator package name and catalog source", lightweightOp.PackageName))
		lightweightOperatorName, lwCatalogSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			lightweightOp.PackageName, randomTargetingNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+lightweightOp.PackageName+" operator")
		Expect(lightweightOperatorName).ToNot(Equal("not found"), lightweightOp.PackageName+" operator package not found")
		Expect(lwCatalogSource).ToNot(Equal("not found"), lightweightOp.PackageName+" operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + lightweightOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			lightweightOperatorName, randomTargetingNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+lightweightOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		By(fmt.Sprintf("Deploy %s operator for testing", lightweightOp.PackageName))
		err = tshelper.DeployOperatorSubscription(
			lightweightOperatorName+"-targeting",
			lightweightOperatorName,
			channel,
			randomTargetingNamespace,
			lwCatalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			lightweightOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(lightweightOp.CSVPrefix,
			randomTargetingNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")

		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				lightweightOp.CSVPrefix,
				randomTargetingNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+lightweightOperatorName)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				"grafana-operator",
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+"grafana-operator")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// negative
	It("operator namespace contains single namespaced operator with operators not labelled", func() {
		// Note: This test uses grafana-operator and a lightweight operator that varies by OCP version
		// See issue #1283 and operatorversions package for operator catalog availability

		createTestOperatorGroup(randomNamespace, tsparams.SingleOrMultiNamespacedOperatorGroup, []string{randomNamespace + "-one"})
		installAndLabelOperator(randomNamespace)

		// Get the lightweight operator based on OCP version
		ocpVersion := "4.19"
		if !globalhelper.IsKindCluster() {
			var verr error
			ocpVersion, verr = globalhelper.GetClusterVersion()
			Expect(verr).ToNot(HaveOccurred(), "Error getting cluster version")
		}
		lightweightOp := operatorversions.GetLightweightOperator(ocpVersion)

		By(fmt.Sprintf("Query the packagemanifest for %s operator package name and catalog source", lightweightOp.PackageName))
		lightweightOperatorName, lwCatalogSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			lightweightOp.PackageName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+lightweightOp.PackageName+" operator")
		Expect(lightweightOperatorName).ToNot(Equal("not found"), lightweightOp.PackageName+" operator package not found")
		Expect(lwCatalogSource).ToNot(Equal("not found"), lightweightOp.PackageName+" operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + lightweightOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			lightweightOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+lightweightOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

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

		err = tshelper.WaitUntilOperatorIsReady(lightweightOp.CSVPrefix, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")

		// NOTE: Intentionally NOT labeling the lightweight operator - this should cause the test to fail

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)

		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces,
			globalparameters.TestCaseFailed, randomReportDir)

		Expect(err).ToNot(HaveOccurred())
	})
})

func installClusterWideOperator() {
	const (
		clusterLoggingOperatorName = "cluster-logging"
		openshiftLoggingNamespace  = "cluster-logging"
	)

	By("Preemptively delete the namespace if it already exists")
	err := globalhelper.DeleteNamespaceAndWait(openshiftLoggingNamespace, tsparams.Timeout)
	Expect(err).ToNot(HaveOccurred(), "Error deleting namespace "+openshiftLoggingNamespace)

	By("Create openshift-logging namespace")
	err = globalhelper.CreateNamespace(openshiftLoggingNamespace)
	Expect(err).ToNot(HaveOccurred())

	DeferCleanup(func() {
		err := globalhelper.DeleteNamespaceAndWait(openshiftLoggingNamespace, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred(), "Error deleting namespace "+openshiftLoggingNamespace)
	})

	By("Create fake operator group for cluster-logging operator")
	err = tshelper.DeployTestOperatorGroup(openshiftLoggingNamespace, true)
	Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

	By("Query the packagemanifest for defaultChannel for " + clusterLoggingOperatorName)
	channel, err := globalhelper.QueryPackageManifestForDefaultChannel(clusterLoggingOperatorName, openshiftLoggingNamespace)
	Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+clusterLoggingOperatorName)

	fmt.Printf("CHANNEL FOUND: %s\n", channel)

	By("Query the packagemanifest for the " + clusterLoggingOperatorName)
	version, err := globalhelper.QueryPackageManifestForVersion(clusterLoggingOperatorName, openshiftLoggingNamespace, channel)
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
	err = tshelper.WaitUntilOperatorIsReady(clusterLoggingOperatorName, openshiftLoggingNamespace)
	Expect(err).ToNot(HaveOccurred(), "Operator "+clusterLoggingOperatorName+" is not ready")

	By("Label operators")
	Eventually(func() error {
		return tshelper.AddLabelToInstalledCSV(
			clusterLoggingOperatorName,
			openshiftLoggingNamespace,
			tsparams.OperatorLabel)
	}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
		ErrorLabelingOperatorStr+clusterLoggingOperatorName)
}

func createTestOperatorGroup(namespace, operatorGroupName string, targetNamespaces []string) {
	err := globalhelper.CreateNamespace(namespace)
	Expect(err).ToNot(HaveOccurred(), "Error creating namespace "+namespace)

	By("Create target namespaces")

	for _, targetNamespace := range targetNamespaces {
		err := globalhelper.CreateNamespace(targetNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace "+targetNamespace)
	}

	DeferCleanup(func() {
		for _, targetNamespace := range targetNamespaces {
			err := globalhelper.DeleteNamespaceAndWait(targetNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace "+targetNamespace)
		}
	})

	By("Deploy operator group for namespace " + namespace)
	err = tshelper.DeployTestOperatorGroupWithTargetNamespace(operatorGroupName, namespace, targetNamespaces)
	Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")
}

func installAndLabelOperator(operatorNamespace string) {
	By("Query the packagemanifest for Grafana operator package name and catalog source")
	grafanaOperatorName, catalogSource := globalhelper.CheckOperatorExistsOrFail("grafana", operatorNamespace)

	By("Query the packagemanifest for available channel, version and CSV for " + grafanaOperatorName)
	channel, version, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(grafanaOperatorName, operatorNamespace)

	By(fmt.Sprintf("Deploy Grafana operator (channel %s, version %s) for testing", channel, version))
	err := tshelper.DeployOperatorSubscription(
		grafanaOperatorName,
		grafanaOperatorName,
		channel,
		operatorNamespace,
		catalogSource,
		tsparams.OperatorSourceNamespace,
		csvName,
		v1alpha1.ApprovalAutomatic,
	)
	Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
		grafanaOperatorName)

	err = tshelper.WaitUntilOperatorIsReady(grafanaOperatorName,
		operatorNamespace)
	Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
		" is not ready")

	By("Label operator")
	Eventually(func() error {
		return tshelper.AddLabelToInstalledCSV(
			grafanaOperatorName,
			operatorNamespace,
			tsparams.OperatorLabel)
	}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
		ErrorLabelingOperatorStr+grafanaOperatorName)
}

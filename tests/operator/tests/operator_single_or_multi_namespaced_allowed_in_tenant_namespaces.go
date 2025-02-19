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
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
)

var _ = Describe("Operator single-or-multi-namespaced-allowed-in-tenant-namespaces", func() {
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
			globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

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

	// negative
	It("operator namespace contains single namespaced operator with operators targeting this namespace", func() {
		createTestOperatorGroup(randomNamespace, tsparams.SingleOrMultiNamespacedOperatorGroup, []string{randomNamespace + "-one"})
		installAndLabelOperator(randomNamespace)

		By("Deploy anchore-engine operator for testing")
		err := tshelper.DeployOperatorSubscription(
			tsparams.OperatorPrefixAnchore,
			tsparams.OperatorPrefixAnchore,
			"alpha",
			randomTargetingNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPrefixAnchore)

		err = waitUntilOperatorIsReady(tsparams.OperatorPrefixAnchore,
			randomTargetingNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixAnchore+
			" is not ready")

		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				randomTargetingNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixAnchore)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixAnchore)

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
		createTestOperatorGroup(randomNamespace, tsparams.SingleOrMultiNamespacedOperatorGroup, []string{randomNamespace + "-one"})
		installAndLabelOperator(randomNamespace)

		By("Deploy anchore-engine operator for testing")
		err := tshelper.DeployOperatorSubscription(
			tsparams.OperatorPrefixAnchore,
			tsparams.OperatorPrefixAnchore,
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPrefixAnchore)

		err = waitUntilOperatorIsReady(tsparams.OperatorPrefixAnchore, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixAnchore+
			" is not ready")

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

	By("Create openshift-logging namespace")
	err := globalhelper.CreateNamespace(openshiftLoggingNamespace)
	Expect(err).ToNot(HaveOccurred())

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
	By("Query the packagemanifest for the default channel")
	channel, err := globalhelper.QueryPackageManifestForDefaultChannel(
		tsparams.CertifiedOperatorPrefixNginx,
		operatorNamespace,
	)
	Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for nginx-ingress-operator")

	By("Query the packagemanifest for the " + tsparams.CertifiedOperatorPrefixNginx)
	version, err := globalhelper.QueryPackageManifestForVersion(tsparams.CertifiedOperatorPrefixNginx, operatorNamespace, channel)
	Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for nginx-ingress-operator")

	By(fmt.Sprintf("Deploy nginx-ingress-operator%s for testing", "."+version))
	err = tshelper.DeployOperatorSubscription(
		tsparams.CertifiedOperatorPrefixNginx,
		tsparams.CertifiedOperatorPrefixNginx,
		channel,
		operatorNamespace,
		tsparams.CertifiedOperatorGroup,
		tsparams.OperatorSourceNamespace,
		tsparams.CertifiedOperatorPrefixNginx+".v"+version,
		v1alpha1.ApprovalAutomatic,
	)
	Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
		tsparams.CertifiedOperatorPrefixNginx)

	err = waitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixNginx,
		operatorNamespace)
	Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixNginx+".v"+version+
		" is not ready")

	By("Label operator")
	Eventually(func() error {
		return tshelper.AddLabelToInstalledCSV(
			tsparams.CertifiedOperatorPrefixNginx,
			operatorNamespace,
			tsparams.OperatorLabel)
	}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
		ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixNginx)
}

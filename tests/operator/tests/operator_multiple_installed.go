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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var _ = Describe("Operator multiple installed,", Serial, func() {
	var randomNamespace string
	var secondNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.OperatorNamespace)

		secondNamespace = randomNamespace + "-second"

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			tsparams.CertsuiteTargetCrdFilters, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("Deploy the same operator (and version) twice in the different namespaces", func() {
		// This is a positive test case to verify that the same operator can be deployed
		// in different namespaces.  This is a valid use case.

		By("Deploy operator group for namespace " + randomNamespace)
		err := tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Create second namespace")
		err = globalhelper.CreateNamespace(secondNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err := globalhelper.DeleteNamespaceAndWait(secondNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Deploy operator group for namespace " + secondNamespace)
		err = tshelper.DeployTestOperatorGroup(secondNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Query the packagemanifest for postgresql operator package name and catalog source")
		postgresqlOperatorName, catalogSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			tsparams.OperatorPackageNamePrefixLightweight, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for postgresql operator")
		Expect(postgresqlOperatorName).ToNot(Equal("not found"), "postgresql operator package not found")
		Expect(catalogSource).ToNot(Equal("not found"), "postgresql operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + postgresqlOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			postgresqlOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+postgresqlOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		// Note: The key to this setup is that the subscriptions can be named separately/uniquely.
		// This is because the operator/csv name is the same, but the subscription name is different.
		// The subscription name cannot be the same, as it is a unique identifier in the namespace.

		By(fmt.Sprintf("Deploy first operator (postgresql-operator) version %s for testing", version))
		err = tshelper.DeployOperatorSubscription(
			"operator1",
			postgresqlOperatorName,
			channel,
			randomNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			postgresqlOperatorName)

		By(fmt.Sprintf("Deploy second operator (postgresql-operator) version %s for testing", version))
		err = tshelper.DeployOperatorSubscription(
			"operator2",
			postgresqlOperatorName,
			channel,
			secondNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			postgresqlOperatorName)

		By("Wait until the first postgresql operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixLightweight, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+postgresqlOperatorName+
			" is not ready")

		By("Wait until the second postgresql operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixLightweight, secondNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+postgresqlOperatorName+
			" is not ready")

		// Note: No need to label these operators as we are testing all operators in the cluster.
		// At this point, two subscriptions, two installplans, and two CSVs should be present in the cluster.

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorMultipleInstalled,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorMultipleInstalled,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Deploy the same operator (different versions) different namespaces [negative]", Serial, func() {
		// This is a negative test case to verify that the same operator cannot be deployed
		// in different namespaces with different versions. This is an invalid use case.

		// We want to create a custom catalog source for this test.
		// This means we will have access to a "new" and "old" channel for the postgresql.
		// We will deploy the "new" channel in the first namespace and the "old" channel in the second namespace.

		By("Create custom-operator catalog source")
		err := globalhelper.DeployCustomOperatorSource("quay.io/redhat-best-practices-for-k8s/qe-custom-catalog")
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			err := globalhelper.DeleteCustomOperatorSource()
			Expect(err).ToNot(HaveOccurred())
		})

		By("Deploy operator group for namespace " + randomNamespace)
		err = tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Create second namespace")
		err = globalhelper.CreateNamespace(secondNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err := globalhelper.DeleteNamespaceAndWait(secondNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Deploy operator group for namespace " + secondNamespace)
		err = tshelper.DeployTestOperatorGroup(secondNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		// Define resource limits to reduce memory footprint
		resourceLimits := &corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("512Mi"),
				corev1.ResourceCPU:    resource.MustParse("500m"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("256Mi"),
				corev1.ResourceCPU:    resource.MustParse("250m"),
			},
		}

		By(fmt.Sprintf("Deploy first operator (nginx-ingress-operator) for testing - channel %s with resource limits", "new"))
		err = tshelper.DeployOperatorSubscriptionWithResourceLimits(
			"operator1",
			tsparams.OperatorPackageNamePrefixLightweightCustomCatalog,
			"new",
			randomNamespace,
			"custom-catalog",
			tsparams.OperatorSourceNamespace,
			"nginx-ingress-operator.v3.0.1",
			v1alpha1.ApprovalAutomatic,
			resourceLimits,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPackageNamePrefixLightweightCustomCatalog)

		By(fmt.Sprintf("Deploy second operator (nginx-ingress-operator) for testing - channel %s with resource limits", "old"))
		err = tshelper.DeployOperatorSubscriptionWithResourceLimits(
			"operator2",
			tsparams.OperatorPackageNamePrefixLightweightCustomCatalog,
			"old",
			secondNamespace,
			"custom-catalog",
			tsparams.OperatorSourceNamespace,
			"nginx-ingress-operator.v3.0.0",
			v1alpha1.ApprovalAutomatic,
			resourceLimits,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPackageNamePrefixLightweightCustomCatalog)

		By("Wait until the first nginx-ingress-operator operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPackageNamePrefixLightweightCustomCatalog, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPackageNamePrefixLightweightCustomCatalog+
			" is not ready")

		By("Wait until the second nginx-ingress-operator operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPackageNamePrefixLightweightCustomCatalog, secondNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPackageNamePrefixLightweightCustomCatalog+
			" is not ready")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorMultipleInstalled,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorMultipleInstalled,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

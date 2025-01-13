package operator

import (
	"fmt"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
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
			[]string{randomNamespace, secondNamespace},
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

		By("Query the packagemanifest for the " + tsparams.CertifiedOperatorPrefixNginx + " default channel")
		channel, err := globalhelper.QueryPackageManifestForDefaultChannel(tsparams.CertifiedOperatorPrefixNginx, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for nginx-ingress-operator")
		Expect(channel).ToNot(Equal("not found"), "Channel not found")

		By("Query the packagemanifest for the " + tsparams.CertifiedOperatorPrefixNginx + " version")
		version, err := globalhelper.QueryPackageManifestForVersion(tsparams.CertifiedOperatorPrefixNginx, randomNamespace, channel)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for nginx-ingress-operator")
		Expect(version).ToNot(Equal("not found"), "Version not found")

		// Note: The key to this setup is that the subscriptions can be named separately/uniquely.
		// This is because the operator/csv name is the same, but the subscription name is different.
		// The subscription name cannot be the same, as it is a unique identifier in the namespace.

		By(fmt.Sprintf("Deploy first operator (nginx-ingress-operator) version %s for testing", version))
		err = tshelper.DeployOperatorSubscription(
			"operator1",
			tsparams.CertifiedOperatorPrefixNginx,
			channel,
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorPrefixNginx+".v"+version,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.CertifiedOperatorPrefixNginx)

		By(fmt.Sprintf("Deploy second operator (nginx-ingress-operator) version %s for testing", version))
		err = tshelper.DeployOperatorSubscription(
			"operator2",
			tsparams.CertifiedOperatorPrefixNginx,
			channel,
			secondNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorPrefixNginx+".v"+version,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.CertifiedOperatorPrefixNginx)

		By("Wait until the first operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixNginx, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixNginx+
			" is not ready")

		By("Wait until the second operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixNginx, secondNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixNginx+
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

	It("Deploy the same operator (different versions) different namespaces [negative]", func() {
		// This is a negative test case to verify that the same operator cannot be deployed
		// in different namespaces with different versions.  This is an invalid use case.

		// We want to create a custom catalog source for this test.
		// This means we will have access to a "new" and "old" channel for the nginx-ingress-operator.
		// We will deploy the "new" channel in the first namespace and the "old" channel in the second namespace.

		// Skip if the host is not x86 based
		if runtime.GOARCH != "amd64" {
			Skip("This test is only supported on x86 based hosts")
		}

		By("Create custom-operator catalog source")
		err := globalhelper.DeployCustomOperatorSource("quay.io/deliedit/test:catalog-index-test")
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

		By(fmt.Sprintf("Deploy first operator (nginx-ingress-operator) for testing - channel %s", "new"))
		err = tshelper.DeployOperatorSubscription(
			"operator1",
			tsparams.CertifiedOperatorPrefixNginx,
			"new",
			randomNamespace,
			"custom-catalog",
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorPrefixNginx+".v2.4.0",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.CertifiedOperatorPrefixNginx)

		By(fmt.Sprintf("Deploy second operator (nginx-ingress-operator) for testing - channel %s", "new"))
		err = tshelper.DeployOperatorSubscription(
			"operator2",
			tsparams.CertifiedOperatorPrefixNginx,
			"old",
			secondNamespace,
			"custom-catalog",
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorPrefixNginx+".v2.3.2",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.CertifiedOperatorPrefixNginx)

		By("Wait until the first operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixNginx, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixNginx+
			" is not ready")

		By("Wait until the second operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixNginx, secondNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixNginx+
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
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

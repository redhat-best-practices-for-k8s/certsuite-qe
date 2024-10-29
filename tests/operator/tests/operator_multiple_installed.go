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

var _ = Describe("Operator multiple installed,", func() {
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

		By(fmt.Sprintf("Deploy first operator (nginx-ingress-operator) for testing"))
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

		By("Deploy second operator (nginx-ingress-operator) for testing")
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
})

// TODO: Add a negative case here with two operators same name and different versions [negative]

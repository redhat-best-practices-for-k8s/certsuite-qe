package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
)

var _ = Describe("Affiliated-certification invalid operator certification,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		if globalhelper.IsKindCluster() {
			Skip("Skipping operator certification tests on kind cluster")
		}

		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestCertificationNameSpace)

		preConfigureAffiliatedCertificationEnvironment(
			randomNamespace,
			randomCertsuiteConfigDir,
		)

		By("Query the packagemanifest for the " + tsparams.CertifiedOperatorPrefixNginx)
		version, err := globalhelper.QueryPackageManifestForVersion(tsparams.CertifiedOperatorPrefixNginx, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for nginx-ingress-operator")

		By(fmt.Sprintf("Deploy nginx-ingress-operator%s for testing", "."+version))
		// nginx-ingress-operator: in certified-operators group and version is certified
		err = tshelper.DeployOperatorSubscription(
			tsparams.CertifiedOperatorPrefixNginx,
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorPrefixNginx+".v"+version,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.CertifiedOperatorPrefixNginx)

		err = waitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixNginx,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixNginx+".v"+version+
			" is not ready")

		// sriov-fec.v1.1.0 operator : in certified-operators group, version is not certified
		By("Deploy alternate operator catalog source")
		err = globalhelper.DisableCatalogSource(tsparams.CertifiedOperatorGroup)
		Expect(err).ToNot(HaveOccurred(), "Error disabling "+
			tsparams.CertifiedOperatorGroup+" catalog source")
		Eventually(func() bool {
			stillEnabled, err := globalhelper.IsCatalogSourceEnabled(
				tsparams.CertifiedOperatorGroup,
				tsparams.OperatorSourceNamespace,
				tsparams.CertifiedOperatorDisplayName)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

			return !stillEnabled
		}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
			"Default catalog source is still enabled")

		// Deploying certified operator with invalid catalog version is necessary in order to cover negative scenarios
		err = globalhelper.DeployRHCertifiedOperatorSource("4.7")
		Expect(err).ToNot(HaveOccurred(), "Error deploying catalog source")

		By("Deploy sriov-fec operator with uncertified version")
		err = tshelper.DeployOperatorSubscription(
			tsparams.UncertifiedOperatorPrefixSriov,
			"stable",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.UncertifiedOperatorFullSriov,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.UncertifiedOperatorPrefixSriov)

		approveInstallPlanWhenReady(tsparams.UncertifiedOperatorFullSriov,
			randomNamespace)

		err = waitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixSriov,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.UncertifiedOperatorPrefixSriov+
			" is not ready")

		By("Re-enable default catalog source")
		err = globalhelper.DeleteCatalogSource(tsparams.CertifiedOperatorGroup,
			randomNamespace,
			"redhat-certified")
		Expect(err).ToNot(HaveOccurred(), "Error removing alternate catalog source")

		err = globalhelper.EnableCatalogSource(tsparams.CertifiedOperatorGroup)
		Expect(err).ToNot(HaveOccurred(), "Error enabling default catalog source")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 46695
	It("one operator to test, operator is in certified-operators organization but its version"+
		" is not certified [negative]", func() {

		if globalhelper.IsKindCluster() {
			Skip("Skip on kind cluster")
		}

		By("Label operator to be certified")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixSriov,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixSriov)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46700
	It("two operators to test, both are in certified-operators organization,"+
		" one’s version is certified, the other’s is not [negative]", func() {

		if globalhelper.IsKindCluster() {
			Skip("Skip on kind cluster")
		}

		By("Label operators to be certified")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixSriov,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixSriov)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixNginx,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixNginx)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})
})

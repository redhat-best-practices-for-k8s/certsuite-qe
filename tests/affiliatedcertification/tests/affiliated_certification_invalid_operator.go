package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
)

var _ = Describe("Affiliated-certification invalid operator certification,", Serial, func() {

	var (
		installedLabeledOperators []tsparams.OperatorLabelInfo
	)

	execute.BeforeAll(func() {
		preConfigureAffiliatedCertificationEnvironment(
			tsparams.TestCertificationNameSpace,
			globalhelper.GetConfiguration().General.TnfConfigDir,
		)

		By("Deploy nginx-ingress-operator for testing")
		// nginx-ingress-operator: in certified-operators group and version is certified
		err := tshelper.DeployOperatorSubscription(
			"nginx-ingress-operator",
			"alpha",
			tsparams.TestCertificationNameSpace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorFullNginx,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.CertifiedOperatorPrefixNginx)

		err = waitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixNginx,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixNginx+
			" is not ready")

		// add nginx-ingress-operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.CertifiedOperatorPrefixNginx,
			Namespace:      tsparams.TestCertificationNameSpace,
			Label:          tsparams.OperatorLabel,
		})

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
			tsparams.TestCertificationNameSpace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.UncertifiedOperatorFullSriov,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.UncertifiedOperatorPrefixSriov)

		approveInstallPlanWhenReady(tsparams.UncertifiedOperatorFullSriov,
			tsparams.TestCertificationNameSpace)

		err = waitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixSriov,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.UncertifiedOperatorPrefixSriov+
			" is not ready")

		By("Re-enable default catalog source")
		err = globalhelper.DeleteCatalogSource(tsparams.CertifiedOperatorGroup,
			tsparams.TestCertificationNameSpace,
			"redhat-certified")
		Expect(err).ToNot(HaveOccurred(), "Error removing alternate catalog source")

		err = globalhelper.EnableCatalogSource(tsparams.CertifiedOperatorGroup)
		Expect(err).ToNot(HaveOccurred(), "Error enabling default catalog source")

		// add sriov-fec operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.UncertifiedOperatorPrefixSriov,
			Namespace:      tsparams.TestCertificationNameSpace,
			Label:          tsparams.OperatorLabel,
		})

	})

	AfterEach(func() {
		By("Remove labels from operators")
		for _, info := range installedLabeledOperators {
			err := tshelper.DeleteLabelFromInstalledCSV(
				info.OperatorPrefix,
				info.Namespace,
				info.Label)
			Expect(err).ToNot(HaveOccurred(), "Error removing label from operator "+info.OperatorPrefix)
		}

		By("Remove reports from report directory")
		err := globalhelper.RemoveContentsFromReportDir(globalhelper.GetConfiguration().General.TnfReportDir)
		Expect(err).ToNot(HaveOccurred())
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
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixSriov)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			globalhelper.GetConfiguration().General.TnfReportDir,
			globalhelper.GetConfiguration().General.TnfConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, globalhelper.GetConfiguration().General.TnfReportDir)
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
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.UncertifiedOperatorPrefixSriov)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixNginx,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixNginx)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			globalhelper.GetConfiguration().General.TnfReportDir,
			globalhelper.GetConfiguration().General.TnfConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed, globalhelper.GetConfiguration().General.TnfReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})
})

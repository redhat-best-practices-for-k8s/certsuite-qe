package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	operatorutils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

var _ = Describe("Affiliated-certification invalid operator certification,", func() {

	var (
		installedLabeledOperators []tsparams.OperatorLabelInfo
	)

	execute.BeforeAll(func() {
		preConfigureAffiliatedCertificationEnvironment()

		By("Deploy openshiftartifactoryha-operator for testing")
		// openshiftartifactoryha-operator: in certified-operators group and version is certified
		err := operatorutils.DeployOperatorSubscription(
			"openshiftartifactoryha-operator",
			"alpha",
			tsparams.TestCertificationNameSpace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorFullArtifactoryHa,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.CertifiedOperatorPrefixArtifactoryHa)

		approveInstallPlanWhenReady(tsparams.CertifiedOperatorFullArtifactoryHa,
			tsparams.TestCertificationNameSpace)

		err = operatorutils.WaitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixArtifactoryHa,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixArtifactoryHa+
			" is not ready")

		// add openshiftartifactoryha operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.CertifiedOperatorPrefixArtifactoryHa,
			Namespace:      tsparams.TestCertificationNameSpace,
			Label:          tsparams.OperatorLabel,
		})

		// sriov-fec.v1.1.0 operator : in certified-operators group, version is not certified
		By("Deploy alternate operator catalog source")
		err = operatorutils.DisableCatalogSource(tsparams.CertifiedOperatorGroup)
		Expect(err).ToNot(HaveOccurred(), "Error disabling "+
			tsparams.CertifiedOperatorGroup+" catalog source")
		Eventually(func() bool {
			stillEnabled, err := operatorutils.IsCatalogSourceEnabled(
				tsparams.CertifiedOperatorGroup,
				tsparams.OperatorSourceNamespace,
				tsparams.CertifiedOperatorDisplayName)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

			return !stillEnabled
		}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
			"Default catalog source is still enabled")

		// Deploying certified operator with invalid catalog version is necessary in order to cover negative scenarios
		err = operatorutils.DeployRHCertifiedOperatorSource("4.5")
		Expect(err).ToNot(HaveOccurred(), "Error deploying catalog source")

		By("Deploy sriov-fec operator with uncertified version")
		err = operatorutils.DeployOperatorSubscription(
			tsparams.UncertifiedOperatorPrefixSriov,
			"stable",
			tsparams.TestCertificationNameSpace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.UncertifiedOperatorFullSriov,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.UncertifiedOperatorPrefixSriov)

		approveInstallPlanWhenReady(tsparams.UncertifiedOperatorFullSriov,
			tsparams.TestCertificationNameSpace)

		err = operatorutils.WaitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixSriov,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.UncertifiedOperatorPrefixSriov+
			" is not ready")

		By("Re-enable default catalog source")
		err = operatorutils.DeleteCatalogSource(tsparams.CertifiedOperatorGroup,
			tsparams.TestCertificationNameSpace,
			"redhat-certified")
		Expect(err).ToNot(HaveOccurred(), "Error removing alternate catalog source")

		err = operatorutils.EnableCatalogSource(tsparams.CertifiedOperatorGroup)
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
			err := operatorutils.DeleteLabelFromInstalledCSV(
				info.OperatorPrefix,
				info.Namespace,
				info.Label)
			Expect(err).ToNot(HaveOccurred(), "Error removing label from operator "+info.OperatorPrefix)
		}
	})

	// 46695
	It("one operator to test, operator is in certified-operators organization but its version"+
		" is not certified [negative]", func() {

		By("Label operator to be certified")
		Eventually(func() error {
			return operatorutils.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixSriov,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.UncertifiedOperatorPrefixSriov)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46700
	It("two operators to test, both are in certified-operators organization,"+
		" one’s version is certified, the other’s is not [negative]", func() {

		By("Label operators to be certified")
		Eventually(func() error {
			return operatorutils.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixSriov,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.UncertifiedOperatorPrefixSriov)

		Eventually(func() error {
			return operatorutils.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixArtifactoryHa,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.CertifiedOperatorPrefixArtifactoryHa)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")

	})
})

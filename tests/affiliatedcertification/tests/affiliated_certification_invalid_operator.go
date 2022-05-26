package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcerthelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcertparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Affiliated-certification invalid operator certification,", func() {

	var (
		installedLabeledOperators []affiliatedcertparameters.OperatorLabelInfo
	)

	execute.BeforeAll(func() {
		preConfigureAffiliatedCertificationEnvironment()

		By("Deploy operators for testing")
		// openshiftartifactoryha-operator: in certified-operators group and version is certified
		err := affiliatedcerthelper.DeployOperatorSubscription(
			"openshiftartifactoryha-operator",
			"alpha",
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.CertifiedOperatorGroup,
			affiliatedcertparameters.OperatorSourceNamespace,
			affiliatedcertparameters.CertifiedOperatorFullArtifactoryHa,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa)

		approveInstallPlanWhenReady(affiliatedcertparameters.CertifiedOperatorFullArtifactoryHa,
			affiliatedcertparameters.TestCertificationNameSpace)

		err = waitUntilOperatorIsReady(affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa,
			affiliatedcertparameters.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa+
			" is not ready")

		// add openshiftartifactoryha operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		// sriov-fec.v1.1.0 operator : in certified-operators group, version is not certified
		By("Deploy alternate operator catalog source")
		err = affiliatedcerthelper.DisableCatalogSource(affiliatedcertparameters.CertifiedOperatorGroup)
		Expect(err).ToNot(HaveOccurred(), "Error disabling "+
			affiliatedcertparameters.CertifiedOperatorGroup+" catalog source")
		Eventually(func() bool {
			stillEnabled, err := affiliatedcerthelper.IsCatalogSourceEnabled(
				affiliatedcertparameters.CertifiedOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
				affiliatedcertparameters.CertifiedOperatorDisplayName)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

			return !stillEnabled
		}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
			"Default catalog source is still enabled")

		err = affiliatedcerthelper.DeployRHCertifiedOperatorSource("4.5")
		Expect(err).ToNot(HaveOccurred(), "Error deploying catalog source")

		By("Deploy operator with uncertified version")
		err = affiliatedcerthelper.DeployOperatorSubscription(
			affiliatedcertparameters.UncertifiedOperatorPrefixSriov,
			"stable",
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.CertifiedOperatorGroup,
			affiliatedcertparameters.OperatorSourceNamespace,
			affiliatedcertparameters.UncertifiedOperatorFullSriov,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			affiliatedcertparameters.UncertifiedOperatorPrefixSriov)

		approveInstallPlanWhenReady(affiliatedcertparameters.UncertifiedOperatorFullSriov,
			affiliatedcertparameters.TestCertificationNameSpace)

		err = waitUntilOperatorIsReady(affiliatedcertparameters.UncertifiedOperatorPrefixSriov,
			affiliatedcertparameters.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+affiliatedcertparameters.UncertifiedOperatorPrefixSriov+
			" is not ready")

		By("Re-enable default catalog source")
		err = affiliatedcerthelper.DeleteCatalogSource(affiliatedcertparameters.CertifiedOperatorGroup,
			affiliatedcertparameters.TestCertificationNameSpace,
			"redhat-certified")
		Expect(err).ToNot(HaveOccurred(), "Error removing alternate catalog source")

		err = affiliatedcerthelper.EnableCatalogSource(affiliatedcertparameters.CertifiedOperatorGroup)
		Expect(err).ToNot(HaveOccurred(), "Error enabling default catalog source")

		// add sriov-fec operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.UncertifiedOperatorPrefixSriov,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

	})

	AfterEach(func() {
		By("Remove labels from operators")
		for _, info := range installedLabeledOperators {
			err := affiliatedcerthelper.DeleteLabelFromInstalledCSV(
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
			return affiliatedcerthelper.AddLabelToInstalledCSV(
				affiliatedcertparameters.UncertifiedOperatorPrefixSriov,
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.OperatorLabel)
		}, affiliatedcertparameters.TimeoutLabelCsv, affiliatedcertparameters.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+affiliatedcertparameters.UncertifiedOperatorPrefixSriov)

		By("Start test")
		err := globalhelper.LaunchTests(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46700
	It("two operators to test, both are in certified-operators organization,"+
		" one’s version is certified, the other’s is not [negative]", func() {

		By("Label operators to be certified")
		Eventually(func() error {
			return affiliatedcerthelper.AddLabelToInstalledCSV(
				affiliatedcertparameters.UncertifiedOperatorPrefixSriov,
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.OperatorLabel)
		}, affiliatedcertparameters.TimeoutLabelCsv, affiliatedcertparameters.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+affiliatedcertparameters.UncertifiedOperatorPrefixSriov)

		Eventually(func() error {
			return affiliatedcerthelper.AddLabelToInstalledCSV(
				affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa,
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.OperatorLabel)
		}, affiliatedcertparameters.TimeoutLabelCsv, affiliatedcertparameters.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa)

		By("Start test")
		err := globalhelper.LaunchTests(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")

	})
})

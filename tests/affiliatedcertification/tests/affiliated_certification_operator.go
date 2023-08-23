package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

var _ = Describe("Affiliated-certification operator certification,", Serial, func() {

	var (
		installedLabeledOperators []tsparams.OperatorLabelInfo
	)

	execute.BeforeAll(func() {
		preConfigureAffiliatedCertificationEnvironment()

		By("Deploy falcon-operator for testing")
		// falcon-operator: not in certified-operators group in catalog, for negative test cases
		err := tshelper.DeployOperatorSubscription(
			"falcon-operator",
			"alpha",
			tsparams.TestCertificationNameSpace,
			tsparams.CommunityOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.UncertifiedOperatorPrefixFalcon)

		err = waitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixFalcon,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.UncertifiedOperatorPrefixFalcon+
			" is not ready")

		// add falcon operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.UncertifiedOperatorPrefixFalcon,
			Namespace:      tsparams.TestCertificationNameSpace,
			Label:          tsparams.OperatorLabel,
		})

		By("Deploy federatorai operator for testing")
		// federatorai operator: in certified-operators group and version is certified
		err = tshelper.DeployOperatorSubscription(
			"federatorai-certified",
			"stable",
			tsparams.TestCertificationNameSpace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorFullFederatorai,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.CertifiedOperatorPrefixFederatorai)

		approveInstallPlanWhenReady(tsparams.CertifiedOperatorFullFederatorai,
			tsparams.TestCertificationNameSpace)

		err = waitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixFederatorai,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixFederatorai+
			" is not ready")

		// add federatorai operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.CertifiedOperatorPrefixFederatorai,
			Namespace:      tsparams.TestCertificationNameSpace,
			Label:          tsparams.OperatorLabel,
		})

		By("Deploy instana-agent-operator for testing")
		// instana-agent-operator: in certified-operators group and version is certified
		err = tshelper.DeployOperatorSubscription(
			"instana-agent-operator",
			"stable",
			tsparams.TestCertificationNameSpace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorFullInstana,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.CertifiedOperatorPrefixInstana)

		approveInstallPlanWhenReady(tsparams.CertifiedOperatorFullInstana,
			tsparams.TestCertificationNameSpace)

		err = waitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixInstana,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixInstana+
			" is not ready")

		// add instana-agent-operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.CertifiedOperatorPrefixInstana,
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
	})

	// 46699
	It("one operator to test, operator is not in certified-operators organization [negative]",
		func() {
			By("Label operator to be certified")
			Eventually(func() error {
				return tshelper.AddLabelToInstalledCSV(
					tsparams.UncertifiedOperatorPrefixFalcon,
					tsparams.TestCertificationNameSpace,
					tsparams.OperatorLabel)
			}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
				"Error labeling operator "+tsparams.UncertifiedOperatorPrefixFalcon)

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

	// 46697
	It("two operators to test, one is in certified-operators organization and its version is certified,"+
		" one is not in certified-operators organization [negative]", func() {
		By("Label operators to be certified")

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixFederatorai,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.CertifiedOperatorPrefixFederatorai)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixFalcon,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.UncertifiedOperatorPrefixFalcon)

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

	// 46582
	It("one operator to test, operator is in certified-operators organization"+
		" and its version is certified", func() {
		By("Label operator to be certified")

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixInstana,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.CertifiedOperatorPrefixInstana)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46696
	It("two operators to test, both are in certified-operators organization and their"+
		" versions are certified", func() {
		By("Label operators to be certified")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixInstana,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.CertifiedOperatorPrefixInstana)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixFederatorai,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.CertifiedOperatorPrefixFederatorai)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46698
	It("no operators are labeled for testing [skip]", func() {
		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

})

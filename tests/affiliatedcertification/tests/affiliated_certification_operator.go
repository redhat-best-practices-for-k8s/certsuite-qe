package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	operatorutils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"
)

var _ = Describe("Affiliated-certification operator certification,", func() {

	var (
		installedLabeledOperators []tsparams.OperatorLabelInfo
	)

	execute.BeforeAll(func() {
		preConfigureAffiliatedCertificationEnvironment()

		By("Deploy falcon-operator for testing")
		// falcon-operator: not in certified-operators group in catalog, for negative test cases
		err := operatorutils.DeployOperatorSubscription(
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

		err = operatorutils.WaitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixFalcon,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.UncertifiedOperatorPrefixFalcon+
			" is not ready")

		// add falcon operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.UncertifiedOperatorPrefixFalcon,
			Namespace:      tsparams.TestCertificationNameSpace,
			Label:          tsparams.OperatorLabel,
		})

		By("Deploy infinibox-operator for testing")
		// infinibox-operator: in certified-operators group and version is certified
		err = operatorutils.DeployOperatorSubscription(
			"infinibox-operator-certified",
			"stable",
			tsparams.TestCertificationNameSpace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorFullInfinibox,
			v1alpha1.ApprovalManual,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.CertifiedOperatorPrefixInfinibox)

		approveInstallPlanWhenReady(tsparams.CertifiedOperatorFullInfinibox,
			tsparams.TestCertificationNameSpace)

		err = operatorutils.WaitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixInfinibox,
			tsparams.TestCertificationNameSpace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixInfinibox+
			" is not ready")

		// add infinibox operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.CertifiedOperatorPrefixInfinibox,
			Namespace:      tsparams.TestCertificationNameSpace,
			Label:          tsparams.OperatorLabel,
		})

		By("Deploy openshiftartifactoryha-operator for testing")
		// openshiftartifactoryha-operator: in certified-operators group and version is certified
		err = operatorutils.DeployOperatorSubscription(
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

	// 46699
	It("one operator to test, operator is not in certified-operators organization [negative]",
		func() {
			By("Label operator to be certified")
			Eventually(func() error {
				return operatorutils.AddLabelToInstalledCSV(
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
			return operatorutils.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixInfinibox,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.CertifiedOperatorPrefixInfinibox)

		Eventually(func() error {
			return operatorutils.AddLabelToInstalledCSV(
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
			return operatorutils.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixArtifactoryHa,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.CertifiedOperatorPrefixArtifactoryHa)

		Eventually(func() error {
			return operatorutils.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixInfinibox,
				tsparams.TestCertificationNameSpace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.CertifiedOperatorPrefixInfinibox)

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

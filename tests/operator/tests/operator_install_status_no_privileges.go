package operator

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/operator/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/operator/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Operator install-status-no-privileges,", func() {

	var (
		installedLabeledOperators []tsparams.OperatorLabelInfo
	)

	execute.BeforeAll(func() {
		By("Deploy operator group")
		err := tshelper.DeployTestOperatorGroup()
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Deploy cloudbees-ci operator for testing")
		err = tshelper.DeployOperatorSubscription(
			"cloudbees-ci",
			"alpha",
			tsparams.OperatorNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.OperatorPrefixCloudbees)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixCloudbees,
			tsparams.OperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixCloudbees+
			" is not ready")

		// add cloudbees operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.OperatorPrefixCloudbees,
			Namespace:      tsparams.OperatorNamespace,
			Label:          tsparams.OperatorLabel,
		})

		By("Deploy anchore-engine operator for testing")
		err = tshelper.DeployOperatorSubscription(
			"anchore-engine",
			"alpha",
			tsparams.OperatorNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.OperatorPrefixAnchore)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixAnchore,
			tsparams.OperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixAnchore+
			" is not ready")

		// add anchore-engine operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.OperatorPrefixAnchore,
			Namespace:      tsparams.OperatorNamespace,
			Label:          tsparams.OperatorLabel,
		})

		By("Deploy openvino operator for testing")
		err = tshelper.DeployOperatorSubscription(
			"ovms-operator",
			"alpha",
			tsparams.OperatorNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			tsparams.OperatorPrefixOpenvino)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixOpenvino,
			tsparams.OperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixOpenvino+
			" is not ready")

		// add openvino operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.OperatorPrefixOpenvino,
			Namespace:      tsparams.OperatorNamespace,
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

	// 66142
	FIt("one operator with no clusterPermissions", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixCloudbees)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66143
	It("one operator not installed with OLM [negative]", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixOpenvino)

		By("Delete operator's subscription")
		err := globalhelper.DeleteSubscription(tsparams.OperatorNamespace,
			tsparams.SubscriptionNameOpenvino,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66144
	FIt("two operators, both with no clusterPermissions", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixCloudbees)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixAnchore)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66145
	It("two operators, one not installed with OLM [negative]", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixAnchore)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.OperatorPrefixOpenvino)

		By("Delete operator's subscription")
		err := globalhelper.DeleteSubscription(tsparams.OperatorNamespace,
			tsparams.SubscriptionNameOpenvino,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

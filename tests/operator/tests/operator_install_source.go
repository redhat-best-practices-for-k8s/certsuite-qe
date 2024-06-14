package operator

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/operator/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/operator/parameters"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
	ErrorRemovingLabelStr    = "Error removing label from operator "
)

var _ = Describe("Operator install-source,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.OperatorNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			tsparams.TnfTargetCrdFilters, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		// Install 3 separate operators for testing

		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup(randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Deploy openvino operator for testing")
		err = tshelper.DeployOperatorSubscription(
			"ovms-operator",
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPrefixOpenvino)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixOpenvino,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixOpenvino+
			" is not ready")

		By("Deploy cloudbees-ci operator for testing")
		err = tshelper.DeployOperatorSubscription(
			"cloudbees-ci",
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)

		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPrefixCloudbees)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixCloudbees,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixCloudbees+
			" is not ready")

		By("Deploy anchore-engine operator for testing")
		err = tshelper.DeployOperatorSubscription(
			"anchore-engine",
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPrefixAnchore)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixAnchore,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixAnchore+
			" is not ready")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.Timeout)
	})

	// 66142
	It("one operator installed with OLM", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixCloudbees)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallSource,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66143
	It("one operator not installed with OLM [negative]", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixOpenvino)

		By("Delete operator's subscription")
		err := globalhelper.DeleteSubscription(randomNamespace,
			tsparams.SubscriptionNameOpenvino)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallSource,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66144
	It("two operators, both installed with OLM", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixCloudbees)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixAnchore)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallSource,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66145
	It("two operators, one not installed with OLM [negative]", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixAnchore)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixOpenvino)

		By("Delete operator's subscription")
		err := globalhelper.DeleteSubscription(randomNamespace,
			tsparams.SubscriptionNameOpenvino)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallSource,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

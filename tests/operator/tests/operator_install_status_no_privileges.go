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

var _ = Describe("Operator install-status-no-privileges,", Serial, func() {

	var (
		installedLabeledOperators []tsparams.OperatorLabelInfo
	)

	execute.BeforeAll(func() {
		By("Clean namespace")
		err := globalhelper.CleanNamespace(tsparams.OperatorNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup()
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		// cloudbees operator has clusterPermissions but no resourceNames
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
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
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

		// quay operator has no clusterPermissions
		By("Deploy quay operator for testing")
		err = tshelper.DeployOperatorSubscription(
			"project-quay",
			"stable-3.7",
			tsparams.OperatorNamespace,
			tsparams.CommunityOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPrefixQuay)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixQuay,
			tsparams.OperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixQuay+
			" is not ready")

		// add quay operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.OperatorPrefixQuay,
			Namespace:      tsparams.OperatorNamespace,
			Label:          tsparams.OperatorLabel,
		})

		// kiali operator has resourceNames under its rules
		By("Deploy kiali operator for testing")
		err = tshelper.DeployOperatorSubscription(
			"kiali",
			"alpha",
			tsparams.OperatorNamespace,
			tsparams.CommunityOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.OperatorPrefixKiali)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixKiali,
			tsparams.OperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.OperatorPrefixKiali+
			" is not ready")

		// add kiali operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.OperatorPrefixKiali,
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

	// 66381
	It("one operator with no clusterPermissions", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixQuay,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixQuay)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			globalhelper.GetConfiguration().General.TnfReportDir,
			globalhelper.GetConfiguration().General.TnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, globalhelper.GetConfiguration().General.TnfReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66382
	It("one operator with clusterPermissions but no resourceNames", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixCloudbees)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			globalhelper.GetConfiguration().General.TnfReportDir,
			globalhelper.GetConfiguration().General.TnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, globalhelper.GetConfiguration().General.TnfReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66383
	It("one operator with clusterPermissions and resourceNames [negative]", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixKiali,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixKiali)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			globalhelper.GetConfiguration().General.TnfReportDir,
			globalhelper.GetConfiguration().General.TnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCaseFailed, globalhelper.GetConfiguration().General.TnfReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66384
	It("two operators, one with no clusterPermissions and one with clusterPermissions but no resourceNames", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixQuay,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixQuay)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixCloudbees)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			globalhelper.GetConfiguration().General.TnfReportDir,
			globalhelper.GetConfiguration().General.TnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, globalhelper.GetConfiguration().General.TnfReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66385
	It("two operators, one with clusterPermissions and resourceNames [negative]", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixKiali,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixKiali)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixQuay,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixQuay)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			globalhelper.GetConfiguration().General.TnfReportDir,
			globalhelper.GetConfiguration().General.TnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallStatusNoPrivileges,
			globalparameters.TestCaseFailed, globalhelper.GetConfiguration().General.TnfReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

})

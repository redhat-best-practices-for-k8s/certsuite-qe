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

var _ = Describe("Operator install-source, ", func() {

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
			tsparams.UncertifiedOperatorPrefixCloudbees)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixCloudbees,
			tsparams.OperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.UncertifiedOperatorPrefixCloudbees+
			" is not ready")

		// add cloudbees operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, tsparams.OperatorLabelInfo{
			OperatorPrefix: tsparams.UncertifiedOperatorPrefixCloudbees,
			Namespace:      tsparams.OperatorNamespace,
			Label:          tsparams.OperatorLabel,
		})

	})

	BeforeEach(func() {
		// By("Clean namespace before each test")
		// err := namespaces.Clean(tsparams.OperatorNamespace, globalhelper.GetAPIClient())
		// Expect(err).ToNot(HaveOccurred())

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

	// 53939
	It("one operator installed with OLM", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.UncertifiedOperatorPrefixCloudbees,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			"Error labeling operator "+tsparams.UncertifiedOperatorPrefixCloudbees)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallSource,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53940
	It("one operator not installed with OLM [negative]", func() {
		By("Label operator")

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallSource,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53941
	It("two operators, both installed with OLM", func() {
		By("Label operators")

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallSource,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53946
	It("two operators, one not installed with OLM [negative]", func() {
		By("Label operators")

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorInstallSource,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

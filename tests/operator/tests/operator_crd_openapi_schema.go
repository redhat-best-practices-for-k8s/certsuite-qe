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

var _ = Describe("Operator crd-openapi-schema", Serial, func() {

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
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
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
			Expect(err).ToNot(HaveOccurred(), ErrorRemovingLabelStr+info.OperatorPrefix)
		}
	})

	It("operator crd is defined with openapi schema", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				tsparams.OperatorNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixCloudbees)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TnfOperatorCrdOpenAPISchema,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			globalhelper.GetConfiguration().General.TnfReportDir,
			globalhelper.GetConfiguration().General.TnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorCrdOpenAPISchema,
			globalparameters.TestCasePassed, globalhelper.GetConfiguration().General.TnfReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

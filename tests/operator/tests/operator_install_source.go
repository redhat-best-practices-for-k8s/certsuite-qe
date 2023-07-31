package operator

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/operator/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Operator install-source, ", func() {

	execute.BeforeAll(func() {

	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.OperatorNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

	})

	// 53939
	It("one operator installed with OLM", func() {
		By("Label operator")

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

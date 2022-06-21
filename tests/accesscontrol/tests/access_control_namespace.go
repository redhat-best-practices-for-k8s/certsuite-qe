package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control namespace, ", func() {

	execute.BeforeAll(func() {
		By("Create additional namespaces for testing")
		err := namespaces.Create(parameters.AdditionalValidNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Create(parameters.InvalidNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

	})

	BeforeEach(func() {

	})

	// 51860
	It("one namespace, no invalid prefixes", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			parameters.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 51862
	It("one namespace, namespace has invalid prefix [negative]", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.InvalidNamespace},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred(), "Error running "+
			parameters.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlNamespace,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")

	})

	// 51863
	It("two namespaces, no invalid prefixes", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace, parameters.AdditionalValidNamespace},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			parameters.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")

	})

	// 51864
	It("two namespaces, one has invalid prefix [negative]", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.InvalidNamespace, parameters.AdditionalValidNamespace},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred(), "Error running "+
			parameters.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlNamespace,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 51971
	It("one custom resource in a valid namespace", func() {
		Skip("Under development")
	})

	// 52058
	It("one custom resource in an invalid namespace [negative]", func() {
		Skip("Under development")
	})

	// 52069
	It("two custom resources, both in valid namespaces", func() {
		Skip("Under development")
	})

	// 52070
	It("two custom resources, one in invalid namespace [negative]", func() {
		Skip("Under development")
	})

	// 52073
	It("two custom resources of different CRDs, both in valid namespace", func() {
		Skip("Under development")
	})

	// 52098
	It("two custom resources of different CRDs, one in invalid namespace [negative]", func() {
		Skip("Under development")
	})

})

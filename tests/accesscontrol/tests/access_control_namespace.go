package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control namespace, ", func() {

	execute.BeforeAll(func() {
		By("Clean test suite namespace before tests")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Create additional namespaces for testing")
		// these namespaces will only be used for the access-control-namespace tests
		err = namespaces.Create(tsparams.AdditionalValidNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Create(tsparams.InvalidNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

	})

	BeforeEach(func() {
		By("Clean namespaces before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Clean(tsparams.AdditionalValidNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Clean(tsparams.InvalidNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

	})

	// 51860
	It("one namespace, no invalid prefixes", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.TestAccessControlNameSpace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 51862
	It("one namespace, namespace has invalid prefix [negative]", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.InvalidNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")

	})

	// 51863
	It("two namespaces, no invalid prefixes", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.TestAccessControlNameSpace, tsparams.AdditionalValidNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
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
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace, "tnf"},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{"installplans.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resource")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", tsparams.TestAccessControlNameSpace,
			globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

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

	// 52058
	It("one custom resource in an invalid namespace [negative]", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace, "tnf"},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{"installplans.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resource")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", tsparams.InvalidNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

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

	// 52069
	It("two custom resources, both in valid namespaces", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace, tsparams.AdditionalValidNamespace, "tnf"},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{"installplans.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resources")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", tsparams.TestAccessControlNameSpace,
			globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		err = tshelper.DefineAndCreateInstallPlan("test-plan-2", tsparams.AdditionalValidNamespace,
			globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

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

	// 52070
	It("two custom resources, one in invalid namespace [negative]", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace, "tnf"},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{"installplans.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resources")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		err = tshelper.DefineAndCreateInstallPlan("test-plan-2", tsparams.InvalidNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

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

	// 52073
	It("two custom resources of different CRDs, both in valid namespace", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace, "tnf"},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{"installplans.operators.coreos.com", "subscriptions.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resources")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", tsparams.TestAccessControlNameSpace,
			globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		err = tshelper.DefineAndCreateSubscription("test-sub", tsparams.TestAccessControlNameSpace,
			globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating subscription")

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

	// 52098
	It("two custom resources of different CRDs, one in invalid namespace [negative]", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace, "tnf"},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{"installplans.operators.coreos.com", "subscriptions.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resources")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", tsparams.TestAccessControlNameSpace,
			globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		err = tshelper.DefineAndCreateSubscription("test-sub", tsparams.AdditionalValidNamespace,
			globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating subscription")

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

})

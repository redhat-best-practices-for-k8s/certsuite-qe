package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control namespace, ", Serial, func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})

	// 51860
	It("one namespace, no invalid prefixes", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
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
		By("Create Invalid Namespace")
		invalidNamespace := tsparams.InvalidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := namespaces.Create(invalidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, invalidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{invalidNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
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
		By("Create additional valid namespace")
		additionalValidNamespace := tsparams.AdditionalValidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := namespaces.Create(additionalValidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, additionalValidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, additionalValidNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
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
		By("Create additional valid namespace")
		additionalValidNamespace := tsparams.AdditionalValidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := namespaces.Create(additionalValidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, additionalValidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Create Invalid Namespace")
		invalidNamespace := tsparams.InvalidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err = namespaces.Create(invalidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, invalidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{invalidNamespace, additionalValidNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
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

	// 51971
	It("one custom resource in a valid namespace", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resource")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", randomNamespace,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

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

	// 52058
	It("one custom resource in an invalid namespace [negative]", func() {
		By("Create Invalid Namespace")
		invalidNamespace := tsparams.InvalidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := namespaces.Create(invalidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, invalidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resource")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", invalidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

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

	// 52069
	It("two custom resources, both in valid namespaces", func() {
		By("Create additional valid namespace")
		additionalValidNamespace := tsparams.AdditionalValidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := namespaces.Create(additionalValidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, additionalValidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, additionalValidNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resources")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", randomNamespace,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		err = tshelper.DefineAndCreateInstallPlan("test-plan-2", additionalValidNamespace,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

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

	// 52070
	It("two custom resources, one in invalid namespace [negative]", func() {
		By("Create Invalid Namespace")
		invalidNamespace := tsparams.InvalidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := namespaces.Create(invalidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, invalidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resources")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", randomNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		err = tshelper.DefineAndCreateInstallPlan("test-plan-2", invalidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

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

	// 52073
	It("two custom resources of different CRDs, both in valid namespace", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com", "subscriptions.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resources")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", randomNamespace,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		err = tshelper.DefineAndCreateSubscription("test-sub", randomNamespace,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating subscription")

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

	// 52098
	It("two custom resources of different CRDs, one in invalid namespace [negative]", func() {
		By("Create Additional Valid Namespace")
		additionalValidNamespace := tsparams.AdditionalValidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := namespaces.Create(additionalValidNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = namespaces.DeleteAndWait(globalhelper.GetAPIClient().CoreV1Interface, additionalValidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com", "subscriptions.operators.coreos.com"})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create custom resources")
		err = tshelper.DefineAndCreateInstallPlan("test-plan", randomNamespace,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		err = tshelper.DefineAndCreateSubscription("test-sub", additionalValidNamespace,
			globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating subscription")

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
})

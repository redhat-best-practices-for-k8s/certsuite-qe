package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/installplan"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/subscription"
)

const (
	CreateInstallPlanInNamespaceStr = "Create Install Plan in Namespace: "
)

var _ = Describe("Access-control namespace, ", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 51860
	It("one namespace, no invalid prefixes", func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 51862
	It("one namespace, namespace has invalid prefix [negative]", func() {
		By("Create Invalid Namespace")
		invalidNamespace := tsparams.InvalidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(invalidNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(invalidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{invalidNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 51863
	It("two namespaces, no invalid prefixes", func() {
		By("Create additional valid namespace")
		additionalValidNamespace := tsparams.AdditionalValidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(additionalValidNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(additionalValidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, additionalValidNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 51864
	It("two namespaces, one has invalid prefix [negative]", func() {
		By("Create additional valid namespace")
		additionalValidNamespace := tsparams.AdditionalValidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(additionalValidNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(additionalValidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Create Invalid Namespace")
		invalidNamespace := tsparams.InvalidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err = globalhelper.CreateNamespace(invalidNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(invalidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{invalidNamespace, additionalValidNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCaseFailed, randomReportDir)
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
			[]string{"installplans.operators.coreos.com"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Define Install Plan")
		plan := installplan.DefineInstallPlan("test-plan", randomNamespace)

		By(CreateInstallPlanInNamespaceStr + plan.Namespace)
		err = globalhelper.CreateInstallPlan(plan)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 52058
	It("one custom resource in an invalid namespace [negative]", func() {
		By("Create Invalid Namespace")
		invalidNamespace := tsparams.InvalidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(invalidNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(invalidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Define Install Plan")
		plan := installplan.DefineInstallPlan("test-plan", invalidNamespace)

		By(CreateInstallPlanInNamespaceStr + plan.Namespace)
		err = globalhelper.CreateInstallPlan(plan)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 52069
	It("two custom resources, both in valid namespaces", func() {
		By("Create additional valid namespace")
		additionalValidNamespace := tsparams.AdditionalValidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(additionalValidNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(additionalValidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, additionalValidNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Define Install Plan")
		plan := installplan.DefineInstallPlan("test-plan", randomNamespace)

		By(CreateInstallPlanInNamespaceStr + plan.Namespace)
		err = globalhelper.CreateInstallPlan(plan)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		By("Define Install Plan")
		plan2 := installplan.DefineInstallPlan("test-plan-2", additionalValidNamespace)

		By(CreateInstallPlanInNamespaceStr + plan2.Namespace)
		err = globalhelper.CreateInstallPlan(plan2)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 52070
	It("two custom resources, one in invalid namespace [negative]", func() {
		By("Create Invalid Namespace")
		invalidNamespace := tsparams.InvalidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(invalidNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(invalidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Define Install Plan")
		plan := installplan.DefineInstallPlan("test-plan", randomNamespace)

		By(CreateInstallPlanInNamespaceStr + plan.Namespace)
		err = globalhelper.CreateInstallPlan(plan)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		By("Define Install Plan")
		plan2 := installplan.DefineInstallPlan("test-plan-2", invalidNamespace)

		By(CreateInstallPlanInNamespaceStr + plan2.Namespace)
		err = globalhelper.CreateInstallPlan(plan2)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCaseFailed, randomReportDir)
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
			[]string{"installplans.operators.coreos.com", "subscriptions.operators.coreos.com"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Define Install Plan")
		plan := installplan.DefineInstallPlan("test-plan", randomNamespace)

		By(CreateInstallPlanInNamespaceStr + plan.Namespace)
		err = globalhelper.CreateInstallPlan(plan)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		By("Define and create subscription")
		testSub := subscription.DefineSubscription("test-sub", randomNamespace)
		err = globalhelper.CreateSubscription(randomNamespace, testSub)
		Expect(err).ToNot(HaveOccurred(), "Error creating subscription")

		DeferCleanup(func() {
			err = globalhelper.DeleteSubscription(randomNamespace, testSub.Name)
			Expect(err).ToNot(HaveOccurred(), "Error deleting subscription")
		})

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 52098
	It("two custom resources of different CRDs, one in invalid namespace [negative]", func() {
		By("Create Additional Valid Namespace")
		additionalValidNamespace := tsparams.AdditionalValidNamespace + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(additionalValidNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		DeferCleanup(func() {
			err = globalhelper.DeleteNamespaceAndWait(additionalValidNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting namespace")
		})

		By("Define tnf config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace, "tnf"},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"installplans.operators.coreos.com", "subscriptions.operators.coreos.com"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Define Install Plan")
		plan := installplan.DefineInstallPlan("test-plan", randomNamespace)

		By(CreateInstallPlanInNamespaceStr + plan.Namespace)
		err = globalhelper.CreateInstallPlan(plan)
		Expect(err).ToNot(HaveOccurred(), "Error creating installplan")

		By("Define and create subscription")
		testSub := subscription.DefineSubscription("test-sub", additionalValidNamespace)
		err = globalhelper.CreateSubscription(additionalValidNamespace, testSub)
		Expect(err).ToNot(HaveOccurred(), "Error creating subscription")

		DeferCleanup(func() {
			err = globalhelper.DeleteSubscription(additionalValidNamespace, testSub.Name)
			Expect(err).ToNot(HaveOccurred(), "Error deleting subscription")
		})

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespace,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestCaseNameAccessControlNamespace+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespace,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})
})

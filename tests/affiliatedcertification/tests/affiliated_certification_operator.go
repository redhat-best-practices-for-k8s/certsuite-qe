package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcerthelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcertparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	utils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"
)

var _ = Describe("Affiliated-certification operator certification,", func() {

	var (
		installedLabeledOperators []affiliatedcertparameters.OperatorLabelInfo
		csvsToDelete              []affiliatedcertparameters.CsvInfo
	)

	execute.BeforeAll(func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)

		err := globalhelper.DefineTnfConfig(
			[]string{affiliatedcertparameters.TestCertificationNameSpace},
			[]string{affiliatedcertparameters.TestPodLabel},
			[]string{},
			[]string{})

		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Create namespace")
		err = namespaces.Create(affiliatedcertparameters.TestCertificationNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		By("Deploy OperatorGroup")
		err = affiliatedcerthelper.DeployOperatorGroup(affiliatedcertparameters.TestCertificationNameSpace,
			utils.DefineOperatorGroup("affiliatedcert-test-operator-group",
				affiliatedcertparameters.TestCertificationNameSpace,
				[]string{affiliatedcertparameters.TestCertificationNameSpace}),
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operatorgroup")

	})

	AfterEach(func() {
		By("Remove labels from operators")
		for _, info := range installedLabeledOperators {
			err := affiliatedcerthelper.DeleteLabelFromInstalledCSV(
				info.OperatorPrefix,
				info.Namespace,
				info.Label)
			Expect(err).ToNot(HaveOccurred(), "Error removing label from operator "+info.OperatorPrefix)
		}
		installedLabeledOperators = nil

		By("Delete csvs marked for deletion")
		for _, csv := range csvsToDelete {
			err := affiliatedcerthelper.DeleteCsv(csv.OperatorPrefix, csv.Namespace)
			Expect(err).ToNot(HaveOccurred(), "Error deleting csv "+csv.OperatorPrefix)
		}
		csvsToDelete = nil

	})

	// 46699
	It("one operator to test, operator does not belong to certified-operators organization in Red Hat catalog [skip]",
		func() {
			// operator is already installed.
			// not deleting csv yet because it is also used in the next test case

			By("Label operator to be certified")

			err := affiliatedcerthelper.AddLabelToInstalledCSV(
				affiliatedcertparameters.UncertifiedOperatorPrefixNginx,
				affiliatedcertparameters.ExistingOperatorNamespace,
				affiliatedcertparameters.OperatorLabel)
			Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
				affiliatedcertparameters.UncertifiedOperatorPrefixNginx)

			By("Start test")

			err = globalhelper.LaunchTests(
				affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
				affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
				globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
				affiliatedcertparameters.TestCaseOperatorSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

			By("Verify test case status in Junit and Claim reports")

			err = globalhelper.ValidateIfReportsAreValid(
				affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
				globalparameters.TestCaseSkipped)
			Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
		})

	// 46697
	It("two operators to test, one belongs to certified-operators organization and its version is certified,"+
		" one does not belong to certified-operators organization in Red Hat catalog [skip]", func() {
		By("Deploy operator to test")

		certifiedOperatorPostgresSubscription := utils.DefineSubscription("crunchy-postgres-operator-subscription",
			affiliatedcertparameters.TestCertificationNameSpace, "v5", "crunchy-postgres-operator",
			affiliatedcertparameters.CertifiedOperatorGroup, affiliatedcertparameters.OperatorSourceNamespace)

		err := affiliatedcerthelper.DeployOperator(affiliatedcertparameters.TestCertificationNameSpace,
			certifiedOperatorPostgresSubscription)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres)

		By("Confirm that operator is installed and ready")

		Eventually(func() bool {
			err = affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace, "pgo")

			return err == nil
		}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
			"Operator "+affiliatedcertparameters.CertifiedOperatorPrefixPostgres+" is not ready")

		By("Label operator to be certified")

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres)

		// add nginx csv to list to be deleted after test case
		// (only needed for nginx because label removal is not working for this csv)
		csvsToDelete = append(csvsToDelete, affiliatedcertparameters.CsvInfo{
			OperatorPrefix: affiliatedcertparameters.UncertifiedOperatorPrefixNginx,
			Namespace:      affiliatedcertparameters.ExistingOperatorNamespace,
		})

		// add postgres operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46582
	It("one operator to test, operator belongs to certified-operators organization in Red Hat catalog"+
		" and its version is certified", func() {
		// postgres operator is already installed
		By("Label operator to be certified")

		err := affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres)

		// add postgres operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46696
	It("two operators to test, both belong to certified-operators organization in Red Hat catalog and their"+
		" versions are certified", func() {
		By("Deploy additional operator to test")

		certifiedOperatorDatadogSubscription := utils.DefineSubscription("datadog-operator-subscription",
			affiliatedcertparameters.TestCertificationNameSpace, "alpha", "datadog-operator-certified",
			affiliatedcertparameters.CertifiedOperatorGroup, affiliatedcertparameters.OperatorSourceNamespace)

		err := affiliatedcerthelper.DeployOperator(affiliatedcertparameters.TestCertificationNameSpace,
			certifiedOperatorDatadogSubscription)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixDatadog)

		By("Confirm that operator is installed and ready")

		Eventually(func() bool {
			err = affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
				"datadog-operator-manager")

			return err == nil
		}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
			"Operator "+affiliatedcertparameters.CertifiedOperatorPrefixDatadog+" is not ready")

		By("Label operators to be certified")

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixDatadog,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixDatadog)

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres)

		// add datadog and postgres operators info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixDatadog,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().TestText),
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46695
	It("one operator to test, operator is not certified [negative]", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.UncertifiedOperatorBarFoo}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46698
	It("no operators are labeled for testing [skip]", func() {
		By("Start test")

		err := globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46700
	It("name and organization fields exist in certifiedoperatorinfo but are empty [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.EmptyFieldsContainerOrOperator}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

})

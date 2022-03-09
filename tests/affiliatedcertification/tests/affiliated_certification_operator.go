package tests

import (
	"time"

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

		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file: %w", err)

		By("Create namespace")

		err = namespaces.Create(affiliatedcertparameters.TestCertificationNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

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
			// operator is already installed
			// add csv to list to be deleted after test case
			// (only needed for this test case because label removal is not working for this csv)
			csvsToDelete = append(csvsToDelete, affiliatedcertparameters.CsvInfo{
				OperatorPrefix: affiliatedcertparameters.UncertifiedOperatorPrefixNginx,
				Namespace:      affiliatedcertparameters.ExistingOperatorNamespace,
			})

			By("Label operator to be certified")

			err := affiliatedcerthelper.AddLabelToInstalledCSV(
				affiliatedcertparameters.UncertifiedOperatorPrefixNginx,
				affiliatedcertparameters.ExistingOperatorNamespace,
				affiliatedcertparameters.OperatorLabel)
			Expect(err).ToNot(HaveOccurred(), "Error labeling operator")

			By("Start test")

			err = globalhelper.LaunchTests(
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

	// 46582
	It("one operator to test, operator belongs to certified-operators organization in Red Hat catalog"+
		" and its version is certified", func() {
		By("Deploy operators to test")

		operatorGroup := utils.DefineOperatorGroup("affiliatedcert-test-operator-group",
			affiliatedcertparameters.TestCertificationNameSpace,
			[]string{affiliatedcertparameters.TestCertificationNameSpace})

		certifiedOperatorPostgresSubscription := utils.DefineSubscription("crunchy-postgres-operator-subscription",
			affiliatedcertparameters.TestCertificationNameSpace, "v5", "crunchy-postgres-operator",
			affiliatedcertparameters.CertifiedOperatorGroup, affiliatedcertparameters.OperatorSourceNamespace)

		err := affiliatedcerthelper.DeployOperator(affiliatedcertparameters.TestCertificationNameSpace,
			operatorGroup,
			certifiedOperatorPostgresSubscription)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator")

		By("Confirm that operator is installed and ready")

		Eventually(func() bool {
			err = affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace, "pgo")

			return err == nil
		}, 5*time.Minute, 5*time.Second).Should(Equal(true), "Operator is not ready")

		By("Label operator to be certified")

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator")

		// add info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		By("Start test")

		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
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

	// 46696
	It("two operators to test, both are certified", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.CertifiedOperatorApicast,
				affiliatedcertparameters.CertifiedOperatorKubeturbo}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46697
	It("two operators to test, one is certified, one is not [negative]", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.CertifiedOperatorApicast,
				affiliatedcertparameters.UncertifiedOperatorBarFoo}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46698
	It("certifiedoperatorinfo field exists in tnf_config but has no value [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{""}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46700
	It("name and organization fields exist in certifiedoperatorinfo but are empty [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.EmptyFieldsContainerOrOperator}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46702
	It("name field in certifiedoperatorinfo field is populated but organization field is not [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.OperatorNameOnlyKubeturbo}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46704
	It("organization field in certifiedoperatorinfo field is populated but name field is not [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.OperatorOrgOnlyCertifiedOperators}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46706
	It("two operators to test, one is certified, one has empty name and organization fields", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.CertifiedOperatorApicast,
				affiliatedcertparameters.EmptyFieldsContainerOrOperator}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

})

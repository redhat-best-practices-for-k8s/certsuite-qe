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
)

var _ = Describe("Affiliated-certification operator certification,", func() {

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

	BeforeEach(func() {

	})

	// 46699
	It("one operator to test, operator does not belong to certified-operators organization in Red Hat catalog [skip]",
		func() {
			// operator is already installed
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

			By("Remove label from operator")
			err = affiliatedcerthelper.DeleteLabelFromInstalledCSV(
				affiliatedcertparameters.UncertifiedOperatorPrefixNginx,
				affiliatedcertparameters.ExistingOperatorNamespace,
				affiliatedcertparameters.OperatorLabel)
			Expect(err).ToNot(HaveOccurred(), "Error removing label from operator")
		})

	// 46582
	It("one operator to test, operator belongs to certified-operators organization in Red Hat catalog"+
		"and its version is certified", func() {
		By("Deploy operators to test")

		err := affiliatedcerthelper.DeployOperator(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorGroup,
			affiliatedcertparameters.CertifiedOperatorPostgresSubscription)
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

		By("Remove label from operator")
		err = affiliatedcerthelper.DeleteLabelFromInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error removing label from operator")
	})

	// 46695
	It("one operator to test, operator is not certified [negative]", func() {
		Skip("Under development to match new functionality")
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.UncertifiedOperatorBarFoo}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46696
	It("two operators to test, both are certified", func() {
		Skip("Under development to match new functionality")
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.CertifiedOperatorApicast,
				affiliatedcertparameters.CertifiedOperatorKubeturbo}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46697
	It("two operators to test, one is certified, one is not [negative]", func() {
		Skip("Under development to match new functionality")
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
		Skip("Under development to match new functionality")
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{affiliatedcertparameters.CertifiedOperatorApicast,
				affiliatedcertparameters.EmptyFieldsContainerOrOperator}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

})

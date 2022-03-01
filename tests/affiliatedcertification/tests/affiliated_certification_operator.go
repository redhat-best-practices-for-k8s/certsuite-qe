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

	// 46582
	It("one operator to test, operator belongs to certified-operators organization in Red Hat catalog"+
		"and its version is certified", func() {
		By("Deploy operator to be certified")

		err := globalhelper.DeployOperator(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorGroup,
			affiliatedcertparameters.CertifiedOperatorPostgresSubscription)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator")

		By("Confirm that operator is installed and ready")

		Eventually(globalhelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace, "pgo"),
			5*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())

		By("Label operator to be certified")

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			"postgresoperator",
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)

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

	// 46699
	It("certifiedoperatorinfo field does not exist in tnf_config [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunOperatorCertTest(
			[]string{}, globalparameters.TestCaseSkipped)
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

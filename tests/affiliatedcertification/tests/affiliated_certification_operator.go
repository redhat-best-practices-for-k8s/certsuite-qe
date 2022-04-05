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
		if affiliatedcerthelper.IsOperatorGroupInstalled(affiliatedcertparameters.OperatorGroupName,
			affiliatedcertparameters.TestCertificationNameSpace) != nil {
			err = affiliatedcerthelper.DeployOperatorGroup(affiliatedcertparameters.TestCertificationNameSpace,
				utils.DefineOperatorGroup(affiliatedcertparameters.OperatorGroupName,
					affiliatedcertparameters.TestCertificationNameSpace,
					[]string{affiliatedcertparameters.TestCertificationNameSpace}),
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operatorgroup")
		}

		By("Deploy operators for testing")
		// falcon-operator: not in certified-operators group in catalog, for negative test cases
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.UncertifiedOperatorDeploymentFalcon) != nil {
			err = affiliatedcerthelper.DeployAndVerifyOperatorSubscription(
				"falcon-operator",
				"alpha",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CommunityOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
				affiliatedcertparameters.UncertifiedOperatorDeploymentFalcon,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.UncertifiedOperatorPrefixFalcon)
		}

		// crunchy-postgres-operator: in certified-operators group and version is certified
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.CertifiedOperatorDeploymentPostgres) != nil {
			err = affiliatedcerthelper.DeployAndVerifyOperatorSubscription(
				"crunchy-postgres-operator",
				"v5",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CertifiedOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
				affiliatedcertparameters.CertifiedOperatorDeploymentPostgres,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.CertifiedOperatorPrefixPostgres)
		}

		// datadog-operator-certified: in certified-operatos group and version is certified
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.CertifiedOperatorDeploymentDatadog) != nil {
			err = affiliatedcerthelper.DeployAndVerifyOperatorSubscription(
				"datadog-operator-certified",
				"alpha",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CertifiedOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
				affiliatedcertparameters.CertifiedOperatorDeploymentDatadog,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.CertifiedOperatorPrefixDatadog)
		}
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
	})

	// 46699
	It("one operator to test, operator does not belong to certified-operators organization in Red Hat catalog [negative]",
		func() {
			By("Label operator to be certified")

			err := affiliatedcerthelper.AddLabelToInstalledCSV(
				affiliatedcertparameters.UncertifiedOperatorPrefixFalcon,
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.OperatorLabel)
			Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
				affiliatedcertparameters.UncertifiedOperatorPrefixFalcon)

			// add falcon operator info to array for cleanup in AfterEach
			installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
				OperatorPrefix: affiliatedcertparameters.UncertifiedOperatorPrefixFalcon,
				Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
				Label:          affiliatedcertparameters.OperatorLabel,
			})

			By("Start test")

			err = globalhelper.LaunchTests(
				affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
				globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
				affiliatedcertparameters.TestCaseOperatorSkipRegEx,
			)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

			By("Verify test case status in Junit and Claim reports")

			err = globalhelper.ValidateIfReportsAreValid(
				affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
				globalparameters.TestCaseFailed)
			Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
		})

	// 46697
	It("two operators to test, one belongs to certified-operators organization and its version is certified,"+
		" one does not belong to certified-operators organization in Red Hat catalog [negative]", func() {
		By("Label operators to be certified")

		err := affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixPostgres)

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.UncertifiedOperatorPrefixFalcon,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.UncertifiedOperatorPrefixFalcon)

		// add postgres and falcon operators info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixPostgres,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.UncertifiedOperatorPrefixFalcon,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46582
	It("one operator to test, operator belongs to certified-operators organization in Red Hat catalog"+
		" and its version is certified", func() {
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
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
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
		By("Label operators to be certified")

		err := affiliatedcerthelper.AddLabelToInstalledCSV(
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
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
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
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.UncertifiedOperatorBarFoo}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46698
	It("no operators are labeled for testing [skip]", func() {
		By("Start test")

		err := globalhelper.LaunchTests(
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
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
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.EmptyFieldsContainerOrOperator}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

})

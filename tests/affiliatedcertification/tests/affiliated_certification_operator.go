package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcerthelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcertparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
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

		By("Deploy OperatorGroup if not already deployed")
		if affiliatedcerthelper.IsOperatorGroupInstalled(affiliatedcertparameters.OperatorGroupName,
			affiliatedcertparameters.TestCertificationNameSpace) != nil {
			err = affiliatedcerthelper.DeployOperatorGroup(affiliatedcertparameters.TestCertificationNameSpace,
				utils.DefineOperatorGroup(affiliatedcertparameters.OperatorGroupName,
					affiliatedcertparameters.TestCertificationNameSpace,
					[]string{affiliatedcertparameters.TestCertificationNameSpace}),
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operatorgroup")
		}

		By("Deploy operators for testing if not already deployed")
		// falcon-operator: not in certified-operators group in catalog, for negative test cases
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.UncertifiedOperatorPrefixFalcon) != nil {
			err = affiliatedcerthelper.DeployOperatorSubscription(
				"falcon-operator",
				"alpha",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CommunityOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.UncertifiedOperatorPrefixFalcon)
			// confirm that operator is installed and ready
			Eventually(func() bool {
				err = affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
					affiliatedcertparameters.UncertifiedOperatorPrefixFalcon)

				return err == nil
			}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
				affiliatedcertparameters.UncertifiedOperatorPrefixFalcon+" is not ready.")
		}
		// add falcon operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.UncertifiedOperatorPrefixFalcon,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		// kubeturbo-certified: in certified-operators group and version is certified
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo) != nil {
			err = affiliatedcerthelper.DeployOperatorSubscription(
				"kubeturbo-certified",
				"stable",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CertifiedOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo)
			// confirm that operator is installed and ready
			Eventually(func() bool {
				err = affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
					affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo)

				return err == nil
			}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
				affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo+" is not ready.")
		}
		// add kubeturbo operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		// openshiftartifactoryha-operator: in certified-operators group and version is certified
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa) != nil {
			err = affiliatedcerthelper.DeployOperatorSubscription(
				"openshiftartifactoryha-operator",
				"alpha",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CertifiedOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa)
			// confirm that operator is installed and ready
			Eventually(func() bool {
				err = affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
					affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa)

				return err == nil
			}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
				affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa+" is not ready.")
		}
		// add openshiftartifactoryha operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})

		By("Deploy operator with uncertified version if not already deployed")
		// k10-kasten-operator: in certified-operators group, version is not certified
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.UncertifiedOperatorPrefixK10) != nil {

			By("Deploy alternate operator catalog source")

			err = affiliatedcerthelper.DisableCatalogSource(affiliatedcertparameters.CertifiedOperatorGroup)
			Expect(err).ToNot(HaveOccurred(), "Error disabling "+
				affiliatedcertparameters.CertifiedOperatorGroup+" catalog source")
			Eventually(func() bool {
				stillEnabled := affiliatedcerthelper.IsCatalogSourceEnabled(
					affiliatedcertparameters.CertifiedOperatorGroup,
					"openshift-marketplace")

				return !stillEnabled
			}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
				"Default catalog source is still enabled")

			err = affiliatedcerthelper.DeployRHCertifiedOperatorSource("4.5")
			Expect(err).ToNot(HaveOccurred(), "Error deploying catalog source")

			err = affiliatedcerthelper.DeployOperatorSubscription(
				affiliatedcertparameters.UncertifiedOperatorPrefixK10,
				"stable",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CertifiedOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.UncertifiedOperatorPrefixK10)
			// confirm that operator is installed and ready
			Eventually(func() bool {
				err = affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
					affiliatedcertparameters.UncertifiedOperatorPrefixK10)

				return err == nil
			}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
				affiliatedcertparameters.UncertifiedOperatorPrefixK10+" is not ready.")

			By("Re-enable default catalog source")
			err = affiliatedcerthelper.DisableCatalogSource(affiliatedcertparameters.CertifiedOperatorGroup)
			Expect(err).ToNot(HaveOccurred(), "Error disabling catalog source "+affiliatedcertparameters.CertifiedOperatorGroup)
			err = affiliatedcerthelper.EnableCatalogSource(affiliatedcertparameters.CertifiedOperatorGroup)
			Expect(err).ToNot(HaveOccurred(), "Error enabling default catalog source")
		}
		// add kasten-k10 operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.UncertifiedOperatorPrefixK10,
			Namespace:      affiliatedcertparameters.TestCertificationNameSpace,
			Label:          affiliatedcertparameters.OperatorLabel,
		})
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
	})

	// 46699
	It("one operator to test, operator is not in certified-operators organization [negative]",
		func() {
			By("Label operator to be certified")

			err := affiliatedcerthelper.AddLabelToInstalledCSV(
				affiliatedcertparameters.UncertifiedOperatorPrefixFalcon,
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.OperatorLabel)
			Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
				affiliatedcertparameters.UncertifiedOperatorPrefixFalcon)

			By("Start test")

			err = globalhelper.LaunchTests(
				affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
				globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
			Expect(err).To(HaveOccurred(), "Error running "+
				affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

			By("Verify test case status in Junit and Claim reports")

			err = globalhelper.ValidateIfReportsAreValid(
				affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
				globalparameters.TestCaseFailed)
			Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
		})

	// 46697
	It("two operators to test, one is in certified-operators organization and its version is certified,"+
		" one is not in certified-operators organization [negative]", func() {
		By("Label operators to be certified")

		err := affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo)

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.UncertifiedOperatorPrefixFalcon,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.UncertifiedOperatorPrefixFalcon)

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46582
	It("one operator to test, operator is in certified-operators organization"+
		" and its version is certified", func() {
		By("Label operator to be certified")

		err := affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo)

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46696
	It("two operators to test, both are in certified-operators organization and their"+
		" versions are certified", func() {
		By("Label operators to be certified")

		err := affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa)

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixKubeturbo)

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46695
	It("one operator to test, operator is in certified-operators organization but its version"+
		" is not certified [negative]", func() {

		By("Label operator to be certified")

		err := affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.UncertifiedOperatorPrefixK10,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.UncertifiedOperatorPrefixK10)

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46700
	It("two operators to test, both are in certified-operators organization,"+
		" one’s version is certified, the other’s is not [negative]", func() {
		By("Label operators to be certified")

		err := affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.UncertifiedOperatorPrefixK10,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.UncertifiedOperatorPrefixK10)

		err = affiliatedcerthelper.AddLabelToInstalledCSV(
			affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa)

		By("Start test")

		err = globalhelper.LaunchTests(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")

	})

	// 46698
	It("no operators are labeled for testing [skip]", func() {
		By("Start test")

		err := globalhelper.LaunchTests(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName+" test")

		By("Verify test case status in Junit and Claim reports")

		err = globalhelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

})

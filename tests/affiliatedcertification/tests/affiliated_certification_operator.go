package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
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

		By("Clean test namespace")
		err := namespaces.Clean(affiliatedcertparameters.TestCertificationNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred(),
			"Error cleaning namespace "+affiliatedcertparameters.TestCertificationNameSpace)

		By("Ensure default catalog source is enabled")
		Expect(affiliatedcerthelper.IsCatalogSourceEnabled(affiliatedcertparameters.CertifiedOperatorGroup,
			"openshift-marketplace",
			"Certified Operators")).To(BeTrue(), "Default catalog source "+
			affiliatedcertparameters.CertifiedOperatorGroup+" is not enabled")

		By("Define config file " + globalparameters.DefaultTnfConfigFileName)
		err = globalhelper.DefineTnfConfig(
			[]string{affiliatedcertparameters.TestCertificationNameSpace},
			[]string{affiliatedcertparameters.TestPodLabel},
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

		By("Deploy operators for testing")
		// falcon-operator: not in certified-operators group in catalog, for negative test cases
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.UncertifiedOperatorPrefixFalcon) != nil {
			err = affiliatedcerthelper.DeployOperatorSubscription(
				"falcon-operator",
				"alpha",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CommunityOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
				"",
				v1alpha1.ApprovalAutomatic,
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

		// infinibox-operator: in certified-operators group and version is certified
		if affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.CertifiedOperatorPrefixInfinibox) != nil {
			err = affiliatedcerthelper.DeployOperatorSubscription(
				"infinibox-operator-certified",
				"stable",
				affiliatedcertparameters.TestCertificationNameSpace,
				affiliatedcertparameters.CertifiedOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
				affiliatedcertparameters.CertifiedOperatorFullInfinibox,
				v1alpha1.ApprovalManual,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.CertifiedOperatorPrefixInfinibox)

			var installPlan *v1alpha1.InstallPlan
			Eventually(func() bool {
				installPlan, err = affiliatedcerthelper.GetInstallPlanByCSV(affiliatedcertparameters.TestCertificationNameSpace,
					affiliatedcertparameters.CertifiedOperatorFullInfinibox)
				if err == nil {
					return installPlan.Status.Phase != "" &&
						installPlan.Status.Phase != v1alpha1.InstallPlanPhasePlanning
				}

				return false
			}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
				affiliatedcertparameters.CertifiedOperatorPrefixInfinibox+" install plan is not ready.")

			err = affiliatedcerthelper.ApproveInstallPlan(affiliatedcertparameters.TestCertificationNameSpace,
				installPlan)

			Expect(err).ToNot(HaveOccurred(), "Error approving installplan for "+
				affiliatedcertparameters.CertifiedOperatorPrefixInfinibox)
			// confirm that operator is installed and ready
			Eventually(func() bool {
				err = affiliatedcerthelper.IsOperatorInstalled(affiliatedcertparameters.TestCertificationNameSpace,
					affiliatedcertparameters.CertifiedOperatorPrefixInfinibox)

				return err == nil
			}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
				affiliatedcertparameters.CertifiedOperatorPrefixInfinibox+" is not ready.")
		}
		// add infinibox operator info to array for cleanup in AfterEach
		installedLabeledOperators = append(installedLabeledOperators, affiliatedcertparameters.OperatorLabelInfo{
			OperatorPrefix: affiliatedcertparameters.CertifiedOperatorPrefixInfinibox,
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
				affiliatedcertparameters.CertifiedOperatorFullArtifactoryHa,
				v1alpha1.ApprovalManual,
			)
			Expect(err).ToNot(HaveOccurred(), "Error deploying operator "+
				affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa)

			var installPlan *v1alpha1.InstallPlan
			Eventually(func() bool {
				installPlan, err = affiliatedcerthelper.GetInstallPlanByCSV(affiliatedcertparameters.TestCertificationNameSpace,
					affiliatedcertparameters.CertifiedOperatorFullArtifactoryHa)

				if err == nil {
					return installPlan.Status.Phase != "" &&
						installPlan.Status.Phase != v1alpha1.InstallPlanPhasePlanning
				}

				return false
			}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
				affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa+" install plan is not ready.")

			err = affiliatedcerthelper.ApproveInstallPlan(affiliatedcertparameters.TestCertificationNameSpace,
				installPlan)

			Expect(err).ToNot(HaveOccurred(), "Error approving installplan for "+
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
			affiliatedcertparameters.CertifiedOperatorPrefixInfinibox,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixInfinibox)

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
			affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixArtifactoryHa)

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
			affiliatedcertparameters.CertifiedOperatorPrefixInfinibox,
			affiliatedcertparameters.TestCertificationNameSpace,
			affiliatedcertparameters.OperatorLabel)
		Expect(err).ToNot(HaveOccurred(), "Error labeling operator "+
			affiliatedcertparameters.CertifiedOperatorPrefixInfinibox)

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

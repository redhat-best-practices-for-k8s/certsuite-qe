package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcertparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/nethelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Affiliated-certification container certification,", func() {

	execute.BeforeAll(func() {

	})

	BeforeEach(func() {

	})

	// 46562
	It("one container to test, container is certified", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.CertifiedContainerNodeJsUbi})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46563
	It("one container to test, container is not certified [negative]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.UncertifiedContainerFooBar})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46564
	It("two containers to test, both are certified", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.CertifiedContainerNodeJsUbi,
				affiliatedcertparameters.CertifiedContainerRhel7OpenJdk})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46565
	It("two containers to test, one is certified, one is not [negative]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.UncertifiedContainerFooBar,
				affiliatedcertparameters.CertifiedContainerNodeJsUbi})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46566
	It("certifiedcontainerinfo field exists in tnf_config but has no value [skip]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.ContainerFieldNoSubFields})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46567
	It("certifiedcontainerinfo field does not exist in tnf_config [skip]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46578
	It("name and repository fields exist in certifiedcontainerinfo field but are empty [skip]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.EmptyFieldsContainer})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46579
	It("name field in certifiedcontainerinfo field is populated but repository field is not [skip]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.ContainerNameOnlyRhel7OpenJdk})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46580
	It("repository field in certifiedcontainerinfo field is populated but name field is not [skip]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.ContainerRepoOnlyOpenJdk})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46581
	It("two containers listed in certifiedcontainerinfo field, one is certified, one has empty name and "+
		"repository fields", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.CertifiedContainerNodeJsUbi,
				affiliatedcertparameters.EmptyFieldsContainer})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

})

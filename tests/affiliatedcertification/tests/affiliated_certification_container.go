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
		By("Add container information to tnf_config.yml")
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.CertifiedContainer1})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46563
	It("one container to test, container is not certified [negative]", func() {
		By("Add container information to tnf_config.yml")
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.UncertifiedContainer1})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46564
	It("two containers to test, both are certified", func() {
		By("Add container information to tnf_config.yml")
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.CertifiedContainer1, affiliatedcertparameters.CertifiedContainer2})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46565
	It("two containers to test, one is certified, one is not [negative]", func() {
		By("Add container information to tnf_config.yml")
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{affiliatedcertparameters.UncertifiedContainer1, affiliatedcertparameters.CertifiedContainer1})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseContainerSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46566
	It("certifiedcontainerinfo field exists in tnf_config but has no value [skip]", func() {
		Skip("Under development")
	})

	// 46567
	It("certifiedcontainerinfo field does not exist in tnf_config [skip]", func() {
		Skip("Under development")
	})

	// 46578
	It("name and repository fields exist in certifiedcontainerinfo field but are empty [skip]", func() {
		Skip("Under development")
	})

	// 46579
	It("name field in certifiedcontainerinfo field is populated but repository field is not [skip]", func() {
		Skip("Under development")
	})

	// 46580
	It("repository field in certifiedcontainerinfo field is populated but name field is not [skip]", func() {
		Skip("Under development")
	})

	// 46581
	It("two containers listed in certifiedcontainerinfo field, one is certified, one has empty name and "+
		"repository fields", func() {
		Skip("Under development")
	})

})

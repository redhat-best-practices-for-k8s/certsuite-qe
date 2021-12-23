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

var _ = Describe("Affiliated-certification operator certification,", func() {

	execute.BeforeAll(func() {

	})

	BeforeEach(func() {

	})

	// 46582
	It("one operator to test, operator is certified", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{},
			[]string{affiliatedcertparameters.CertifiedOperatorApicast})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46695
	It("one operator to test, operator is not certified [negative]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{},
			[]string{affiliatedcertparameters.UncertifiedOperatorBarFoo})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46696
	It("two operators to test, both are certified", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{},
			[]string{affiliatedcertparameters.CertifiedOperatorApicast,
				affiliatedcertparameters.CertifiedOperatorKubeturbo})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46697
	It("two operators to test, one is certified, one is not [negative]", func() {
		By("Add container information to " + globalparameters.DefaultTnfConfigFileName)
		err := globalhelper.DefineTnfConfig(
			[]string{netparameters.TestNetworkingNameSpace},
			[]string{netparameters.TestPodLabel},
			[]string{},
			[]string{affiliatedcertparameters.CertifiedOperatorApicast,
				affiliatedcertparameters.UncertifiedOperatorBarFoo})
		Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")

		By("Start test")
		err = globalhelper.LaunchTests(
			[]string{affiliatedcertparameters.AffiliatedCertificationTestSuiteName},
			affiliatedcertparameters.TestCaseOperatorSkipRegEx,
		)
		Expect(err).To(HaveOccurred(), "Error running "+
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName+" test")

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred(), "Error validating test reports")
	})

	// 46698
	It("certifiedoperatorinfo field exists in tnf_config but has no value [skip]", func() {
		Skip("Under development")
	})

	// 46699
	It("certifiedoperatorinfo field does not exist in tnf_config [skip]", func() {
		Skip("Under development")
	})

	// 46700
	It("name and organization fields exist in certifiedoperatorinfo but are empty [skip]", func() {
		Skip("Under development")
	})

	// 46702
	It("name field in certifiedoperatorinfo field is populated but organization field is not [skip]", func() {
		Skip("Under development")
	})

	// 46704
	It("organization field in certifiedoperatorinfo field is populated but name field is not [skip]", func() {
		Skip("Under development")
	})

	// 46706
	It("two operators to test, one is certified, one has empty name and organization fields", func() {
		Skip("Under development")
	})

})

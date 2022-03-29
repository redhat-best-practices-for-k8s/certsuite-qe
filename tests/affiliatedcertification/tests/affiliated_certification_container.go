package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcerthelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcertparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Affiliated-certification container certification,", func() {

	execute.BeforeAll(func() {

	})

	BeforeEach(func() {

	})

	// 46562
	It("one container to test, container is certified", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.CertifiedContainerNodeJsUbi}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46563
	It("one container to test, container is not certified [negative]", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.UncertifiedContainerFooBar}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46564
	It("two containers to test, both are certified", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.CertifiedContainerNodeJsUbi,
				affiliatedcertparameters.CertifiedContainerRhel7OpenJdk}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46565
	It("two containers to test, one is certified, one is not [negative]", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.UncertifiedContainerFooBar,
				affiliatedcertparameters.CertifiedContainerNodeJsUbi}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46566
	It("certifiedcontainerinfo field exists in tnf_config but has no value [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{""}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46567
	It("certifiedcontainerinfo field does not exist in tnf_config [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46578
	It("name and repository fields exist in certifiedcontainerinfo field but are empty [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.EmptyFieldsContainerOrOperator}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46579
	It("name field in certifiedcontainerinfo field is populated but repository field is not [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.ContainerNameOnlyRhel7OpenJdk}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46580
	It("repository field in certifiedcontainerinfo field is populated but name field is not [skip]", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.ContainerRepoOnlyOpenJdk}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46581
	It("two containers listed in certifiedcontainerinfo field, one is certified, one has empty name and "+
		"repository fields", func() {
		err := affiliatedcerthelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText),
			[]string{affiliatedcertparameters.CertifiedContainerNodeJsUbi,
				affiliatedcertparameters.EmptyFieldsContainerOrOperator}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

})

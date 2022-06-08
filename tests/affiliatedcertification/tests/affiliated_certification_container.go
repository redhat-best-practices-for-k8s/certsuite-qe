package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

var _ = Describe("Affiliated-certification container certification,", func() {

	execute.BeforeAll(func() {

	})

	BeforeEach(func() {

	})

	// 46562
	It("one container to test, container is certified", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{tsparams.CertifiedContainerCockroachDB}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46563
	It("one container to test, container is not certified [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{tsparams.UncertifiedContainerNodeJs12}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46564
	It("two containers to test, both are certified", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{tsparams.CertifiedContainerCockroachDB,
				tsparams.CertifiedContainer5gc}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46565
	It("two containers to test, one is certified, one is not [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{tsparams.UncertifiedContainerNodeJs12,
				tsparams.CertifiedContainerCockroachDB}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46566
	It("certifiedcontainerinfo field exists in tnf_config but has no value [skip]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{""}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46567
	It("certifiedcontainerinfo field does not exist in tnf_config [skip]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46578
	It("name and repository fields exist in certifiedcontainerinfo field but are empty [skip]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{tsparams.EmptyFieldsContainer}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46579
	It("name field in certifiedcontainerinfo field is populated but repository field is not [skip]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{tsparams.ContainerNameOnlyCockroachDB}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46580
	It("repository field in certifiedcontainerinfo field is populated but name field is not [skip]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{tsparams.ContainerRepoOnlyRedHatRegistry}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46581
	It("two containers listed in certifiedcontainerinfo field, one is certified, one has empty name and "+
		"repository fields", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			[]string{tsparams.CertifiedContainer5gc,
				tsparams.EmptyFieldsContainer}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

})

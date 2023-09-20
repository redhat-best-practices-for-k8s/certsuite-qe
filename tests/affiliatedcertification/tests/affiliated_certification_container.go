package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

var _ = Describe("Affiliated-certification container certification,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestCertificationNameSpace)
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})

	// 46562
	It("one container to test, container is certified", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{tsparams.CertifiedContainerCockroachDB}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46563
	It("one container to test, container is not certified [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{tsparams.UncertifiedContainerNodeJs12}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46564
	It("two containers to test, both are certified", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{tsparams.CertifiedContainerCockroachDB,
				tsparams.CertifiedContainer5gc}, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46565
	It("two containers to test, one is certified, one is not [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{tsparams.UncertifiedContainerNodeJs12,
				tsparams.CertifiedContainerCockroachDB}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46566
	It("certifiedcontainerinfo field exists in tnf_config but has no value [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{""}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46567
	It("certifiedcontainerinfo field does not exist in tnf_config [skip]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{}, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46578
	It("name and repository fields exist in certifiedcontainerinfo field but are empty [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{tsparams.EmptyFieldsContainer}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46579
	It("name field in certifiedcontainerinfo field is populated but repository field is not [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{tsparams.ContainerNameOnlyCockroachDB}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46580
	It("repository field in certifiedcontainerinfo field is populated but name field is not [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{tsparams.ContainerRepoOnlyRedHatRegistry}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 46581
	It("two containers listed in certifiedcontainerinfo field, one is certified, one has empty name and "+
		"repository fields [negative]", func() {
		err := tshelper.SetUpAndRunContainerCertTest(
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomNamespace,
			[]string{tsparams.CertifiedContainer5gc,
				tsparams.EmptyFieldsContainer}, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

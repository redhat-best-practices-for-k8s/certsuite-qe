package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

var _ = Describe("Affiliated-certification container-is-certified-digest,", func() {

	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.TestCertificationNameSpace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestCertificationNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.TestCertificationNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 66765
	It("one container to test, container is certified", func() {
		By("Define deployment with certified container")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", tsparams.TestCertificationNameSpace,
			tsparams.CertifiedContainerUrlNodeJs, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66766
	It("one container to test, container is not certified [negative]", func() {
		By("Define deployment with uncertified container")

		dep := deployment.DefineDeployment("affiliated-cert-deployment", tsparams.TestCertificationNameSpace,
			tsparams.UncertifiedContainerUrlCnfTest, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66767
	It("two containers to test, both are certified", func() {
		By("Define deployments with certified containers")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", tsparams.TestCertificationNameSpace,
			tsparams.CertifiedContainerUrlNodeJs, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2 := deployment.DefineDeployment("affiliated-cert-deployment-2", tsparams.TestCertificationNameSpace,
			tsparams.CertifiedContainerUrlCockroachDB, tsparams.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66768
	It("two containers to test, one is certified, one is not [negative]", func() {
		By("Define deployments with different container certification statuses")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", tsparams.TestCertificationNameSpace,
			tsparams.UncertifiedContainerUrlCnfTest, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2 := deployment.DefineDeployment("affiliated-cert-deployment-2", tsparams.TestCertificationNameSpace,
			tsparams.CertifiedContainerUrlCockroachDB, tsparams.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

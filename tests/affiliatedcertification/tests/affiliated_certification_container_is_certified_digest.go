package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

var _ = Describe("Affiliated-certification container-is-certified-digest,", Serial, func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestCertificationNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})

	// 66765
	It("one container to test, container is certified digest", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skip test due to image pull missing credentials in Kind")
		}

		By("Define deployment with certified container")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
			tsparams.CertifiedContainerURLNodeJs, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66766
	It("one container to test, container is not certified digest [negative]", func() {
		By("Define deployment with uncertified container")

		dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
			tsparams.UncertifiedContainerURLCnfTest, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66767
	It("two containers to test, both are certified digest", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skip test due to image pull missing credentials in Kind")
		}

		By("Define deployments with certified containers")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
			tsparams.CertifiedContainerURLNodeJs, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2 := deployment.DefineDeployment("affiliated-cert-deployment-2", randomNamespace,
			tsparams.CertifiedContainerURLCockroachDB, tsparams.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66768
	It("two containers to test, one is certified, one is not digest [negative]", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skip test due to image pull missing credentials in Kind")
		}

		By("Define deployments with different container certification statuses")
		dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
			tsparams.UncertifiedContainerURLCnfTest, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2 := deployment.DefineDeployment("affiliated-cert-deployment-2", randomNamespace,
			tsparams.CertifiedContainerURLCockroachDB, tsparams.TestDeploymentLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameContainerDigest,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameContainerDigest,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

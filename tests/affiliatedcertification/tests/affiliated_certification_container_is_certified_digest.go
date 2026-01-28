package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"

	"github.com/redhat-best-practices-for-k8s/oct/pkg/certdb/onlinecheck"
)

var _ = Describe("Affiliated-certification container-is-certified-digest,", Serial,
	Label("affiliatedcertification", "ocp-required"), func() {
		var randomNamespace string
		var randomReportDir string
		var randomCertsuiteConfigDir string

		BeforeEach(func() {
			if globalhelper.IsKindCluster() {
				Skip("Skip test due to image pull missing credentials in Kind")
			}

			// Create random namespace and keep original report and certsuite config directories
			randomNamespace, randomReportDir, randomCertsuiteConfigDir =
				globalhelper.BeforeEachSetupWithRandomNamespace(
					tsparams.TestCertificationNameSpace)

			By("Define certsuite config file")
			err := globalhelper.DefineCertsuiteConfig(
				[]string{randomNamespace},
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")

			By("Check if the test image is certified prior to deployment")
			// Using the 'oct' repo, we should do a quick assertion to see if the image is available
			// and certified.
			onlineValidator := onlinecheck.NewOnlineValidator()

			// The information for this is gathered from:
			//nolint:lll
			// https://catalog.redhat.com/api/containers/v1/images?filter=image_id==sha256:41bc5b622db8b5e0d608e7524c39928b191270666252578edbf1e0f84a9e3cab
			//nolint:lll
			Expect(onlineValidator.IsContainerCertified("registry.access.redhat.com", "ubi8/nodejs-12", "latest", "sha256:41bc5b622db8b5e0d608e7524c39928b191270666252578edbf1e0f84a9e3cab")).To(BeTrue())
		})

		AfterEach(func() {
			globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
				randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
		})

		// 66765
		It("one container to test, container is certified digest", func() {
			By("Define deployment with certified container")
			dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
				tsparams.CertifiedContainerURLNodeJs, tsparams.TestDeploymentLabels)

			By("Create deployment")
			err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())

			By("Assert deployment is ready")
			runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
			Expect(err).ToNot(HaveOccurred())
			Expect(runningDeployment).ToNot(BeNil())

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameContainerDigest,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameContainerDigest,
				globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

		// 66766
		It("one container to test, container is not certified digest [negative]", func() {
			By("Define deployment with uncertified container")

			dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
				tsparams.UncertifiedContainerURLCnfTest, tsparams.TestDeploymentLabels)

			By("Create deployment")
			err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())

			By("Assert deployment is ready")
			runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
			Expect(err).ToNot(HaveOccurred())
			Expect(runningDeployment).ToNot(BeNil())

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameContainerDigest,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameContainerDigest,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

		// 66767
		It("two containers to test, both are certified digest", func() {
			By("Define deployments with certified containers")
			dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
				tsparams.CertifiedContainerURLNodeJs, tsparams.TestDeploymentLabels)

			By("Create deployment 1")
			err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())

			dep2 := deployment.DefineDeployment("affiliated-cert-deployment-2", randomNamespace,
				tsparams.CertifiedContainerURLCockroachDB, tsparams.TestDeploymentLabels)

			By("Create deployment 2")
			err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameContainerDigest,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameContainerDigest,
				globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

		// 66768
		It("two containers to test, one is certified, one is not digest [negative]", func() {
			// Note: This test uses container images, not operators, so it's not affected by
			// the OCP 4.20 certified-operators catalog issue (see issue #1283)

			By("Define deployments with different container certification statuses")
			dep := deployment.DefineDeployment("affiliated-cert-deployment", randomNamespace,
				tsparams.UncertifiedContainerURLCnfTest, tsparams.TestDeploymentLabels)

			By("Create deployment 1")
			err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())

			By("Assert deployment1 is ready")
			runningDeployment1, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
			Expect(err).ToNot(HaveOccurred())
			Expect(runningDeployment1).ToNot(BeNil())

			dep2 := deployment.DefineDeployment("affiliated-cert-deployment-2", randomNamespace,
				tsparams.CertifiedContainerURLCockroachDB, tsparams.TestDeploymentLabels)

			By("Create deployment 2")
			err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())

			By("Assert deployment2 is ready")
			runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
			Expect(err).ToNot(HaveOccurred())
			Expect(runningDeployment2).ToNot(BeNil())

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameContainerDigest,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameContainerDigest,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

	})

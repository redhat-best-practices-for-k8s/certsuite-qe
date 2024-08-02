package operator

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/operator/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/operator/parameters"
)

var _ = Describe("Operator pods have runAs userid", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.OperatorNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			tsparams.TnfTargetCrdFilters, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("Operator pods should have runAs userid", func() {
		// Deploy an operator that has runAs userid
		By("Deploy operator group")
		err := tshelper.DeployTestOperatorGroup(randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Query the packagemanifest for the " + tsparams.CertifiedOperatorPrefixNginx)
		version, err := globalhelper.QueryPackageManifestForVersion(tsparams.CertifiedOperatorPrefixNginx, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for nginx-ingress-operator")

		By(fmt.Sprintf("Deploy nginx-ingress-operator%s for testing", "."+version))
		// nginx-ingress-operator: in certified-operators group and version is certified
		err = tshelper.DeployOperatorSubscription(
			tsparams.CertifiedOperatorPrefixNginx,
			"alpha",
			randomNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorPrefixNginx+".v"+version,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			tsparams.CertifiedOperatorPrefixNginx)

		err = waitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixNginx,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+tsparams.CertifiedOperatorPrefixNginx+".v"+version+
			" is not ready")

		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.CertifiedOperatorPrefixNginx,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.CertifiedOperatorPrefixNginx)

		By("Assert that the manager pod has runAs userid")
		controllerPod, err := globalhelper.GetControllerPodFromOperator(randomNamespace, tsparams.CertifiedOperatorPrefixNginx)
		Expect(err).ToNot(HaveOccurred(), "Error getting controller pod")

		for _, container := range controllerPod.Spec.Containers {
			Expect(container.SecurityContext).ToNot(BeNil())
			Expect(container.SecurityContext.RunAsUser).ToNot(BeNil())
			Expect(*container.SecurityContext.RunAsUser).ToNot(Equal(0))
		}

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfOperatorPodRunAsUserID,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfOperatorPodRunAsUserID,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Operator pods do not have runAs userid [negative]", func() {
		// Deploy an operator that has runAs userid

		// TODO: Find an operator that does not have runAs userid
	})
})

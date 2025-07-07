package operator

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
)

var _ = Describe("Operator bundle count,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.OperatorNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{"nginxingresses.charts.nginx.org"}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("CatalogSource with less than 1000 bundle images", func() {
		By("Create custom-operator catalog source")
		err := globalhelper.DeployCustomOperatorSource("quay.io/redhat-best-practices-for-k8s/qe-custom-catalog")
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			err := globalhelper.DeleteCustomOperatorSource()
			Expect(err).ToNot(HaveOccurred())
		})

		By("Deploy operator group for namespace " + randomNamespace)
		err = tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By(fmt.Sprintf("Deploy first operator (nginx-ingress-operator) for testing - channel %s", "new"))
		err = tshelper.DeployOperatorSubscription(
			"operator1",
			"nginx-ingress-operator",
			"new",
			randomNamespace,
			"custom-catalog",
			tsparams.OperatorSourceNamespace,
			"nginx-ingress-operator.v3.0.1",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			"nginx-ingress-operator")

		By("Wait until operator is ready")
		err = tshelper.WaitUntilOperatorIsReady("nginx-ingress-operator", randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+"nginx-ingress-operator"+
			" is not ready")

		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				"nginx-ingress-operator",
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+"nginx-ingress-operator")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorBundleCount,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorBundleCount,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

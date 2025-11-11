package tests

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"

	crdutils "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/crd"
)

var _ = Describe("lifecycle-crd-scaling", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		if globalhelper.IsVanillaK8sCluster() {
			By("Make masters schedulable")
			err := nodes.EnableMasterScheduling(globalhelper.GetAPIClient().Nodes(), true)
			Expect(err).ToNot(HaveOccurred())

			By("Enable intrusive tests")
			err = os.Setenv("CERTSUITE_NON_INTRUSIVE_ONLY", "false")
			Expect(err).ToNot(HaveOccurred())
		}

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.CertsuiteTargetOperatorLabels},
			[]string{},
			[]string{tsparams.CertsuiteTargetCrdFilters}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")

		if globalhelper.GetConfiguration().General.DisableIntrusiveTests == strings.ToLower("true") {
			Skip("Intrusive tests are disabled via config")
		}
	})

	AfterEach(func() {
		By("Disable intrusive tests")
		err := os.Setenv("CERTSUITE_NON_INTRUSIVE_ONLY", "true")
		Expect(err).ToNot(HaveOccurred())

		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	It("Custom resource is deployed, scale in and out", func() {
		// We have to pre-install the crd-operator-scaling resources prior to running these tests.
		By("Check if cr-scale-operator is installed")
		exists, err := globalhelper.NamespaceExists(tsparams.CertsuiteTargetOperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "error checking if cr-scaling-operator is installed")
		if !exists {
			// Skip the test if cr-scaling-operator is not installed
			Skip("cr-scale-operator is not installed, skipping test")
		}

		By("Create a scale custom resource")
		_, err = crdutils.CreateCustomResourceScale(tsparams.CertsuiteCustomResourceName, randomNamespace,
			tsparams.CertsuiteTargetOperatorLabels, tsparams.CertsuiteTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdScaling,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdScaling, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

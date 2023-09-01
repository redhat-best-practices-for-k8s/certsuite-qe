package tests

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"

	crdutils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/crd"
)

var _ = Describe("lifecycle-crd-scaling", Serial, func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		if globalhelper.IsKindCluster() {
			By("Make masters schedulable")
			err := nodes.EnableMasterScheduling(globalhelper.GetAPIClient().CoreV1Interface, true)
			Expect(err).ToNot(HaveOccurred())

			By("Enable intrusive tests")
			err = os.Setenv("TNF_NON_INTRUSIVE_ONLY", "false")
			Expect(err).ToNot(HaveOccurred())
		}

		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			[]string{tsparams.TnfTargetCrdFilters})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		if globalhelper.GetConfiguration().General.DisableIntrusiveTests == strings.ToLower("true") {
			Skip("Intrusive tests are disabled via config")
		}
	})

	AfterEach(func() {
		By("Disable intrusive tests")
		err := os.Setenv("TNF_NON_INTRUSIVE_ONLY", "true")
		Expect(err).ToNot(HaveOccurred())

		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	It("Custom resource is deployed, scale in and out", func() {
		// We have to pre-install the crd-operator-scaling resources prior to running these tests.
		By("Check if cr-scale-operator is installed")
		exists, err := namespaces.Exists(tsparams.TnfTargetOperatorNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "error checking if cr-scaling-operator is installed")
		if !exists {
			// Skip the test if cr-scaling-operator is not installed
			Skip("cr-scale-operator is not installed, skipping test")
		}

		By("Create a scale custom resource")
		_, err = crdutils.CreateCustomResourceScale(tsparams.TnfCustomResourceName, randomNamespace,
			tsparams.TnfTargetOperatorLabels, tsparams.TnfTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.TnfCrdScaling,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdScaling, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

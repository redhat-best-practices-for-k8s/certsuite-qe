package tests

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-statefulset-scaling", Serial, func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		By("Enable intrusive tests")
		err := os.Setenv("TNF_NON_INTRUSIVE_ONLY", "false")
		Expect(err).ToNot(HaveOccurred())

		if globalhelper.IsKindCluster() {
			By("Make masters schedulable")
			err := nodes.EnableMasterScheduling(globalhelper.GetAPIClient().Nodes(), true)
			Expect(err).ToNot(HaveOccurred())
		}

		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define TNF config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())

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

	// 45439
	It("One statefulSet, one pod", func() {
		By("Define statefulSet")
		statefulset := tshelper.DefineStatefulSet(tsparams.TestStatefulSetName, randomNamespace)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("start lifecycle-statefulset-scaling test")
		err = globalhelper.LaunchTests(tsparams.TnfStatefulSetScalingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfStatefulSetScalingTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

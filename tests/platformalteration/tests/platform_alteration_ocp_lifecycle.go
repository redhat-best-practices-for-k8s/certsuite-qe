package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
)

var _ = Describe("platform-alteration-ocp-lifecycle", func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.PlatformAlterationNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.WaitingTime)
	})

	It("OCP version should be supported", func() {
		if globalhelper.IsKindCluster() {
			Skip("OCP version is not applicable for Kind cluster")
		}

		By("Start platform-alteration-ocp-lifecycle test")
		err := globalhelper.LaunchTests(tsparams.TnfOCPLifecycleName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfOCPLifecycleName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

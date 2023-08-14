package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-crd-scaling", func() {
	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.LifecycleNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			[]string{tsparams.TnfTargetCrdFilters})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Check if crd-scaling-operator is installed")
		exists, err := namespaces.Exists(tsparams.TnfTargetOperatorNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error installing crd-scaling-operator")
		Expect(exists).To(BeTrue())
	})

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	It("Crd deployed, scale in and out", func() {
		By("Create a scale custom resource")
		_, err := tshelper.CreateCustomResourceScale(tsparams.TnfCustomResourceName, tsparams.LifecycleNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.TnfCrdScaling,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdScaling, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})
})

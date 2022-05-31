package tests

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle-statefulset-scaling", func() {

	BeforeEach(func() {
		err := helper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(parameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Enable intrusive tests")
		err = os.Setenv("TNF_NON_INTRUSIVE_ONLY", "false")
		Expect(err).ToNot(HaveOccurred())
	})

	// 45439
	It("One statefulSet, one pod", func() {
		By("Define statefulSet")
		statefulset := helper.DefineStatefulSet("lifecyclesf")
		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulset, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("start lifecycle-statefulset-scaling test")
		err = globalhelper.LaunchTests(
			parameters.TnfStatefulSetScalingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TnfStatefulSetScalingTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

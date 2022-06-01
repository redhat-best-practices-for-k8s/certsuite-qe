package tests

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle-deployment-scaling", func() {

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

	// 47398
	It("One deployment, one pod, one container, scale in and out", func() {

		By("Define Deployment")
		deploymenta, err := helper.DefineDeployment(1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
			deploymenta, parameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-deployment-scaling test")
		err = globalhelper.LaunchTests(parameters.TnfDeploymentScalingTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(parameters.TnfDeploymentScalingTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

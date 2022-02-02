package tests

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle lifecycle-scaling", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv("TNF_NON_INTRUSIVE_ONLY", "false")
		Expect(err).ToNot(HaveOccurred())
	})

	// 47398
	It("One deployment, one pod, one container, scale in & out", func() {

		By("Define Deployment")
		deploymentStruct, err := lifehelper.DefineDeployment(false, 1, 1, "lifecycleput")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("start lifecycle lifecycle-scaling")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButScalingRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ScalingDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	AfterEach(func() {
		err := os.Unsetenv("TNF_NON_INTRUSIVE_ONLY")
		Expect(err).ToNot(HaveOccurred())
	})

})

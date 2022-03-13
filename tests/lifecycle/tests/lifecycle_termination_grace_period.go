package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"k8s.io/utils/pointer"
)

var _ = Describe("lifecycle lifecycle-termination-grace-period", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47399
	It("One deployment, one pod, terminationGracePeriodSeconds is not set", func() {

		By("Define deployment without terminationGracePeriodSeconds")
		deploymentStruct := lifehelper.RemoveterminationGracePeriod(
			lifehelper.DefineDeployment(1, 1, "lifecycleput"))

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-termination-grace-period test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			globalhelper.GetStringOfSkipTcs(lifeparameters.SkipsTcsSlice, lifeparameters.TerminationGracePeriodName))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TerminationGracePeriodName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47400
	It("One deployment, one pod, terminationGracePeriodSeconds is set", func() {

		By("Define deployment with terminationGracePeriodSeconds specified")
		deploymentStruct := deployment.RedefineWithTerminationGracePeriod(
			lifehelper.DefineDeployment(1, 1, "lifecycleput"), pointer.Int64Ptr(45))

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-termination-grace-period test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			globalhelper.GetStringOfSkipTcs(lifeparameters.SkipsTcsSlice, lifeparameters.TerminationGracePeriodName))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TerminationGracePeriodName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47401
	It("Multiple deployments, replica > 1, terminationGracePeriodSeconds is set on one deployment", func() {

		By("Define first deployment with terminationGracePeriodSeconds specified")
		firstDeployment := deployment.RedefineWithTerminationGracePeriod(
			lifehelper.DefineDeployment(3, 1, "lifecycleputone"), pointer.Int64Ptr(45))

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(firstDeployment, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment without terminationGracePeriodSeconds")
		secondDeployment := lifehelper.RemoveterminationGracePeriod(
			lifehelper.DefineDeployment(3, 1, "lifecycleputtwo"))

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(secondDeployment, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-pod-termination-grace-period test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			globalhelper.GetStringOfSkipTcs(lifeparameters.SkipsTcsSlice, lifeparameters.TerminationGracePeriodName))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.TerminationGracePeriodName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})
})

package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Access-control pod-automount-service-account-token , ", func() {

	execute.BeforeAll(func() {

	})

	// 53033
	It("one deployment, one pod, token false", func() {
		By("Define deployment with automountServiceAccountToken set to false")
		dep, err := helper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		dep = deployment.RedefineWithAutomountServiceAccountToken(dep, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53034
	It("one deployment, one pod, token true [negative]", func() {
		By("Define deployment with automountServiceAccountToken set to false")
		dep, err := helper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		dep = deployment.RedefineWithAutomountServiceAccountToken(dep, true)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodAutomountToken,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlNamespace,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53035
	It("one deployment, one pod, token not set, service account's token false", func() {

	})

	// 53036
	It("one deployment, one pod, token not set, service account's token true [negative]", func() {

	})

	// 53040
	It("one deployment, one pod, token not set, service account's token not set [negative]", func() {

	})

	// 53054
	It("one deployment, one pod, token false, service account's token true", func() {

	})

	// 53036
	It("two deployments, one pod each, tokens false", func() {

	})

	// 53057
	It("two deployments, one pod each, one token true [negative]", func() {

	})

})

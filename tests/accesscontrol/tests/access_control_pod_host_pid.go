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
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control pod-host-pid ", func() {

	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(parameters.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

	})

	// 53140
	It("one deployment, one pod, HostPid false", func() {
		By("Define deployment with hostPid set to false")
		dep, err := helper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPid(dep, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostPid,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostPid,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53141
	It("one deployment, one pod, HostPid true [negative]", func() {
		By("Define deployment with hostPid set to true")
		dep, err := helper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPid(dep, true)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostPid,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostPid,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53142
	It("two deployments, one pod each, HostPids false", func() {
		By("Define deployments with hostPid set to false")
		dep, err := helper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPid(dep, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := helper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPid(dep2, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostPid,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostPid,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53143
	It("two deployments, one pod each, one HostPid true [negative]", func() {
		By("Define deployments with hostPid set to different values")
		dep, err := helper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPid(dep, true)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := helper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPid(dep2, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostPid,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostPid,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

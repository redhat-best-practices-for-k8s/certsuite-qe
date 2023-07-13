package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control pod-host-ipc, ", func() {

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

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(parameters.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53140
	It("one deployment, one pod, HostIpc false", func() {
		By("Define deployment with hostIPC set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostIpc(dep, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostIpc,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostIpc,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53141
	It("one deployment, one pod, HostIpc true [negative]", func() {
		By("Define deployment with hostIPC set to true")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostIpc(dep, true)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostIpc,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostIpc,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53142
	It("two deployments, one pod each, HostIpcs false", func() {
		By("Define deployments with hostIPC set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostIpc(dep, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostIpc(dep2, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostIpc,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostIpc,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53143
	It("two deployments, one pod each, one HostIpc true [negative]", func() {
		By("Define deployments with hostIPC set to different values")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostIpc(dep, true)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostIpc(dep2, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostIpc,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostIpc,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

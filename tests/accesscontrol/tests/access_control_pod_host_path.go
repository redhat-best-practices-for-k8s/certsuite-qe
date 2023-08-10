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

var _ = Describe("Access-control pod-host-path, ", func() {

	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{parameters.TestAccessControlNameSpace},
			[]string{parameters.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(parameters.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

	})

	// 53939
	It("one deployment, one pod, HostPath not set", func() {
		By("Define deployment with hostPath set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostPath,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostPath,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53940
	It("one deployment, one pod, HostPath set [negative]", func() {
		By("Define deployment with hostPath set to true")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPath(dep, "volume", "mnt/data")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostPath,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostPath,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53941
	It("two deployments, one pod each, HostPath not set", func() {
		By("Define deployments with hostPath set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostPath,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostPath,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53946
	It("two deployments, one pod each, one HostPath set [negative]", func() {
		By("Define deployments with hostPath set to different values")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPath(dep2, "volume", "mnt/data")

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, parameters.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			parameters.TestCaseNameAccessControlPodHostPath,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			parameters.TestCaseNameAccessControlPodHostPath,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

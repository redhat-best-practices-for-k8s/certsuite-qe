package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control requests-and-limits,", Serial, func() {

	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.TestAccessControlNameSpace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 55021
	It("one deployment, one container with all requests and limits set", func() {
		By("Define deployment with requests and limits set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithAllRequestsAndLimits(dep, tsparams.MemoryLimit, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55022
	It("one deployment, one container with no limits set [negative]", func() {
		By("Define deployment with no limits set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithResourceRequests(dep, tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55023
	It("one deployment, one container with no limits or requests set [negative]", func() {
		By("Define deployment with no limits or requests set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55025
	It("one deployment, one container with CPU limits not set [negative]", func() {
		By("Define deployment with CPU limits not set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithMemoryRequestsAndLimitsAndCPURequest(dep, tsparams.MemoryLimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55026
	It("one deployment, one container with memory limits not set [negative]", func() {
		By("Define deployment with memory limits not set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithMemoryRequestAndCPURequestsAndLimits(dep, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55027
	It("two deployments, one container each with all limits and requests set", func() {
		By("Define deployments with requests and limits set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithAllRequestsAndLimits(dep, tsparams.MemoryLimit, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithAllRequestsAndLimits(dep2, tsparams.MemoryLimit, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55028
	It("two deployments, one container each, one with memory limits not set [negative]", func() {
		By("Define deployments with memory limits not set on one")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithAllRequestsAndLimits(dep, tsparams.MemoryLimit, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithMemoryRequestAndCPURequestsAndLimits(dep2, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

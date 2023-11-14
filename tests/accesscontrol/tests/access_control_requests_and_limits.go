package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
)

var _ = Describe("Access-control requests-and-limits,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})

	// 55021
	It("one deployment, one container with all requests and limits set", func() {
		By("Define deployment with requests and limits set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithAllRequestsAndLimits(dep, tsparams.MemoryLimit, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment has requests and limits set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal(tsparams.CPULimit))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal(tsparams.MemoryLimit))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal(tsparams.CPURequest))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal(tsparams.MemoryRequest))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55022
	It("one deployment, one container with no limits set [negative]", func() {
		By("Define deployment with no limits set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithResourceRequests(dep, tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment has no limits set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal("0"))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal("0"))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55023
	It("one deployment, one container with no limits or requests set [negative]", func() {
		By("Define deployment with no limits or requests set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has no limits or requests set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal("0"))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal("0"))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal("0"))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal("0"))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55025
	It("one deployment, one container with CPU limits not set [negative]", func() {
		By("Define deployment with CPU limits not set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithMemoryRequestsAndLimitsAndCPURequest(dep, tsparams.MemoryLimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment has CPU limits not set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal("0"))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal(tsparams.CPURequest))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55026
	It("one deployment, one container with memory limits not set [negative]", func() {
		By("Define deployment with memory limits not set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithMemoryRequestAndCPURequestsAndLimits(dep, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment has memory limits not set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal("0"))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal(tsparams.MemoryRequest))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55027
	It("two deployments, one container each with all limits and requests set", func() {
		By("Define deployments with requests and limits set")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithAllRequestsAndLimits(dep, tsparams.MemoryLimit, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment1 has requests and limits set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal(tsparams.CPULimit))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal(tsparams.MemoryLimit))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal(tsparams.CPURequest))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal(tsparams.MemoryRequest))

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithAllRequestsAndLimits(dep2, tsparams.MemoryLimit, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment2 has requests and limits set")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal(tsparams.CPULimit))
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal(tsparams.MemoryLimit))
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal(tsparams.CPURequest))
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal(tsparams.MemoryRequest))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55028
	It("two deployments, one container each, one with memory limits not set [negative]", func() {
		By("Define deployments with memory limits not set on one")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithAllRequestsAndLimits(dep, tsparams.MemoryLimit, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment1 has memory limits set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal(tsparams.MemoryLimit))
		Expect(runningDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal(tsparams.MemoryRequest))

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithMemoryRequestAndCPURequestsAndLimits(dep2, tsparams.CPULimit,
			tsparams.MemoryRequest, tsparams.CPURequest)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert that deployment2 has memory limits not set")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal("0"))
		Expect(runningDeployment2.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal(tsparams.MemoryRequest))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRequestsAndLimits,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

package accesscontrol

import (
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	klog "k8s.io/klog/v2"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/statefulset"
)

//
// These test cases assume that the cluster worker nodes have a non-realtime kernel.
// Thus, they need to run in an orderly and serial fashion. Rationale:
//   - ordered: needed so BeforeAll and AfterAll funcs can be used to switch the kernel type.
//   - serial: switching the kernel type of nodes makes MachineConfigPool to reboot them one by one, which may affect other test cases
//     if they're running in parallel.
//

var _ = Describe("Access-control sys-nice_capability", Ordered, Serial,
	Label("rt-kernel"), Label("accesscontrol13", "ocp-required"), func() {
		var randomNamespace string
		var randomReportDir string
		var randomCertsuiteConfigDir string

		// Save the MCP name that manages worker nodes, so:
		// - re-check the kernel type of the first node on every test case.
		// - calculate the dynamic timeout to be used to wait for the nodes to be rebooted.
		var mcpName string
		var mcpWorkerNames []string

		// This flag is needed to avoid the AfterAll/AfterEach blocks if the whole test suite was skipped in BeforeAll.
		skipTestSuite := false

		BeforeAll(func() {
			if globalhelper.GetConfiguration().General.DisableIntrusiveTests == strings.ToLower("true") {
				skipTestSuite = true
				Skip("Intrusive tests are disabled via config")
			}

			// Skip all if running in a kind cluster
			if globalhelper.IsKindCluster() {
				skipTestSuite = true
				Skip("Kind cluster detected")
			}

			By("Getting workers MCP")
			mcpNodes, err := tshelper.GetWorkersMCPs()
			Expect(err).ToNot(HaveOccurred(), "failed to get MCP for workers")
			Expect(len(mcpNodes)).To(BeNumerically(">", 0), "no workers MCP found")

			// This test suite does not support more than one workers MCP.
			if len(mcpNodes) > 1 {
				skipTestSuite = true
				Skip(fmt.Sprintf("more than one workers MCP found: %v", mcpNodes))
			}

			// There should be only one element in the map, so let's get it:
			for name, nodes := range mcpNodes {
				mcpName = name
				mcpWorkerNames = nodes

				break
			}

			// We need to deploy a custom MachineConfig so the MachineConfig Operator can install the realtime version
			// of the current kernel in the worker nodes. But first, we need to know which labels to put on our MC so
			// the workers MCP can use it.
			By(fmt.Sprintf("Selected MCP=%s (workers=%v) Getting target MachineConfig label/s.", mcpNodes, mcpWorkerNames))
			mcLabels, err := tshelper.GetMachineConfigTargetLabels(mcpName)
			Expect(err).ToNot(HaveOccurred(), "failed to get MachineConfig labels")

			timeOut := tsparams.NodeRebootTimeout * time.Duration(len(mcpWorkerNames))
			mcName := tsparams.RelatimeKernelMachineConfigName

			By(fmt.Sprintf("Deploying MachineConfig %q for realtime kernel. Labels=%v. Total nodes reboot timeout=%s, (MCP startup timeout=%s",
				mcName, mcLabels, timeOut, tsparams.McpStartTimeout))
			err = tshelper.DeployRTKernelMachineConfig(mcpName, mcName, mcLabels, timeOut)
			Expect(err).ToNot(HaveOccurred(), "failed to deploy MachineConfig for realtime kernel in node "+mcpWorkerNames[0])
		})

		AfterAll(func() {
			if skipTestSuite {
				return
			}
			// We need to remove the MC that installs the realtime kernel in the MCP workers.
			timeout := tsparams.NodeRebootTimeout * time.Duration(len(mcpWorkerNames))
			mcName := tsparams.RelatimeKernelMachineConfigName
			By("Removing MachineConfig " + mcName + " so workers of MCP " + mcpName + " can switch to default kernel version.")
			err := tshelper.RemoveRTKernelMachineConfig(mcpName, tsparams.RelatimeKernelMachineConfigName, timeout)
			Expect(err).ToNot(HaveOccurred(), "failed to remove MachineConfig "+tsparams.RelatimeKernelMachineConfigName)
		})

		BeforeEach(func() {
			// Create random namespace and keep original report and certsuite config directories
			randomNamespace, randomReportDir, randomCertsuiteConfigDir =
				globalhelper.BeforeEachSetupWithRandomPrivilegedNamespace(
					tsparams.TestAccessControlNameSpace)

			By("Define certsuite config file")
			err := globalhelper.DefineCertsuiteConfig(
				[]string{randomNamespace},
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")
		})

		AfterEach(func() {
			if skipTestSuite {
				return
			}

			globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
				randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
		})

		//
		// Realtime kernel test cases block.
		// The BeforeAll must have switched a suitable workers MCP to use kerneltype "realtime"
		//

		It("One deployment with one pod with sys-nice cap at container level, realtime kernel", func() {
			// Double-check that we have realtime kernel. Using the first node would suffice.
			Expect(mcpWorkerNames).NotTo(BeEmpty(), "worker names slice must not be empty")

			workersUseRealtimeKernel, err := tshelper.HasNodeRtKernel(mcpWorkerNames[0])
			Expect(err).ToNot(HaveOccurred(), "failed to get kernel type from node "+mcpWorkerNames[0])
			Expect(workersUseRealtimeKernel).To(BeTrue(), "no workers with realtime kernel found")

			By("Define deployment with cap sys admin added")
			dep, err := tshelper.DefineDeployment(1, 1, "acdeployment", randomNamespace)
			Expect(err).ToNot(HaveOccurred(), "failed to define deployment")

			deployment.RedefineWithContainersSecurityContextCaps(dep, []string{"SYS_NICE"}, nil)

			By("Create deployment")
			err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "failed to deploy deployment")

			By("Ensure all running pods are adding the SYS_NICE cap")
			podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed getting list of pods in namespace "+randomNamespace)

			Expect(len(podList.Items)).To(Equal(1), "Invalid number of pods in namespace "+randomNamespace)

			pod := &podList.Items[0]
			By("Ensure pod " + pod.Name + " has added SYS_NICE cap")
			Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(Equal([]corev1.Capability{"SYS_NICE"}))

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameAccessControlRtSysNiceCapability,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameAccessControlRtSysNiceCapability,
				globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("One statefulset with two pods with two containers with sys-nice cap at container level, realtime kernel", func() {
			// Double-check that we have realtime kernel. Using the first node would suffice.
			Expect(mcpWorkerNames).NotTo(BeEmpty(), "worker names slice must not be empty")

			workersUseRealtimeKernel, err := tshelper.HasNodeRtKernel(mcpWorkerNames[0])
			Expect(err).ToNot(HaveOccurred(), "failed to get kernel type from node "+mcpWorkerNames[0])
			Expect(workersUseRealtimeKernel).To(BeTrue(), "no workers with realtime kernel found")

			By("Define statefulset with cap sys admin added")
			sts, err := tshelper.DefineStatefulSet(2, 2, "sts-2-2", randomNamespace)
			Expect(err).ToNot(HaveOccurred(), "failed to define statefulset")

			statefulset.RedefineWithContainersSecurityContextCaps(sts, []string{"SYS_NICE"}, nil)

			By("Create statefulset")
			err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(sts, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "failed to deploy statefulset")

			By("Ensure all the statefulset's pods exist in namespace " + randomNamespace)
			var podList *corev1.PodList
			Eventually(func() bool {
				podList, err = globalhelper.GetListOfPodsInNamespace(randomNamespace)
				if err != nil {
					klog.Info(fmt.Sprintf("Failed to list pods in ns %s: %v", randomNamespace, err))

					return false
				}

				if len(podList.Items) == 2 {

					return true
				}

				return false
			}).WithTimeout(5*time.Minute).WithPolling(5*time.Second).
				Should(BeTrue(), "timeout waiting statefulset pods to be running")

			Expect(len(podList.Items)).To(Equal(2), "Invalid number of pods in namespace "+randomNamespace)

			pod := &podList.Items[0]
			By("Ensure pod " + pod.Name + " has added SYS_NICE cap")
			Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(Equal([]corev1.Capability{"SYS_NICE"}))
			Expect(pod.Spec.Containers[1].SecurityContext.Capabilities.Add).To(Equal([]corev1.Capability{"SYS_NICE"}))

			pod = &podList.Items[1]
			By("Ensure pod " + pod.Name + " has added SYS_NICE cap")
			Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(Equal([]corev1.Capability{"SYS_NICE"}))
			Expect(pod.Spec.Containers[1].SecurityContext.Capabilities.Add).To(Equal([]corev1.Capability{"SYS_NICE"}))

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameAccessControlRtSysNiceCapability,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameAccessControlRtSysNiceCapability,
				globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("One deployment with one pod without sys-nice cap at container level, realtime kernel [negative]", func() {
			// Double-check that we have realtime kernel. Using the first node would suffice.
			Expect(mcpWorkerNames).NotTo(BeEmpty(), "worker names slice must not be empty")

			workersUseRealtimeKernel, err := tshelper.HasNodeRtKernel(mcpWorkerNames[0])
			Expect(err).ToNot(HaveOccurred(), "failed to get kernel type from node "+mcpWorkerNames[0])
			Expect(workersUseRealtimeKernel).To(BeTrue(), "no workers with realtime kernel found")

			By("Define deployment with cap sys admin added")
			dep, err := tshelper.DefineDeployment(1, 1, "dep-1-1", randomNamespace)
			Expect(err).ToNot(HaveOccurred(), "failed to define deployment")

			By("Create deployment")
			err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "failed to deploy deployment")

			By("Ensure the pod is not adding the SYS_NICE cap")
			podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed getting list of pods in namespace "+randomNamespace)

			Expect(len(podList.Items)).To(Equal(1), "Invalid number of pods in namespace "+randomNamespace)

			pod := &podList.Items[0]
			By("Ensure pod " + pod.Name + " has not added SYS_NICE cap")
			Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(BeNil())

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameAccessControlRtSysNiceCapability,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameAccessControlRtSysNiceCapability,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("One statefulset with two pods with two containers without sys-nice cap, realtime kernel [negative]", func() {
			// Double-check that we have realtime kernel. Using the first node would suffice.
			Expect(mcpWorkerNames).NotTo(BeEmpty(), "worker names slice must not be empty")

			workersUseRealtimeKernel, err := tshelper.HasNodeRtKernel(mcpWorkerNames[0])
			Expect(err).ToNot(HaveOccurred(), "failed to get kernel type of node "+mcpWorkerNames[0])

			Expect(workersUseRealtimeKernel).To(BeTrue(), "no workers with realtime kernel found")

			By("Define statefulset with cap sys admin added")
			sts, err := tshelper.DefineStatefulSet(2, 2, "sts-2-2", randomNamespace)
			Expect(err).ToNot(HaveOccurred(), "failed to define statefulset")

			By("Create statefulset")
			err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(sts, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "failed to deploy statefulset")

			By("Ensure all the statefulset's pods exist in namespace " + randomNamespace)
			var podList *corev1.PodList
			Eventually(func() bool {
				podList, err = globalhelper.GetListOfPodsInNamespace(randomNamespace)
				if err != nil {
					klog.Info(fmt.Sprintf("Failed to list pods in ns %s: %v", randomNamespace, err))

					return false
				}

				if len(podList.Items) == 2 {

					return true
				}

				return false
			}).WithTimeout(5*time.Minute).WithPolling(5*time.Second).
				Should(BeTrue(), "timeout waiting statefulset pods to be running")

			Expect(len(podList.Items)).To(Equal(2), "Invalid number of pods in namespace "+randomNamespace)

			pod := &podList.Items[0]
			By("Ensure pod " + pod.Name + " has not added SYS_NICE cap")
			Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(BeNil())
			Expect(pod.Spec.Containers[1].SecurityContext.Capabilities.Add).To(BeNil())

			pod = &podList.Items[1]
			By("Ensure pod " + pod.Name + " has not added SYS_NICE cap")
			Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(BeNil())
			Expect(pod.Spec.Containers[1].SecurityContext.Capabilities.Add).To(BeNil())

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameAccessControlRtSysNiceCapability,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameAccessControlRtSysNiceCapability,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

	})

//
// Non-realtime kernel test cases block.
// They can run unordered and in parallel with other tcs.
//

var _ = Describe("Access-control sys-nice_capability check, non-realtime kernel", Label("accesscontrol13"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomPrivilegedNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("No workload [skip]", func() {
		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRtSysNiceCapability,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRtSysNiceCapability,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Testing whether the tc passes since non-realtime kernel is detected, no
	// matter whether the containers have the sys-nice cap or not.
	It("One deployment with one pod with sys-nice cap at container level", func() {
		By("Define deployment with cap sys admin added")
		dep, err := tshelper.DefineDeployment(1, 1, "acdeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextCaps(dep, []string{"SYS_NICE"}, nil)

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all running pods are adding the SYS_NICE cap")
		podList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Failed getting list of pods in namespace "+randomNamespace)

		Expect(len(podList.Items)).To(Equal(1), "Invalid number of pods in namespace "+randomNamespace)

		pod := &podList.Items[0]
		By("Ensure pod " + pod.Name + " has added SYS_NICE cap")
		Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(Equal([]corev1.Capability{"SYS_NICE"}))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRtSysNiceCapability,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRtSysNiceCapability,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Testing whether the tc passes since non-realtime kernel is detected, no
	// matter whether the containers have the sys-nice cap or not.
	It("One statefulset with two pods with two containers without sys-nice cap, non-realtime kernel", func() {
		By("Define statefulset with cap sys admin added")
		sts, err := tshelper.DefineStatefulSet(2, 2, "sts-2-2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create statefulset")
		err = globalhelper.CreateAndWaitUntilStatefulSetIsReady(sts, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all the statefulset's pods exist in namespace " + randomNamespace)
		var podList *corev1.PodList
		Eventually(func() bool {
			podList, err = globalhelper.GetListOfPodsInNamespace(randomNamespace)
			if err != nil {
				klog.Info(fmt.Sprintf("Failed to list pods in ns %s: %v", randomNamespace, err))
				klog.Info(fmt.Sprintf("Failed to list pods in ns %s: %v", randomNamespace, err))

				return false
			}

			if len(podList.Items) == 2 {
				return true
			}

			return false
		}).WithTimeout(5*time.Minute).WithPolling(5*time.Second).
			Should(BeTrue(), "timeout waiting statefulset pods to be running")

		Expect(len(podList.Items)).To(Equal(2), "Invalid number of pods in namespace "+randomNamespace)

		pod := &podList.Items[0]
		By("Ensure pod " + pod.Name + " has not added SYS_NICE cap")
		if globalhelper.IsKindCluster() {
			Expect(pod.Spec.Containers[0].SecurityContext).To(BeNil())
			Expect(pod.Spec.Containers[1].SecurityContext).To(BeNil())
		} else {
			Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(BeNil())
			Expect(pod.Spec.Containers[1].SecurityContext.Capabilities.Add).To(BeNil())
		}

		pod = &podList.Items[1]
		By("Ensure pod " + pod.Name + " has not added SYS_NICE cap")
		if globalhelper.IsKindCluster() {
			Expect(pod.Spec.Containers[0].SecurityContext).To(BeNil())
			Expect(pod.Spec.Containers[1].SecurityContext).To(BeNil())
		} else {
			Expect(pod.Spec.Containers[0].SecurityContext.Capabilities.Add).To(BeNil())
			Expect(pod.Spec.Containers[1].SecurityContext.Capabilities.Add).To(BeNil())
		}

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlRtSysNiceCapability,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlRtSysNiceCapability,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

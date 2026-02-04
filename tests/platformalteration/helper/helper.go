package helper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/client"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/statefulset"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const WaitingTime = 5 * time.Minute

// IsRealOCPCluster checks if the cluster is a real OCP cluster (not CRC/SNO development cluster).
// Real clusters typically have multiple nodes or specific configurations that indicate
// they are production-like environments.
func IsRealOCPCluster() (bool, string) {
	// Check if MCO is available
	mcpList, err := globalhelper.GetAPIClient().MachineConfigPools().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, fmt.Sprintf("cannot access MachineConfigPools: %v", err)
	}

	// Check for worker-cnf pool which indicates a real telco/CNF cluster
	for _, mcp := range mcpList.Items {
		if mcp.Name == "worker-cnf" {
			return true, "found worker-cnf MachineConfigPool"
		}
	}

	// Check number of nodes - CRC typically has only 1 node
	nodesList, err := globalhelper.GetAPIClient().K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, fmt.Sprintf("cannot list nodes: %v", err)
	}

	if len(nodesList.Items) > 1 {
		// Check if we have dedicated worker nodes (not just control-plane nodes)
		workerCount := 0

		for _, node := range nodesList.Items {
			_, isControlPlane := node.Labels["node-role.kubernetes.io/control-plane"]

			_, isMaster := node.Labels["node-role.kubernetes.io/master"]
			if !isControlPlane && !isMaster {
				workerCount++
			}
		}

		if workerCount > 0 {
			return true, fmt.Sprintf("cluster has %d worker nodes", workerCount)
		}
	}

	return false, "appears to be a CRC/SNO development cluster"
}

// DetectBaseImageAlterations checks if the cluster has conditions that would cause
// the certsuite base-image test to fail. This includes checking for:
// - Custom RPM packages installed on nodes via MachineConfig extensions
// - Real OCP clusters (not CRC/SNO) which often have customizations
// Note: Kernel arguments are NOT checked here as they don't affect base image tests.
// Returns true if alterations are detected, along with details.
func DetectBaseImageAlterations() (bool, string) {
	// First check if this is a real OCP cluster - these often have custom packages
	isReal, realDetails := IsRealOCPCluster()
	if isReal {
		return true, "running on real OCP cluster: " + realDetails
	}

	// Check for MachineConfigs with extensions (custom RPMs)
	mcList, err := globalhelper.GetAPIClient().MachineConfigs().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, fmt.Sprintf("cannot access MachineConfigs: %v", err)
	}

	var alterationDetails []string

	for _, machineConfig := range mcList.Items {
		// Only check for extensions (custom RPMs) - these indicate base image modifications
		// Kernel arguments are standard OCP configuration and don't affect base image tests
		if len(machineConfig.Spec.Extensions) > 0 {
			alterationDetails = append(alterationDetails,
				fmt.Sprintf("MachineConfig %s has extensions: %v", machineConfig.Name, machineConfig.Spec.Extensions))
		}
	}

	if len(alterationDetails) > 0 {
		return true, strings.Join(alterationDetails, "; ")
	}

	return false, "no base image alterations detected"
}

// DetectBootParamsAlterations checks if the cluster has conditions that would cause
// the certsuite boot-params test to fail. This includes checking for:
// - Custom kernel arguments in MachineConfigs
// - Performance profiles with custom isolcpus, nohz, etc.
// Returns true if alterations are detected, along with details.
func DetectBootParamsAlterations() (bool, string) {
	mcList, err := globalhelper.GetAPIClient().MachineConfigs().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, fmt.Sprintf("cannot access MachineConfigs: %v", err)
	}

	var alterationDetails []string

	// Known kernel arguments that indicate boot params customization
	customKernelArgs := []string{
		"isolcpus", "nohz", "nohz_full", "rcu_nocbs", "kthread_cpus",
		"irqaffinity", "skew_tick", "intel_pstate", "nosmt", "hugepages",
		"default_hugepagesz", "tsc", "clocksource", "processor.max_cstate",
		"intel_idle.max_cstate", "mce", "audit", "systemd.cpu_affinity",
	}

	for _, machineConfig := range mcList.Items {
		if len(machineConfig.Spec.KernelArguments) > 0 {
			for _, arg := range machineConfig.Spec.KernelArguments {
				for _, customArg := range customKernelArgs {
					if strings.HasPrefix(arg, customArg) {
						alterationDetails = append(alterationDetails,
							fmt.Sprintf("MachineConfig %s has custom kernel arg: %s", machineConfig.Name, arg))

						break
					}
				}
			}
		}
	}

	// Check if this is a real cluster
	isReal, realDetails := IsRealOCPCluster()
	if isReal {
		alterationDetails = append(alterationDetails, "running on real OCP cluster: "+realDetails)
	}

	if len(alterationDetails) > 0 {
		return true, strings.Join(alterationDetails, "; ")
	}

	return false, "no boot params alterations detected"
}

// WaitForSpecificNodeCondition waits for a given node to become ready or not.
func WaitForSpecificNodeCondition(clients *client.ClientSet, timeout, interval time.Duration, nodeName string,
	ready bool) error {
	return wait.PollUntilContextTimeout(context.TODO(), interval, timeout, true,
		func(ctx context.Context) (bool, error) {
			nodesList, err := clients.Nodes().List(ctx, metav1.ListOptions{})
			if err != nil {
				return false, err
			}

			// Verify the node condition
			for _, node := range nodesList.Items {
				if node.Name == nodeName && nodes.IsNodeInCondition(&node, corev1.NodeReady) == ready {
					return true, nil
				}
			}

			return false, nil
		})
}

// UpdateAndVerifyHugePagesConfig updates the hugepages file with a given number,
// and verify that it was updated successfully.
func UpdateAndVerifyHugePagesConfig(updatedHugePagesNumber int, filePath string, pod *corev1.Pod) error {
	cmd := fmt.Sprintf("echo %d > %s", updatedHugePagesNumber, filePath)

	_, err := globalhelper.ExecCommand(*pod, []string{"/bin/bash", "-c", cmd})
	if err != nil {
		return fmt.Errorf("failed to execute command %s: %w", cmd, err)
	}

	// loop to wait until the file has been actually updated.
	timeout := time.Now().Add(5 * time.Minute)

	for {
		currentHugepagesNumber, err := GetHugePagesConfigNumber(filePath, pod)
		if err != nil {
			return fmt.Errorf("failed to get hugepages number: %w", err)
		}

		if updatedHugePagesNumber == currentHugepagesNumber {
			return nil
		}

		if time.Now().After(timeout) {
			return fmt.Errorf("timedout waiting for hugepages to be updated, currently: %d, expected: %d",
				currentHugepagesNumber, updatedHugePagesNumber)
		}

		time.Sleep(tsparams.RetryInterval * time.Second)
	}
}

// GetHugePagesConfigNumber returns hugepages config number from a given file.
func GetHugePagesConfigNumber(file string, pod *corev1.Pod) (int, error) {
	cmd := fmt.Sprintf("cat %s", file)

	buf, err := globalhelper.ExecCommand(*pod, []string{"/bin/bash", "-c", cmd})
	if err != nil {
		return -1, err
	}

	hugepagesNumber, err := strconv.Atoi(strings.Split(buf.String(), "\r\n")[0])
	if err != nil {
		return -1, err
	}

	return hugepagesNumber, nil
}

// ArgListToMap takes a list of strings of the form "key=value" and translate it into a map
// of the form {key: value}.
func ArgListToMap(lst []string) map[string]string {
	retval := make(map[string]string)

	for _, arg := range lst {
		splitArgs := strings.Split(arg, "=")
		if len(splitArgs) == 1 {
			retval[splitArgs[0]] = ""
		} else {
			retval[splitArgs[0]] = splitArgs[1]
		}
	}

	return retval
}

// AppendIstioContainerToPod appends istio-proxy container to a pod.
func AppendIstioContainerToPod(pod *corev1.Pod, image string) {
	pod.Spec.Containers = append(
		pod.Spec.Containers, corev1.Container{
			Name:    "istio-proxy",
			Image:   image,
			Command: []string{"/bin/bash", "-c", "sleep INF"},
		})
}

// Creates deployment with one pod with one non-UBI based container.
func DefineDeploymentWithNonUBIContainer(namespace string) *appsv1.Deployment {
	dep := deployment.DefineDeployment(tsparams.TestDeploymentName, namespace,
		tsparams.NotRedHatRelease, tsparams.CertsuiteTargetPodLabels)

	// Workaround as this non-ubi test image needs /bin/sh (busybox) instead of /bin/bash.
	deployment.RedefineWithContainerSpecs(dep, []corev1.Container{
		{
			Name:    "test",
			Image:   tsparams.NotRedHatRelease,
			Command: []string{"/bin/sh", "-c", "sleep INF"},
		},
	})

	return dep
}

// Creates statefulset with one pod with one non-UBI based container.
func DefineStatefulSetWithNonUBIContainer(namespace string) *appsv1.StatefulSet {
	sts := statefulset.DefineStatefulSet(tsparams.TestStatefulSetName, namespace,
		tsparams.NotRedHatRelease, tsparams.CertsuiteTargetPodLabels)

	// Workaround as this non-ubi test image needs /bin/sh (busybox) instead of /bin/bash.
	statefulset.RedefineWithContainerSpecs(sts, []corev1.Container{
		{
			Name:    "test",
			Image:   tsparams.NotRedHatRelease,
			Command: []string{"/bin/sh", "-c", "sleep INF"},
		},
	})

	return sts
}

// HasMachineConfigKernelArguments checks if any MachineConfig in the cluster has custom kernel arguments.
// Returns true if custom kernel arguments are found, along with a description of what was found.
func HasMachineConfigKernelArguments() (bool, string, error) {
	machineConfigList, err := globalhelper.GetAPIClient().MachineConfigs().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, "", fmt.Errorf("failed to list MachineConfigs: %w", err)
	}

	for _, mc := range machineConfigList.Items {
		if len(mc.Spec.KernelArguments) > 0 {
			return true, fmt.Sprintf("MachineConfig %s has kernelArguments: %v", mc.Name, mc.Spec.KernelArguments), nil
		}
	}

	return false, "No MachineConfigs with custom kernelArguments found", nil
}

// DetectBootParamsAlteration checks if the actual kernel cmdline contains parameters
// that indicate the boot configuration has been altered from defaults.
// It checks for common performance tuning parameters that wouldn't be in a default RHCOS config.
func DetectBootParamsAlteration(kernelCmdline string) (bool, string) {
	// Common parameters that indicate intentional boot param customization
	customParams := []string{
		"nohz=",
		"nohz_full=",
		"rcu_nocbs=",
		"isolcpus=",
		"intel_pstate=",
		"processor.max_cstate=",
		"intel_idle.max_cstate=",
		"skew_tick=",
		"nosoftlockup",
		"tsc=",
		"cgroup_no_v1=",
		"psi=",
		"hugepagesz=",
		"hugepages=",
		"default_hugepagesz=",
	}

	var foundParams []string

	for _, param := range customParams {
		if strings.Contains(kernelCmdline, param) {
			// Extract the full parameter
			for _, part := range strings.Fields(kernelCmdline) {
				if strings.HasPrefix(part, strings.TrimSuffix(param, "=")) {
					foundParams = append(foundParams, part)
				}
			}
		}
	}

	if len(foundParams) > 0 {
		return true, fmt.Sprintf("Detected custom boot parameters: %v", foundParams)
	}

	return false, "No custom boot parameters detected"
}

// IsNodeControlPlane checks if a node is a control plane node based on its name or labels.
func IsNodeControlPlane(nodeName string) bool {
	// Check common control plane naming patterns
	controlPlanePatterns := []string{
		"master",
		"control-plane",
		"ctlplane",
		"controlplane",
	}

	nodeNameLower := strings.ToLower(nodeName)
	for _, pattern := range controlPlanePatterns {
		if strings.Contains(nodeNameLower, pattern) {
			return true
		}
	}

	return false
}

// CheckPodsOnControlPlaneNodes checks if any pods in the list are running on control plane nodes.
// Returns true if any pod is on a control plane node, along with details.
func CheckPodsOnControlPlaneNodes(pods []corev1.Pod) (bool, string) {
	var controlPlanePods []string

	for _, pod := range pods {
		if IsNodeControlPlane(pod.Spec.NodeName) {
			controlPlanePods = append(controlPlanePods, fmt.Sprintf("%s (node: %s)", pod.Name, pod.Spec.NodeName))
		}
	}

	if len(controlPlanePods) > 0 {
		return true, fmt.Sprintf("Pods running on control plane nodes: %v", controlPlanePods)
	}

	return false, "No pods running on control plane nodes"
}

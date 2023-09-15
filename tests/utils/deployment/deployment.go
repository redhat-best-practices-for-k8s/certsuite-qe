package deployment

import (
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type MultusAnnotation struct {
	Name string `json:"name"`
}

// DefineDeployment returns deployment struct.
func DefineDeployment(deploymentName string, namespace string, image string, label map[string]string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:        pointer.Int32(1),
			MinReadySeconds: 30,
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testpod-",
					Labels:    label,
					Namespace: namespace,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: pointer.Int64(0),
					Containers: []corev1.Container{
						{
							Name:    "test",
							Image:   image,
							Command: []string{"/bin/bash", "-c", "sleep INF"},
						},
					},
				},
			},
		},
	}
}

// RedefineAllContainersWithPreStopSpec redefines deployment with requested lifecycle/preStop spec.
func RedefineAllContainersWithPreStopSpec(deployment *appsv1.Deployment, command []string) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].Lifecycle = &corev1.Lifecycle{
			PreStop: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: command,
				},
			},
		}
	}
}

// RedefineWithLabels redefines deployment with additional label.
func RedefineWithLabels(deployment *appsv1.Deployment, label map[string]string) {
	newMap := make(map[string]string)
	for k, v := range deployment.Spec.Template.Labels {
		newMap[k] = v
	}

	for k, v := range label {
		newMap[k] = v
	}

	deployment.Spec.Template.Labels = newMap
}

// RedefineWithMultus redefines deployment with additional labels.
func RedefineWithMultus(deployment *appsv1.Deployment, nadNames []string) *appsv1.Deployment {
	if len(nadNames) < 1 {
		return deployment
	}

	var nadAnnotations []MultusAnnotation

	for _, nadName := range nadNames {
		nadAnnotations = append(nadAnnotations, MultusAnnotation{Name: nadName})
	}

	bString, _ := json.Marshal(nadAnnotations)

	deployment.Spec.Template.Annotations = map[string]string{
		"k8s.v1.cni.cncf.io/networks": string(bString),
	}

	return deployment
}

// RedefineWithReplicaNumber redefines deployment with requested replica number.
func RedefineWithReplicaNumber(deployment *appsv1.Deployment, replicasNumber int32) {
	deployment.Spec.Replicas = pointer.Int32(replicasNumber)
}
func AppendServiceAccount(deployment *appsv1.Deployment, serviceAccountName string) {
	deployment.Spec.Template.Spec.ServiceAccountName = serviceAccountName
}

// RedefineFirstContainerWithPreStopSpec redefines deployment first container with lifecycle/preStop spec.
func RedefineFirstContainerWithPreStopSpec(deployment *appsv1.Deployment, command []string) error {
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		deployment.Spec.Template.Spec.Containers[0].Lifecycle = &corev1.Lifecycle{
			PreStop: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: command}}}

		return nil
	}

	return fmt.Errorf("deployment %s does not have any containers", deployment.Name)
}

// RedefineWithTerminationGracePeriod redefines deployment with terminationGracePeriod spec.
func RedefineWithTerminationGracePeriod(deployment *appsv1.Deployment, terminationGracePeriod *int64) {
	deployment.Spec.Template.Spec.TerminationGracePeriodSeconds = terminationGracePeriod
}

// RedefineWithPodAntiAffinity redefines deployment with podAntiAffinity spec.
func RedefineWithPodAntiAffinity(deployment *appsv1.Deployment, label map[string]string) {
	deployment.Spec.Template.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: label,
					},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		}}
}

func RedefineWithImagePullPolicy(deployment *appsv1.Deployment, pullPolicy corev1.PullPolicy) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].ImagePullPolicy = pullPolicy
	}
}

func RedefineWithNodeSelector(deployment *appsv1.Deployment, nodeSelector map[string]string) {
	deployment.Spec.Template.Spec.NodeSelector = nodeSelector
}

func RedefineWithNodeAffinity(deployment *appsv1.Deployment, key string) {
	deployment.Spec.Template.Spec.Affinity = &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      key,
								Operator: corev1.NodeSelectorOpExists,
							},
						},
					},
				},
			},
		}}
}

func RedefineWithReadinessProbe(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

func RedefineWithLivenessProbe(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

// RedefineWithStartUpProbe adds startup probe to deployment manifest.
func RedefineWithStartUpProbe(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].StartupProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

func RedefineWithContainerSpecs(deployment *appsv1.Deployment, containerSpecs []corev1.Container) {
	deployment.Spec.Template.Spec.Containers = containerSpecs
}

func RedefineWithPrivilegedContainer(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"ALL"}},
		}
	}
}

func RedefineWithHostPid(deployment *appsv1.Deployment, hostPid bool) {
	deployment.Spec.Template.Spec.HostPID = hostPid
}

func RedefineWithHostIpc(deployment *appsv1.Deployment, hostIpc bool) {
	deployment.Spec.Template.Spec.HostIPC = hostIpc
}

func RedefineWithAutomountServiceAccountToken(deployment *appsv1.Deployment, token bool) {
	deployment.Spec.Template.Spec.AutomountServiceAccountToken = &token
}

func RedefineWithHostNetwork(deployment *appsv1.Deployment, hostNetwork bool) {
	deployment.Spec.Template.Spec.HostNetwork = hostNetwork
}

func RedefineWithPVC(deployment *appsv1.Deployment, name string, claim string) {
	deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: claim,
				},
			},
		},
	}
}

func RedefineWithHostPath(deployment *appsv1.Deployment, name string, path string) {
	deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: path,
				},
			},
		},
	}
}

func RedefineWithCPUResources(deployment *appsv1.Deployment, limit string, req string) {
	for i := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[i].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(limit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(req),
			},
		}
	}
}

func RedefineWithAllRequestsAndLimits(deployment *appsv1.Deployment, memoryLimit string, cpuLimit string,
	memoryRequest string, cpuRequest string) {
	for i := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[i].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(memoryLimit),
				corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(memoryRequest),
				corev1.ResourceCPU:    resource.MustParse(cpuRequest),
			},
		}
	}
}

func RedefineWithMemoryRequestsAndLimitsAndCPURequest(deployment *appsv1.Deployment, memoryLimit string,
	memoryRequest string, cpuRequest string) {
	for i := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[i].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(memoryLimit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(memoryRequest),
				corev1.ResourceCPU:    resource.MustParse(cpuRequest),
			},
		}
	}
}

func RedefineWithMemoryRequestAndCPURequestsAndLimits(deployment *appsv1.Deployment, cpuLimit string,
	memoryRequest string, cpuRequest string) {
	for i := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[i].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(cpuLimit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(memoryRequest),
				corev1.ResourceCPU:    resource.MustParse(cpuRequest),
			},
		}
	}
}

func RedefineWithResourceRequests(deployment *appsv1.Deployment, memory string, cpu string) {
	for i := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[i].Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(memory),
				corev1.ResourceCPU:    resource.MustParse(cpu),
			},
		}
	}
}

func RedefineWithRunTimeClass(deployment *appsv1.Deployment, rtcName string) {
	deployment.Spec.Template.Spec.RuntimeClassName = pointer.String(rtcName)
}

func RedefineWithShareProcessNamespace(deployment *appsv1.Deployment, shareProcessNamespace bool) {
	deployment.Spec.Template.Spec.ShareProcessNamespace = &shareProcessNamespace
}

func RedefineWithSysPtrace(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"SYS_PTRACE"}},
		}
	}
}

func RedefineWith2MiHugepages(deployment *appsv1.Deployment, hugepages int) {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Mi"))

	for i := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[i].Resources.Requests[corev1.ResourceHugePagesPrefix+"2Mi"] = hugepagesVal
		deployment.Spec.Template.Spec.Containers[i].Resources.Limits[corev1.ResourceHugePagesPrefix+"2Mi"] = hugepagesVal
	}
}

func RedefineWith1GiHugepages(deployment *appsv1.Deployment, hugepages int) {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Gi"))

	for i := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[i].Resources.Requests[corev1.ResourceHugePagesPrefix+"1Gi"] = hugepagesVal
		deployment.Spec.Template.Spec.Containers[i].Resources.Limits[corev1.ResourceHugePagesPrefix+"1Gi"] = hugepagesVal
	}
}

func RedefineWithNoExecuteToleration(deployment *appsv1.Deployment) {
	tol := corev1.Toleration{
		Effect:            "NoExecute",
		Key:               "node.kubernetes.io/not-ready",
		Operator:          "Exists",
		TolerationSeconds: pointer.Int64(365),
	}
	deployment.Spec.Template.Spec.Tolerations = append(deployment.Spec.Template.Spec.Tolerations, tol)
}

func RedefineWithPreferNoScheduleToleration(deployment *appsv1.Deployment) {
	tol := corev1.Toleration{
		Effect:   "PreferNoSchedule",
		Key:      "node.kubernetes.io/memory-pressure",
		Operator: "Equal",
		Value:    "value1",
	}
	deployment.Spec.Template.Spec.Tolerations = append(deployment.Spec.Template.Spec.Tolerations, tol)
}

func RedefineWithNoScheduleToleration(deployment *appsv1.Deployment) {
	tol := corev1.Toleration{
		Effect:   "NoSchedule",
		Key:      "node.kubernetes.io/memory-pressure",
		Operator: "Equal",
		Value:    "value2",
	}
	deployment.Spec.Template.Spec.Tolerations = append(deployment.Spec.Template.Spec.Tolerations, tol)
}

func RedefineWithServiceAccount(deployment *appsv1.Deployment, serviceAccountName string) {
	deployment.Spec.Template.Spec.ServiceAccountName = serviceAccountName
}

// RedefineWithPostStart adds postStart to deployment manifest.
func RedefineWithPostStart(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].Lifecycle = &corev1.Lifecycle{
			PostStart: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

func RedefineWithPodSecurityContextRunAsUser(deployment *appsv1.Deployment, uid int64) {
	deployment.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{
		RunAsUser: pointer.Int64(uid),
	}
}

// RedefineWithContainersSecurityContextAll redefines deployment with extended permissions.
func RedefineWithContainersSecurityContextAll(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"ALL"}},
		}
	}
}

func RedefineWithContainersSecurityContextIpcLock(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"IPC_LOCK"}},
		}
	}
}

func RedefineWithContainersSecurityContextNetAdmin(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"NET_ADMIN"}},
		}
	}
}

func RedefineWithContainersSecurityContextNetRaw(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"NET_RAW"}},
		}
	}
}

func RedefineWithContainersSecurityContextSysAdmin(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"SYS_ADMIN"}},
		}
	}
}

func RedefineWithContainersSecurityContextBpf(deployment *appsv1.Deployment) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"BPF"}},
		}
	}
}

func RedefineWithContainersSecurityContextAllowPrivilegeEscalation(deployment *appsv1.Deployment,
	allowPrivilegeEscalation bool) {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			RunAsUser:                pointer.Int64(0),
			AllowPrivilegeEscalation: &allowPrivilegeEscalation,
		}
	}
}

func RedefineContainerCommand(deployment *appsv1.Deployment, index int, command []string) error {
	if len(deployment.Spec.Template.Spec.Containers) > index {
		deployment.Spec.Template.Spec.Containers[index].Command = command

		return nil
	}

	return fmt.Errorf("deployment %s does not have container index %d", deployment.Name, index)
}

func RedefineContainerEnvVarList(deployment *appsv1.Deployment, index int, envVars []corev1.EnvVar) error {
	if len(deployment.Spec.Template.Spec.Containers) > index {
		deployment.Spec.Template.Spec.Containers[index].Env = envVars

		return nil
	}

	return fmt.Errorf("deployment %s does not have container index %d", deployment.Name, index)
}

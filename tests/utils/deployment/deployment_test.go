package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func TestDefineDeployment(t *testing.T) {
	testCases := []struct {
		name      string
		namespace string
		image     string
		labels    map[string]string
	}{
		{
			name:      "test-deployment",
			namespace: "test-namespace",
			image:     "test-image",
			labels: map[string]string{
				"app": "test",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			deployment := DefineDeployment(testCase.name, testCase.namespace, testCase.image, testCase.labels)
			assert.NotNil(t, deployment)

			// Assert that parameters passed in are in the deployment.
			assert.Equal(t, testCase.name, deployment.Name)
			assert.Equal(t, testCase.namespace, deployment.Namespace)
			assert.Equal(t, testCase.image, deployment.Spec.Template.Spec.Containers[0].Image)
			assert.Equal(t, testCase.labels, deployment.Spec.Template.Labels)

			// Assert hardcoded default values are in the deployment.
			assert.Equal(t, int32(1), *deployment.Spec.Replicas)
			assert.Equal(t, int32(5), deployment.Spec.MinReadySeconds)
		})
	}
}

func TestRedefineAllContainersWithPreStopSpec(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineAllContainersWithPreStopSpec(deployment, []string{"sleep 10"})
	assert.NotNil(t, deployment)

	// Assert preStop hook is in the deployment.
	assert.Equal(t, "sleep 10", deployment.Spec.Template.Spec.Containers[0].Lifecycle.PreStop.Exec.Command[0])
}

func TestRedefineWithLabels(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineWithLabels(deployment, map[string]string{"test": "test"})
	assert.NotNil(t, deployment)

	// Assert labels are in the deployment.
	assert.Equal(t, "test", deployment.Spec.Template.Labels["test"])
}

func TestRedefineWithMultus(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineWithMultus(deployment, []string{"test-network"})
	assert.NotNil(t, deployment)

	// Assert multus is in the deployment.
	assert.Equal(t, "[{\"name\":\"test-network\"}]", deployment.Spec.Template.Annotations["k8s.v1.cni.cncf.io/networks"])
}

func TestRedefineWithReplicaNumber(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineWithReplicaNumber(deployment, 2)
	assert.NotNil(t, deployment)

	// Assert replica number is in the deployment.
	assert.Equal(t, int32(2), *deployment.Spec.Replicas)
}

func TestRedefineFirstContainerWithPreStopSpec(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	assert.Nil(t, RedefineFirstContainerWithPreStopSpec(deployment, []string{"sleep 10"}))
	assert.NotNil(t, deployment)

	// Assert preStop hook is in the deployment.
	assert.Equal(t, "sleep 10", deployment.Spec.Template.Spec.Containers[0].Lifecycle.PreStop.Exec.Command[0])

	// Delete the containers and assert an error
	deployment.Spec.Template.Spec.Containers = []corev1.Container{}
	err := RedefineFirstContainerWithPreStopSpec(deployment, []string{"sleep 10"})
	assert.NotNil(t, err)
}

func TestRedefineWithTerminationGracePeriod(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	int64Val := int64(10)

	RedefineWithTerminationGracePeriod(deployment, &int64Val)
	assert.NotNil(t, deployment)

	// Assert terminationGracePeriod is in the deployment.
	assert.Equal(t, int64(10), *deployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
}

func TestRedefineWithPodAntiAffinity(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineWithPodAntiAffinity(deployment, map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	// Assert podAntiAffinity is in the deployment.
	//nolint:lll
	assert.Equal(t, deployment.Spec.Template.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchLabels["app"], deployment.Spec.Template.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchLabels["app"])
}

func TestRedefineWithImagePullPolicy(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineWithImagePullPolicy(deployment, corev1.PullAlways)
	assert.NotNil(t, deployment)

	// Assert imagePullPolicy is in the deployment.
	assert.Equal(t, corev1.PullAlways, deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy)
}

func TestRedefineWithNodeSelector(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineWithNodeSelector(deployment, map[string]string{"test": "test"})
	assert.NotNil(t, deployment)

	// Assert nodeSelector is in the deployment.
	assert.Equal(t, "test", deployment.Spec.Template.Spec.NodeSelector["test"])
}

func TestRedefineWithNodeAffinity(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineWithNodeAffinity(deployment, "test-key")
	assert.NotNil(t, deployment)

	// Assert nodeAffinity is in the deployment.
	//nolint:lll
	assert.Equal(t, "test-key", deployment.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Key)
}

func TestRedefineWithReadinessProbe(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have a readiness probe.
	assert.Nil(t, deployment.Spec.Template.Spec.Containers[0].ReadinessProbe)

	RedefineWithReadinessProbe(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has a readiness probe.
	assert.NotNil(t, deployment.Spec.Template.Spec.Containers[0].ReadinessProbe)
	assert.Equal(t, "ls", deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Exec.Command[0])
}

func TestRedefineWithLivenessProbe(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have a liveness probe.
	assert.Nil(t, deployment.Spec.Template.Spec.Containers[0].LivenessProbe)

	RedefineWithLivenessProbe(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has a liveness probe.
	assert.NotNil(t, deployment.Spec.Template.Spec.Containers[0].LivenessProbe)
	assert.Equal(t, "ls", deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Exec.Command[0])
}

func TestRedefineWithStartUpProbe(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have a startup probe.
	assert.Nil(t, deployment.Spec.Template.Spec.Containers[0].StartupProbe)

	RedefineWithStartUpProbe(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has a startup probe.
	assert.NotNil(t, deployment.Spec.Template.Spec.Containers[0].StartupProbe)
	assert.Equal(t, "ls", deployment.Spec.Template.Spec.Containers[0].StartupProbe.Exec.Command[0])
}

func TestRedefineWithContainerSpecs(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	// Clear the containers
	deployment.Spec.Template.Spec.Containers = []corev1.Container{}

	RedefineWithContainerSpecs(deployment, []corev1.Container{
		{
			Name: "test-container",
		},
	})
	assert.NotNil(t, deployment)

	// Assert that the deployment has container specs.
	assert.Equal(t, "test-container", deployment.Spec.Template.Spec.Containers[0].Name)
}

func TestRedefineWithPrivilegedContainer(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	assert.NotNil(t, deployment)

	RedefineWithPrivilegedContainer(deployment)
	assert.NotNil(t, deployment)

	securityContext := corev1.SecurityContext{
		Privileged: ptr.To[bool](true),
		RunAsUser:  ptr.To[int64](0),
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				"ALL",
			},
		},
	}

	// Assert that the deployment has securityContext set.
	assert.Equal(t, securityContext.Privileged, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Privileged)
	assert.Equal(t, securityContext.RunAsUser, deployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser)
	assert.Equal(t, securityContext.Capabilities, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities)
}

func TestRedefineWithHostPid(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have hostPID set.
	assert.False(t, deployment.Spec.Template.Spec.HostPID)

	RedefineWithHostPid(deployment, true)
	assert.NotNil(t, deployment)

	// Assert that the deployment has hostPID set.
	assert.True(t, deployment.Spec.Template.Spec.HostPID)
}

func TestRedefineWithHostIpc(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have hostIPC set.
	assert.False(t, deployment.Spec.Template.Spec.HostIPC)

	RedefineWithHostIpc(deployment, true)
	assert.NotNil(t, deployment)

	// Assert that the deployment has hostIPC set.
	assert.True(t, deployment.Spec.Template.Spec.HostIPC)
}

func TestRedefineWithAutomountServiceAccountToken(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have automountServiceAccountToken set.
	assert.Nil(t, deployment.Spec.Template.Spec.AutomountServiceAccountToken)

	RedefineWithAutomountServiceAccountToken(deployment, true)
	assert.NotNil(t, deployment)

	// Assert that the deployment has automountServiceAccountToken set.
	assert.True(t, *deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}

func TestRedefineWithHostNetwork(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have hostNetwork set.
	assert.False(t, deployment.Spec.Template.Spec.HostNetwork)

	RedefineWithHostNetwork(deployment, true)
	assert.NotNil(t, deployment)

	// Assert that the deployment has hostNetwork set.
	assert.True(t, deployment.Spec.Template.Spec.HostNetwork)
}

func TestRedefineWithPVC(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have volumes set.
	assert.Equal(t, 0, len(deployment.Spec.Template.Spec.Volumes))

	RedefineWithPVC(deployment, "test-volume", "test-pvc")
	assert.NotNil(t, deployment)

	// Assert that the deployment has volumes set.
	assert.Equal(t, 1, len(deployment.Spec.Template.Spec.Volumes))
	assert.Equal(t, "test-volume", deployment.Spec.Template.Spec.Volumes[0].Name)
	assert.Equal(t, "test-pvc", deployment.Spec.Template.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim.ClaimName)
}

func TestRedefineWithHostPath(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have volumes set.
	assert.Equal(t, 0, len(deployment.Spec.Template.Spec.Volumes))

	RedefineWithHostPath(deployment, "test-volume", "test-path")
	assert.NotNil(t, deployment)

	// Assert that the deployment has volumes set.
	assert.Equal(t, 1, len(deployment.Spec.Template.Spec.Volumes))
	assert.Equal(t, "test-volume", deployment.Spec.Template.Spec.Volumes[0].Name)
	assert.Equal(t, "test-path", deployment.Spec.Template.Spec.Volumes[0].VolumeSource.HostPath.Path)
}

func TestRedefineWithCPUResources(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	RedefineWithCPUResources(deployment, "150m", "100m")
	assert.NotNil(t, deployment)

	// Assert that the deployment has resources set.
	assert.Equal(t, "100m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	assert.Equal(t, "150m", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
}

func TestRedefineWithAllRequestsAndLimits(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	RedefineWithAllRequestsAndLimits(deployment, "150m", "100m", "150Mi", "100Mi")
	assert.NotNil(t, deployment)

	// Assert that the deployment has resources set.
	assert.Equal(t, "100Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	assert.Equal(t, "100m", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	assert.Equal(t, "150Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
	assert.Equal(t, "150m", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
}

func TestRedefineWithMemoryRequestsAndLimitsAndCPURequest(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	RedefineWithMemoryRequestsAndLimitsAndCPURequest(deployment, "150m", "100Mi", "150m")
	assert.NotNil(t, deployment)

	// Assert that the deployment has resources set.
	assert.Equal(t, "100Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
	assert.Equal(t, "150m", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
	assert.Equal(t, "150m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	assert.Equal(t, "0", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()) // purposefully set to 0
}

func TestRedefineWithMemoryRequestAndCPURequestsAndLimits(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	RedefineWithMemoryRequestAndCPURequestsAndLimits(deployment, "150Mi", "150m", "100m")
	assert.NotNil(t, deployment)

	// Assert that the deployment has resources set.
	assert.Equal(t, "150m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
	assert.Equal(t, "100m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	assert.Equal(t, "150Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	assert.Equal(t, "0", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()) // purposefully set to 0
}

func TestRedefineWithResourceRequests(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	RedefineWithResourceRequests(deployment, "150m", "100Mi")
	assert.NotNil(t, deployment)

	// Assert that the deployment has resources set.
	assert.Equal(t, "150m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
	assert.Equal(t, "100Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
}

func TestRedefineWithRunTimeClass(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	RedefineWithRunTimeClass(deployment, "test-runtime-class")
	assert.NotNil(t, deployment)

	// Assert that the deployment has runtime class set.
	assert.Equal(t, "test-runtime-class", *deployment.Spec.Template.Spec.RuntimeClassName)
}

func TestRedefineWithShareProcessNamespace(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Assert that the deployment does not have shareProcessNamespace set.
	assert.Nil(t, deployment.Spec.Template.Spec.ShareProcessNamespace)

	RedefineWithShareProcessNamespace(deployment, true)
	assert.NotNil(t, deployment)

	// Assert that the deployment has shareProcessNamespace set.
	assert.True(t, *deployment.Spec.Template.Spec.ShareProcessNamespace)
}

func TestRedefineWithSysPtrace(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	securityContext := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				"SYS_PTRACE",
			},
		},
	}

	// Assert that the deployment does not have securityContext set.
	assert.Nil(t, deployment.Spec.Template.Spec.Containers[0].SecurityContext)

	RedefineWithSysPtrace(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has securityContext set.
	assert.Equal(t, securityContext.Capabilities, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities)
}

func TestRedefineWithNoExecuteToleration(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Clear the tolerations
	deployment.Spec.Template.Spec.Tolerations = []corev1.Toleration{}

	RedefineWithNoExecuteToleration(deployment)
	assert.NotNil(t, deployment)

	tol := corev1.Toleration{
		Effect:            "NoExecute",
		Key:               "node.kubernetes.io/not-ready",
		Operator:          "Exists",
		TolerationSeconds: ptr.To[int64](365),
	}

	// Assert that the deployment has tolerations set.
	assert.Equal(t, tol.Key, deployment.Spec.Template.Spec.Tolerations[0].Key)
	assert.Equal(t, tol.Effect, deployment.Spec.Template.Spec.Tolerations[0].Effect)
	assert.Equal(t, tol.Operator, deployment.Spec.Template.Spec.Tolerations[0].Operator)
	assert.Equal(t, tol.TolerationSeconds, deployment.Spec.Template.Spec.Tolerations[0].TolerationSeconds)
}

func TestRedefineWithPreferNoScheduleToleration(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Clear the tolerations
	deployment.Spec.Template.Spec.Tolerations = []corev1.Toleration{}

	RedefineWithPreferNoScheduleToleration(deployment)
	assert.NotNil(t, deployment)

	tol := corev1.Toleration{
		Effect:   "PreferNoSchedule",
		Key:      "node.kubernetes.io/memory-pressure",
		Operator: "Equal",
		Value:    "value1",
	}

	// Assert that the deployment has tolerations set.
	assert.Equal(t, tol.Key, deployment.Spec.Template.Spec.Tolerations[0].Key)
	assert.Equal(t, tol.Effect, deployment.Spec.Template.Spec.Tolerations[0].Effect)
	assert.Equal(t, tol.Operator, deployment.Spec.Template.Spec.Tolerations[0].Operator)
	assert.Equal(t, tol.Value, deployment.Spec.Template.Spec.Tolerations[0].Value)
}

func TestRedefineWithNoScheduleToleration(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	// Clear the tolerations
	deployment.Spec.Template.Spec.Tolerations = []corev1.Toleration{}

	RedefineWithNoScheduleToleration(deployment)
	assert.NotNil(t, deployment)

	tol := corev1.Toleration{
		Effect:   "NoSchedule",
		Key:      "node.kubernetes.io/memory-pressure",
		Operator: "Equal",
		Value:    "value2",
	}

	// Assert that the deployment has tolerations set.
	assert.Equal(t, tol.Key, deployment.Spec.Template.Spec.Tolerations[0].Key)
	assert.Equal(t, tol.Effect, deployment.Spec.Template.Spec.Tolerations[0].Effect)
	assert.Equal(t, tol.Operator, deployment.Spec.Template.Spec.Tolerations[0].Operator)
	assert.Equal(t, tol.Value, deployment.Spec.Template.Spec.Tolerations[0].Value)
}

func TestRedefineWithServiceAccount(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	RedefineWithServiceAccount(deployment, "test-service-account")
	assert.NotNil(t, deployment)

	// Assert that the deployment has serviceAccount set.
	assert.Equal(t, "test-service-account", deployment.Spec.Template.Spec.ServiceAccountName)
}

func TestAppendServiceAccount(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	AppendServiceAccount(deployment, "test-service-account")
	assert.NotNil(t, deployment)

	// Assert that the deployment has serviceAccount set.
	assert.Equal(t, "test-service-account", deployment.Spec.Template.Spec.ServiceAccountName)
}

func TestRedefineWithPostStart(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	RedefineWithPostStart(deployment)
	assert.NotNil(t, deployment)

	// Assert postStart hook is in the deployment.
	assert.Equal(t, "ls", deployment.Spec.Template.Spec.Containers[0].Lifecycle.PostStart.Exec.Command[0])
}

func TestRedefineWithPodSecurityContextRunAsUser(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	RedefineWithPodSecurityContextRunAsUser(deployment, 1000)
	assert.NotNil(t, deployment)

	// Assert that the deployment has podSecurityContext set.
	assert.Equal(t, int64(1000), *deployment.Spec.Template.Spec.SecurityContext.RunAsUser)
}

func TestRedefineWithContainersSecurityContextAll(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

	securityContext := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				"ALL",
			},
		},
	}

	RedefineWithContainersSecurityContextAll(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has securityContext set.
	assert.Equal(t, securityContext.Capabilities, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities)
}

func TestRedefineWithContainersSecurityContextIpcLock(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	securityContext := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				"IPC_LOCK",
			},
		},
	}

	RedefineWithContainersSecurityContextIpcLock(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has securityContext set.
	assert.Equal(t, securityContext.Capabilities, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities)
}

func TestRedefineWithContainersSecurityContextNetAdmin(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	securityContext := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				"NET_ADMIN",
			},
		},
	}

	RedefineWithContainersSecurityContextNetAdmin(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has securityContext set.
	assert.Equal(t, securityContext.Capabilities, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities)
}

func TestRedefineWithContainersSecurityContextNetRaw(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	securityContext := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				"NET_RAW",
			},
		},
	}

	RedefineWithContainersSecurityContextNetRaw(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has securityContext set.
	assert.Equal(t, securityContext.Capabilities, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities)
}

func TestRedefineWithContainersSecurityContextSysAdmin(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	securityContext := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				"SYS_ADMIN",
			},
		},
	}

	RedefineWithContainersSecurityContextSysAdmin(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has securityContext set.
	assert.Equal(t, securityContext.Capabilities, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities)
}

func TestRedefineWithContainersSecurityContextBpf(t *testing.T) {
	deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
	securityContext := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				"BPF",
			},
		},
	}

	RedefineWithContainersSecurityContextBpf(deployment)
	assert.NotNil(t, deployment)

	// Assert that the deployment has securityContext set.
	assert.Equal(t, securityContext.Capabilities, deployment.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities)
}

func TestRedefineWithContainersSecurityContextAllowPrivilegeEscalation(t *testing.T) {
	testCases := []struct {
		name                     string
		allowPrivilegeEscalation bool
	}{
		{
			name:                     "true",
			allowPrivilegeEscalation: true,
		},
		{
			name:                     "false",
			allowPrivilegeEscalation: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})
			securityContext := corev1.SecurityContext{
				AllowPrivilegeEscalation: &testCase.allowPrivilegeEscalation,
			}

			RedefineWithContainersSecurityContextAllowPrivilegeEscalation(deployment, testCase.allowPrivilegeEscalation)
			assert.NotNil(t, deployment)

			// Assert that the deployment has securityContext set.
			assert.Equal(t, securityContext.AllowPrivilegeEscalation,
				deployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation)
		})
	}
}

func TestRedefineContainerCommand(t *testing.T) {
	testCases := []struct {
		name     string
		commands []string
	}{
		{
			name:     "single command",
			commands: []string{"ls"},
		},
		{
			name:     "multiple commands",
			commands: []string{"ls", "-l"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

			assert.Nil(t, RedefineContainerCommand(deployment, 0, testCase.commands))
			assert.NotNil(t, deployment)

			// Assert that the deployment has commands set.
			assert.Equal(t, testCase.commands, deployment.Spec.Template.Spec.Containers[0].Command)
		})
	}
}

func TestRedefineContainerEnvVarList(t *testing.T) {
	testCases := []struct {
		name string
		envs []corev1.EnvVar
	}{
		{
			name: "single env",
			envs: []corev1.EnvVar{
				{
					Name:  "test",
					Value: "test",
				},
			},
		},
		{
			name: "multiple envs",
			envs: []corev1.EnvVar{
				{
					Name:  "test",
					Value: "test",
				},
				{
					Name:  "test2",
					Value: "test2",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			deployment := DefineDeployment("test-deployment", "test-namespace", "test-image", map[string]string{"app": "test"})

			assert.Nil(t, RedefineContainerEnvVarList(deployment, 0, testCase.envs))
			assert.NotNil(t, deployment)

			// Assert that the deployment has envs set.
			assert.Equal(t, testCase.envs, deployment.Spec.Template.Spec.Containers[0].Env)
		})
	}
}

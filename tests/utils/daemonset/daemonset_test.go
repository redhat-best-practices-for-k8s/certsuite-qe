package daemonset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestDefineDaemonSet(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	assert.NotNil(t, ds)
	assert.Equal(t, "nginx", ds.Name)
	assert.Equal(t, "default", ds.Namespace)
	assert.Equal(t, int32(3), ds.Spec.MinReadySeconds)
	assert.Equal(t, "testpod-", ds.Spec.Template.Name)
}

func TestDefineDaemonSetWithContainerSpecs(t *testing.T) {
	testDS := DefineDaemonSetWithContainerSpecs("nginx", "default", map[string]string{"app": "nginx"}, []corev1.Container{
		{
			Name:    "test",
			Image:   "nginx",
			Command: []string{"/bin/bash", "-c", "sleep INF"}}})
	assert.NotNil(t, testDS)
	assert.Equal(t, "nginx", testDS.Name)
	assert.Equal(t, "default", testDS.Namespace)
	assert.Equal(t, int32(3), testDS.Spec.MinReadySeconds)
}

func TestRedefineDaemonSetWithNodeSelector(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineDaemonSetWithNodeSelector(ds, map[string]string{"node-role.kubernetes.io/master": ""})
	assert.Equal(t, "", ds.Spec.Template.Spec.NodeSelector["node-role.kubernetes.io/master"])
}

func TestRedefineWithLabel(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithLabel(ds, map[string]string{"app": "nginx"})
	assert.Equal(t, "nginx", ds.Spec.Template.Labels["app"])
}

func TestRedefineWithPrivilegeAndHostNetwork(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithPrivilegeAndHostNetwork(ds)
	assert.Equal(t, true, *ds.Spec.Template.Spec.Containers[0].SecurityContext.Privileged)
	assert.Equal(t, true, ds.Spec.Template.Spec.HostNetwork)
	assert.NotNil(t, ds.Spec.Template.Spec.Containers[0].SecurityContext)
	assert.Equal(t, int64(0), *ds.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser)
}

func TestRedefineWithMultus(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithMultus(ds, "multus")
	assert.Equal(t, "[ { \"name\": \"multus\" } ]", ds.Spec.Template.Annotations["k8s.v1.cni.cncf.io/networks"])
}

func TestRedefineWithImagePullPolicy(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithImagePullPolicy(ds, "Always")
	assert.Equal(t, corev1.PullAlways, ds.Spec.Template.Spec.Containers[0].ImagePullPolicy)
}

func TestRedefineWithContainerSpecs(t *testing.T) {
	testDS := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithContainerSpecs(testDS, []corev1.Container{
		{
			Name:    "test",
			Image:   "nginx",
			Command: []string{"/bin/bash", "-c", "sleep INF"}}})
	assert.Equal(t, "test", testDS.Spec.Template.Spec.Containers[0].Name)
	assert.Equal(t, "nginx", testDS.Spec.Template.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/bin/bash", "-c", "sleep INF"}, testDS.Spec.Template.Spec.Containers[0].Command)
}

func TestRedefineWithPrivilegedContainer(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithPrivilegedContainer(ds)
	assert.Equal(t, true, *ds.Spec.Template.Spec.Containers[0].SecurityContext.Privileged)
	assert.Equal(t, int64(0), *ds.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser)
}

func TestRedefineWithVolumeMount(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithVolumeMount(ds)
	assert.Equal(t, "host", ds.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name)
	assert.Equal(t, "/host", ds.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath)
}

func TestRedefineWithCPUResources(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithCPUResources(ds, "100m", "101m")
	assert.Equal(t, "101m", ds.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	assert.Equal(t, "100m", ds.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
}

func TestRedefineWithRunTimeClass(t *testing.T) {
	ds := DefineDaemonSet("default", "nginx", map[string]string{"app": "nginx"}, "nginx")
	RedefineWithRunTimeClass(ds, "test")
	assert.Equal(t, "test", *ds.Spec.Template.Spec.RuntimeClassName)
}

package helper

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/utils/pointer"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/performance/parameters"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineExclusivePod(podName string, namespace string, image string, label map[string]string) *corev1.Pod {
	cpuLimit := "1"
	memoryLimit := "512Mi"
	containerCommand := []string{"/bin/bash", "-c", "sleep INF"}

	containerResource := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    label},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: pointer.Int64(0),
			ServiceAccountName:            tsparams.PriviledgedRoleName,
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    pointer.Int64(1000),
				RunAsGroup:   pointer.Int64(1000),
				RunAsNonRoot: pointer.Bool(true)},
			Containers: []corev1.Container{
				{
					Name:      "shared",
					Image:     image,
					Command:   containerCommand,
					Resources: containerResource},
				{
					Name:      "exclusive",
					Image:     image,
					Command:   containerCommand,
					Resources: containerResource},
			},
		},
	}
}

func DefineRtPod(podName string, namespace string, image string, label map[string]string) *corev1.Pod {
	cpuLimit := "1"
	memoryLimit := "512Mi"
	containerCommand := []string{"/bin/bash", "-c", "sleep INF"}

	containerResource := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
	}

	containerSecurityContext := &corev1.SecurityContext{
		Privileged: pointer.Bool(true),
		RunAsUser:  pointer.Int64(0),
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    label},
		Spec: corev1.PodSpec{
			ServiceAccountName:            tsparams.PriviledgedRoleName,
			TerminationGracePeriodSeconds: pointer.Int64(0),
			Containers: []corev1.Container{
				{
					Name:            "rt-app",
					Image:           image,
					Command:         containerCommand,
					Resources:       containerResource,
					SecurityContext: containerSecurityContext},
			},
		},
	}
}

func RedefinePodWithSharedContainer(pod *corev1.Pod, containerIndex int) {
	totalContainers := len(pod.Spec.Containers)
	limit := "1"
	req := "250m"

	if containerIndex >= 0 && containerIndex < totalContainers {
		pod.Spec.Containers[containerIndex].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(limit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(req),
			},
		}
	}
}

func ChangeSchedulingPolicy(pod *corev1.Pod, command string) (outStr, errStr string, err error) {
	outStr, errStr, err = ExecCommandContainer(pod, command)
	if err != nil {
		return "", "", fmt.Errorf("cannot execute command: \" %s \"  on %s err:%w", command, pod.Spec.Containers[0].Name, err)
	}

	return outStr, errStr, err
}

// ExecCommand runs command in the pod and returns buffer output.
func ExecCommandContainer(
	pod *corev1.Pod, command string) (stdout, stderr string, err error) {
	commandStr := []string{"sh", "-c", command}

	var buffOut bytes.Buffer

	var buffErr bytes.Buffer

	podName := pod.Name
	podNamespace := pod.Namespace
	container := pod.Spec.Containers[0].Name

	logrus.Trace(fmt.Sprintf("execute command on ns=%s, pod=%s container=%s, cmd: %s",
		podNamespace, podName, container, strings.Join(commandStr, " ")))

	req := globalhelper.APIClient.CoreV1Interface.RESTClient().
		Post().
		Namespace(podNamespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   commandStr,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(globalhelper.APIClient.Config, "POST", req.URL())
	if err != nil {
		logrus.Error(err)

		return stdout, stderr, err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: &buffOut,
		Stderr: &buffErr,
	})

	stdout, stderr = buffOut.String(), buffErr.String()

	if err != nil {
		logrus.Error(err)
		logrus.Error(req.URL())
		logrus.Error("command: ", command)
		logrus.Error("stderr: ", stderr)
		logrus.Error("stdout: ", stdout)

		return stdout, stderr, err
	}

	return stdout, stderr, err
}

func DeleteRunTimeClass(rtcName string) error {
	err := globalhelper.APIClient.RuntimeClasses().Delete(context.Background(), rtcName,
		metav1.DeleteOptions{GracePeriodSeconds: pointer.Int64(0)})
	if err != nil {
		return fmt.Errorf("failed to delete RunTimeClasses %w", err)
	}

	return nil
}

func ConfigurePrivilegedServiceAccount(namespace string) error {
	aRole, aRoleBinding, aServiceAccount := getPrivilegedServiceAccountObjects(namespace)
	// create role
	_, err := globalhelper.APIClient.RbacV1Interface.Roles(namespace).Create(context.TODO(), &aRole, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating role, err=%w", err)
	}

	// create rolebinding
	_, err = globalhelper.APIClient.RbacV1Interface.RoleBindings(namespace).Create(context.TODO(), &aRoleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating role bindings, err=%w", err)
	}

	// create service account
	_, err = globalhelper.APIClient.CoreV1Interface.ServiceAccounts(namespace).Create(context.TODO(),
		&aServiceAccount, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating service account, err=%w", err)
	}

	return nil
}

func getPrivilegedServiceAccountObjects(namespace string) (aRole rbacv1.Role,
	aRoleBinding rbacv1.RoleBinding, aServiceAccount corev1.ServiceAccount) {
	aRole = rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsparams.PriviledgedRoleName,
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{{
			APIGroups: []string{"*"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		},
		},
	}

	aRoleBinding = rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsparams.PriviledgedRoleName,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      tsparams.PriviledgedRoleName,
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     tsparams.PriviledgedRoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	aServiceAccount = corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsparams.PriviledgedRoleName,
			Namespace: namespace,
		},
	}

	return aRole, aRoleBinding, aServiceAccount
}
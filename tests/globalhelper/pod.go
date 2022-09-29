package globalhelper

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/gomega"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecCommand runs command in the pod and returns buffer output.
func ExecCommand(pod corev1.Pod, command []string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	req := APIClient.CoreV1Interface.RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: pod.Spec.Containers[0].Name,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(APIClient.Config, "POST", req.URL())
	if err != nil {
		return buf, err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: &buf,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return buf, err
	}

	return buf, nil
}

// GetListOfPodsInNamespace returns list of pods for given namespace.
func GetListOfPodsInNamespace(namespace string) (*corev1.PodList, error) {
	runningPods, err := APIClient.Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return runningPods, nil
}

// CreateAndWaitUntilPodIsReady creates a pod and waits until all it's containers are ready.
func CreateAndWaitUntilPodIsReady(pod *corev1.Pod, timeout time.Duration) error {
	createdPod, err := APIClient.Pods(pod.Namespace).Create(
		context.Background(),
		pod,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isPodReady(createdPod.Namespace, createdPod.Name)
		if err != nil {

			glog.V(5).Info(fmt.Sprintf(
				"Pod %s is not ready, retry in %d seconds", createdPod.Name, retryInterval))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "Pod is not ready")

	return nil
}

// isPodReady checks if a pod is ready.
func isPodReady(namespace string, podName string) (bool, error) {
	podObject, err := APIClient.Pods(namespace).Get(
		context.Background(),
		podName,
		metav1.GetOptions{},
	)

	if err != nil {
		return false, err
	}

	numContainers := len(podObject.Spec.Containers)

	if len(podObject.Status.ContainerStatuses) != numContainers {
		return false, nil
	}

	for index := range podObject.Spec.Containers {
		if !podObject.Status.ContainerStatuses[index].Ready {
			return false, nil
		}
	}

	return true, nil
}

// AppendContainersToPod appends containers to a pod.
func AppendContainersToPod(pod *corev1.Pod, containersNum int, image string) {
	containerList := &pod.Spec.Containers

	for i := 0; i < containersNum; i++ {
		*containerList = append(
			*containerList, corev1.Container{
				Name:    fmt.Sprintf("container%d", i+1),
				Image:   image,
				Command: []string{"/bin/bash", "-c", "sleep INF"},
			})
	}
}

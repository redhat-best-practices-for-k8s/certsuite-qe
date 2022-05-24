package globalhelper

import (
	"bytes"
	"context"
	"os"
	"time"

	. "github.com/onsi/gomega"
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

// CreateAndWaitUntilPodIsReady creates a pod and waits until all its containers are ready.
func CreateAndWaitUntilPodIsReady(pod *corev1.Pod, timeout time.Duration) error {
	_, err := APIClient.Pods(pod.Namespace).Create(
		context.Background(),
		pod,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	numContainers := len(pod.Spec.Containers)

	Eventually(func() bool {
		runningPod, err := APIClient.Pods(pod.Namespace).Get(
			context.Background(),
			pod.Name,
			metav1.GetOptions{})
		if err != nil {
			return false
		}

		// We need to wait until all the containers have an status entry.
		if len(runningPod.Status.ContainerStatuses) != numContainers {
			return false
		}

		for index := range runningPod.Spec.Containers {
			if !runningPod.Status.ContainerStatuses[index].Ready {
				return false
			}
		}

		return true
	}, timeout, 5*time.Second).Should(Equal(true), "Pod is not ready")

	return nil
}

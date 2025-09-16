package globalhelper

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/remotecommand"
	klog "k8s.io/klog/v2"
)

// ExecCommand runs command in the pod and returns buffer output.
func ExecCommand(pod corev1.Pod, command []string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	req := GetAPIClient().CoreV1Interface.RESTClient().
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

	exec, err := remotecommand.NewSPDYExecutor(GetAPIClient().Config, "POST", req.URL())
	if err != nil {
		return buf, err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
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
	runningPods, err := GetAPIClient().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return runningPods, nil
}

// CreateAndWaitUntilPodIsReady creates a pod and waits until all it's containers are ready.
func CreateAndWaitUntilPodIsReady(pod *corev1.Pod, timeout time.Duration) error {
	createdPod, err := GetAPIClient().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Info(fmt.Sprintf("pod %s already exists", pod.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create pod %q (ns %s): %w", pod.Name, pod.Namespace, err)
	}

	// Pod isn't running, we need to wait for it to be scheduled
	// Loop through pod conditions to check if pod is schedulable
	// If it is not schedulable, we return an error
	podUnschedulable := false

	Eventually(func() bool {
		runningPod, err := GetAPIClient().Pods(pod.Namespace).Get(
			context.TODO(), pod.Name, metav1.GetOptions{})

		if err != nil {
			klog.V(5).Info(fmt.Sprintf(
				"Pod %s is not running, retry in 1 second", createdPod.Name))

			return false
		}

		// If it is running, we can break the loop
		if runningPod.Status.Phase == corev1.PodRunning {
			return true
		}

		// print the conditions
		fmt.Printf("Pod %s conditions: %v\n", runningPod.Name, runningPod.Status.Conditions)

		for _, condition := range runningPod.Status.Conditions {
			if condition.Type == corev1.PodScheduled && condition.Status == corev1.ConditionFalse && condition.Reason == "Unschedulable" {
				podUnschedulable = true

				break
			}
		}

		return podUnschedulable
	}, timeout, 1*time.Second).Should(Equal(true), "Pod is not running")

	if podUnschedulable {
		return fmt.Errorf("pod %s is not schedulable", createdPod.Name)
	}

	Eventually(func() bool {
		status, err := isPodReady(createdPod.Namespace, createdPod.Name)
		if err != nil {
			klog.V(5).Info(fmt.Sprintf(
				"Pod %s is not ready, retry in 1 second", createdPod.Name))

			return false
		}

		return status
	}, timeout, 1*time.Second).Should(Equal(true), "Pod is not ready")

	return nil
}

// isPodReady checks if a pod is ready.
func isPodReady(namespace string, podName string) (bool, error) {
	podObject, err := GetAPIClient().Pods(namespace).Get(
		context.TODO(),
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

// AppendLabelsToPod appends labels to given pod manifest.
// from go documentation : "if you pass a map to a function that changes the contents of the map,
// the changes will be visible in the caller" - thats the reason for copying the map.
func AppendLabelsToPod(pod *corev1.Pod, labels map[string]string) {
	newMap := make(map[string]string)
	for k, v := range pod.Labels {
		newMap[k] = v
	}

	for k, v := range labels {
		newMap[k] = v
	}

	pod.Labels = newMap
}

func GetRunningPod(namespace, name string) (*corev1.Pod, error) {
	return getRunningPod(GetAPIClient().K8sClient.CoreV1(), namespace, name)
}

func getRunningPod(client typedcorev1.CoreV1Interface, namespace, name string) (*corev1.Pod, error) {
	pod, err := client.Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %q (ns %s): %w", name, namespace, err)
	}

	return pod, nil
}

func GetControllerPodFromOperator(namespace, operatorName string) (*corev1.Pod, error) {
	// Wait for the controller manager pod to come up
	podsFound := false

	var (
		pods *corev1.PodList
		err  error
	)

	// Try 10 times to find the pod
	for i := 0; i < 10; i++ {
		pods, err = GetListOfPodsInNamespace(namespace)
		if err != nil {
			return nil, err
		}

		if len(pods.Items) == 0 {
			fmt.Println("No pods found, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)

			continue
		} else {
			podsFound = true

			break
		}
	}

	if !podsFound {
		return nil, fmt.Errorf("no pods found in namespace %s", namespace)
	}

	for _, pod := range pods.Items {
		fmt.Printf("Checking pod %s\n", pod.Name)

		if strings.Contains(pod.Name, operatorName) {
			return &pod, nil
		}
	}

	return nil, fmt.Errorf("pod for operator %s not found", operatorName)
}

package helper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const WaitingTime = 5 * time.Minute

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

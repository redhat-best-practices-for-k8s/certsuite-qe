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
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const WaitingTime = 5 * time.Minute

// WaitForSpecificNodeCondition waits for a given node to become ready or not.
func WaitForSpecificNodeCondition(clients *client.ClientSet, timeout, interval time.Duration, nodeName string,
	ready bool) error {
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		nodesList, err := clients.Nodes().List(context.Background(), metav1.ListOptions{})
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
func UpdateAndVerifyHugePagesConfig(updatedHugePagesNumber int, filePath string, pod *corev1.Pod) (bool, error) {
	cmd := fmt.Sprintf("echo %d > %s", updatedHugePagesNumber, filePath)
	_, err := globalhelper.ExecCommand(*pod, []string{"/bin/bash", "-c", cmd})
	if err != nil {
		return false, err
	}

	// loop to wait until the file has been actually updated.
	timeout := time.Now().Add(5 * time.Minute)

	for {
		currentHugepagesNumber, err := GetHugePagesConfigNumber(filePath, pod)
		if err != nil {
			return false, err
		}

		if updatedHugePagesNumber == currentHugepagesNumber {
			return true, nil
		}

		if time.Now().After(timeout) {
			return false, nil
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

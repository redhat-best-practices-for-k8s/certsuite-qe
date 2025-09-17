package globalhelper

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/gomega"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiStatefulSet "github.com/openshift-kni/eco-goinfra/pkg/statefulset"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klog "k8s.io/klog/v2"
)

// CreateAndWaitUntilStatefulSetIsReady creates statefulset and waits until all it's replicas are ready.
func CreateAndWaitUntilStatefulSetIsReady(statefulSet *appsv1.StatefulSet, timeout time.Duration) error {
	return createAndWaitUntilStatefulSetIsReady(egiClients.New(""), statefulSet, timeout)
}

func createAndWaitUntilStatefulSetIsReady(client *egiClients.Settings, statefulSet *appsv1.StatefulSet, timeout time.Duration) error {
	statefulSet, err := client.StatefulSets(statefulSet.Namespace).Create(
		context.TODO(), statefulSet, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Info(fmt.Sprintf("statefulSet %s already exists", statefulSet.Name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create statefulSet %q (ns %s): %w", statefulSet.Name, statefulSet.Namespace, err)
	}

	Eventually(func() bool {
		status, err := isStatefulSetReady(client, statefulSet.Namespace, statefulSet.Name)
		if err != nil {
			klog.V(5).Info(fmt.Sprintf(
				"statefulSet %s is not ready, retry in 1 second", statefulSet.Name))

			return false
		}

		return status
	}, timeout, 1*time.Second).Should(Equal(true), "statefulSet is not ready")

	return nil
}

// isStatefulSetReady checks if a statefulset is ready.
func isStatefulSetReady(client *egiClients.Settings, namespace, statefulSetName string) (bool, error) {
	testStatefulSet, err := egiStatefulSet.Pull(client, statefulSetName, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get statefulSet %q (ns %s): %w", statefulSetName, namespace, err)
	}

	return testStatefulSet.IsReady(1 * time.Second), nil
}

func GetRunningStatefulSet(namespace, statefulSetName string) (*appsv1.StatefulSet, error) {
	return getRunningStatefulSet(egiClients.New(""), namespace, statefulSetName)
}

func getRunningStatefulSet(client *egiClients.Settings, namespace, statefulSetName string) (*appsv1.StatefulSet, error) {
	runningStatefulSet, err := egiStatefulSet.Pull(client, statefulSetName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get statefulSet %q (ns %s): %w", statefulSetName, namespace, err)
	}

	return runningStatefulSet.Object, nil
}

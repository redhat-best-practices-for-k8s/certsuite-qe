package namespaces

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	k8sv1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/pointer"
)

// WaitForDeletion waits until the namespace will be removed from the cluster.
func WaitForDeletion(cs *testclient.ClientSet, nsName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		_, err := cs.Namespaces().Get(context.Background(), nsName, metav1.GetOptions{})
		if k8serrors.IsNotFound(err) {
			glog.V(5).Info(fmt.Sprintf("namespaces %s is not found", nsName))

			return true, nil
		}

		return false, nil
	})
}

// Create creates a new namespace with the given name.
// If the namespace exists, it returns.
func Create(namespace string, cs *testclient.ClientSet) error {
	_, err := cs.Namespaces().Create(context.Background(), &k8sv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		}}, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("namespaces %s already installed", namespace))

		return nil
	}

	return err
}

// DeleteAndWait deletes a namespace and waits until delete.
func DeleteAndWait(clientSet *testclient.ClientSet, namespace string, timeout time.Duration) error {
	err := clientSet.Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return WaitForDeletion(clientSet, namespace, timeout)
}

func Exists(namespace string, cs *testclient.ClientSet) (bool, error) {
	_, err := cs.Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err == nil {
		return true, nil
	}

	if k8serrors.IsNotFound(err) {
		return false, nil
	}

	return false, err
}

// CleanPods deletes all pods in namespace.
func CleanPods(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.Pods(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete pods %w", err)
	}

	return err
}

// CleanDeployments deletes all deployments in namespace.
func CleanDeployments(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.Deployments(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete deployment %w", err)
	}

	return err
}

// CleanDaemonSets deletes all daemonsets in namespace.
func CleanDaemonSets(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.DaemonSets(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete daemonSet %w", err)
	}

	return err
}

func CleanNetworkAttachmentDefinition(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	nadList := &v1.NetworkAttachmentDefinitionList{}
	err = clientSet.List(context.Background(), nadList)

	if err != nil {
		return err
	}

	if len(nadList.Items) > 1 {
		for _, nad := range nadList.Items {
			if nad.Name != "dummy-dhcp-network" {
				err = clientSet.Delete(context.Background(), &nad)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// CleanReplicaSets deletes all ReplicaSets in namespace.
func CleanReplicaSets(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.ReplicaSets(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete ReplicaSet %w", err)
	}

	return err
}

// CleanStatefulSets deletes all StatefulSets in namespace.
func CleanStatefulSets(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.StatefulSets(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete StatefulSet %w", err)
	}

	return err
}

// CleanServices deletes all service in namespace.
func CleanServices(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	serviceList, err := clientSet.Services(namespace).List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, s := range serviceList.Items {
		err = clientSet.Services(namespace).Delete(context.Background(), s.Name, metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64Ptr(0),
		})
		if err != nil {
			return fmt.Errorf("failed to delete service %w", err)
		}
	}

	return err
}

func CleanSubscriptions(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.Subscriptions(namespace).DeleteCollection(context.Background(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64Ptr(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete subscriptions %w", err)
	}

	return err
}

func CleanCSVs(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.ClusterServiceVersions(namespace).DeleteCollection(context.Background(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64Ptr(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete CSVs %w", err)
	}

	return err
}

func CleanInstallPlans(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.InstallPlans(namespace).DeleteCollection(context.Background(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64Ptr(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete installplans %w", err)
	}

	return err
}

func CleanPVCs(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.PersistentVolumeClaims(namespace).DeleteCollection(context.Background(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64Ptr(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete persistent volume claim %w", err)
	}

	return err
}

// Clean cleans all dangling objects from the given namespace.
func Clean(namespace string, clientSet *testclient.ClientSet) error {
	err := CleanDeployments(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanDaemonSets(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanPods(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanReplicaSets(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanStatefulSets(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanNetworkAttachmentDefinition(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanServices(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanSubscriptions(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanCSVs(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanPVCs(namespace, clientSet)
	if err != nil {
		return err
	}

	return CleanInstallPlans(namespace, clientSet)
}

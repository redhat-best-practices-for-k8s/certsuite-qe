package namespaces

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/pointer"
)

// WaitForDeletion waits until the namespace will be removed from the cluster.
func WaitForDeletion(cs corev1Typed.CoreV1Interface, nsName string, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(context.TODO(), time.Second, timeout, true,
		func(ctx context.Context) (bool, error) {
			_, err := cs.Namespaces().Get(ctx, nsName, metav1.GetOptions{})
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
	_, err := cs.Namespaces().Create(context.TODO(), &corev1.Namespace{
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
func DeleteAndWait(clientSet corev1Typed.CoreV1Interface, namespace string, timeout time.Duration) error {
	err := clientSet.Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("namespaces %s is not found", namespace))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	return WaitForDeletion(clientSet, namespace, timeout)
}

func Exists(namespace string, cs *testclient.ClientSet) (bool, error) {
	_, err := cs.Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
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

	err = clientSet.Pods(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete pods %w", err)
	}

	return nil
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

	err = clientSet.Deployments(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete deployment %w", err)
	}

	return nil
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

	err = clientSet.DaemonSets(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete daemonSet %w", err)
	}

	return nil
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

	err = clientSet.List(context.TODO(), nadList)
	if err != nil {
		return err
	}

	if len(nadList.Items) > 1 {
		for _, nad := range nadList.Items {
			if nad.Name != "dummy-dhcp-network" {
				err = clientSet.Delete(context.TODO(), &nad)
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

	err = clientSet.ReplicaSets(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete ReplicaSet %w", err)
	}

	return nil
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

	err = clientSet.StatefulSets(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64(0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete StatefulSet %w", err)
	}

	return nil
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

	serviceList, err := clientSet.Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, s := range serviceList.Items {
		err = clientSet.Services(namespace).Delete(context.TODO(), s.Name, metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64(0),
		})
		if err != nil {
			return fmt.Errorf("failed to delete service %w", err)
		}
	}

	return nil
}

func CleanSubscriptions(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.Subscriptions(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete subscriptions %w", err)
	}

	return nil
}

func CleanCSVs(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.ClusterServiceVersions(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete CSVs %w", err)
	}

	return nil
}

func CleanInstallPlans(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.InstallPlans(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete installplans %w", err)
	}

	return nil
}

func CleanPVCs(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.PersistentVolumeClaims(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete persistent volume claim %w", err)
	}

	return nil
}

func CleanPodDistruptionBudget(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.PolicyV1Interface.PodDisruptionBudgets(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete pod distruption budget %w", err)
	}

	return nil
}

func CleanResourceQuotas(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.ResourceQuotas(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete resource quotas %w", err)
	}

	return nil
}

func CleanNetworkPolicies(namespace string, clientSet *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = clientSet.NetworkPolicies(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64(0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete network policies %w", err)
	}

	return nil
}

// Clean cleans all dangling objects from the given namespace.
func Clean(namespace string, clientSet *testclient.ClientSet) error {
	// check if the namespace exists first
	if _, err := clientSet.Namespaces().Get(context.Background(), namespace, metav1.GetOptions{}); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
	}

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

	err = CleanPodDistruptionBudget(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanResourceQuotas(namespace, clientSet)
	if err != nil {
		return err
	}

	err = CleanNetworkPolicies(namespace, clientSet)
	if err != nil {
		return err
	}

	return CleanInstallPlans(namespace, clientSet)
}

func ApplyResourceQuota(namespace string, clientSet *testclient.ClientSet, quota *corev1.ResourceQuota) error {
	nsExist, err := Exists(namespace, clientSet)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	_, err1 := clientSet.ResourceQuotas(namespace).Create(context.TODO(), quota, metav1.CreateOptions{})

	return err1
}

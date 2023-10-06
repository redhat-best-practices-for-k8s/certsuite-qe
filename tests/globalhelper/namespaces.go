package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	v1alpha1typed "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/ptr"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// WaitForDeletion waits until the namespace will be removed from the cluster.
func WaitForNamespaceDeletion(cs corev1Typed.CoreV1Interface, nsName string, timeout time.Duration) error {
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

func CreateNamespace(namespace string) error {
	return createNamespace(namespace, GetAPIClient().K8sClient.CoreV1())
}

// Create creates a new namespace with the given name.
// If the namespace exists, it returns.
func createNamespace(namespace string, client corev1Typed.CoreV1Interface) error {
	_, err := client.Namespaces().Create(context.TODO(), &corev1.Namespace{
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
func DeleteNamespaceAndWait(namespace string, timeout time.Duration) error {
	return deleteNamespaceAndWait(GetAPIClient().K8sClient.CoreV1(), namespace, timeout)
}

func deleteNamespaceAndWait(clientSet corev1Typed.CoreV1Interface, namespace string, timeout time.Duration) error {
	err := clientSet.Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("namespaces %s is not found", namespace))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	return WaitForNamespaceDeletion(clientSet, namespace, timeout)
}

func NamespaceExists(namespace string) (bool, error) {
	return namespaceExists(namespace, GetAPIClient().K8sClient)
}

func namespaceExists(namespace string, client kubernetes.Interface) (bool, error) {
	_, err := client.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err == nil {
		return true, nil
	}

	if k8serrors.IsNotFound(err) {
		return false, nil
	}

	return false, err
}

// CleanPods deletes all pods in namespace.
func CleanPods(namespace string) error {
	return cleanPods(namespace, GetAPIClient().K8sClient)
}

// CleanPods deletes all pods in namespace.
func cleanPods(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.CoreV1().Pods(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To[int64](0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete pods %w", err)
	}

	return nil
}

// CleanDeployments deletes all deployments in namespace.
func CleanDeployments(namespace string) error {
	return cleanDeployments(namespace, GetAPIClient().K8sClient)
}

func cleanDeployments(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.AppsV1().Deployments(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To[int64](0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete deployment %w", err)
	}

	return nil
}

func CleanDaemonSets(namespace string) error {
	return cleanDaemonSets(namespace, GetAPIClient().K8sClient)
}

// CleanDaemonSets deletes all daemonsets in namespace.
func cleanDaemonSets(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.AppsV1().DaemonSets(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To[int64](0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete daemonSet %w", err)
	}

	return nil
}

// CleanNetworkAttachmentDefinition deletes all network attachment definitions in namespace.
func CleanNetworkAttachmentDefinition(namespace string) error {
	return cleanNetworkAttachmentDefinition(namespace,
		GetAPIClient().K8sClient, GetAPIClient().Client)
}

func cleanNetworkAttachmentDefinition(namespace string, client kubernetes.Interface, rtc runtimeclient.Client) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	nadList := &v1.NetworkAttachmentDefinitionList{}

	err = rtc.List(context.TODO(), nadList)
	if err != nil {
		return err
	}

	if len(nadList.Items) > 1 {
		for _, nad := range nadList.Items {
			if nad.Name != "dummy-dhcp-network" {
				err = rtc.Delete(context.TODO(), &nad)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// CleanReplicaSets deletes all ReplicaSets in namespace.
func CleanReplicaSets(namespace string) error {
	return cleanReplicaSets(namespace, GetAPIClient().K8sClient)
}

func cleanReplicaSets(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.AppsV1().ReplicaSets(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To[int64](0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete ReplicaSet %w", err)
	}

	return nil
}

// CleanStatefulSets deletes all StatefulSets in namespace.
func CleanStatefulSets(namespace string) error {
	return cleanStatefulSets(namespace, GetAPIClient().K8sClient)
}

func cleanStatefulSets(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.AppsV1().StatefulSets(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To[int64](0),
	}, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete StatefulSet %w", err)
	}

	return nil
}

// CleanServices deletes all service in namespace.
func CleanServices(namespace string) error {
	return cleanServices(namespace, GetAPIClient().K8sClient)
}

func cleanServices(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	serviceList, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, s := range serviceList.Items {
		err = client.CoreV1().Services(namespace).Delete(context.TODO(), s.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil {
			return fmt.Errorf("failed to delete service %w", err)
		}
	}

	return nil
}

// CleanSubscriptions deletes all subscriptions in namespace.
func CleanSubscriptions(namespace string) error {
	return cleanSubscriptions(namespace, GetAPIClient().K8sClient, GetAPIClient().OperatorsV1alpha1Interface)
}

func cleanSubscriptions(namespace string, k8sclient kubernetes.Interface, subclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	nsExist, err := namespaceExists(namespace, k8sclient)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = subclient.Subscriptions(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete subscriptions %w", err)
	}

	return nil
}

// CleanCSVs deletes all CSVs in namespace.
func CleanCSVs(namespace string) error {
	return cleanCSVs(namespace, GetAPIClient().K8sClient, GetAPIClient().OperatorsV1alpha1Interface)
}

func cleanCSVs(namespace string, k8sclient kubernetes.Interface, opclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	nsExist, err := namespaceExists(namespace, k8sclient)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = opclient.ClusterServiceVersions(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete CSVs %w", err)
	}

	return nil
}

func CleanInstallPlans(namespace string) error {
	return cleanInstallPlans(namespace, GetAPIClient().K8sClient, GetAPIClient().OperatorsV1alpha1Interface)
}

func cleanInstallPlans(namespace string, k8sclient kubernetes.Interface, opclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	nsExist, err := namespaceExists(namespace, k8sclient)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = opclient.InstallPlans(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete installplans %w", err)
	}

	return nil
}

// CleanPVCs deletes all persistent volume claims in namespace.
func CleanPVCs(namespace string) error {
	return cleanPVCs(namespace, GetAPIClient().K8sClient)
}

func cleanPVCs(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.CoreV1().PersistentVolumeClaims(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete persistent volume claim %w", err)
	}

	return nil
}

// CleanPodDistruptionBudget deletes all pod disruption budget in namespace.
func CleanPodDistruptionBudget(namespace string) error {
	return cleanPodDistruptionBudget(namespace, GetAPIClient().K8sClient)
}

func cleanPodDistruptionBudget(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.PolicyV1().PodDisruptionBudgets(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete pod disruption budget %w", err)
	}

	return nil
}

// CleanResourceQuotas deletes all resource quotas in namespace.
func CleanResourceQuotas(namespace string) error {
	return cleanResourceQuotas(namespace, GetAPIClient().K8sClient)
}

func cleanResourceQuotas(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.CoreV1().ResourceQuotas(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete resource quotas %w", err)
	}

	return nil
}

// CleanNetworkPolicies deletes all network policies in namespace.
func CleanNetworkPolicies(namespace string) error {
	return cleanNetworkPolicies(namespace, GetAPIClient().K8sClient)
}

func cleanNetworkPolicies(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.NetworkingV1().NetworkPolicies(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		},
		metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete network policies %w", err)
	}

	return nil
}

// Clean cleans all dangling objects from the given namespace.
func CleanNamespace(namespace string) error {
	clientSet := GetAPIClient()
	// check if the namespace exists first
	if _, err := clientSet.Namespaces().Get(context.Background(), namespace, metav1.GetOptions{}); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
	}

	err := CleanDeployments(namespace)
	if err != nil {
		return err
	}

	err = CleanDaemonSets(namespace)
	if err != nil {
		return err
	}

	err = CleanPods(namespace)
	if err != nil {
		return err
	}

	err = CleanReplicaSets(namespace)
	if err != nil {
		return err
	}

	err = CleanStatefulSets(namespace)
	if err != nil {
		return err
	}

	err = CleanNetworkAttachmentDefinition(namespace)
	if err != nil {
		return err
	}

	err = CleanServices(namespace)
	if err != nil {
		return err
	}

	err = CleanSubscriptions(namespace)
	if err != nil {
		return err
	}

	err = CleanCSVs(namespace)
	if err != nil {
		return err
	}

	err = CleanPVCs(namespace)
	if err != nil {
		return err
	}

	err = CleanPodDistruptionBudget(namespace)
	if err != nil {
		return err
	}

	err = CleanResourceQuotas(namespace)
	if err != nil {
		return err
	}

	err = CleanNetworkPolicies(namespace)
	if err != nil {
		return err
	}

	return CleanInstallPlans(namespace)
}

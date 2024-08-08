package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/client"

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
func cleanPods(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	// Delete all pods in namespace
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, p := range pods.Items {
		err = client.CoreV1().Pods(namespace).Delete(context.TODO(), p.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete pod %w", err)
		}
	}

	return nil
}

func cleanDeployments(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	deployments, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, d := range deployments.Items {
		err = client.AppsV1().Deployments(namespace).Delete(context.TODO(), d.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete deployment %w", err)
		}
	}

	return nil
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

	daemoneSets, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ds := range daemoneSets.Items {
		err = client.AppsV1().DaemonSets(namespace).Delete(context.TODO(), ds.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete daemonset %w", err)
		}
	}

	return nil
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

func cleanReplicaSets(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	replicaSets, err := client.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, rs := range replicaSets.Items {
		err = client.AppsV1().ReplicaSets(namespace).Delete(context.TODO(), rs.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete replicaSet %w", err)
		}
	}

	return nil
}

func cleanStatefulSets(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	statefulSets, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ss := range statefulSets.Items {
		err = client.AppsV1().StatefulSets(namespace).Delete(context.TODO(), ss.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete statefulSet %w", err)
		}
	}

	return nil
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
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete service %w", err)
		}
	}

	return nil
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

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete subscriptions %w", err)
	}

	return nil
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

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete CSVs %w", err)
	}

	return nil
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

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete installplans %w", err)
	}

	return nil
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

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete persistent volume claim %w", err)
	}

	return nil
}

func cleanPodDisruptionBudget(namespace string, client kubernetes.Interface) error {
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

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete pod disruption budget %w", err)
	}

	return nil
}

func cleanResourceQuotas(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	rqs, err := client.CoreV1().ResourceQuotas(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, rq := range rqs.Items {
		err = client.CoreV1().ResourceQuotas(namespace).Delete(context.TODO(), rq.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete resource quota %w", err)
		}
	}

	return nil
}

func cleanNetworkPolicies(namespace string, client kubernetes.Interface) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	nps, err := client.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, np := range nps.Items {
		err = client.NetworkingV1().NetworkPolicies(namespace).Delete(context.TODO(), np.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete network policy %w", err)
		}
	}

	return nil
}

func CleanNamespace(namespace string) error {
	return cleanNamespace(namespace, GetAPIClient())
}

// Clean cleans all dangling objects from the given namespace.
func cleanNamespace(namespace string, clientSet *client.ClientSet) error {
	// check if the namespace exists first
	exists, err := namespaceExists(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	err = cleanDeployments(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanDaemonSets(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanPods(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanReplicaSets(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanStatefulSets(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanNetworkAttachmentDefinition(namespace, clientSet.K8sClient, clientSet.Client)
	if err != nil {
		return err
	}

	err = cleanServices(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanSubscriptions(namespace, clientSet.K8sClient, clientSet.OperatorsV1alpha1Interface)
	if err != nil {
		return err
	}

	err = cleanCSVs(namespace, clientSet.K8sClient, clientSet.OperatorsV1alpha1Interface)
	if err != nil {
		return err
	}

	err = cleanPVCs(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanPodDisruptionBudget(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanResourceQuotas(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	err = cleanNetworkPolicies(namespace, clientSet.K8sClient)
	if err != nil {
		return err
	}

	return cleanInstallPlans(namespace, clientSet.K8sClient, clientSet.OperatorsV1alpha1Interface)
}

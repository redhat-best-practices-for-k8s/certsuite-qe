package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/client"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiDeployment "github.com/openshift-kni/eco-goinfra/pkg/deployment"
	egiNamespaces "github.com/openshift-kni/eco-goinfra/pkg/namespace"
	egiPod "github.com/openshift-kni/eco-goinfra/pkg/pod"
	egiStatefulSet "github.com/openshift-kni/eco-goinfra/pkg/statefulset"
	egiStorage "github.com/openshift-kni/eco-goinfra/pkg/storage"
	v1alpha1typed "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
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
	return createNamespace(namespace, egiClients.New(""))
}

// Create creates a new namespace with the given name.
// If the namespace exists, it returns.
func createNamespace(namespace string, client *egiClients.Settings) error {
	_, err := egiNamespaces.NewBuilder(client, namespace).Create()

	return err
}

// DeleteAndWait deletes a namespace and waits until delete.
func DeleteNamespaceAndWait(namespace string, timeout time.Duration) error {
	return deleteNamespaceAndWait(egiClients.New(""), namespace, timeout)
}

func deleteNamespaceAndWait(client *egiClients.Settings, namespace string, timeout time.Duration) error {
	builder := egiNamespaces.NewBuilder(client, namespace)

	// Issue the delete request
	err := builder.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	// Wait for the namespace to be deleted
	return WaitForNamespaceDeletion(client.K8sClient.CoreV1(), namespace, timeout)
}

func NamespaceExists(namespace string) (bool, error) {
	return namespaceExists(namespace, egiClients.New(""))
}

func namespaceExists(namespace string, client *egiClients.Settings) (bool, error) {
	builder := egiNamespaces.NewBuilder(client, namespace)

	return builder.Exists(), nil
}

// CleanPods deletes all pods in namespace.
func cleanPods(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	egiPodBuilders, err := egiPod.List(client, namespace, metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Loop through the builders and delete the pods
	for _, podBuilder := range egiPodBuilders {
		_, err = podBuilder.Delete()
		if err != nil {
			return fmt.Errorf("failed to delete pod %w", err)
		}
	}

	return nil
}

func cleanDeployments(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	egiDeploymentBuilders, err := egiDeployment.List(client, namespace, metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Loop through the builders and delete the deployments
	for _, deploymentBuilder := range egiDeploymentBuilders {
		err = deploymentBuilder.Delete()
		if err != nil {
			return fmt.Errorf("failed to delete deployment %w", err)
		}
	}

	return nil
}

// CleanDaemonSets deletes all daemonsets in namespace.
func cleanDaemonSets(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	daemoneSets, err := client.K8sClient.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ds := range daemoneSets.Items {
		err = client.K8sClient.AppsV1().DaemonSets(namespace).Delete(context.TODO(), ds.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete daemonset %w", err)
		}
	}

	return nil
}

func cleanNetworkAttachmentDefinition(namespace string, client *egiClients.Settings, rtc runtimeclient.Client) error {
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

func cleanReplicaSets(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	replicaSets, err := client.K8sClient.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, rs := range replicaSets.Items {
		err = client.K8sClient.AppsV1().ReplicaSets(namespace).Delete(context.TODO(), rs.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete replicaSet %w", err)
		}
	}

	return nil
}

func cleanStatefulSets(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	statefulSets, err := client.K8sClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ss := range statefulSets.Items {
		err = client.K8sClient.AppsV1().StatefulSets(namespace).Delete(context.TODO(), ss.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete statefulSet %w", err)
		}
	}

	return nil
}

func cleanServices(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	serviceList, err := client.K8sClient.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, s := range serviceList.Items {
		err = client.K8sClient.CoreV1().Services(namespace).Delete(context.TODO(), s.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete service %w", err)
		}
	}

	return nil
}

func cleanSubscriptions(namespace string, client *egiClients.Settings, subclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	nsExist, err := namespaceExists(namespace, client)
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

func cleanCSVs(namespace string, client *egiClients.Settings, opclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	nsExist, err := namespaceExists(namespace, client)
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

func cleanInstallPlans(namespace string, client *egiClients.Settings, opclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	nsExist, err := namespaceExists(namespace, client)
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

func cleanPVCs(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	pvcBuilder, err := egiStorage.ListPVC(client, namespace, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pvc := range pvcBuilder {
		err = pvc.Delete()
		if err != nil {
			return fmt.Errorf("failed to delete persistent volume claim %w", err)
		}
	}

	return nil
}

func cleanPodDisruptionBudget(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	err = client.K8sClient.PolicyV1().PodDisruptionBudgets(namespace).DeleteCollection(context.TODO(),
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		},
		metav1.ListOptions{})

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete pod disruption budget %w", err)
	}

	return nil
}

func cleanResourceQuotas(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	rqs, err := client.K8sClient.CoreV1().ResourceQuotas(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, rq := range rqs.Items {
		err = client.K8sClient.CoreV1().ResourceQuotas(namespace).Delete(context.TODO(), rq.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete resource quota %w", err)
		}
	}

	return nil
}

func cleanNetworkPolicies(namespace string, client *egiClients.Settings) error {
	nsExist, err := namespaceExists(namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	nps, err := client.K8sClient.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, np := range nps.Items {
		err = client.K8sClient.NetworkingV1().NetworkPolicies(namespace).Delete(context.TODO(), np.Name, metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete network policy %w", err)
		}
	}

	return nil
}

func CleanNamespace(namespace string) error {
	return cleanNamespace(namespace, GetAPIClient(), egiClients.New(""))
}

// Clean cleans all dangling objects from the given namespace.
func cleanNamespace(namespace string, clientSet *client.ClientSet, egiClient *egiClients.Settings) error {
	// check if the namespace exists first
	exists, err := namespaceExists(namespace, egiClient)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	builder := egiNamespaces.NewBuilder(egiClient, namespace)

	err = builder.CleanObjects(5*time.Second,
		egiDeployment.GetGVR(),
		// egiDaemonSet.GetGVR(),
		// egiPods.GetGVR(),
		// egiReplicaSet.GetGVR(),
		egiStatefulSet.GetGVR(),
	)

	if err != nil {
		return err
	}

	err = cleanPods(namespace, egiClient)
	if err != nil {
		return err
	}

	err = cleanReplicaSets(namespace, egiClient)
	if err != nil {
		return err
	}

	err = cleanStatefulSets(namespace, egiClient)
	if err != nil {
		return err
	}

	err = cleanNetworkAttachmentDefinition(namespace, egiClient, clientSet.Client)
	if err != nil {
		return err
	}

	err = cleanServices(namespace, egiClient)
	if err != nil {
		return err
	}

	err = cleanSubscriptions(namespace, egiClient, clientSet.OperatorsV1alpha1Interface)
	if err != nil {
		return err
	}

	err = cleanCSVs(namespace, egiClient, clientSet.OperatorsV1alpha1Interface)
	if err != nil {
		return err
	}

	err = cleanPVCs(namespace, egiClient)
	if err != nil {
		return err
	}

	err = cleanPodDisruptionBudget(namespace, egiClient)
	if err != nil {
		return err
	}

	err = cleanResourceQuotas(namespace, egiClient)
	if err != nil {
		return err
	}

	err = cleanNetworkPolicies(namespace, egiClient)
	if err != nil {
		return err
	}

	return cleanInstallPlans(namespace, egiClient, clientSet.OperatorsV1alpha1Interface)
}

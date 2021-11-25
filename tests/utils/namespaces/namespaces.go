package namespaces

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"

	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/pointer"
)

// WaitForDeletion waits until the namespace will be removed from the cluster
func WaitForDeletion(cs *testclient.ClientSet, nsName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		_, err := cs.Namespaces().Get(context.Background(), nsName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
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

// DeleteAndWait deletes a namespace and waits until delete
func DeleteAndWait(cs *testclient.ClientSet, namespace string, timeout time.Duration) error {
	err := cs.Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return WaitForDeletion(cs, namespace, timeout)
}

func Exists(namespace string, cs *testclient.ClientSet) (bool, error) {
	_, err := cs.Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err == nil {
		return true, nil
	} else {
		if k8serrors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}
}

// CleanPods deletes all pods in namespace
func CleanPods(namespace string, cs *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, cs)
	if err != nil {
		return err
	}
	if !nsExist {
		return nil
	}
	err = cs.Pods(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pods %v", err)
	}
	return err
}

// CleanDeployments deletes all deployments in namespace
func CleanDeployments(namespace string, cs *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, cs)
	if err != nil {
		return err
	}
	if !nsExist {
		return nil
	}
	err = cs.Deployments(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete deployment %v", err)
	}
	return err
}

// CleanDaemonSets deletes all daemonsets in namespace
func CleanDaemonSets(namespace string, cs *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, cs)
	if err != nil {
		return err
	}
	if !nsExist {
		return nil
	}
	err = cs.DaemonSets(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete daemonSet %v", err)
	}
	return err
}

// CleanServices deletes all service in namespace
func CleanServices(namespace string, cs *testclient.ClientSet) error {
	nsExist, err := Exists(namespace, cs)
	if err != nil {
		return err
	}
	if !nsExist {
		return nil
	}
	serviceList, err := cs.Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, s := range serviceList.Items {
		err = cs.Services(namespace).Delete(context.Background(), s.Name, metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64Ptr(0),
		})
		if err != nil {
			return fmt.Errorf("failed to delete service %v", err)
		}
	}
	return err
}

// Clean cleans all dangling objects from the given namespace.
func Clean(namespace string, cs *testclient.ClientSet) error {
	err := CleanDeployments(namespace, cs)
	if err != nil {
		return err
	}
	err = CleanDaemonSets(namespace, cs)
	if err != nil {
		return err
	}
	err = CleanPods(namespace, cs)
	if err != nil {
		return err
	}
	return CleanServices(namespace, cs)
}

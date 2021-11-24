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

type Namespace struct {
	k8sv1.Namespace
}

// DefineNamespace return namespace struct
func DefineNamespace(name string) *Namespace {
	return &Namespace{
		k8sv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			}}}
}

// WaitForDeletion waits until the namespace will be removed from the cluster
func (namespace Namespace) WaitForDeletion(cs *testclient.ClientSet, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		_, err := cs.Namespaces().Get(context.Background(), namespace.Name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			glog.V(5).Info(fmt.Sprintf("namespaces %s is not found", namespace.Name))
			return true, nil
		}
		return false, nil
	})
}

// Create creates a new namespace with the given name.
// If the namespace exists, it returns.
func (namespace *Namespace) Create(cs *testclient.ClientSet) error {
	_, err := cs.Namespaces().Create(context.Background(), &namespace.Namespace, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("namespaces %s already installed", namespace.Name))
		return nil
	}
	return err
}

// DeleteAndWait deletes a namespace and waits until delete
func (namespace Namespace) DeleteAndWait(cs *testclient.ClientSet, timeout time.Duration) error {
	err := cs.Namespaces().Delete(context.Background(), namespace.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return namespace.WaitForDeletion(cs, timeout)
}

func (namespace *Namespace) Exists(cs *testclient.ClientSet) (bool, error) {
	_, err := cs.Namespaces().Get(context.Background(), namespace.Name, metav1.GetOptions{})
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
func (namespace *Namespace) CleanPods(cs *testclient.ClientSet) error {
	nsExist, err := namespace.Exists(cs)
	if err != nil {
		return err
	}
	if !nsExist {
		return nil
	}
	err = cs.Pods(namespace.Name).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pods %v", err)
	}
	return err
}

// CleanDeployments deletes all deployments in namespace
func (namespace *Namespace) CleanDeployments(cs *testclient.ClientSet) error {
	nsExist, err := namespace.Exists(cs)
	if err != nil {
		return err
	}
	if !nsExist {
		return nil
	}
	err = cs.Deployments(namespace.Name).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete deployment %v", err)
	}
	return err
}

// CleanDaemonSets deletes all daemonsets in namespace
func (namespace *Namespace) CleanDaemonSets(cs *testclient.ClientSet) error {
	nsExist, err := namespace.Exists(cs)
	if err != nil {
		return err
	}
	if !nsExist {
		return nil
	}
	err = cs.DaemonSets(namespace.Name).DeleteCollection(context.Background(), metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	}, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete daemonSet %v", err)
	}
	return err
}

// Clean cleans all dangling objects from the given namespace.
func (namespace Namespace) Clean(cs *testclient.ClientSet) error {
	err := namespace.CleanPods(cs)
	if err != nil {
		return err
	}
	err = namespace.CleanDeployments(cs)
	if err != nil {
		return err
	}
	err = namespace.CleanDaemonSets(cs)
	return err
}

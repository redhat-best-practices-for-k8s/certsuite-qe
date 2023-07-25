package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreatePersistentVolumeClaim(pvcName, namespace string) error {
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: namespace,
		},
	}

	_, err :=
		GetAPIClient().CoreV1Interface.PersistentVolumeClaims(namespace).Create(context.Background(), &pvc, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("persistentvolumeclaim %s already installed", pvcName))

		return nil
	}

	return err
}

func DeletePersistentVolumeClaim(pvcName, namespace string) error {
	pvc, err := GetAPIClient().ServiceAccounts(namespace).Get(context.Background(), pvcName, metav1.GetOptions{})

	if k8serrors.IsNotFound(err) {
		return nil
	}

	if pvc == nil {
		return err
	}

	err = GetAPIClient().PersistentVolumes().Delete(context.Background(), pvcName,
		metav1.DeleteOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete persistentvolumeclaim %w", err)
	}

	return nil
}

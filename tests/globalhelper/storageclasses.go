package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	storagev1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	storagev1typed "k8s.io/client-go/kubernetes/typed/storage/v1"
	"k8s.io/utils/ptr"
)

func CreateStorageClass(storageClassName string, defaultSC bool) error {
	return createStorageClass(GetAPIClient().K8sClient.StorageV1(), storageClassName, defaultSC)
}

func createStorageClass(client storagev1typed.StorageV1Interface, storageClassName string, defaultSC bool) error {
	storageClassTemplate := storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: storageClassName,
		},
		Provisioner: "kubernetes.io/no-provisioner",
	}

	// Set the storageclass as default if needed.
	if defaultSC {
		storageClassTemplate.Annotations = map[string]string{
			"storageclass.kubernetes.io/is-default-class": "true",
		}
	}

	_, err := client.StorageClasses().Create(context.Background(),
		&storageClassTemplate, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("storageclass %s already installed", storageClassName))

		return nil
	}

	return err
}

func DeleteStorageClass(storageClassName string) error {
	return deleteStorageClass(GetAPIClient().K8sClient.StorageV1(), storageClassName)
}

func deleteStorageClass(client storagev1typed.StorageV1Interface, storageClassName string) error {
	err := client.StorageClasses().Delete(context.Background(),
		storageClassName, metav1.DeleteOptions{GracePeriodSeconds: ptr.To[int64](0)})

	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("storageclass %s already deleted", storageClassName))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete storageclass %w", err)
	}

	return nil
}

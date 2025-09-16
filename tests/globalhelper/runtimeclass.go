package globalhelper

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/gomega"
	nodev1 "k8s.io/api/node/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	klog "k8s.io/klog/v2"
)

func CreateRunTimeClass(rtc *nodev1.RuntimeClass) error {
	return createRunTimeClass(GetAPIClient().K8sClient, rtc)
}

func createRunTimeClass(client kubernetes.Interface, rtc *nodev1.RuntimeClass) error {
	rtc, err := client.NodeV1().RuntimeClasses().Create(context.TODO(), rtc, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Info(fmt.Sprintf("runtimeclass %s already created", rtc.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create runtimeclass %q (ns %s): %w", rtc.Name, rtc.Namespace, err)
	}

	Eventually(func() bool {
		rtcCreated, err := isRtcCreated(client, rtc)
		if err != nil {
			klog.V(5).Info(fmt.Sprintf("rtc %s was not created, retry in %d seconds", rtc.Name, retryInterval))

			return false
		}

		return rtcCreated
	}, retryInterval*time.Minute, retryInterval*time.Second).Should(Equal(true), "rtc was not created")

	return nil
}

func isRtcCreated(client kubernetes.Interface, rtc *nodev1.RuntimeClass) (bool, error) {
	rtc, err := client.NodeV1().RuntimeClasses().Get(context.TODO(), rtc.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return len(rtc.UID) != 0, nil
}

func DeleteRunTimeClass(rtc *nodev1.RuntimeClass) error {
	return deleteRunTimeClass(GetAPIClient().K8sClient, rtc)
}

func deleteRunTimeClass(client kubernetes.Interface, rtc *nodev1.RuntimeClass) error {
	err := client.NodeV1().RuntimeClasses().Delete(context.TODO(), rtc.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete runtimeclass %q (ns %s): %w", rtc.Name, rtc.Namespace, err)
	}

	Eventually(func() bool {
		rtcDeleted, err := isRtcDeleted(client, rtc)
		if err != nil {
			klog.V(5).Info(fmt.Sprintf("rtc %s was not deleted, retry in %d seconds", rtc.Name, retryInterval))

			return false
		}

		return rtcDeleted
	}, retryInterval*time.Minute, retryInterval*time.Second).Should(Equal(true), "rtc was not deleted")

	return nil
}

func isRtcDeleted(client kubernetes.Interface, rtc *nodev1.RuntimeClass) (bool, error) {
	_, err := client.NodeV1().RuntimeClasses().Get(context.TODO(), rtc.Name, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return true, nil
	}

	return false, err
}

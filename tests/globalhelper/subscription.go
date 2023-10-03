package globalhelper

import (
	"context"
	"fmt"

	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func DeleteSubscription(namespace string, subscriptionName string, clientSet *testclient.ClientSet) error {
	subscription, err := clientSet.Subscriptions(namespace).Get(context.TODO(),
		subscriptionName,
		metav1.GetOptions{})

	if k8serrors.IsNotFound(err) {
		return nil
	}

	if subscription == nil {
		return err
	}

	err = clientSet.Subscriptions(namespace).Delete(context.TODO(),
		subscriptionName,
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})

	if err != nil {
		return fmt.Errorf("failed to delete subscription %w", err)
	}

	return nil
}

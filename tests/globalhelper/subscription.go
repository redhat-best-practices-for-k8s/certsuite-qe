package globalhelper

import (
	"context"
	"fmt"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	v1alpha1typed "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// DeleteSubscription deletes a subscription in the given namespace.
func DeleteSubscription(namespace string, subscriptionName string) error {
	return deleteSubscription(namespace, subscriptionName, GetAPIClient().OperatorsV1alpha1Interface)
}

func deleteSubscription(namespace string, subscriptionName string, client v1alpha1typed.OperatorsV1alpha1Interface) error {
	subscription, err := client.Subscriptions(namespace).Get(context.TODO(),
		subscriptionName,
		metav1.GetOptions{})

	if k8serrors.IsNotFound(err) {
		return nil
	}

	if subscription == nil {
		return err
	}

	err = client.Subscriptions(namespace).Delete(context.TODO(),
		subscriptionName,
		metav1.DeleteOptions{
			GracePeriodSeconds: ptr.To[int64](0),
		})

	if err != nil {
		return fmt.Errorf("failed to delete subscription %w", err)
	}

	return nil
}

// CreateSubscription creates a subscription in the given namespace.
func CreateSubscription(namespace string, subscription *v1alpha1.Subscription) error {
	return createSubscription(namespace, subscription, GetAPIClient().OperatorsV1alpha1Interface)
}

func createSubscription(namespace string, subscription *v1alpha1.Subscription, client v1alpha1typed.OperatorsV1alpha1Interface) error {
	_, err := client.Subscriptions(namespace).Create(context.TODO(),
		subscription,
		metav1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("failed to create subscription %w", err)
	}

	return nil
}

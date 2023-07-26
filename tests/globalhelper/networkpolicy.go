package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
)

func CreateAndWaitUntilNetworkPolicyIsReady(networkPolicy *networkingv1.NetworkPolicy, timeout time.Duration) error {
	policy, err := GetAPIClient().NetworkPolicies(networkPolicy.Namespace).Create(
		context.TODO(), networkPolicy, metav1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("failed to create Network Policy %q (ns %s): %w", networkPolicy.Name, networkPolicy.Namespace, err)
	}

	Eventually(func() bool {
		status, err := doesNetworkPolicyExist(policy.Namespace, policy.Name)
		if err != nil {
			glog.Fatal(fmt.Sprintf(
				"Network Policy %s is not ready.", policy.Name))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "Network Policy is not ready")

	return nil
}

func doesNetworkPolicyExist(namespace, name string) (bool, error) {
	if _, err := GetAPIClient().NetworkPolicies(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	); err != nil {
		return false, err
	}

	return true, nil
}

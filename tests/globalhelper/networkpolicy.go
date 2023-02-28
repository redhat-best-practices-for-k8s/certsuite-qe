package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
)

func CreateAndWaitUntilNetworkPolicyIsReady(networkPolicy *v1.NetworkPolicy, timeout time.Duration) error {
	policy, err := APIClient.NetworkPolicies(networkPolicy.Namespace).Create(
		context.Background(), networkPolicy, metav1.CreateOptions{})

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
	_, err := APIClient.NetworkPolicies(namespace).Get(
		context.Background(),
		name,
		metav1.GetOptions{},
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

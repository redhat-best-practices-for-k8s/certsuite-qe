package globalhelper

import (
	"context"
	"time"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
)

func CreateNetworkPolicy(networkPolicy *networkingv1.NetworkPolicy, timeout time.Duration) error {
	client := egiClients.New("")
	_, err := client.NetworkPolicies(networkPolicy.Namespace).Create(
		context.TODO(),
		networkPolicy,
		metav1.CreateOptions{})

	return err
}

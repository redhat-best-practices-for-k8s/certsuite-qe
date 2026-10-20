package globalhelper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestGetServicesFromNamespace(t *testing.T) {
	client := k8sfake.NewClientset(
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "in-namespace", Namespace: "target"}},
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "elsewhere", Namespace: "other"}},
	)

	services, err := getServicesFromNamespace("target", client)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Len(t, services, 1) {
		return
	}

	assert.Equal(t, "in-namespace", services[0].Name)
}

func TestGetServicesFromNamespaceListError(t *testing.T) {
	client := k8sfake.NewClientset()
	client.PrependReactor("list", "services", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("service list failed")
	})

	services, err := getServicesFromNamespace("target", client)
	assert.Nil(t, services)
	assert.EqualError(t, err, "service list failed")
}

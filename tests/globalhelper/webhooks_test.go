package globalhelper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestDeleteWebhookConfigurations(t *testing.T) {
	testCases := []struct {
		name       string
		resource   string
		webhook    runtime.Object
		deleteFunc func(*k8sfake.Clientset, string) error
	}{
		{
			name:     "mutating",
			resource: "mutatingwebhookconfigurations",
			webhook: &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: metav1.ObjectMeta{Name: "test-webhook"},
			},
			deleteFunc: func(client *k8sfake.Clientset, name string) error {
				return deleteMutatingWebhookConfiguration(client.AdmissionregistrationV1(), name)
			},
		},
		{
			name:     "validating",
			resource: "validatingwebhookconfigurations",
			webhook: &admissionregistrationv1.ValidatingWebhookConfiguration{
				ObjectMeta: metav1.ObjectMeta{Name: "test-webhook"},
			},
			deleteFunc: func(client *k8sfake.Clientset, name string) error {
				return deleteValidatingWebhookConfiguration(client.AdmissionregistrationV1(), name)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Run("deletes existing configuration", func(t *testing.T) {
				client := k8sfake.NewClientset(testCase.webhook)

				err := testCase.deleteFunc(client, "test-webhook")
				assert.NoError(t, err)

				actions := client.Actions()
				if !assert.Len(t, actions, 1) {
					return
				}

				assert.Equal(t, testCase.resource, actions[0].GetResource().Resource)
				assert.Equal(t, "delete", actions[0].GetVerb())

				deleteAction, ok := actions[0].(k8stesting.DeleteAction)
				if !assert.True(t, ok) {
					return
				}

				assert.Equal(t, "test-webhook", deleteAction.GetName())
			})

			t.Run("missing configuration is ignored", func(t *testing.T) {
				client := k8sfake.NewClientset()

				assert.NoError(t, testCase.deleteFunc(client, "missing-webhook"))
			})

			t.Run("wraps API errors", func(t *testing.T) {
				client := k8sfake.NewClientset()
				client.PrependReactor("delete", testCase.resource, func(k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, errors.New("API unavailable")
				})

				err := testCase.deleteFunc(client, "test-webhook")
				if !assert.Error(t, err) {
					return
				}

				assert.Contains(t, err.Error(), "failed to delete")
				assert.Contains(t, err.Error(), "API unavailable")
			})
		})
	}
}

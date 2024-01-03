package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/resourcequota"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestCreateResourceQuota(t *testing.T) {
	testCases := []struct {
		testResourceQuota  *corev1.ResourceQuota
		quotaAlreadyExists bool
	}{
		{
			testResourceQuota:  resourcequota.DefineResourceQuota("quota1", "default", "1", "1", "1", "1"),
			quotaAlreadyExists: false,
		},
		{
			testResourceQuota:  resourcequota.DefineResourceQuota("quota2", "default", "1", "1", "1", "1"),
			quotaAlreadyExists: true,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		// Create namespace object
		runtimeObjects = append(runtimeObjects, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: testCase.testResourceQuota.Namespace,
			},
		})

		// Cover the scenario where the resource quota already exists
		if testCase.quotaAlreadyExists {
			runtimeObjects = append(runtimeObjects, testCase.testResourceQuota)
		}

		// Add the runtime objects to the fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		err := createResourceQuota(client, testCase.testResourceQuota)
		assert.Nil(t, err)
	}
}

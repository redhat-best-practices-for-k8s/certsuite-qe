package globalhelper

import (
	"testing"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGetRunningStatefulSet(t *testing.T) {
	generateStatefulSet := func(name, namespace string) *appsv1.StatefulSet {
		return &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Status: appsv1.StatefulSetStatus{
				ReadyReplicas: 1,
			},
		}
	}

	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, generateStatefulSet("testStatefulSet", "testNamespace"))
	fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
		K8sMockObjects: runtimeObjects,
	})
	statefulSet, err := getRunningStatefulSet(fakeClient, "testNamespace", "testStatefulSet")
	assert.Nil(t, err)
	assert.Equal(t, "testStatefulSet", statefulSet.Name)
	assert.Equal(t, "testNamespace", statefulSet.Namespace)
}

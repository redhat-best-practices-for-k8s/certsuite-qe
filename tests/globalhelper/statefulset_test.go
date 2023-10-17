package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
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
	client := k8sfake.NewSimpleClientset(runtimeObjects...)
	statefulSet, err := getRunningStatefulSet(client.AppsV1(), "testNamespace", "testStatefulSet")
	assert.Nil(t, err)
	assert.Equal(t, "testStatefulSet", statefulSet.Name)
	assert.Equal(t, "testNamespace", statefulSet.Namespace)
}

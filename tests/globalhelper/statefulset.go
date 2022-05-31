package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateAndWaitUntilStatefulSetIsReady creates statefulset and waits until all its replicas are up and running.
func CreateAndWaitUntilStatefulSetIsReady(statefulSet *v1.StatefulSet, timeout time.Duration) error {
	runningStatefulSet, err := APIClient.StatefulSets(statefulSet.Namespace).Create(
		context.Background(),
		statefulSet,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		st, err := APIClient.StatefulSets(statefulSet.Namespace).Get(
			context.Background(),
			statefulSet.Name,
			metav1.GetOptions{})
		if err != nil || *st.Spec.Replicas != st.Status.ReadyReplicas {
			glog.V(5).Info(fmt.Sprintf(
				"statefulSet %s is not ready, retry in 5 seconds", runningStatefulSet.Name))

			return false
		}

		return true
	}, timeout, 5*time.Second).Should(Equal(true), "statefulSet is not ready")

	return nil
}

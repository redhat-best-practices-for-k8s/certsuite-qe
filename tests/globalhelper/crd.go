package globalhelper

import (
	"context"
	"fmt"
	"time"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"
)

const (
	crdRetryIntervalSecs = 5
)

func CreateAndWaitUntilCrdIsReady(crd *apiextv1.CustomResourceDefinition, timeout time.Duration) error {
	_, err := APIClient.CustomResourceDefinitions().Create(
		context.Background(),
		crd,
		metav1.CreateOptions{},
	)
	if err != nil {
		return err
	}

	Eventually(func() bool {
		runningCrd, err := APIClient.CustomResourceDefinitions().Get(
			context.Background(),
			crd.Name,
			metav1.GetOptions{},
		)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"crd %s is not ready, retry in 5 seconds", runningCrd.Name))

			return false
		}

		for _, condition := range runningCrd.Status.Conditions {
			if condition.Type == apiextv1.Established {
				return true
			}
		}

		return false
	}, timeout, crdRetryIntervalSecs*time.Second).Should(Equal(true), "CRD is not ready")

	return nil
}

func DeleteCrdAndWaitUntilIsRemoved(crd string, timeout time.Duration) {
	err := APIClient.CustomResourceDefinitions().Delete(
		context.Background(),
		crd,
		metav1.DeleteOptions{})
	Expect(err).ToNot(HaveOccurred())

	Eventually(func() bool {
		_, err := APIClient.CustomResourceDefinitions().Get(
			context.Background(),
			crd,
			metav1.GetOptions{})

		// If the CRD was already removed, we'll get an error.
		return err != nil
	}, timeout, crdRetryIntervalSecs*time.Second).Should(Equal(true), "CRD is not removed yet")
}

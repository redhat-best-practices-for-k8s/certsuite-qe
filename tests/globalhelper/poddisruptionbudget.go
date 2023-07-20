package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
)

// IsPodDisruptionBudgetCreated checks if a pod disruption budget is created.
func IsPodDisruptionBudgetCreated(pdbName string, namespace string) (bool, error) {
	pdb, err := GetAPIClient().PolicyV1Interface.PodDisruptionBudgets(namespace).Get(context.Background(), pdbName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return len(pdb.UID) != 0, nil
}

// CreatPodDisruptionBudget creates Pod Disruption Budget and wait until pdb is created.
func CreatePodDisruptionBudget(pdb *policyv1.PodDisruptionBudget, timeout time.Duration) error {
	poddisruptionbudget, err := GetAPIClient().PolicyV1Interface.PodDisruptionBudgets(pdb.Namespace).Create(
		context.Background(),
		pdb,
		metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Pod Disruption Budget %q (ns %s): %w", pdb.Name, pdb.Namespace, err)
	}

	Eventually(func() bool {
		status, err := IsPodDisruptionBudgetCreated(poddisruptionbudget.Name, poddisruptionbudget.Namespace)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"Pod Disruption Budget %s is not ready, retry in 5 seconds", poddisruptionbudget.Name))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "Pod Disruption Budget is not ready")

	return nil
}

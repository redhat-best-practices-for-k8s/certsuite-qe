package globalhelper

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnableMasterScheduling enables/disables master nodes scheduling.
func EnableMasterScheduling(scheduleable bool) error {
	scheduler, err := APIClient.ConfigV1Interface.Schedulers().Get(context.TODO(), "cluster", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get schedulers: %w", err)
	}

	scheduler.Spec.MastersSchedulable = scheduleable

	_, err = APIClient.ConfigV1Interface.Schedulers().Update(context.TODO(), scheduler, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update schedulers: %w", err)
	}

	return nil
}

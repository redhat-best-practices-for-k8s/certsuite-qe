package globalhelper

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnableMasterScheduling enables/disables master nodes scheduling.
func EnableMasterScheduling(scheduleable bool) error {
	scheduler, err := APIClient.ConfigV1Interface.Schedulers().Get(
		context.TODO(), "cluster", metav1.GetOptions{})
	if err != nil {
		return err
	}

	scheduler.Spec.MastersSchedulable = scheduleable
	_, err = APIClient.ConfigV1Interface.Schedulers().Update(context.TODO(),
		scheduler, metav1.UpdateOptions{})

	return err
}

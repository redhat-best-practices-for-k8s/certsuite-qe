package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/node/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateRunTimeClass(rtc *v1.RuntimeClass) error {
	rtc, err := APIClient.RuntimeClasses().Create(context.Background(), rtc, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		rtcCreated, err := isRtcCreated(rtc)
		if err != nil {

			glog.V(5).Info(fmt.Sprintf(
				"rtc %s was not created, retry in %d seconds", rtc.Name, retryInterval))

			return false
		}

		return rtcCreated
	}, retryInterval*time.Minute, retryInterval*time.Second).Should(Equal(true), "rtc was not created")

	return nil
}

func isRtcCreated(rtc *v1.RuntimeClass) (bool, error) {
	rtc, err := APIClient.RuntimeClasses().Get(context.Background(), rtc.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return len(rtc.UID) != 0, nil
}

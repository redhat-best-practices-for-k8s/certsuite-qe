package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"
	nodev1 "k8s.io/api/node/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateRunTimeClass(rtc *nodev1.RuntimeClass) error {
	rtc, err := GetAPIClient().RuntimeClasses().Create(context.TODO(), rtc, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("runtimeclass %s already created", rtc.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create runtimeclass %q (ns %s): %w", rtc.Name, rtc.Namespace, err)
	}

	Eventually(func() bool {
		rtcCreated, err := isRtcCreated(rtc)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf("rtc %s was not created, retry in %d seconds", rtc.Name, retryInterval))

			return false
		}

		return rtcCreated
	}, retryInterval*time.Minute, retryInterval*time.Second).Should(Equal(true), "rtc was not created")

	return nil
}

func isRtcCreated(rtc *nodev1.RuntimeClass) (bool, error) {
	rtc, err := GetAPIClient().RuntimeClasses().Get(context.TODO(), rtc.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return len(rtc.UID) != 0, nil
}

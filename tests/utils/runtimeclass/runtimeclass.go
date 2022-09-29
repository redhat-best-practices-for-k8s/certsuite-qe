package runtimeclass

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	v1 "k8s.io/api/node/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	retryInterval = 5
)

func DefineRunTimeClass(name string) *v1.RuntimeClass {
	return &v1.RuntimeClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Handler: "runc",
	}
}

func CreateRunTimeClass(rtc *v1.RuntimeClass) error {
	rtc, err := globalhelper.APIClient.RuntimeClasses().Create(context.Background(), rtc, metav1.CreateOptions{})
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
	rtc, err := globalhelper.APIClient.RuntimeClasses().Get(context.Background(), rtc.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return len(rtc.UID) != 0, nil
}

package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime   = 5 * time.Minute
	RetryInterval = 5
)

var (
	LifecycleNamespace     = "lifecycle-tests"
	testPodLabelPrefixName = "test-network-function.com/lifecycle"
	testPodLabelValue      = "testing"
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentLabels   = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
	PreStopCommand = []string{"/bin/sh", "-c", "killall -0 tail"}

	TnfShutdownTcName            = "lifecycle-container-shutdown"
	TnfDeploymentScalingTcName   = "lifecycle-deployment-scaling"
	TnfPodOwnerTypeTcName        = "lifecycle-pod-owner-type"
	TnfPodRecreationTcName       = "lifecycle-pod-recreation"
	TnfPodHighAvailabilityTcName = "lifecycle-pod-high-availability"
	TnfPodSchedulingTcName       = "lifecycle-pod-scheduling"
	TnfLivenessTcName            = "lifecycle-liveness-probe"
	TnfReadinessTcName           = "lifecycle-readiness-probe"
	TnfStatefulSetScalingTcName  = "lifecycle-statefulset-scaling"
	TnfImagePullPolicyTcName     = "lifecycle-image-pull-policy"
)

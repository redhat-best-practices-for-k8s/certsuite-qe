package lifeparameters

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
	PreStopCommand         = []string{"/bin/sh", "-c", "killall -0 tail"}
	LifecycleTestSuiteName = "lifecycle"

	SkipsTcsSlice = []string{"lifecycle-pod-recreation", "lifecycle-pod-scheduling",
		"lifecycle-deployment-scaling", "lifecycle-pod-termination-grace-period",
		"lifecycle-pod-owner-type", "lifecycle-container-shutdown", "lifecycle-image-pull-policy",
		"lifecycle-liveness", "lifecycle-readiness", "lifecycle-statefulset-scaling",
		"lifecycle-pod-high-availability"}

	ShutdownName               = "lifecycle-container-shutdown"
	DeploymentScalingName      = "lifecycle-deployment-scaling"
	TerminationGracePeriodName = "lifecycle-pod-termination-grace-period"
	PodOwnerTypeName           = "lifecycle-pod-owner-type"
	PodRecreationName          = "lifecycle-pod-recreation"
	PodHighAvailabilityName    = "lifecycle-pod-high-availability"
)

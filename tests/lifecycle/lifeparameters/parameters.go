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
	PreStopCommand          = []string{"/bin/sh", "-c", "killall -0 tail"}
	LifecycleTestSuiteName  = "lifecycle"
	SkipAllButShutdownRegex = "lifecycle-pod-high-availability lifecycle-pod-scheduling" +
		" lifecycle-pod-termination-grace-period lifecycle-pod-owner-type" +
		" lifecycle-pod-recreation lifecycle-scaling lifecycle-image-pull-policy"

	SkipAllButScalingRegex = "lifecycle-pod-high-availability lifecycle-pod-scheduling" +
		" lifecycle-pod-termination-grace-period lifecycle-pod-owner-type" +
		" lifecycle-pod-recreation lifecycle-container-shutdown lifecycle-image-pull-policy"

	SkipAllButTerminationGracePeriodRegex = "lifecycle-pod-high-availability lifecycle-pod-scheduling" +
		" lifecycle-scaling lifecycle-pod-owner-type" +
		" lifecycle-pod-recreation lifecycle-container-shutdown lifecycle-image-pull-policy"

	SkipAllButPodOwnerTypeRegex = "lifecycle-pod-high-availability lifecycle-pod-scheduling" +
		" lifecycle-scaling lifecycle-pod-termination-grace-period" +
		" lifecycle-pod-recreation lifecycle-container-shutdown lifecycle-image-pull-policy"

	SkipAllButPodRecreationRegex = "lifecycle-pod-high-availability lifecycle-pod-scheduling" +
		" lifecycle-scaling lifecycle-pod-termination-grace-period" +
		" lifecycle-pod-owner-type lifecycle-container-shutdown lifecycle-image-pull-policy"

	ShutdownDefaultName               = "lifecycle lifecycle-container-shutdown"
	ScalingDefaultName                = "lifecycle lifecycle-deployment-scaling"
	TerminationGracePeriodDefaultName = "lifecycle lifecycle-pod-termination-grace-period"
	PodOwnerTypeDefaultName           = "lifecycle lifecycle-pod-owner-type"
	PodRecreationDefaultName          = "lifecycle lifecycle-pod-recreation"
)

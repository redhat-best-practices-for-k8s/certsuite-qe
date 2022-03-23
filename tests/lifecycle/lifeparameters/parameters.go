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

	TnfShutdownTcName               = "lifecycle-container-shutdown"
	TnfDeploymentScalingTcName      = "lifecycle-deployment-scaling"
	TnfTerminationGracePeriodTcName = "lifecycle-pod-termination-grace-period"
	TnfPodOwnerTypeTcName           = "lifecycle-pod-owner-type"
	TnfPodRecreationTcName          = "lifecycle-pod-recreation"
	TnfPodHighAvailabilityTcName    = "lifecycle-pod-high-availability"
	TnfPodSchedulingTcName          = "lifecycle-pod-scheduling"
	TnfLivenessTcName               = "lifecycle-liveness"
	TnfReadinessTcName              = "lifecycle-readiness"
	TnfStatefulSetScalingTcName     = "lifecycle-statefulset-scaling"
	TnfImagePullPolicyTcName        = "lifecycle-image-pull-policy"

	TnfTestCases = []string{TnfPodRecreationTcName, TnfPodSchedulingTcName,
		TnfDeploymentScalingTcName, TnfTerminationGracePeriodTcName,
		TnfPodOwnerTypeTcName, TnfShutdownTcName, TnfImagePullPolicyTcName,
		TnfLivenessTcName, TnfReadinessTcName, TnfStatefulSetScalingTcName,
		TnfPodHighAvailabilityTcName}
)

package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime = 5 * time.Minute
)

var (
	testPodLabelPrefixName = "test-network-function.com/performance"
	testPodLabelValue      = "testing"
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestPodName            = "performance-pod"
	TnfTargetPodLabels     = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
)

const (
	TnfTestSuiteName     = "performance"
	PerformanceNamespace = "performance-ns"

	RtImageName = "quay.io/testnetworkfunction/debug-partner:latest"

	// TNF test cases names.
	TnfExclusiveCPUPool                  = "performance-exclusive-cpu-pool"
	TnfSharedCPUPoolSchedulingPolicy     = "performance-shared-cpu-pool-non-rt-scheduling-policy"
	TnfRtIsolatedCPUPoolSchedulingPolicy = "performance-isolated-cpu-pool-rt-scheduling-policy"
	TnfRtAppsNoExecProbes                = "performance-rt-apps-no-exec-probes"

	PriviledgedRoleName = "privileged-role"
	TnfRunTimeClass     = "performance-rtc"
)

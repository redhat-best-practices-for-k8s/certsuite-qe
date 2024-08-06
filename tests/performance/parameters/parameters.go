package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime = 5 * time.Minute
)

var (
	testPodLabelPrefixName   = "test-network-function.com/performance"
	testPodLabelValue        = "testing"
	TestPodLabel             = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestPodName              = "performance-pod"
	CertsuiteTargetPodLabels = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
)

const (
	CertsuiteTestSuiteName = "performance"
	PerformanceNamespace   = "performance-ns"

	RtImageName = "quay.io/testnetworkfunction/k8s-best-practices-debug:latest"

	// Certsuite test case names.
	CertsuiteExclusiveCPUPool                   = "performance-exclusive-cpu-pool"
	CertsuiteSharedCPUPoolSchedulingPolicy      = "performance-shared-cpu-pool-non-rt-scheduling-policy"
	CertsuiteRtIsolatedCPUPoolSchedulingPolicy  = "performance-isolated-cpu-pool-rt-scheduling-policy"
	CertsuiteRtExclusiveCPUPoolSchedulingPolicy = "performance-exclusive-cpu-pool-rt-scheduling-policy"
	CertsuiteRtAppsNoExecProbes                 = "performance-rt-apps-no-exec-probes"

	PrivilegedRoleName    = "privileged-role"
	CertsuiteRunTimeClass = "performance-rtc"

	DisableStr = "disable"
)

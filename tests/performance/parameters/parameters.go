package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime = 5 * time.Minute
)

var (
	testPodLabelPrefixName   = "redhat-best-practices-for-k8s.com/performance"
	testPodLabelValue        = "testing"
	TestPodLabel             = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestPodName              = "performance-pod"
	DpdkPodName              = "dpdk-pod"
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
	CertsuiteCPUPinningNoExecProbes             = "performance-cpu-pinning-no-exec-probes"

	PrivilegedRoleName    = "privileged-role"
	CertsuiteRunTimeClass = "performance-rtc"

	DisableStr = "disable"

	SampleWorkloadImage = "registry.access.redhat.com/ubi8/ubi-micro:latest"
)

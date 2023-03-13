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
	testPodLabelPrefixName = "test-network-function.com/platform-alteration"
	testPodLabelValue      = "testing"
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentName     = "platform-alteration-dpa"
	TestDaemonSetName      = "platform-alteration-dsa"
	TestStatefulSetName    = "platform-alteration-sfa"
	TestPodName            = "platform-alteration-pod"
	TnfTargetPodLabels     = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}

	NotRedHatRelease = "quay.io/jitesoft/nginx:1.23.3"
)

const (
	TnfTestSuiteName            = "platform-alteration"
	PlatformAlterationNamespace = "platform-alteration-ns"

	// TNF test cases names.
	TnfBaseImageName          = "platform-alteration-base-image"
	TnfIsSelinuxEnforcingName = "platform-alteration-is-selinux-enforcing"
	TnfIsRedHatReleaseName    = "platform-alteration-isredhat-release"
	TnfTaintedNodeKernelName  = "platform-alteration-tainted-node-kernel"
	TnfHugePagesConfigName    = "platform-alteration-hugepages-config"
	TnfBootParamsName         = "platform-alteration-boot-params"
	TnfSysctlConfigName       = "platform-alteration-sysctl-config"
	TnfHugePages2mOnlyName    = "platform-alteration-hugepages-2m-only"
	TnfOCPLifecycleName       = "platform-alteration-ocp-lifecycle"
	TnfOCPNodeOsName          = "platform-alteration-ocp-node-os-lifecycle"
	TnfServiceMeshUsageName   = "platform-alteration-service-mesh-usage"

	Getenforce    = `chroot /host getenforce`
	Enforcing     = "Enforcing"
	SetPermissive = `chroot /host setenforce 0`
	SetEnforce    = `chroot /host setenforce 1`

	RebootWaitingTime     = 10 * time.Minute
	Reboot                = `chroot /host systemctl reboot`
	FindHugePagesFiles    = "find /host/sys/devices/system/node/ -name nr_hugepages"
	PerformanceProfileCrd = "performanceprofiles.performance.openshift.io"
)

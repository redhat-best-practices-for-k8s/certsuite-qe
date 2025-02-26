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
	testPodLabelPrefixName   = "redhat-best-practices-for-k8s.com/platform-alteration"
	testPodLabelValue        = "testing"
	TestPodLabel             = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentName       = "platform-alteration-dpa"
	TestDaemonSetName        = "platform-alteration-dsa"
	TestStatefulSetName      = "platform-alteration-sfa"
	TestPodName              = "platform-alteration-pod"
	CertsuiteTargetPodLabels = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}

	NotRedHatRelease = "quay.io/jitesoft/nginx:stable"
)

const (
	CertsuiteTestSuiteName      = "platform-alteration"
	PlatformAlterationNamespace = "platform-alteration-ns"

	// Certsuite test case names.
	CertsuiteBaseImageName          = "platform-alteration-base-image"
	CertsuiteIsSelinuxEnforcingName = "platform-alteration-is-selinux-enforcing"
	CertsuiteIsRedHatReleaseName    = "platform-alteration-isredhat-release"
	CertsuiteTaintedNodeKernelName  = "platform-alteration-tainted-node-kernel"
	CertsuiteHugePagesConfigName    = "platform-alteration-hugepages-config"
	CertsuiteBootParamsName         = "platform-alteration-boot-params"
	CertsuiteSysctlConfigName       = "platform-alteration-sysctl-config"
	CertsuiteHugePages2mOnlyName    = "platform-alteration-hugepages-2m-only"
	CertsuiteHugePages1gOnlyName    = "platform-alteration-hugepages-1g-only"
	CertsuiteOCPLifecycleName       = "platform-alteration-ocp-lifecycle"
	CertsuiteOCPNodeOsName          = "platform-alteration-ocp-node-os-lifecycle"
	CertsuiteServiceMeshUsageName   = "platform-alteration-service-mesh-usage"
	CertsuiteClusterOperatorHealth  = "platform-alteration-cluster-operator-health"

	Getenforce    = `chroot /host getenforce`
	Enforcing     = "Enforcing"
	SetPermissive = `chroot /host setenforce 0`
	SetEnforce    = `chroot /host setenforce 1`

	RebootWaitingTime     = 10 * time.Minute
	Reboot                = `chroot /host systemctl reboot`
	FindHugePagesFiles    = "find /host/sys/devices/system/node/ -name nr_hugepages"
	PerformanceProfileCrd = "performanceprofiles.performance.openshift.io"
)

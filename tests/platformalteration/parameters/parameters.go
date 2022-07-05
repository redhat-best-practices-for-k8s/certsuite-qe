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
	}

	NotRedHatRelease = "ubuntu:20.04"
)

const (
	TnfTestSuiteName            = "platform-alteration"
	PlatformAlterationNamespace = "platform-alteration-ns"

	// TNF test cases names.
	TnfBaseImageName          = "platform-alteration-base-image"
	TnfIsSelinuxEnforcingName = "platform-alteration-is-selinux-enforcing"
	TnfIsRedHatReleaseName    = "platform-alteration-isredhat-release"
)

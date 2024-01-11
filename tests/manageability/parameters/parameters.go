package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime = 5 * time.Minute
)

var (
	testPodLabelPrefixName = "test-network-function.com/manageability"
	testPodLabelValue      = "testing"
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestPodName            = "manageability-pod"
	TnfTargetPodLabels     = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
	TestImageWithValidTag = "quay.io/bapalm/httpd:2.4.58"
	InvalidPortName       = "sftp"
)

const (
	TnfTestSuiteName       = "manageability"
	ManageabilityNamespace = "manageability-ns"

	// TNF test cases names.
	TnfContainerPortName = "manageability-container-port-name-format"
	TnfContainerImageTag = "manageability-containers-image-tag"
)

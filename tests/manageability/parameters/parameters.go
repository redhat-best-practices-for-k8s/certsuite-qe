package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime = 5 * time.Minute
)

var (
	testPodLabelPrefixName   = "redhat-best-practices-for-k8s.com/manageability"
	testPodLabelValue        = "testing"
	TestPodLabel             = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestPodName              = "manageability-pod"
	CertsuiteTargetPodLabels = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
	TestImageWithValidTag = "quay.io/bapalm/httpd:2.4.58"
	InvalidPortName       = "sftp"
)

const (
	CertsuiteTestSuiteName = "manageability"
	ManageabilityNamespace = "manageability-ns"

	// Certsuite test case names.
	CertsuiteContainerPortName = "manageability-container-port-name-format"
	CertsuiteContainerImageTag = "manageability-containers-image-tag"
)

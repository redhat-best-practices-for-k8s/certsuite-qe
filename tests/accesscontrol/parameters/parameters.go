package parameters

import (
	"fmt"
	"time"
)

const (
	Timeout = 5 * time.Minute
)

var (
	TestAccessControlNameSpace = "accesscontrol-tests"
	testPodLabelPrefixName     = "accesscontrol-test/test"
	testPodLabelValue          = "testing"
	TestPodLabel               = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	InvalidNamespace           = "openshift-test"
	AdditionalValidNamespace   = "ac-test"

	TestCaseNameAccessControlNamespace         = "access-control-namespace"
	TestCaseNameAccessControlPodHostIpc        = "access-control-pod-host-ipc"
	TestCaseNameAccessControlPodAutomountToken = "access-control-pod-automount-service-account-token"

	TestDeploymentLabels = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
)

package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime = 5 * time.Minute
	Timeout     = 5 * time.Minute
)

var (
	testPodLabelPrefixName = "test-network-function.com/operator"
	testPodLabelValue      = "testing"
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TnfTargetPodLabels     = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
)

const (
	OperatorNamespace = "operator-ns"

	// TNF test cases names.
	TnfOperatorInstallSource = "operator-install-source"
)

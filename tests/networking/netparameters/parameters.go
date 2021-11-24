package netparameters

import (
	"fmt"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"time"
)

const (
	WaitingTime   = 5 * time.Minute
	RetryInterval = 5
)

var (
	testPodLabelPrefixName = "networking-test/test"
	testPodLabelValue      = "testing"
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentLabels   = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "networkingput"}
	TestCaseDefaultNetworkName  = "networking Both Pods are on the Default network Testing network connectivity networking-icmpv4-connectivity"
	TestCaseDefaultSkipRegEx    = "nodePort|Multus"
	TestCaseNodePortNetworkName = "networking Should not have type of nodePort networking-service-type"
	TestCaseNodePortSkipRegEx   = "Default|Multus"
	NetworkingTestSkipLabel     = map[string]string{"test-network-function.com/skip_connectivity_tests": ""}
	NetworkingTestSuiteName     = "networking"
	DefaultPartnerPodNamespace  = "default"
	DefaultPartnerPodPrefixName = "tnfpartner-"
	TestNamespace               = namespaces.DefineNamespace("networking-tests")
)

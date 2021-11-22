package netparameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime   = 5 * time.Minute
	RetryInterval = 5
)

var (
	TestNetworkingNameSpace = "networking-tests"
	testPodLabelPrefixName  = "networking-test/test"
	testPodLabelValue       = "testing"
	TestPodLabel            = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentLabels    = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "networkingput"}
	TestCaseDefaultNetworkName  = "networking Both Pods are on the Default network Testing network connectivity networking-icmpv4-connectivity"
	TestCaseDefaultSkipRegEx    = "nodePort|Multus"
	NetworkingTestSuiteName     = "networking"
	DefaultPartnerPodNamespace  = "default"
	DefaultPartnerPodPrefixName = "tnfpartner-"
)

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
	TestNadNameA               = "networking-nada"
	TestIPamIPNetworkA         = "10.255.255.0/25"
	TestDeploymentAName        = "networkingputa"
	TestNadNameB               = "networking-nadb"
	TestIPamIPNetworkB         = "10.255.128.0/25"
	TestDeploymentBName        = "networkingputb"
	TestCaseDefaultNetworkName = "networking Both Pods are on the Default network Testing Default network connectivity " +
		"networking-icmpv4-connectivity"
	TestCaseMultusConnectivityName = "networking Both Pods are connected via a Multus Overlay Network Testing Multus " +
		"network connectivity networking-icmpv4-connectivity"
	TestCaseDefaultSkipRegEx    = "nodePort|Multus"
	TestCaseNodePortNetworkName = "networking Should not have type of nodePort networking-service-type"
	TestCaseNodePortSkipRegEx   = "Default|Multus"
	TestCaseMultusSkipRegEx     = "nodePort|Default"
	NetworkingTestSkipLabel     = map[string]string{"test-network-function.com/skip_connectivity_tests": ""}
	NetworkingTestSuiteName     = "networking"
)

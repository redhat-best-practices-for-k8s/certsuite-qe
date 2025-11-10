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
	TestNetworkingNameSpace       = "networking-tests"
	AdditionalNetworkingNamespace = "net-tests"
	testPodLabelPrefixName        = "networking-test/test"
	testPodLabelValue             = "testing"
	TestPodLabel                  = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentLabels          = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "networkingput"}
	TestNadNameA                                 = "networking-nada"
	TestIPamIPNetworkA                           = "10.255.255.0/25"
	TestDeploymentAName                          = "networkingputa"
	TestNadNameB                                 = "networking-nadb"
	TestIPamIPNetworkB                           = "10.255.128.0/25"
	TestDeploymentBName                          = "networkingputb"
	CertsuiteDefaultNetworkTcName                = "networking-icmpv4-connectivity"
	CertsuiteMultusIpv4TcName                    = "networking-icmpv4-connectivity-multus"
	CertsuiteNodePortTcName                      = "access-control-service-type"
	CertsuiteNetworkPolicyDenyAllTcName          = "networking-network-policy-deny-all"
	CertsuiteOcpReservedPortsUsageTcName         = "networking-ocp-reserved-ports-usage"
	CertsuiteUndeclaredContainerPortsUsageTcName = "networking-undeclared-container-ports-usage"
	CertsuiteReservedPartnerPortsTcName          = "networking-reserved-partner-ports"
	CertsuiteDualStackServiceTcName              = "networking-dual-stack-service"
	NetworkingTestSkipLabel                      = map[string]string{"redhat-best-practices-for-k8s.com/skip_connectivity_tests": ""}
	NetworkingTestMultusSkipLabel                = map[string]string{"redhat-best-practices-for-k8s.com/skip_multus_connectivity_tests": ""}
)

type IPOutputInterface struct {
	IfIndex   uint     `json:"ifindex"`
	IfName    string   `json:"ifname"`
	Flags     []string `json:"flags"`
	Mtu       uint     `json:"mtu"`
	Qdisc     string   `json:"qdisc"`
	Master    string   `json:"master"`
	Operstate string   `json:"operstate"`
	Linkmode  string   `json:"linkmode"`
	Group     string   `json:"group"`
	Txqlen    int      `json:"txqlen"`
	LinkType  string   `json:"link_type"`
	Address   string   `json:"address"`
	Broadcast string   `json:"broadcast"`
}

const (
	CertsuiteTestSuiteName = "networking"
	NetworkingNamespace    = "networking-ns"

	// Certsuite test case names.
	CertsuiteNetworkingIcmpv4TcName               = "networking-icmpv4-connectivity"
	CertsuiteNetworkingIcmpv6TcName               = "networking-icmpv6-connectivity"
	CertsuiteNetworkingOcpReservedPortsTcName     = "networking-ocp-reserved-ports"
	CertsuiteNetworkingDefaultNetworkTcName       = "networking-network-policy-deny-all"
	CertsuiteNetworkingUnderTestContainersTcName  = "networking-undeclared-container-ports"
	CertsuiteNetworkingUnderTestPodsTcName        = "networking-dual-stack-service"
	CertsuiteNetworkingReservedPartnerPortsTcName = "networking-reserved-partner-ports"
	CertsuiteNetworkingPtpDaemonTcName            = "networking-ptp-daemon"
	CertsuiteNetworkingRestartUnderTestTcName     = "networking-restart-on-reboot"
	CertsuiteNetworkingDpdkCPUPinningTcName       = "networking-dpdk-cpu-pinning"
	CertsuiteNetworkingMultipleIPTcName           = "networking-multiple-ip-families"
	CertsuiteNetworkingMultusBridgeTcName         = "networking-multus-bridge"
	CertsuiteNetworkingMultusIpamTcName           = "networking-multus-ipam"
	CertsuiteNetworkingMultusNodeSelectorTcName   = "networking-multus-node-selector"

	SampleWorkloadImage = "registry.access.redhat.com/ubi8/ubi-micro:latest"
)

package parameters

import (
	"fmt"
	"time"
)

type (
	OperatorLabelInfo struct {
		OperatorPrefix string
		Namespace      string
		Label          map[string]string
	}

	CsvInfo struct {
		OperatorPrefix string
		Namespace      string
	}
)

const (
	WaitingTime     = 5 * time.Minute
	Timeout         = 5 * time.Minute
	TimeoutLabelCsv = 2 * time.Minute
	PollingInterval = 5 * time.Second
)

var (
	testPodLabelPrefixName = "test-network-function.com/operator"
	testPodLabelValue      = "testing"
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TnfTargetPodLabels     = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
	TnfTargetOperatorLabels  = fmt.Sprintf("%s: %s", "test-network-function.com/operator", "target")
	OperatorGroupName        = "operator-test-operator-group"
	OperatorLabel            = map[string]string{"test-network-function.com/operator": "target"}
	CertifiedOperatorGroup   = "certified-operators"
	CommunityOperatorGroup   = "community-operators"
	OperatorSourceNamespace  = "openshift-marketplace"
	OperatorPrefixCloudbees  = "cloudbees-ci"
	OperatorPrefixAnchore    = "anchore-engine"
	OperatorPrefixQuay       = "quay-operator"
	OperatorPrefixKiali      = "kiali-operator"
	OperatorPrefixOpenvino   = "openvino-operator"
	SubscriptionNameOpenvino = "ovms-operator-subscription"
)

const (
	OperatorNamespace = "operator-ns"

	// TNF test cases names.
	TnfOperatorInstallSource             = "operator-install-source"
	TnfOperatorInstallStatusNoPrivileges = "operator-install-status-no-privileges"
	TnfOperatorInstallStatus             = "operator-install-status-succeeded"
	TnfOperatorSemanticVersioning        = "operator-semantic-versioning"
	TnfOperatorCrdVersioning             = "operator-crd-versioning"
	TnfOperatorCrdOpenAPISchema          = "operator-crd-openapi-schema"
)

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
	testPodLabelPrefixName   = "redhat-best-practices-for-k8s.com/operator"
	testPodLabelValue        = "testing"
	TestPodLabel             = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	CertsuiteTargetPodLabels = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
	CertsuiteTargetOperatorLabels = fmt.Sprintf("%s: %s", "redhat-best-practices-for-k8s.com/operator", "target")
	CertsuiteTargetCrdFilters     = []string{"charts.operatorhub.io"}
	OperatorGroupName             = "operator-test-operator-group"
	OperatorLabel                 = map[string]string{"redhat-best-practices-for-k8s.com/operator": "target"}
	CertifiedOperatorGroup        = "certified-operators"
	RedhatOperatorGroup           = "redhat-operators"
	CommunityOperatorGroup        = "community-operators"
	OperatorSourceNamespace       = "openshift-marketplace"
	OperatorPrefixCloudbees       = "cloudbees-ci"
	OperatorPrefixAnchore         = "anchore-engine"
	OperatorPrefixQuay            = "quay-operator"
	OperatorPrefixKiali           = "kiali-operator"
	OperatorPrefixOpenvino        = "openvino-operator"
	CertifiedOperatorPrefixNginx  = "nginx-ingress-operator"
	SubscriptionNameOpenvino      = "ovms-operator-subscription"
)

const (
	OperatorNamespace = "operator-ns"

	// Certsuite test case names.
	CertsuiteOperatorInstallSource             = "operator-install-source"
	CertsuiteOperatorInstallStatusNoPrivileges = "operator-install-status-no-privileges"
	CertsuiteOperatorInstallStatus             = "operator-install-status-succeeded"
	CertsuiteOperatorSemanticVersioning        = "operator-semantic-versioning"
	CertsuiteOperatorCrdVersioning             = "operator-crd-versioning"
	CertsuiteOperatorCrdOpenAPISchema          = "operator-crd-openapi-schema"
	CertsuiteOperatorNonRoot                   = "operator-run-as-non-root"
	CertsuiteOperatorReadOnlyFilesystem        = "operator-read-only-file-system"
	CertsuiteOperatorPodAutomountToken         = "operator-automount-tokens"
	CertsuiteOperatorPodRunAsUserID            = "operator-run-as-user-id"
)

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
	CertsuiteTargetOperatorLabels             = fmt.Sprintf("%s: %s", "redhat-best-practices-for-k8s.com/operator", "target")
	CertsuiteTargetCrdFilters                 = []string{"grafanadashboards.grafana.integreatly.org"}
	OperatorGroupName                         = "operator-test-operator-group"
	OperatorLabel                             = map[string]string{"redhat-best-practices-for-k8s.com/operator": "target"}
	CertifiedOperatorGroup                    = "certified-operators"
	RedhatOperatorGroup                       = "redhat-operators"
	CommunityOperatorGroup                    = "community-operators"
	OperatorSourceNamespace                   = "openshift-marketplace"
	OperatorPrefixLightweight                 = "postgresoperator"
	OperatorPrefixKiali                       = "kiali-operator"
	CertifiedOperatorPrefix                   = "grafana-operator"
	UncertifiedOperatorPrefixCockroach        = "cockroachdb"
	CertifiedOperatorPrefixCockroachCertified = "cockroach-operator"
	SingleOrMultiNamespacedOperatorGroup      = "single-or-multi-og"

	TestDeploymentLabels = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}

	TestPodName = "operator-test-pod"
)

const (
	OperatorNamespace = "operator-ns"

	// Certsuite test case names.
	CertsuiteOperatorInstallSource                                    = "operator-install-source"
	CertsuiteOperatorInstallStatusNoPrivileges                        = "operator-install-status-no-privileges"
	CertsuiteOperatorInstallStatus                                    = "operator-install-status-succeeded"
	CertsuiteOperatorSemanticVersioning                               = "operator-semantic-versioning"
	CertsuiteOperatorCrdVersioning                                    = "operator-crd-versioning"
	CertsuiteOperatorCrdOpenAPISchema                                 = "operator-crd-openapi-schema"
	CertsuiteOperatorSingleOrMultiNamespacedAllowedInTenantNamespaces = "operator-single-or-multi-namespaced-allowed-in-tenant-namespaces"
	CertsuiteOperatorNonRoot                                          = "operator-run-as-non-root"
	CertsuiteOperatorPodAutomountToken                                = "operator-automount-tokens"
	CertsuiteOperatorPodRunAsUserID                                   = "operator-run-as-user-id"
	CertsuiteOperatorMultipleInstalled                                = "operator-multiple-same-operators"
	CertsuiteOperatorBundleCount                                      = "operator-catalogsource-bundle-count"

	TestOperatorInstallStatusSucceeded      = "operator-install-status-succeeded"
	TestOperatorNoPrivileges                = "operator-no-privileges"
	TestOperatorIsNotUsingDeprecatedAPIs    = "operator-is-not-using-deprecated-apis"
	TestOperatorHubSubscriptionValidChannel = "operator-hub-subscription-is-valid-channel"
	TestOperatorSingleCrdOwner              = "operator-single-crd-owner"

	SampleWorkloadImage = "registry.access.redhat.com/ubi8/ubi-micro:latest"
)

const (
	OperatorPackageNamePrefixLightweight              = "postgresql"
	OperatorPackageNamePrefixLightweightCustomCatalog = "nginx-ingress-operator"
)

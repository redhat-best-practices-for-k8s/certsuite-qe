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
	Timeout         = 5 * time.Minute
	TimeoutLabelCsv = 2 * time.Minute
	PollingInterval = 5 * time.Second

	SampleWorkloadImage = "registry.access.redhat.com/ubi8/ubi-micro:latest"
)

var (
	TestCertificationNameSpace = "affiliatedcert-tests"
	testPodLabelPrefixName     = "affiliatedcert-test/test"
	testPodLabelValue          = "testing"
	TestPodLabel               = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentLabels       = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}

	UnrelatedOperatorPrefixCloudcasa = "cloudcasa"
	UnrelatedNamespace               = "certsuite"

	TestCaseContainerAffiliatedCertName = "affiliated-certification-container-is-certified"
	TestCaseNameContainerDigest         = "affiliated-certification-container-is-certified-digest"
	EmptyFieldsContainer                = ";;;"
	ContainerNameOnlyCockroachDB        = "cockroachdb/cockroach;;;"
	ContainerRepoOnlyRedHatRegistry     = ";registry.connect.redhat.com;;"
	CertifiedContainerURLNodeJs         = "registry.access.redhat.com/ubi8/nodejs-12:latest"
	CertifiedContainerURLCockroachDB    = "registry.connect.redhat.com/cockroachdb/cockroach:v23.1.17" // 'latest' tag is not available
	UncertifiedContainerURLCnfTest      = "quay.io/testnetworkfunction/k8s-best-practices-debug:latest"

	TestCaseOperatorAffiliatedCertName        = "affiliated-certification-operator-is-certified"
	TestHelmChartCertified                    = "affiliated-certification-helmchart-is-certified"
	TestHelmVersion                           = "affiliated-certification-helm-version"
	OperatorGroupName                         = "affiliatedcert-test-operator-group"
	CertifiedOperatorGroup                    = "certified-operators"
	CertifiedOperatorDisplayName              = "Certified Operators"
	CommunityOperatorGroup                    = "community-operators"
	OperatorSourceNamespace                   = "openshift-marketplace"
	OperatorLabel                             = map[string]string{"redhat-best-practices-for-k8s.com/operator": "target"}
	UncertifiedOperatorPrefixCockroach        = "cockroachdb"
	CertifiedOperatorPrefixCockroachCertified = "cockroach-operator"
	CertifiedOperatorPrefix                   = "grafana-operator"
	UncertifiedOperatorPrefixSriov            = "sriov-fec"
	UncertifiedOperatorFullSriov              = "sriov-fec.v1.2.1"
)

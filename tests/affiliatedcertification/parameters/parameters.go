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
	UnrelatedNamespace               = "tnf"

	TestCaseContainerAffiliatedCertName = "affiliated-certification-container-is-certified"
	TestCaseNameContainerDigest         = "affiliated-certification-container-is-certified-digest"
	CertifiedContainerCockroachDB       = "cockroachdb/cockroach;registry.connect.redhat.com;v20.1.8;" +
		"sha256:6667919a41d304d5d4ade3ded4f11b42d722a995a4283e11e15320529f7f9abf"
	CertifiedContainer5gc            = "5gc/vru-sig-mp;artnj.zte.com.cn;2021010510011609815594;"
	UncertifiedContainerNodeJs12     = "nodejs-12/ubi8;registry.connect.redhat.com;latest;"
	EmptyFieldsContainer             = ";;;"
	ContainerNameOnlyCockroachDB     = "cockroachdb/cockroach;;;"
	ContainerRepoOnlyRedHatRegistry  = ";registry.connect.redhat.com;;"
	CertifiedContainerURLNodeJs      = "registry.access.redhat.com/ubi8/nodejs-12:latest"
	CertifiedContainerURLCockroachDB = "registry.connect.redhat.com/cockroachdb/cockroach:latest"
	UncertifiedContainerURLCnfTest   = "quay.io/testnetworkfunction/cnf-test-partner:latest"

	TestCaseOperatorAffiliatedCertName = "affiliated-certification-operator-is-certified"
	TestHelmChartCertified             = "affiliated-certification-helmchart-is-certified"
	OperatorGroupName                  = "affiliatedcert-test-operator-group"
	CertifiedOperatorGroup             = "certified-operators"
	CertifiedOperatorDisplayName       = "Certified Operators"
	CommunityOperatorGroup             = "community-operators"
	OperatorSourceNamespace            = "openshift-marketplace"
	OperatorLabel                      = map[string]string{"test-network-function.com/operator": "target"}
	UncertifiedOperatorPrefixFalcon    = "falcon-operator"
	CertifiedOperatorPrefixFederatorai = "federatorai"
	CertifiedOperatorFullFederatorai   = "federatorai.v4.7.1-1"
	CertifiedOperatorPrefixInstana     = "instana-agent-operator"
	CertifiedOperatorFullInstana       = "instana-agent-operator.v2.0.4"
	UncertifiedOperatorPrefixSriov     = "sriov-fec"
	UncertifiedOperatorFullSriov       = "sriov-fec.v1.2.1"
)

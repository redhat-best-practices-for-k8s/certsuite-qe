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

	UnrelatedOperatorPrefixCloudcasa = "cloudcasa"
	UnrelatedNamespace               = "tnf"

	TestCaseContainerAffiliatedCertName = "affiliated-certification-container-is-certified"
	CertifiedContainerCockroachDB       = "cockroachdb/cockroach;registry.connect.redhat.com;v20.1.8;" +
		"sha256:6667919a41d304d5d4ade3ded4f11b42d722a995a4283e11e15320529f7f9abf"
	CertifiedContainer5gc           = "5gc/vru-sig-mp;artnj.zte.com.cn;2021010510011609815594;"
	UncertifiedContainerNodeJs12    = "nodejs-12/ubi8;registry.connect.redhat.com;latest;"
	EmptyFieldsContainer            = ";;;"
	ContainerNameOnlyCockroachDB    = "cockroachdb/cockroach;;;"
	ContainerRepoOnlyRedHatRegistry = ";registry.connect.redhat.com;;"

	TestCaseOperatorAffiliatedCertName = "affiliated-certification-operator-is-certified"
	OperatorGroupName                  = "affiliatedcert-test-operator-group"
	CertifiedOperatorGroup             = "certified-operators"
	CertifiedOperatorDisplayName       = "Certified Operators"
	CommunityOperatorGroup             = "community-operators"
	OperatorSourceNamespace            = "openshift-marketplace"
	OperatorLabel                      = map[string]string{"test-network-function.com/operator": "target"}
	UncertifiedOperatorPrefixFalcon    = "falcon-operator"
	CertifiedOperatorPrefixInfinibox   = "infinibox-operator"
	CertifiedOperatorFullInfinibox     = "infinibox-operator.v2.4.0"
	CertifiedOperatorPrefixInstana     = "instana-agent-operator"
	CertifiedOperatorFullInstana       = "instana-agent-operator.v2.0.4"
	UncertifiedOperatorPrefixSriov     = "sriov-fec"
	UncertifiedOperatorFullSriov       = "sriov-fec.v1.2.1"
)

package affiliatedcertparameters

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
	Timeout         = 20 * time.Minute
	TimeoutLabelCsv = 2 * time.Minute
	PollingInterval = 5 * time.Second
)

var (
	TestCertificationNameSpace = "affiliatedcert-tests"
	testPodLabelPrefixName     = "affiliatedcert-test/test"
	testPodLabelValue          = "testing"
	TestPodLabel               = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)

	TestCaseContainerAffiliatedCertName = "affiliated-certification-container-is-certified"
	CertifiedContainerCockroachDB       = "cockroachdb/cockroach;registry.connect.redhat.com;latest;" +
		"sha256:5d7061441ac5f16153301346fd3cfd645d976983edd2457a32c4f914407b9957"
	CertifiedContainer5gc           = "5gc/vru-sig-mp;artnj.zte.com.cn;2021010510011609815594;"
	UncertifiedContainerNodeJs12    = "nodejs-12/ubi8;registry.connect.redhat.com;latest;"
	EmptyFieldsContainer            = ";;;"
	ContainerNameOnlyCockroachDB    = "cockroachdb/cockroach;;;"
	ContainerRepoOnlyRedHatRegistry = ";registry.connect.redhat.com;;"

	TestCaseOperatorAffiliatedCertName   = "affiliated-certification-operator-is-certified"
	OperatorGroupName                    = "affiliatedcert-test-operator-group"
	CertifiedOperatorGroup               = "certified-operators"
	CertifiedOperatorDisplayName         = "Certified Operators"
	CommunityOperatorGroup               = "community-operators"
	OperatorSourceNamespace              = "openshift-marketplace"
	OperatorLabel                        = map[string]string{"test-network-function.com/operator": "target"}
	UncertifiedOperatorPrefixFalcon      = "falcon-operator"
	CertifiedOperatorPrefixInfinibox     = "infinibox-operator"
	CertifiedOperatorFullInfinibox       = "infinibox-operator.v2.1.2"
	CertifiedOperatorPrefixArtifactoryHa = "artifactory-ha-operator"
	CertifiedOperatorFullArtifactoryHa   = "artifactory-ha-operator.v1.1.20"
	UncertifiedOperatorPrefixSriov       = "sriov-fec"
	UncertifiedOperatorFullSriov         = "sriov-fec.v1.1.0"
)

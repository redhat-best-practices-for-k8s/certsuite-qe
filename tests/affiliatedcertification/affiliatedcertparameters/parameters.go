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
	Timeout         = 10 * time.Minute
	PollingInterval = 5 * time.Second
)

var (
	AffiliatedCertificationTestSuiteName = "affiliated-certification"

	TestCertificationNameSpace = "affiliatedcert-tests"
	testPodLabelPrefixName     = "affiliatedcert-test/test"
	testPodLabelValue          = "testing"
	TestPodLabel               = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)

	TestCaseContainerSkipRegEx          = "operator-is-certified helmchart-is-certified"
	TestCaseContainerAffiliatedCertName = "affiliated-certification affiliated-certification-container-is-certified"
	CertifiedContainerNodeJsUbi         = "nodejs-12/ubi8"
	CertifiedContainerRhel7OpenJdk      = "openjdk-11-rhel7/openjdk"
	UncertifiedContainerFooBar          = "foo/bar"
	EmptyFieldsContainerOrOperator      = "/"
	ContainerNameOnlyRhel7OpenJdk       = "openjdk-11-rhel7/"
	ContainerRepoOnlyOpenJdk            = "/openjdk"

	TestCaseOperatorSkipRegEx           = "container-is-certified helmchart-is-certified"
	TestCaseOperatorAffiliatedCertName  = "affiliated-certification affiliated-certification-operator-is-certified"
	OperatorGroupName                   = "affiliatedcert-test-operator-group"
	CertifiedOperatorGroup              = "certified-operators"
	CommunityOperatorGroup              = "community-operators"
	OperatorSourceNamespace             = "openshift-marketplace"
	OperatorLabel                       = map[string]string{"test-network-function.com/operator": "target"}
	UncertifiedOperatorPrefixFalcon     = "falcon-operator"
	UncertifiedOperatorDeploymentFalcon = "falcon-operator-controller-manager"
	CertifiedOperatorPrefixPostgres     = "postgresoperator"
	CertifiedOperatorDeploymentPostgres = "pgo"
	CertifiedOperatorPrefixDatadog      = "datadog-operator"
	CertifiedOperatorDeploymentDatadog  = "datadog-operator-manager"
	UncertifiedOperatorBarFoo           = "bar/foo"
)

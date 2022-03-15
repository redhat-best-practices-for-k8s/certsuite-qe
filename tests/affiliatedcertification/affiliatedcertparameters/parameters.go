package affiliatedcertparameters

import (
	"fmt"
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

	TestCaseOperatorSkipRegEx          = "container-is-certified helmchart-is-certified"
	TestCaseOperatorAffiliatedCertName = "affiliated-certification affiliated-certification-operator-is-certified"
	CertifiedOperatorGroup             = "certified-operators"
	OperatorSourceNamespace            = "openshift-marketplace"
	OperatorLabel                      = map[string]string{"test-network-function.com/operator": "target"}
	UncertifiedOperatorPrefixNginx     = "nginx-operator"
	ExistingOperatorNamespace          = "tnf"
	CertifiedOperatorPrefixPostgres    = "postgresoperator"
	CertifiedOperatorPrefixDatadog     = "datadog-operator"
	CertifiedOperatorKubeturbo         = "kubeturbo-certified/certified-operators"
	UncertifiedOperatorBarFoo          = "bar/foo"
	OperatorNameOnlyKubeturbo          = "kubeturbo-certified"
	OperatorOrgOnlyCertifiedOperators  = "certified-operators"
)

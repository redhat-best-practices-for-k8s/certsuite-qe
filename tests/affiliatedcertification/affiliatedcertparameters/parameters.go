package affiliatedcertparameters

import (
	"fmt"

	utils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"
)

var (
	AffiliatedCertificationTestSuiteName = "affiliated-certification"

	TestCertificationNameSpace = "affiliatedcert-tests"
	testPodLabelPrefixName     = "affiliatedcert-test/test"
	testPodLabelValue          = "testing"
	TestPodLabel               = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	certifiedOperatorGroup     = "certified-operators"
	operatorSourceNamespace    = "openshift-marketplace"
	OperatorLabel              = map[string]string{"test-network-function.com/operator": "target"}

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
	OperatorGroup                      = utils.DefineOperatorGroup("affiliatedcert-test-operator-group",
		TestCertificationNameSpace, []string{TestCertificationNameSpace})
	CertifiedOperatorPostgresSubscription = utils.DefineSubscription("crunchy-postgres-operator-subscription",
		TestCertificationNameSpace, "v5", "crunchy-postgres-operator", certifiedOperatorGroup, operatorSourceNamespace)
	CertifiedOperatorApicast          = "apicast-operator/redhat-operators"
	CertifiedOperatorKubeturbo        = "kubeturbo-certified/certified-operators"
	UncertifiedOperatorBarFoo         = "bar/foo"
	OperatorNameOnlyKubeturbo         = "kubeturbo-certified"
	OperatorOrgOnlyCertifiedOperators = "certified-operators"
)

package affiliatedcertparameters

import "fmt"

var (
	AffiliatedCertificationTestSuiteName = "affiliated-certification"

	TestCertificationNameSpace = "cert-tests"
	testPodLabelPrefixName     = "cert-test/test"
	testPodLabelValue          = "testing"
	TestPodLabel               = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)

	TestCaseContainerSkipRegEx          = "operator-is-certified"
	TestCaseContainerAffiliatedCertName = "affiliated-certification affiliated-certification-container-is-certified"
	CertifiedContainerNodeJsUbi         = "nodejs-12/ubi8"
	CertifiedContainerRhel7OpenJdk      = "openjdk-11-rhel7/openjdk"
	UncertifiedContainerFooBar          = "foo/bar"
	EmptyFieldsContainerOrOperator      = "/"
	ContainerNameOnlyRhel7OpenJdk       = "openjdk-11-rhel7/"
	ContainerRepoOnlyOpenJdk            = "/openjdk"

	TestCaseOperatorSkipRegEx          = "container-is-certified"
	TestCaseOperatorAffiliatedCertName = "affiliated-certification affiliated-certification-operator-is-certified"
	CertifiedOperatorApicast           = "apicast-operator/redhat-operators"
	CertifiedOperatorKubeturbo         = "kubeturbo-certified/certified-operators"
	UncertifiedOperatorBarFoo          = "bar/foo"
	OperatorNameOnlyKubeturbo          = "kubeturbo-certified"
	OperatorOrgOnlyCertifiedOperators  = "certified-operators"
)

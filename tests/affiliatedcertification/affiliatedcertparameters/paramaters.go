package affiliatedcertparameters

var (
	AffiliatedCertificationTestSuiteName = "affiliated-certification"
	TestCaseContainerSkipRegEx           = "operator-is-certified"
	TestCaseOperatorSkipRegEx            = "container-is-certified"
	TestCaseContainerAffiliatedCertName  = "affiliated-certification affiliated-certification-container-is-certified"
	CertifiedContainer1                  = "nodejs-12/ubi8"
	CertifiedContainer2                  = "openjdk-11-rhel7/openjdk"
	UncertifiedContainer1                = "foo/bar"
)

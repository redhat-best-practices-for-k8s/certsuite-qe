package affiliatedcertparameters

var (
	AffiliatedCertificationTestSuiteName = "affiliated-certification"
	TestCaseContainerSkipRegEx           = "operator-is-certified"
	TestCaseOperatorSkipRegEx            = "container-is-certified"
	TestCaseContainerAffiliatedCertName  = "affiliated-certification affiliated-certification-container-is-certified"
	CertifiedContainerNodeJsUbi          = "nodejs-12/ubi8"
	CertifiedContainerRhel7OpenJdk       = "openjdk-11-rhel7/openjdk"
	UncertifiedContainerFooBar           = "foo/bar"
	ContainerFieldNoSubFields            = ""
	EmptyFieldsContainer                 = "/"
	ContainerNameOnlyRhel7OpenJdk        = "openjdk-11-rhel7/"
	ContainerRepoOnlyOpenJdk             = "/openjdk"
)

package globalparameters

import "encoding/xml"

type (
	TnfConfig struct {
		TargetNameSpaces         []TargetNameSpace            `yaml:"targetNameSpaces" json:"targetNameSpaces"`
		TargetPodLabels          []PodLabel                   `yaml:"targetPodLabels" json:"targetPodLabels"`
		OperatorsUnderTestLabels []string                     `yaml:"operatorsUnderTestLabels" json:"operatorsUnderTestLabels"`
		Certifiedcontainerinfo   []CertifiedContainerRepoInfo `yaml:"certifiedcontainerinfo" json:"certifiedcontainerinfo"`
		TargetCrdFilters         []TargetCrdFilter            `yaml:"targetCrdFilters" json:"targetCrdFilters"`
	}

	TargetCrdFilter struct {
		NameSuffix string `yaml:"nameSuffix" json:"nameSuffix"`
	}

	TargetNameSpace struct {
		Name string `yaml:"name" json:"name"`
	}

	PodLabel struct {
		Prefix string `yaml:"prefix" json:"prefix"`
		Name   string `yaml:"name" json:"name"`
		Value  string `yaml:"value" json:"value"`
	}

	Label struct {
		Name  string `yaml:"name" json:"name"`
		Value string `yaml:"value" json:"value"`
	}

	CertifiedContainerRepoInfo struct {
		Repository string `yaml:"repository" json:"repository"`
		Registry   string `yaml:"registry" json:"registry"`
		Tag        string `yaml:"tag" json:"tag"`
		Digest     string `yaml:"digest" json:"digest"`
	}

	CertifiedOperatorRepoInfo struct {
		Name         string `yaml:"name" json:"name"`
		Organization string `yaml:"organization" json:"organization"`
	}

	JUnitTestSuites struct {
		XMLName xml.Name    `xml:"testsuites"`
		Suites  []TestSuite `xml:"testsuite"`
	}

	// TestSuite is created from XML output by a Ant JUnit task.
	TestSuite struct {
		XMLName    xml.Name   `xml:"testsuite"`
		Errors     int        `xml:"errors,attr"`
		Failures   int        `xml:"failures,attr"`
		Name       string     `xml:"name,attr"`
		Tests      int        `xml:"tests,attr"`
		Time       float64    `xml:"time,attr"`
		Properties []Property `xml:"properties>property,omitempty"`
		Testcases  []TestCase `xml:"testcase"`
	}

	// TestCase represents a failed testcase.
	TestCase struct {
		XMLName   xml.Name `xml:"testcase"`
		ClassName string   `xml:"classname,attr"`
		Name      string   `xml:"name,attr"`
		Time      float64  `xml:"time,attr"`
		Status    string   `xml:"status,attr"`
		Fail      *Failure `xml:"failure"`
		Skipped   *string  `xml:"skipped"`
	}

	// Failure gives details as to why a test case failed.
	Failure struct {
		Message string `xml:"message,attr"`
		Type    string `xml:"type,attr"`
		Value   string `xml:",innerxml"`
	}
	// Property represents a key/value pair used to define properties.
	Property struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	}
)

var (
	DefaultClaimFileName             = "claim.json"
	DefaultTnfConfigFileName         = "tnf_config.yml"
	DefaultJunitReportName           = "cnf-certification-tests_junit.xml"
	PartnerNamespaceEnvVarName       = "TNF_PARTNER_NAMESPACE"
	TestCasePassed                   = "passed"
	TestCaseFailed                   = "failed"
	TestCaseSkipped                  = "skipped"
	NetworkSuiteName                 = "networking"
	AffiliatedCertificationSuiteName = "affiliated-certification"
	LifecycleSuiteName               = "lifecycle"
	PlatformAlterationSuiteName      = "platform-alteration"
	ObservabilitySuiteName           = "observability"
	AccessControlSuiteName           = "access-control"
	PerformanceSuiteName             = "performance"
	ManageabilitySuiteName           = "manageability"
	OperatorSuiteName                = "operator"
)

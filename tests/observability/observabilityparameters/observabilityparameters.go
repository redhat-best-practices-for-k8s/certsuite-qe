package observabilityparameters

import "time"

const (
	TnfTestSuiteName = "observability"

	// TNF test cases names.
	TnfContainerLoggingTcName = "observability-container-logging"
	TnfCrdStatusTcName        = "observability-crd-status"

	TestNamespace = "observability-ns"

	// Resources base names. In case a TC creates more than one of them,
	// an index will be appended: observability-dp1, observability-dp2...
	TestDeploymentBaseName  = "observability-dp"
	TestStatefulSetBaseName = "observability-st"
	TestDaemonSetBaseName   = "observability-ds"
	TestPodBaseName         = "observability-pod"
	TestContainerBaseName   = "qe-observability-container"

	TestPodLabelKey   = "tnf-qe/observability"
	TestPodLabelValue = "container-logging-tc"
)

// observability-container-logging helper params.
const (
	OneLogLine               = "Hello world line 1\n"
	OneLogLineWithoutNewLine = "Hello world line 1"

	TwoLogLines               = "Hello world line 1\nHello world line 2\n"
	TwoLogLinesWithoutNewLine = "Hello world line 1\nHello world line 2"

	NoLogLines = ""
)

// observability-crd-status helper params.
const (
	CrdSuffix1 = "tnf-qe.suffix1.com"
	CrdSuffix2 = "tnf-qe.suffix2.com"

	NotConfiguredCrdSuffix = "not-configured-suffix.com"
)

var (
	TnfTargetPodLabels = map[string]string{
		TestPodLabelKey: TestPodLabelValue,
	}

	PodDeployTimeoutMins         = 5 * time.Minute
	DeploymentDeployTimeoutMins  = 5 * time.Minute
	StatefulSetDeployTimeoutMins = 5 * time.Minute
	DaemonSetDeployTimeoutMins   = 5 * time.Minute
	CrdDeployTimeoutMins         = 1 * time.Minute

	NsResourcesDeleteTimeoutMins = 5 * time.Minute
)

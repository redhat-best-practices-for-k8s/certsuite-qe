package parameters

import "time"

const (
	TnfTestSuiteName = "observability"

	// TNF test cases names.
	TnfContainerLoggingTcName = "observability-container-logging"

	QeTestNamespace = "observability-ns"

	// Resources base names. In case a TC creates more than one of them,
	// an index will be appended: observability-dp1, observability-dp2...
	QeTestDeploymentBaseName  = "observability-dp"
	QeTestStatefulSetBaseName = "observability-st"
	QeTestDaemonSetBaseName   = "observability-ds"
	QeTestPodBaseName         = "observability-pod"
	QeTestContainerBaseName   = "qe-observability-container"

	QeTestPodLabelKey   = "tnf-qe/observability"
	QeTestPodLabelValue = "container-logging-tc"
)

// observability-container-logging helper params.
const (
	OneLogLine               = "Hello world line 1\n"
	OneLogLineWithoutNewLine = "Hello world line 1"

	TwoLogLines               = "Hello world line 1\nHello world line 2\n"
	TwoLogLinesWithoutNewLine = "Hello world line 1\nHello world line 2"

	NoLogLines = ""
)

var (
	TnfTargetPodLabels = map[string]string{
		QeTestPodLabelKey: QeTestPodLabelValue,
	}

	PodDeployTimeoutMins         = 5 * time.Minute
	DeploymentDeployTimeoutMins  = 5 * time.Minute
	StatefulSetDeployTimeoutMins = 5 * time.Minute
	DaemonSetDeployTimeoutMins   = 5 * time.Minute

	NsResourcesDeleteTimeoutMins = 5 * time.Minute
)

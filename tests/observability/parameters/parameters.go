package parameters

import "time"

const (
	CertsuiteTestSuiteName = "observability"

	// Certsuite test case names.
	CertsuiteContainerLoggingTcName     = "observability-container-logging"
	CertsuiteCrdStatusTcName            = "observability-crd-status"
	CertsuiteTerminationMsgPolicyTcName = "observability-termination-policy"
	CertsuitePodDisruptionBudgetTcName  = "observability-pod-disruption-budget"

	TestNamespace = "observability-ns"

	// Resources base names. In case a TC creates more than one of them,
	// an index will be appended: observability-dp1, observability-dp2...
	TestDeploymentBaseName  = "observability-dp"
	TestStatefulSetBaseName = "observability-st"
	TestDaemonSetBaseName   = "observability-ds"
	TestPodBaseName         = "observability-pod"
	TestContainerBaseName   = "qe-observability-container"
	TestPdbBaseName         = "observability-pdb"

	TestPodLabelKey   = "certsuite-qe/observability"
	TestPodLabelValue = "container-logging-tc"

	TestPodExtraLabelKey   = "extra-key"
	TestPodExtraLabelValue = "extra-value"

	UnknownKey   = "unknown-key"
	UnknownValue = "unknown-value"
)

// observability-container-logging helper params.
const (
	OneLogLine               = "Hello world line 1\n"
	OneLogLineWithoutNewLine = "Hello world line 1"

	TwoLogLines = "Hello world line 1\nHello world line 2\n"

	NoLogLines = ""
)

// observability-crd-status helper params.
const (
	CrdSuffix1 = "certsuite-qe.suffix1.com"
	CrdSuffix2 = "certsuite-qe.suffix2.com"

	NotConfiguredCrdSuffix = "not-configured-suffix.com"

	CrdRetryInterval = 5 * time.Second
)

// observability-termination-policy helper params.
const (
	UseDefaultTerminationMsgPolicy = ""
)

var (
	CertsuiteTargetPodLabels = map[string]string{
		TestPodLabelKey: TestPodLabelValue,
	}

	CertsuiteUnknownPodLabels = map[string]string{
		UnknownKey: UnknownValue,
	}

	TestContainerNormalCommand = []string{"/bin/bash", "-c", "sleep INF"}

	PodDeployTimeoutMins         = 5 * time.Minute
	DeploymentDeployTimeoutMins  = 5 * time.Minute
	StatefulSetDeployTimeoutMins = 5 * time.Minute
	ReplicaSetDeployTimeoutMins  = 5 * time.Minute
	DaemonSetDeployTimeoutMins   = 5 * time.Minute
	CrdDeployTimeoutMins         = 1 * time.Minute
	PdbDeployTimeoutMins         = 5 * time.Minute

	NsResourcesDeleteTimeoutMins = 5 * time.Minute
)

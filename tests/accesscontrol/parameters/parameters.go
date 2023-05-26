package parameters

import (
	"fmt"
	"time"
)

const (
	Timeout = 5 * time.Minute
)

var (
	TestAccessControlNameSpace           = "accesscontrol-tests"
	testPodLabelPrefixName               = "accesscontrol-test/test"
	testPodLabelValue                    = "testing"
	TestPodLabel                         = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	InvalidNamespace                     = "openshift-test"
	AdditionalValidNamespace             = "ac-test"
	AdditionalNamespaceForResourceQuotas = "ac-rq-test"

	TestCaseNameAccessControlNamespace              = "access-control-namespace"
	TestCaseNameAccessControlPodHostIpc             = "access-control-pod-host-ipc"
	TestCaseNameAccessControlPodHostPid             = "access-control-pod-host-pid"
	TestCaseNameAccessControlPodAutomountToken      = "access-control-pod-automount-service-account-token"
	TestCaseNameAccessControlPodHostNetwork         = "access-control-pod-host-network"
	TestCaseNameAccessControlSysPtraceCapability    = "access-control-sys-ptrace-capability"
	TestCaseNameAccessControlRequestsAndLimits      = "access-control-requests-and-limits"
	TestCaseNameAccessControlNo1337Uid              = "access-control-no-1337-uid"
	TestCaseNameAccessControlNamespaceResourceQuota = "access-control-namespace-resource-quota"
	TestCaseNameAccessControlIpcLockCapability      = "access-control-ipc-lock-capability-check"
	TestCaseNameAccessControlNetAdminCapability     = "access-control-net-admin-capability-check"
	TestCaseNameAccessControlNetRawCapability       = "access-control-net-raw-capability-check"
	TestCaseNameAccessControlContainerHostPort      = "access-control-container-host-port"
	TestCaseNameAccessControlSysAdminCapability     = "access-control-sys-admin-capability-check"

	TestDeploymentLabels = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}

	ServiceAccountName = "default"
	MemoryLimit        = "512Mi"
	MemoryRequest      = "500Mi"
	CPULimit           = "1"
	CPURequest         = "1"
)

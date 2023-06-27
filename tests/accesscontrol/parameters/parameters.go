package parameters

import (
	"fmt"
	"time"
)

const (
	Timeout = 5 * time.Minute
)

var (
	testPodLabelPrefixName               = "accesscontrol-test/test"
	testPodLabelValue                    = "testing"
	TestPodName                          = "access-control-pod"
	TestPodLabel                         = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	InvalidNamespace                     = "openshift-test"
	AdditionalValidNamespace             = "ac-test"
	AdditionalNamespaceForResourceQuotas = "ac-rq-test"

	TestDeploymentLabels = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
)

const (
	// TNF Test case names.
	TnfPodRoleBindings                              = "access-control-pod-role-bindings"
	TnfPodServiceAccount                            = "access-control-pod-service-account"
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
	TestCaseNameAccessControlNonRootUser            = "access-control-security-context-non-root-user-check"
	TestCaseNameAccessControlCluterRoleBindings     = "access-control-cluster-role-bindings"
	TnfNodePortTcName                               = "access-control-service-type"
	TestCaseNameAccessControlPrivilegeEscalation    = "access-control-security-context-privilege-escalation"
	TestCaseNameAccessControlPodHostPath            = "access-control-pod-host-path"
	TnfSecurityContextTcName                        = "access-control-security-context"

	TestAccessControlNameSpace = "accesscontrol-tests"

	ServiceAccountName = "automount-test-sa"
	MemoryLimit        = "512Mi"
	MemoryRequest      = "500Mi"
	CPULimit           = "1"
	CPURequest         = "1"

	TestServiceAccount   = "my-sa"
	TestRoleBindingName  = "my-rb"
	TestRoleName         = "my-r"
	TestAnotherNamespace = "my-ns"
)

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

	CertsuiteTargetOperatorLabels    = fmt.Sprintf("%s: %s", "cnf/test", "cr-scale-operator")
	CertsuiteTargetOperatorLabelsMap = map[string]string{
		"cnf/test": "cr-scale-operator",
	}
	CertsuiteTargetCrdFilters        = "memcacheds.cache.example.com"
	CertsuiteTargetOperatorNamespace = "cr-scale-operator-system"
	CertsuiteCustomResourceName      = "memcached-sample"

	CertsuiteCustomResourceAPIGroupName = "cache.example.com"
	CertsuiteCustomResourceResourceName = "memcacheds"

	SSHDaemonStartContainerCommand = []string{"/usr/sbin/sshd", "-f", "/home/tnf-user/sshd/sshd_config", "-D", "-d"}
)

const (
	// Certsuite test case names.
	CertsuiteCrdRoles                               = "access-control-crd-roles"
	CertsuitePodRoleBindings                        = "access-control-pod-role-bindings"
	CertsuitePodServiceAccount                      = "access-control-pod-service-account"
	CertsuiteNoSSHDaemonsAllowed                    = "access-control-ssh-daemons"
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
	TestCaseNameAccessControlBpfCapability          = "access-control-bpf-capability-check"
	TestCaseNameAccessControlContainerHostPort      = "access-control-container-host-port"
	TestCaseNameAccessControlSysAdminCapability     = "access-control-sys-admin-capability-check"
	TestCaseNameAccessControlNonRootUserID          = "access-control-security-context-non-root-user-id-check"
	TestCaseNameAccessControlClusterRoleBindings    = "access-control-cluster-role-bindings"
	CertsuiteNodePortTcName                         = "access-control-service-type"
	TestCaseNameAccessControlPrivilegeEscalation    = "access-control-security-context-privilege-escalation"
	TestCaseNameAccessControlPodHostPath            = "access-control-pod-host-path"
	CertsuiteSecurityContextTcName                  = "access-control-security-context"
	TestCaseNameAccessControlOneProcessPerContainer = "access-control-one-process-per-container"

	TestAccessControlNameSpace = "accesscontrol-tests"

	SSHDaemonImageName = "quay.io/testnetworkfunction/k8s-best-practices-debug:latest"

	ServiceAccountName = "automount-test-sa"
	MemoryLimit        = "112Mi"
	MemoryRequest      = "100Mi"
	CPULimit           = "500m"
	CPURequest         = "500m"

	TestServiceAccount   = "my-sa"
	TestRoleBindingName  = "my-rb"
	TestRoleName         = "my-r"
	TestAnotherNamespace = "my-ns"
)

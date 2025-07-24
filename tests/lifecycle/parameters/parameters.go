package parameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime   = 5 * time.Minute
	RetryInterval = 5
)

var (
	LifecycleNamespace     = "lifecycle-tests"
	testPodLabelPrefixName = "redhat-best-practices-for-k8s.com/lifecycle"
	testPodLabelValue      = "testing"
	PreStopCommand         = []string{"/bin/sh", "-c", "killall -0 tail"}
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestTargetLabels       = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
	AffinityRequiredPodLabels = map[string]string{
		"AffinityRequired": "true",
	}
	CertsuiteTargetOperatorLabels    = fmt.Sprintf("%s: %s", "cnf/test", "cr-scale-operator")
	CertsuiteTargetOperatorLabelsMap = map[string]string{
		"cnf/test": "cr-scale-operator",
	}
	CertsuiteTargetCrdFilters        = "memcacheds.cache.example.com"
	CertsuiteTargetOperatorNamespace = "cr-scale-operator-system"
	CertsuiteCustomResourceName      = "memcached-sample"

	TestLocalStorageClassName = "local-storage"
)

const (
	TestDeploymentName    = "lifecycle-dpa"
	TestDaemonSetName     = "lifecycle-dsa"
	TestStatefulSetName   = "lifecycle-sfa"
	TestPodName           = "lifecycle-pod"
	TestReplicaSetName    = "lifecycle-rsa"
	TestPVName            = "lifecycle-pv"
	TestPVCName           = "lifecycle-pvc"
	TestVolumeName        = "lifecycle-storage"
	CertsuiteRunTimeClass = "lifecycle-rtc"

	// Test Case names.
	CertsuiteCrdScaling                          = "lifecycle-crd-scaling"
	CertsuiteShutdownTcName                      = "lifecycle-container-prestop"
	CertsuiteDeploymentScalingTcName             = "lifecycle-deployment-scaling"
	CertsuitePodOwnerTypeTcName                  = "lifecycle-pod-owner-type"
	CertsuitePodRecreationTcName                 = "lifecycle-pod-recreation"
	CertsuitePodHighAvailabilityTcName           = "lifecycle-pod-high-availability"
	CertsuitePodSchedulingTcName                 = "lifecycle-pod-scheduling"
	CertsuiteLivenessTcName                      = "lifecycle-liveness-probe"
	CertsuiteReadinessTcName                     = "lifecycle-readiness-probe"
	CertsuiteStatefulSetScalingTcName            = "lifecycle-statefulset-scaling"
	CertsuiteImagePullPolicyTcName               = "lifecycle-image-pull-policy"
	CertsuitePersistentVolumeReclaimPolicyTcName = "lifecycle-persistent-volume-reclaim-policy"
	CertsuiteCPUIsolationTcName                  = "lifecycle-cpu-isolation"
	CertsuiteStartUpProbeTcName                  = "lifecycle-startup-probe"
	CertsuiteAffinityRequiredPodsTcName          = "lifecycle-affinity-required-pods"
	CertsuiteContainerStartUpTcName              = "lifecycle-container-poststart"
	CertsuitePodTolerationBypassTcName           = "lifecycle-pod-toleration-bypass"
	CertsuiteStorageProvisioner                  = "lifecycle-storage-provisioner"

	SampleWorkloadImage = "registry.access.redhat.com/ubi8/ubi-micro:latest"
)

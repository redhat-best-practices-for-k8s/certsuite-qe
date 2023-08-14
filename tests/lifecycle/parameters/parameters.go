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
	testPodLabelPrefixName = "test-network-function.com/lifecycle"
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
	TnfTargetOperatorLabels    = fmt.Sprintf("%s: %s", "cnf/test", "cr-scale-operator")
	TnfTargetOperatorLabelsMap = map[string]string{
		"cnf/test": "cr-scale-operator",
	}
	TnfTargetCrdFilters        = "memcacheds.cache.example.com"
	TnfTargetOperatorNamespace = "cr-scale-operator-system"
	TnfCustomResourceName      = "memcached-sample"

	TestLocalStorageClassName = "local-storage"
)

const (
	TestDeploymentName  = "lifecycle-dpa"
	TestDaemonSetName   = "lifecycle-dsa"
	TestStatefulSetName = "lifecycle-sfa"
	TestPodName         = "lifecycle-pod"
	TestReplicaSetName  = "lifecycle-rsa"
	TestPVName          = "lifecycle-pv"
	TestPVCName         = "lifecycle-pvc"
	TestVolumeName      = "lifecycle-storage"
	TnfRunTimeClass     = "lifecycle-rtc"

	// Test Case names.
	TnfCrdScaling                          = "lifecycle-crd-scaling"
	TnfShutdownTcName                      = "lifecycle-container-shutdown"
	TnfDeploymentScalingTcName             = "lifecycle-deployment-scaling"
	TnfPodOwnerTypeTcName                  = "lifecycle-pod-owner-type"
	TnfPodRecreationTcName                 = "lifecycle-pod-recreation"
	TnfPodHighAvailabilityTcName           = "lifecycle-pod-high-availability"
	TnfPodSchedulingTcName                 = "lifecycle-pod-scheduling"
	TnfLivenessTcName                      = "lifecycle-liveness-probe"
	TnfReadinessTcName                     = "lifecycle-readiness-probe"
	TnfStatefulSetScalingTcName            = "lifecycle-statefulset-scaling"
	TnfImagePullPolicyTcName               = "lifecycle-image-pull-policy"
	TnfPersistentVolumeReclaimPolicyTcName = "lifecycle-persistent-volume-reclaim-policy"
	TnfCPUIsolationTcName                  = "lifecycle-cpu-isolation"
	TnfStartUpProbeTcName                  = "lifecycle-startup-probe"
	TnfAffinityRequiredPodsTcName          = "lifecycle-affinity-required-pods"
	TnfContainerStartUpTcName              = "lifecycle-container-startup"
	TnfPodTolerationBypassTcName           = "lifecycle-pod-toleration-bypass"
	TnfStorageRequiredPods                 = "lifecycle-storage-required-pods"
)

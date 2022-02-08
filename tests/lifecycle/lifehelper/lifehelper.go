package lifehelper

import (
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// DefineDeployment defines a deployment.
func DefineDeployment(replica int32, containers int, name string) *v1.Deployment {
	deploymentStruct := globalhelper.AppendContainersToDeployment(
		deployment.RedefineWithReplicaNumber(
			deployment.DefineDeployment(
				name,
				lifeparameters.LifecycleNamespace,
				globalhelper.Configuration.General.TnfImage,
				lifeparameters.TestDeploymentLabels), replica),
		containers,
		globalhelper.Configuration.General.TnfImage)

	return deploymentStruct
}

// RemoveterminationGracePeriod removes terminationGracePeriodSeconds field in a deployment.
func RemoveterminationGracePeriod(deploymentStruct *v1.Deployment) *v1.Deployment {
	return deployment.RedefineWithTerminationGracePeriod(deploymentStruct, nil)
}

func DefineReplicaSet(name string) *v1.ReplicaSet {
	return replicaset.DefineReplicaSet(name,
		lifeparameters.LifecycleNamespace,
		globalhelper.Configuration.General.TnfImage,
		lifeparameters.TestDeploymentLabels)
}

func DefineStatefulSet(name string) *v1.StatefulSet {
	return statefulset.DefineStatefulSet(name,
		lifeparameters.LifecycleNamespace,
		globalhelper.Configuration.General.TnfImage,
		lifeparameters.TestDeploymentLabels)
}

func DefindPod(name string) *corev1.Pod {
	return pod.DefinePod(name, lifeparameters.LifecycleNamespace,
		globalhelper.Configuration.General.TnfImage,
		lifeparameters.TestDeploymentLabels)
}

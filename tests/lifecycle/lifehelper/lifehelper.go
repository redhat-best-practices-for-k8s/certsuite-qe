package lifehelper

import (
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	v1 "k8s.io/api/apps/v1"
)

// DefineLifecycleDeployment Defines basic deployment structure for lifecycle tests.
func DefineLifecycleDeployment(moreThanOneContainer bool, deploymentName string) *v1.Deployment {
	if !moreThanOneContainer {
		return deployment.DefineDeployment(
			deploymentName,
			lifeparameters.LifecycleNamespace,
			globalhelper.Configuration.General.TnfImage,
			lifeparameters.TestDeploymentLabels)
	}

	return deployment.DefineDeploymentWithTwoContainers(
		deploymentName,
		lifeparameters.LifecycleNamespace,
		globalhelper.Configuration.General.TnfImage,
		lifeparameters.TestDeploymentLabels)
}

/* DefineLifecycleDeploymentSeveralPodsWithTwoContainers Defines a deployment with
several pods and 2 containers. */
func DefineLifecycleDeploymentSeveralPodsWithTwoContainers(deploymentName string, numOfPods int32) *v1.Deployment {
	return deployment.RedefineWithReplicaNumber(
		DefineLifecycleDeployment(true, deploymentName),
		numOfPods)
}

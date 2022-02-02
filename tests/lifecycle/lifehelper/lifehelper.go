package lifehelper

import (
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	v1 "k8s.io/api/apps/v1"
)

// DefineDeployment defines a deployment with/out preStop field.
func DefineDeployment(preStop bool, replica int32, containers int, name string) (*v1.Deployment, error) {
	deploymentStruct := globalhelper.AppendContainersToDeployment(
		deployment.RedefineWithReplicaNumber(
			deployment.DefineDeployment(
				name,
				lifeparameters.LifecycleNamespace,
				globalhelper.Configuration.General.TnfImage,
				lifeparameters.TestDeploymentLabels), replica),
		containers,
		globalhelper.Configuration.General.TnfImage)

	if !preStop {
		return deploymentStruct, nil
	}

	return deployment.RedefineAllContainersWithPreStopSpec(
		deploymentStruct, lifeparameters.PreStopCommand)
}

// RemoveterminationGracePeriod removes terminationGracePeriodSeconds field in a deployment.
func RemoveterminationGracePeriod(deploymentStruct *v1.Deployment) *v1.Deployment {
	return deployment.RedefineWithTerminationGracePeriod(deploymentStruct, nil)
}

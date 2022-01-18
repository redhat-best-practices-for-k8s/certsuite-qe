package lifehelper

import (
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	v1 "k8s.io/api/apps/v1"
)

// DefineLifecycleDeployment Defines basic deployment structure for lifecycle tests.
func DefineLifecycleDeployment() *v1.Deployment {
	deploymentStruct := deployment.DefineDeployment(
		"lifecycleput",
		lifeparameters.LifecycleNamespace,
		globalhelper.Configuration.General.TnfImage,
		lifeparameters.TestDeploymentLabels)
	return deploymentStruct

}

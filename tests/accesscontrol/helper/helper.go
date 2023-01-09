package helper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/resourcequota"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteNamespaces(nsToDelete []string, clientSet *client.ClientSet, timeout time.Duration) error {
	failedNs := make(map[string]error)

	for _, namespace := range nsToDelete {
		err := namespaces.DeleteAndWait(
			clientSet,
			namespace,
			timeout,
		)
		if err != nil {
			failedNs[namespace] = err
		}
	}

	if len(failedNs) > 0 {
		return fmt.Errorf("failed to remove the following namespaces: %v", failedNs)
	}

	return nil
}

func DefineDeployment(replica int32, containers int, name string) (*v1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, parameters.TestAccessControlNameSpace,
		globalhelper.Configuration.General.TestImage, parameters.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, globalhelper.Configuration.General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func DefineDeploymentWithNamespace(replica int32, containers int, name string, namespace string) (*v1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		globalhelper.Configuration.General.TestImage, parameters.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, globalhelper.Configuration.General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func SetServiceAccountAutomountServiceAccountToken(namespace, saname, value string) error {
	var boolVal bool

	serviceacct, err := globalhelper.APIClient.ServiceAccounts(namespace).
		Get(context.TODO(), saname, metav1.GetOptions{})

	if err != nil {
		return fmt.Errorf("error getting service account: %w", err)
	}

	switch value {
	case "true":
		boolVal = true
		serviceacct.AutomountServiceAccountToken = &boolVal

	case "false":
		boolVal = false
		serviceacct.AutomountServiceAccountToken = &boolVal

	case "nil":
		serviceacct.AutomountServiceAccountToken = nil

	default:
		return fmt.Errorf("invalid value for token value")
	}

	_, err = globalhelper.APIClient.ServiceAccounts(parameters.TestAccessControlNameSpace).
		Update(context.TODO(), serviceacct, metav1.UpdateOptions{})

	return err
}

func DefineAndCreateResourceQuota(namespace string, clientSet *client.ClientSet) error {
	quota := resourcequota.DefineResourceQuota("quota1", parameters.CPURequest, parameters.MemoryRequest,
		parameters.CPULimit, parameters.MemoryLimit)

	return namespaces.ApplyResourceQuota(namespace, clientSet, quota)

}

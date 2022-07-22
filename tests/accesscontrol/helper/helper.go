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
		return nil, errors.New("invalid containers number")
	}

	deploymentStruct := globalhelper.AppendContainersToDeployment(
		deployment.RedefineWithReplicaNumber(
			deployment.DefineDeployment(
				name,
				parameters.TestAccessControlNameSpace,
				globalhelper.Configuration.General.TestImage,
				parameters.TestDeploymentLabels), replica),
		containers-1,
		globalhelper.Configuration.General.TestImage)

	return deploymentStruct, nil
}

func SetServiceAccountAutomountServiceAccountToken(namespace, saname, value string) error {
	var boolVal bool
	serviceacct, err := globalhelper.APIClient.ServiceAccounts(namespace).
		Get(context.TODO(), saname, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Error getting service account: %w", err)
	}

	if value == "true" {
		boolVal = true
		serviceacct.AutomountServiceAccountToken = &boolVal
	} else if value == "false" {
		boolVal = false
		serviceacct.AutomountServiceAccountToken = &boolVal
	} else if value == "nil" {
		serviceacct.AutomountServiceAccountToken = nil
	} else {
		return fmt.Errorf("Invalid value for token value")
	}

	_, err = globalhelper.APIClient.ServiceAccounts(parameters.TestAccessControlNameSpace).
		Update(context.TODO(), serviceacct, metav1.UpdateOptions{})
	return err
}

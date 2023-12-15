package helper

import (
	"context"
	"errors"
	"fmt"

	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/container"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/service"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineDeployment(replica int32, containers int, name, namespace string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		globalhelper.GetConfiguration().General.TestImage, parameters.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, globalhelper.GetConfiguration().General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func DefineDeploymentWithClusterRoleBindingWithServiceAccount(replica int32,
	containers int, name, namespace, serviceAccountName string) (*appsv1.Deployment, error) {
	deploymentStruct := deployment.DefineDeployment(name, namespace,
		globalhelper.GetConfiguration().General.TestImage, parameters.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, globalhelper.GetConfiguration().General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)
	deployment.AppendServiceAccount(deploymentStruct, serviceAccountName)

	return deploymentStruct, nil
}

func DefineDeploymentWithNamespace(replica int32, containers int, name string, namespace string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		globalhelper.GetConfiguration().General.TestImage, parameters.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, globalhelper.GetConfiguration().General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func DefineDeploymentWithContainerPorts(name, namespace string, replicaNumber int32,
	ports []corev1.ContainerPort) (*appsv1.Deployment, error) {
	if len(ports) < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		globalhelper.GetConfiguration().General.TestImage, parameters.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, len(ports)-1, globalhelper.GetConfiguration().General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replicaNumber)

	portSpecs := container.CreateContainerSpecsFromContainerPorts(ports,
		globalhelper.GetConfiguration().General.TestImage, "test")

	deployment.RedefineWithContainerSpecs(deploymentStruct, portSpecs)

	return deploymentStruct, nil
}

func SetServiceAccountAutomountServiceAccountToken(namespace, saname, value string) error {
	var boolVal bool

	serviceacct, err := globalhelper.GetAPIClient().ServiceAccounts(namespace).
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

	_, err = globalhelper.GetAPIClient().ServiceAccounts(namespace).
		Update(context.TODO(), serviceacct, metav1.UpdateOptions{})

	return err
}

// DefineAndCreateServiceOnCluster defines service resource and creates it on cluster.
func DefineAndCreateServiceOnCluster(name, namespace string, port int32, targetPort int32, withNodePort bool,
	ipFams []corev1.IPFamily, ipFamPolicy string) error {
	var testService *corev1.Service

	if ipFamPolicy == "" {
		testService = service.DefineService(
			name,
			namespace,
			port,
			targetPort,
			corev1.ProtocolTCP,
			parameters.TestDeploymentLabels,
			ipFams,
			nil)
	} else {
		ipPolicy := corev1.IPFamilyPolicy(ipFamPolicy)

		testService = service.DefineService(
			name,
			namespace,
			port,
			targetPort,
			corev1.ProtocolTCP,
			parameters.TestDeploymentLabels,
			ipFams,
			&ipPolicy)
	}

	if withNodePort {
		err := service.RedefineWithNodePort(testService)
		if err != nil {
			return err
		}
	}

	_, err := globalhelper.GetAPIClient().Services(namespace).Create(
		context.TODO(),
		testService, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create service on cluster: %w", err)
	}

	return nil
}

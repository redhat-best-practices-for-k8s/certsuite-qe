package globalhelper

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

func defineTnfNamespaces(config *globalparameters.TnfConfig, namespaces []string) error {
	if len(namespaces) < 1 {
		return fmt.Errorf("target namespaces cannot be empty list")
	}

	if config == nil {
		return fmt.Errorf("config struct cannot be nil")
	}

	for _, namespace := range namespaces {
		config.TargetNameSpaces = append(config.TargetNameSpaces, globalparameters.TargetNameSpace{
			Name: namespace,
		})
	}

	return nil
}

func defineTargetPodLabels(config *globalparameters.TnfConfig, targetPodLabels []string) error {
	if len(targetPodLabels) < 1 {
		return fmt.Errorf("target pod labels cannot be empty list")
	}

	for _, targetPodLabel := range targetPodLabels {
		prefixNameValue := strings.Split(targetPodLabel, "/")
		if len(prefixNameValue) != 2 {
			return fmt.Errorf(fmt.Sprintf("target pod label %s is invalid", targetPodLabel))
		}

		prefix := strings.TrimSpace(prefixNameValue[0])
		nameValue := strings.Split(prefixNameValue[1], ":")

		if len(nameValue) != 2 {
			return fmt.Errorf(fmt.Sprintf("target pod label %s is invalid", targetPodLabel))
		}

		name := strings.TrimSpace(nameValue[0])
		value := strings.TrimSpace(nameValue[1])

		config.TargetPodLabels = append(config.TargetPodLabels, globalparameters.PodLabel{
			Prefix: prefix,
			Name:   name,
			Value:  value,
		})
	}

	return nil
}

// ValidateIfReportsAreValid test if report is valid for given test case.
func ValidateIfReportsAreValid(tcName string, tcExpectedStatus string) error {
	glog.V(5).Info("Verify test case status in Junit report")

	junitTestReport, err := OpenJunitTestReport()

	if err != nil {
		return err
	}

	claimReport, err := OpenClaimReport()

	if err != nil {
		return err
	}

	err = IsExpectedStatusParamValid(tcExpectedStatus)

	if err != nil {
		return err
	}

	isTestCaseInValidStatusInJunitReport := IsTestCasePassedInJunitReport
	isTestCaseInValidStatusInClaimReport := IsTestCasePassedInClaimReport

	if tcExpectedStatus == globalparameters.TestCaseFailed {
		isTestCaseInValidStatusInJunitReport = IsTestCaseFailedInJunitReport
		isTestCaseInValidStatusInClaimReport = IsTestCaseFailedInClaimReport
	}

	if tcExpectedStatus == globalparameters.TestCaseSkipped {
		isTestCaseInValidStatusInJunitReport = IsTestCaseSkippedInJunitReport
		isTestCaseInValidStatusInClaimReport = IsTestCaseSkippedInClaimReport
	}

	if !isTestCaseInValidStatusInJunitReport(junitTestReport, tcName) {
		return fmt.Errorf("test case %s is not in expected %s state in junit report", tcName, tcExpectedStatus)
	}

	glog.V(5).Info("Verify test case status in claim report file")

	testPassed, err := isTestCaseInValidStatusInClaimReport(tcName, *claimReport)

	if err != nil {
		return err
	}

	if !testPassed {
		return fmt.Errorf("test case %s is not in expected %s state in claim report", tcName, tcExpectedStatus)
	}

	return nil
}

func defineCertifiedContainersInfo(config *globalparameters.TnfConfig, certifiedContainerInfo []string) error {
	if len(certifiedContainerInfo) < 1 {
		// do not add certifiedcontainerinfo to tnf_config at all in this case
		return nil
	}

	for _, certifiedContainerFields := range certifiedContainerInfo {
		nameRepository := strings.Split(certifiedContainerFields, "/")

		if len(nameRepository) == 1 {
			// certifiedContainerInfo item does not contain separation character
			// use this to add only the Certifiedcontainerinfo field with no sub fields
			var emptyInfo globalparameters.CertifiedContainerRepoInfo
			config.Certifiedcontainerinfo = append(config.Certifiedcontainerinfo, emptyInfo)

			return nil
		}

		if len(nameRepository) != 2 {
			return fmt.Errorf(fmt.Sprintf("certified container info %s is invalid", certifiedContainerFields))
		}

		name := strings.TrimSpace(nameRepository[0])
		repo := strings.TrimSpace(nameRepository[1])

		glog.V(5).Info(fmt.Sprintf("Adding container name:%s repository:%s to configuration", name, repo))

		config.Certifiedcontainerinfo = append(config.Certifiedcontainerinfo, globalparameters.CertifiedContainerRepoInfo{
			Name:       name,
			Repository: repo,
		})
	}

	return nil
}

func defineCertifiedOperatorsInfo(config *globalparameters.TnfConfig, certifiedOperatorInfo []string) error {
	if len(certifiedOperatorInfo) < 1 {
		// do not add certifiedoperatorinfo to tnf_config at all in this case
		return nil
	}

	for _, certifiedOperatorFields := range certifiedOperatorInfo {
		nameOrganization := strings.Split(certifiedOperatorFields, "/")

		if len(nameOrganization) == 1 {
			// certifiedOperatorInfo item does not contain separation character
			// use this to add only the Certifiedoperatorinfo field with no sub fields
			var emptyInfo globalparameters.CertifiedOperatorRepoInfo
			config.Certifiedoperatorinfo = append(config.Certifiedoperatorinfo, emptyInfo)

			return nil
		}

		if len(nameOrganization) != 2 {
			return fmt.Errorf(fmt.Sprintf("certified operator info %s is invalid", certifiedOperatorFields))
		}

		name := strings.TrimSpace(nameOrganization[0])
		org := strings.TrimSpace(nameOrganization[1])

		glog.V(5).Info(fmt.Sprintf("Adding operator name:%s organization:%s to configuration", name, org))

		config.Certifiedoperatorinfo = append(config.Certifiedoperatorinfo, globalparameters.CertifiedOperatorRepoInfo{
			Name:         name,
			Organization: org,
		})
	}

	return nil
}

// DefineTnfConfig creates tnf_config.yml file under tnf config directory.
func DefineTnfConfig(namespaces []string, targetPodLabels []string, certifiedContainerInfo []string,
	certifiedOperatorsInfo []string) error {
	configFile, err := os.OpenFile(
		path.Join(
			Configuration.General.TnfConfigDir,
			globalparameters.DefaultTnfConfigFileName),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error opening/creating file: %w", err)
	}
	defer configFile.Close()
	configFileEncoder := yaml.NewEncoder(configFile)
	tnfConfig := globalparameters.TnfConfig{}

	err = defineTnfNamespaces(&tnfConfig, namespaces)
	if err != nil {
		return err
	}

	err = defineTargetPodLabels(&tnfConfig, targetPodLabels)
	if err != nil {
		return err
	}

	err = defineCertifiedContainersInfo(&tnfConfig, certifiedContainerInfo)
	if err != nil {
		return err
	}

	err = defineCertifiedOperatorsInfo(&tnfConfig, certifiedOperatorsInfo)
	if err != nil {
		return err
	}

	err = configFileEncoder.Encode(tnfConfig)

	glog.V(5).Info(fmt.Sprintf("%s deployed under %s directory",
		globalparameters.DefaultTnfConfigFileName, Configuration.General.TnfConfigDir))

	return err
}

// IsExpectedStatusParamValid validates if requested test status is valid.
func IsExpectedStatusParamValid(status string) error {
	return validateIfParamInAllowedListOfParams(
		status,
		[]string{globalparameters.TestCaseFailed, globalparameters.TestCasePassed, globalparameters.TestCaseSkipped})
}

func validateIfParamInAllowedListOfParams(parameter string, listOfParameters []string) error {
	for _, allowedParameter := range listOfParameters {
		if allowedParameter == parameter {
			return nil
		}
	}

	return fmt.Errorf("parameter %s is not allowed. List of allowed parameters %s", parameter, listOfParameters)
}

// AppendContainersToDeployment appends containers to a deployment.
func AppendContainersToDeployment(deployment *v1.Deployment, containersNum int, image string) *v1.Deployment {
	containerList := &deployment.Spec.Template.Spec.Containers

	for i := 0; i < containersNum; i++ {
		*containerList = append(
			*containerList, corev1.Container{
				Name:    fmt.Sprintf("container%d", i+1),
				Image:   image,
				Command: []string{"/bin/bash", "-c", "sleep INF"},
			})
	}

	return deployment
}

// CreateAndWaitUntilReplicaSetIsReady creates replicaSet and wait until all replicas are up and running.
func CreateAndWaitUntilStatefulSetIsReady(statefulSet *v1.StatefulSet, timeout time.Duration) error {
	runningReplica, err := APIClient.StatefulSets(statefulSet.Namespace).Create(
		context.Background(),
		statefulSet,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isStatefulSetReady(runningReplica.Namespace, runningReplica.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"statefulSet %s is not ready, retry in 5 seconds", runningReplica.Name))

			return false
		}

		return status
	}, timeout, 5*time.Second).Should(Equal(true), "statefulSet is not ready")

	return nil
}

// CreateAndWaitUntilReplicaSetIsReady creates replicaSet and wait until all replicas are up and running.
func CreateAndWaitUntilReplicaSetIsReady(replicaSet *v1.ReplicaSet, timeout time.Duration) error {
	runningReplica, err := APIClient.ReplicaSets(replicaSet.Namespace).Create(
		context.Background(),
		replicaSet,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isReplicaSetReady(runningReplica.Namespace, runningReplica.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"replicaSet %s is not ready, retry in 5 seconds", runningReplica.Name))

			return false
		}

		return status
	}, timeout, 5*time.Second).Should(Equal(true), "replicaSet is not ready")

	return nil
}

func isReplicaSetReady(namespace string, replicaSetName string) (bool, error) {
	testReplicaSet, err := APIClient.ReplicaSets(namespace).Get(
		context.Background(),
		replicaSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if testReplicaSet.Status.ReadyReplicas > 0 {
		if testReplicaSet.Status.Replicas == testReplicaSet.Status.ReadyReplicas {
			return true, nil
		}
	}

	return false, nil
}

func isStatefulSetReady(namespace string, statefulSetName string) (bool, error) {
	testStatefulSet, err := APIClient.StatefulSets(namespace).Get(
		context.Background(),
		statefulSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if testStatefulSet.Status.ReadyReplicas > 0 {
		if testStatefulSet.Status.Replicas == testStatefulSet.Status.ReadyReplicas {
			return true, nil
		}
	}

	return false, nil
}

func isPodReady(namespace string, podName string) (bool, error) {
	pod, err := APIClient.Pods(namespace).Get(
		context.Background(),
		podName,
		metav1.GetOptions{},
	)

	if err != nil {
		return false, err
	}

	if pod.Status.Phase == "Running" {
		return true, nil
	}

	return false, nil
}

func CreateAndWaitUntilPodIsReady(pod *corev1.Pod, timeout time.Duration) error {
	pod, err := APIClient.Pods(pod.Namespace).Create(
		context.Background(),
		pod,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isPodReady(pod.Namespace, pod.Name)
		if err != nil {

			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in 5 seconds", pod.Name))

			return false
		}

		return status
	}, timeout, 5*time.Second).Should(Equal(true), "Deployment is not ready")

	return nil
}

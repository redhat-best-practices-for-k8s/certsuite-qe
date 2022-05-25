package globalhelper

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
	v1 "k8s.io/api/apps/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

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

// DefineTnfConfig creates tnf_config.yml file under tnf config directory.
func DefineTnfConfig(namespaces []string, targetPodLabels []string,
	certifiedContainerInfo []string, crdFilters []string) error {
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

	// CRD filters is an optional field.
	defineCrdFilters(&tnfConfig, crdFilters)

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

func defineCertifiedContainersInfo(config *globalparameters.TnfConfig, certifiedContainerInfo []string) error {
	if len(certifiedContainerInfo) < 1 {
		// do not add certifiedcontainerinfo to tnf_config at all in this case
		return nil
	}

	for _, certifiedContainerFields := range certifiedContainerInfo {
		nameRepositoryTagDigest := strings.Split(certifiedContainerFields, ";")

		if len(nameRepositoryTagDigest) == 1 {
			// certifiedContainerInfo item does not contain separation character
			// use this to add only the Certifiedcontainerinfo field with no sub fields
			var emptyInfo globalparameters.CertifiedContainerRepoInfo
			config.Certifiedcontainerinfo = append(config.Certifiedcontainerinfo, emptyInfo)

			return nil
		}

		if len(nameRepositoryTagDigest) != 4 {
			return fmt.Errorf(fmt.Sprintf("certified container info %s is invalid", certifiedContainerFields))
		}

		name := strings.TrimSpace(nameRepositoryTagDigest[0])
		repo := strings.TrimSpace(nameRepositoryTagDigest[1])
		tag := strings.TrimSpace(nameRepositoryTagDigest[2])
		digest := strings.TrimSpace(nameRepositoryTagDigest[3])

		glog.V(5).Info(fmt.Sprintf("Adding container name:%s repository:%s to configuration", name, repo))

		config.Certifiedcontainerinfo = append(config.Certifiedcontainerinfo, globalparameters.CertifiedContainerRepoInfo{
			Name:       name,
			Repository: repo,
			Tag:        tag,
			Digest:     digest,
		})
	}

	return nil
}

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

func defineCrdFilters(config *globalparameters.TnfConfig, crdSuffixes []string) {
	for _, crdSuffix := range crdSuffixes {
		glog.V(5).Info(fmt.Sprintf("Adding crd suffix %s to tnf configuration file", crdSuffix))

		config.TargetCrdFilters = append(config.TargetCrdFilters, globalparameters.TargetCrdFilter{
			NameSuffix: crdSuffix,
		})
	}
}

func validateIfParamInAllowedListOfParams(parameter string, listOfParameters []string) error {
	for _, allowedParameter := range listOfParameters {
		if allowedParameter == parameter {
			return nil
		}
	}

	return fmt.Errorf("parameter %s is not allowed. List of allowed parameters %s", parameter, listOfParameters)
}

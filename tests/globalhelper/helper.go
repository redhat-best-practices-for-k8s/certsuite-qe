package globalhelper

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"gopkg.in/yaml.v2"
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

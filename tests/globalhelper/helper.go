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

// DefineTnfConfig creates tnf_config.yml file under tnf config directory
func DefineTnfConfig(namespaces []string, targetPodLabels []string) error {
	configFile, err := os.OpenFile(
		path.Join(
			Configuration.General.TnfConfigDir,
			globalparameters.DefaultTnfConfigFileName),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error opening/creating file: %v", err)
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
	err = configFileEncoder.Encode(tnfConfig)
	glog.V(5).Info(fmt.Sprintf("%s deployed under %s directory",
		globalparameters.DefaultTnfConfigFileName, Configuration.General.TnfConfigDir))
	return err
}

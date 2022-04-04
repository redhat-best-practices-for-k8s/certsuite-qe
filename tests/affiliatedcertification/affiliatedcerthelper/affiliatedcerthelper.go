package affiliatedcerthelper

import (
	"context"
	"fmt"
	"strings"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcertparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	utils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"
)

func SetUpAndRunContainerCertTest(tcName string, containersInfo []string, expectedResult string) error {
	var err error

	ginkgo.By("Add container information to " + globalparameters.DefaultTnfConfigFileName)

	err = globalhelper.DefineTnfConfig(
		[]string{affiliatedcertparameters.TestCertificationNameSpace},
		[]string{affiliatedcertparameters.TestPodLabel},
		containersInfo,
		[]string{})

	if err != nil {
		return fmt.Errorf("error defining tnf config file: %w", err)
	}

	ginkgo.By("Start test")

	err = globalhelper.LaunchTests(
		affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
		tcName,
		affiliatedcertparameters.TestCaseContainerSkipRegEx,
	)

	if strings.Contains(expectedResult, globalparameters.TestCaseFailed) && err == nil {
		return fmt.Errorf("error running %s test",
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName)
	}

	if (strings.Contains(expectedResult, globalparameters.TestCasePassed) ||
		strings.Contains(expectedResult, globalparameters.TestCaseSkipped)) && err != nil {
		return fmt.Errorf("error running %s test: %w",
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName, err)
	}

	ginkgo.By("Verify test case status in Junit and Claim reports")

	err = globalhelper.ValidateIfReportsAreValid(
		affiliatedcertparameters.TestCaseContainerAffiliatedCertName,
		expectedResult)

	if err != nil {
		return fmt.Errorf("error validating test reports: %w", err)
	}

	return nil
}

func SetUpAndRunOperatorCertTest(tcName string, operatorsInfo []string, expectedResult string) error {
	ginkgo.By("Add container information to " + globalparameters.DefaultTnfConfigFileName)

	err := globalhelper.DefineTnfConfig(
		[]string{affiliatedcertparameters.TestCertificationNameSpace},
		[]string{affiliatedcertparameters.TestPodLabel},
		[]string{},
		operatorsInfo)
	if err != nil {
		return fmt.Errorf("error defining tnf config file: %w", err)
	}

	ginkgo.By("Start test")

	err = globalhelper.LaunchTests(
		affiliatedcertparameters.AffiliatedCertificationTestSuiteName,
		tcName,
		affiliatedcertparameters.TestCaseOperatorSkipRegEx,
	)

	if strings.Contains(expectedResult, globalparameters.TestCaseFailed) && err == nil {
		return fmt.Errorf("error running %s test",
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName)
	}

	if (strings.Contains(expectedResult, globalparameters.TestCasePassed) ||
		strings.Contains(expectedResult, globalparameters.TestCaseSkipped)) && err != nil {
		return fmt.Errorf("error running %s test: %w",
			affiliatedcertparameters.AffiliatedCertificationTestSuiteName, err)
	}

	ginkgo.By("Verify test case status in Junit and Claim reports")

	err = globalhelper.ValidateIfReportsAreValid(
		affiliatedcertparameters.TestCaseOperatorAffiliatedCertName,
		expectedResult)

	if err != nil {
		return fmt.Errorf("error validating test reports: %w", err)
	}

	return nil
}

// AddLabelToInstalledCSV adds given label to existing csv object.
func AddLabelToInstalledCSV(prefixCsvName string, namespace string, label map[string]string) error {
	csv, err := getCsvByPrefix(prefixCsvName, namespace)
	if err != nil {
		return err
	}

	newMap := make(map[string]string)
	for k, v := range csv.GetLabels() {
		newMap[k] = v
	}

	for k, v := range label {
		newMap[k] = v
	}

	csv.SetLabels(newMap)

	return updateCsv(namespace, csv)
}

// DeleteLabelFromInstalledCSV removes given label from existing csv object.
func DeleteLabelFromInstalledCSV(prefixCsvName string, namespace string, label map[string]string) error {
	csv, err := getCsvByPrefix(prefixCsvName, namespace)
	if err != nil {
		return err
	}

	newMap := make(map[string]string)

	for k, v := range csv.GetLabels() {
		if _, ok := label[k]; ok {
			continue
		}

		newMap[k] = v
	}

	csv.SetLabels(newMap)

	return updateCsv(namespace, csv)
}

func getCsvByPrefix(prefixCsvName string, namespace string) (*v1alpha1.ClusterServiceVersion, error) {
	csvs, err := globalhelper.APIClient.ClusterServiceVersions(namespace).List(
		context.TODO(), metav1.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	var neededCSV v1alpha1.ClusterServiceVersion

	for _, csv := range csvs.Items {
		if strings.Contains(csv.Name, prefixCsvName) {
			neededCSV = csv
		}
	}

	if neededCSV.Name == "" {
		return nil, fmt.Errorf("failed to detect a needed CSV")
	}

	return &neededCSV, nil
}

func updateCsv(namespace string, csv *v1alpha1.ClusterServiceVersion) error {
	_, err := globalhelper.APIClient.ClusterServiceVersions(namespace).Update(
		context.TODO(), csv, metav1.UpdateOptions{},
	)

	if err != nil {
		return fmt.Errorf("fail to update CSV due to %w", err)
	}

	return nil
}

// Deploys an operator and returns when its deployment is confirmed to exist.
func DeployAndVerifyOperatorSubscription(operatorPackage, chanel, namespace, group,
	sourceNamespace, deploymentName string) error {
	operatorSubscription := utils.DefineSubscription(
		operatorPackage+"-subscription",
		namespace,
		chanel,
		operatorPackage,
		group,
		sourceNamespace)

	err := DeployOperator(namespace, operatorSubscription)

	if err != nil {
		return fmt.Errorf("Error deploying operator "+operatorPackage+": %w", err)
	}

	// confirm that operator is installed and ready
	gomega.Eventually(func() bool {
		err = IsOperatorInstalled(namespace, deploymentName)

		return err == nil
	}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(gomega.Equal(true))

	if err != nil {
		return fmt.Errorf("Operator "+operatorPackage+" is not ready: %w", err)
	}

	return nil
}

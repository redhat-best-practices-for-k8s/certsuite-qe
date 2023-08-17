package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/onsi/ginkgo/v2"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	utils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"
)

func SetUpAndRunContainerCertTest(tcName string, containersInfo []string, expectedResult string) error {
	var err error

	ginkgo.By("Add container information to " + globalparameters.DefaultTnfConfigFileName)

	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.TestCertificationNameSpace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		containersInfo,
		[]string{})

	if err != nil {
		return fmt.Errorf("error defining tnf config file: %w", err)
	}

	ginkgo.By("Start test")

	err = globalhelper.LaunchTests(
		tsparams.TestCaseContainerAffiliatedCertName,
		tcName)

	if strings.Contains(expectedResult, globalparameters.TestCaseFailed) && err == nil {
		return fmt.Errorf("error running %s test",
			tsparams.TestCaseContainerAffiliatedCertName)
	}

	if (strings.Contains(expectedResult, globalparameters.TestCasePassed) ||
		strings.Contains(expectedResult, globalparameters.TestCaseSkipped)) && err != nil {
		return fmt.Errorf("error running %s test: %w",
			tsparams.TestCaseContainerAffiliatedCertName, err)
	}

	ginkgo.By("Verify test case status in Junit and Claim reports")

	err = globalhelper.ValidateIfReportsAreValid(
		tsparams.TestCaseContainerAffiliatedCertName,
		expectedResult)

	if err != nil {
		return fmt.Errorf("error validating test reports: %w", err)
	}

	return nil
}

// AddLabelToInstalledCSV adds given label to existing csv object.
func AddLabelToInstalledCSV(prefixCsvName string, namespace string, label map[string]string) error {
	csv, err := GetCsvByPrefix(prefixCsvName, namespace)
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
	csv, err := GetCsvByPrefix(prefixCsvName, namespace)
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

func DoesOperatorHaveLabels(prefixCsvName string, namespace string, labels map[string]string) (bool, error) {
	csv, err := GetCsvByPrefix(prefixCsvName, namespace)
	if err != nil {
		return false, err
	}

	csvLabels := csv.GetLabels()
	for k, v := range labels {
		csvLabelValue, csvLabelExists := csvLabels[k]
		if !csvLabelExists || csvLabelValue != v {
			return false, nil
		}
	}

	return true, nil
}

// GetCsvByPrefix returns csv object based on given prefix.
func GetCsvByPrefix(prefixCsvName string, namespace string) (*v1alpha1.ClusterServiceVersion, error) {
	csvs, err := globalhelper.GetAPIClient().ClusterServiceVersions(namespace).List(
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

func DeployOperatorSubscription(operatorPackage, channel, namespace, group,
	sourceNamespace, startingCSV string, installApproval v1alpha1.Approval) error {
	operatorSubscription := utils.DefineSubscription(
		operatorPackage+"-subscription",
		namespace,
		channel,
		operatorPackage,
		group,
		sourceNamespace,
		startingCSV,
		installApproval)

	err := globalhelper.DeployOperator(operatorSubscription)

	if err != nil {
		return fmt.Errorf("Error deploying operator "+operatorPackage+": %w", err)
	}

	return nil
}

func updateCsv(namespace string, csv *v1alpha1.ClusterServiceVersion) error {
	_, err := globalhelper.GetAPIClient().ClusterServiceVersions(namespace).Update(
		context.TODO(), csv, metav1.UpdateOptions{},
	)

	if err != nil {
		return fmt.Errorf("fail to update CSV due to %w", err)
	}

	return nil
}

package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/operator/parameters"
	utils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"
)

func DeployTestOperatorGroup() error {
	if globalhelper.IsOperatorGroupInstalled(tsparams.OperatorGroupName,
		tsparams.OperatorNamespace) != nil {
		return globalhelper.DeployOperatorGroup(tsparams.OperatorNamespace,
			utils.DefineOperatorGroup(tsparams.OperatorGroupName,
				tsparams.OperatorNamespace,
				[]string{tsparams.OperatorNamespace}),
		)
	}

	return nil
}

func WaitUntilOperatorIsReady(csvPrefix, namespace string) error {
	var err error

	var csv *v1alpha1.ClusterServiceVersion

	Eventually(func() bool {
		csv, err = GetCsvByPrefix(csvPrefix, namespace)
		if csv != nil && csv.Status.Phase != v1alpha1.CSVPhaseNone {
			return csv.Status.Phase != "InstallReady" &&
				csv.Status.Phase != "Deleting" &&
				csv.Status.Phase != "Installing" &&
				csv.Status.Phase != "Replacing" &&
				csv.Status.Phase != "Unknown"
		}

		return false
	}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
		csvPrefix+" is not ready.")

	return err
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

func DeployOperatorSubscriptionWithNodeSelector(operatorPackage, channel, namespace, group,
	sourceNamespace, startingCSV string, installApproval v1alpha1.Approval, nodeSelector map[string]string) error {
	operatorSubscription := utils.DefineSubscriptionWithNodeSelector(
		operatorPackage+"-subscription",
		namespace,
		channel,
		operatorPackage,
		group,
		sourceNamespace,
		startingCSV,
		installApproval,
		nodeSelector)

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

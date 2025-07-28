package helper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
	utils "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/operator"
)

func DeployTestOperatorGroup(namespace string, clusterWide bool) error {
	if globalhelper.IsOperatorGroupInstalled(tsparams.OperatorGroupName,
		namespace) != nil {
		targetNamespaces := []string{namespace}
		if clusterWide {
			targetNamespaces = []string{}
		}

		return globalhelper.DeployOperatorGroup(namespace,
			utils.DefineOperatorGroup(tsparams.OperatorGroupName,
				namespace,
				targetNamespaces),
		)
	}

	return nil
}

func WaitUntilOperatorIsReady(csvPrefix, namespace string) error {
	var err error

	var csv *v1alpha1.ClusterServiceVersion

	Eventually(func() bool {
		// Fetch the CSV from the cluster
		csv, err = GetCsvByPrefix(csvPrefix, namespace)
		if err != nil {
			glog.V(5).Infof("Failed to get CSV with prefix %s: %v", csvPrefix, err)

			return false
		}

		if csv != nil && csv.Status.Phase == v1alpha1.CSVPhaseSucceeded {
			glog.V(5).Infof("CSV %s is in Succeeded phase", csv.Name)
			checkInstalledCSV(namespace)
			glog.Infof("Operator %s is ready and subscription verification passed in namespace %s", csvPrefix, namespace)

			return true
		}

		displayCSVStatus(csv, 5)

		return false
	}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
		csvPrefix+" is not ready.")

	// Display deployment status in the namespace
	glog.V(4).Infof("Checking deployment status for operator %s in namespace %s", csvPrefix, namespace)
	displayDeploymentStatus(namespace, 4)

	return nil
}

func checkInstalledCSV(namespace string) {
	// Additional check: verify that the subscription's installedcsv matches currentCsv
	// This ensures the operator was installed successfully
	Eventually(func() bool {
		subscriptions, err := globalhelper.GetAPIClient().OperatorsV1alpha1Interface.Subscriptions(namespace).List(
			context.TODO(), metav1.ListOptions{})
		if err != nil {
			glog.V(5).Infof("Failed to list subscriptions in namespace %s: %v", namespace, err)

			return false
		}

		if len(subscriptions.Items) == 0 {
			glog.V(5).Infof("No subscriptions found in namespace %s", namespace)

			return false
		}

		// Assume there's a single subscription in the namespace
		subscription := subscriptions.Items[0]

		if subscription.Status.InstalledCSV == "" || subscription.Status.CurrentCSV == "" {
			glog.V(5).Infof("Subscription %s: InstalledCSV='%s', CurrentCSV='%s' - not ready yet",
				subscription.Name, subscription.Status.InstalledCSV, subscription.Status.CurrentCSV)

			return false
		}

		if subscription.Status.InstalledCSV != subscription.Status.CurrentCSV {
			glog.V(5).Infof("Subscription %s: InstalledCSV='%s' does not match CurrentCSV='%s'",
				subscription.Name, subscription.Status.InstalledCSV, subscription.Status.CurrentCSV)

			return false
		}

		glog.V(5).Infof("Subscription %s: InstalledCSV matches CurrentCSV (%s) - operator installed successfully",
			subscription.Name, subscription.Status.InstalledCSV)

		return true
	}, 10*time.Second, tsparams.PollingInterval).Should(Equal(true),
		"Subscription's installedCSV does not match currentCSV within 10 seconds")
}

func displayCSVStatus(csv *v1alpha1.ClusterServiceVersion, logLevel int) {
	// Display the whole CSV status field
	if csv != nil {
		glog.V(glog.Level(logLevel)).Infof("Waiting for CSV to be ready. CSV: %s", csv.Name)
		glog.V(glog.Level(logLevel)).Infof("CSV Status Phase: %s", csv.Status.Phase)
		glog.V(glog.Level(logLevel)).Infof("CSV Status Message: %s", csv.Status.Message)
		glog.V(glog.Level(logLevel)).Infof("CSV Status Reason: %s", csv.Status.Reason)
		glog.V(glog.Level(logLevel)).Infof("CSV Status LastUpdateTime: %s", csv.Status.LastUpdateTime)
		glog.V(glog.Level(logLevel)).Infof("CSV Status LastTransitionTime: %s", csv.Status.LastTransitionTime)

		// Display conditions if available
		if len(csv.Status.Conditions) > 0 {
			glog.V(glog.Level(logLevel)).Infof("CSV Conditions:")

			for _, condition := range csv.Status.Conditions {
				glog.V(glog.Level(logLevel)).Infof("  Phase: %s, Reason: %s, Message: %s, LastUpdate: %s",
					condition.Phase, condition.Reason, condition.Message, condition.LastUpdateTime)
			}
		}

		// Display requirements if available
		if len(csv.Status.RequirementStatus) > 0 {
			glog.V(glog.Level(logLevel)).Infof("CSV Requirements:")

			for _, req := range csv.Status.RequirementStatus {
				glog.V(glog.Level(logLevel)).Infof("  Group: %s, Version: %s, Kind: %s, Name: %s, Status: %s, Message: %s",
					req.Group, req.Version, req.Kind, req.Name, req.Status, req.Message)
			}
		}
	}
}

func displayDeploymentStatus(namespace string, logLevel int) {
	deployments, err := globalhelper.GetAPIClient().AppsV1Interface.Deployments(namespace).List(
		context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.V(glog.Level(logLevel)).Infof("Failed to list deployments in namespace %s: %v", namespace, err)

		return
	}

	if len(deployments.Items) == 0 {
		glog.V(glog.Level(logLevel)).Infof("No deployments found in namespace %s", namespace)

		return
	}

	allDeploymentsReady := true

	for _, deployment := range deployments.Items {
		glog.V(glog.Level(logLevel)).Infof("Deployment %s status:", deployment.Name)
		glog.V(glog.Level(logLevel)).Infof("  Replicas: %d", deployment.Status.Replicas)
		glog.V(glog.Level(logLevel)).Infof("  ReadyReplicas: %d", deployment.Status.ReadyReplicas)
		glog.V(glog.Level(logLevel)).Infof("  AvailableReplicas: %d", deployment.Status.AvailableReplicas)
		glog.V(glog.Level(logLevel)).Infof("  UnavailableReplicas: %d", deployment.Status.UnavailableReplicas)
		glog.V(glog.Level(logLevel)).Infof("  UpdatedReplicas: %d", deployment.Status.UpdatedReplicas)
		glog.V(glog.Level(logLevel)).Infof("  ObservedGeneration: %d", deployment.Status.ObservedGeneration)

		// Display deployment conditions
		if len(deployment.Status.Conditions) > 0 {
			glog.V(glog.Level(logLevel)).Infof("  Deployment Conditions:")

			for _, condition := range deployment.Status.Conditions {
				glog.V(glog.Level(logLevel)).Infof("    Type: %s, Status: %s, Reason: %s, Message: %s, LastUpdate: %s",
					condition.Type, condition.Status, condition.Reason, condition.Message, condition.LastUpdateTime)
			}
		}

		// Check if deployment is ready
		desiredReplicas := *deployment.Spec.Replicas
		if deployment.Status.ReadyReplicas != desiredReplicas ||
			deployment.Status.AvailableReplicas != desiredReplicas {
			glog.V(glog.Level(logLevel)).Infof("Deployment %s is not ready: ReadyReplicas=%d, AvailableReplicas=%d, DesiredReplicas=%d",
				deployment.Name, deployment.Status.ReadyReplicas, deployment.Status.AvailableReplicas, desiredReplicas)
			allDeploymentsReady = false
		} else {
			glog.V(glog.Level(logLevel)).Infof("Deployment %s is ready", deployment.Name)
		}
	}

	if allDeploymentsReady {
		glog.V(glog.Level(logLevel)).Infof("All deployments are ready in namespace %s", namespace)
	} else {
		glog.V(glog.Level(logLevel)).Infof("Not all deployments are ready in namespace %s", namespace)
	}
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
		if strings.HasPrefix(csv.Name, prefixCsvName) {
			neededCSV = csv
		}
	}

	if neededCSV.Name == "" {
		return nil, fmt.Errorf("failed to detect a needed CSV")
	}

	return &neededCSV, nil
}

func DeployOperatorSubscription(subscriptionName, operatorPackage, channel, namespace, group,
	sourceNamespace, startingCSV string, installApproval v1alpha1.Approval) error {
	operatorSubscription := utils.DefineSubscription(
		subscriptionName+"-subscription",
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

// IsCSVNotSucceeded checks if CSV installation status is not Succeeded.
func IsCSVNotSucceeded(csvPrefix, namespace string) (bool, error) {
	csv, err := GetCsvByPrefix(csvPrefix, namespace)
	if err != nil {
		return false, err
	}

	// Return true if CSV is NOT in Succeeded phase
	return csv.Status.Phase != v1alpha1.CSVPhaseSucceeded, nil
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

func DeployTestOperatorGroupWithTargetNamespace(operatorGroupName, namespace string, targetNamespaces []string) error {
	if globalhelper.IsOperatorGroupInstalled(operatorGroupName,
		namespace) != nil {
		return globalhelper.DeployOperatorGroup(namespace,
			utils.DefineOperatorGroup(operatorGroupName,
				namespace,
				targetNamespaces),
		)
	}

	return nil
}

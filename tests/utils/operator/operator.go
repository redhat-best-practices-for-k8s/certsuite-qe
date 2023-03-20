package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	goclient "sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/gomega"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

// DefineOperatorGroup returns an operator group struct.
func DefineOperatorGroup(groupName string, namespace string, targetNamespace []string) *olmv1.OperatorGroup {
	return &olmv1.OperatorGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      groupName,
			Namespace: namespace},
		Spec: olmv1.OperatorGroupSpec{
			TargetNamespaces: targetNamespace},
	}
}

// DefineSubscription returns a subscription struct.
func DefineSubscription(subName, namespace, channel, operatorName, catalogSource,
	catalogSourceNamespace, startingCSV string, installPlanApproval v1alpha1.Approval) *v1alpha1.Subscription {
	return &v1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subName,
			Namespace: namespace,
		},
		Spec: &v1alpha1.SubscriptionSpec{
			Channel:                channel,
			Package:                operatorName,
			CatalogSource:          catalogSource,
			CatalogSourceNamespace: catalogSourceNamespace,
			StartingCSV:            startingCSV,
			InstallPlanApproval:    installPlanApproval,
		},
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

func DeployOperatorSubscription(operatorPackage, chanel, namespace, group,
	sourceNamespace, startingCSV string, installApproval v1alpha1.Approval) error {
	operatorSubscription := DefineSubscription(
		operatorPackage+"-subscription",
		namespace,
		chanel,
		operatorPackage,
		group,
		sourceNamespace,
		startingCSV,
		installApproval)

	err := DeployOperator(operatorSubscription)

	if err != nil {
		return fmt.Errorf("Error deploying operator "+operatorPackage+": %w", err)
	}

	return nil
}

func DeployOperatorGroup(namespace string, operatorGroup *olmv1.OperatorGroup) error {
	err := namespaces.Create(namespace, globalhelper.APIClient)
	if err != nil {
		return err
	}

	err = globalhelper.APIClient.Create(context.TODO(),
		&olmv1.OperatorGroup{
			ObjectMeta: metav1.ObjectMeta{
				Name:      operatorGroup.Name,
				Namespace: operatorGroup.Namespace},
			Spec: olmv1.OperatorGroupSpec{
				TargetNamespaces: operatorGroup.Spec.TargetNamespaces},
		},
	)

	if err != nil {
		return fmt.Errorf("can not deploy operatorGroup %w", err)
	}

	return nil
}

func DeployOperator(subscription *v1alpha1.Subscription) error {
	err := globalhelper.APIClient.Create(context.TODO(),
		&v1alpha1.Subscription{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subscription.Name,
				Namespace: subscription.Namespace,
			},
			Spec: &v1alpha1.SubscriptionSpec{
				Channel:                subscription.Spec.Channel,
				Package:                subscription.Spec.Package,
				CatalogSource:          subscription.Spec.CatalogSource,
				CatalogSourceNamespace: subscription.Spec.CatalogSourceNamespace,
				StartingCSV:            subscription.Spec.StartingCSV,
				InstallPlanApproval:    subscription.Spec.InstallPlanApproval,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("can not install Subscription %w", err)
	}

	return nil
}

func IsOperatorGroupInstalled(operatorGroupName, namespace string) error {
	var operatorGroup olmv1.OperatorGroup
	err := globalhelper.APIClient.Get(context.TODO(),
		goclient.ObjectKey{Name: operatorGroupName, Namespace: namespace},
		&operatorGroup)

	if err != nil {
		return fmt.Errorf(operatorGroupName+" operatorGroup resource not found: %w", err)
	}

	return nil
}

func GetInstallPlanByCSV(namespace string, csvName string) (*v1alpha1.InstallPlan, error) {
	installPlans, err := globalhelper.APIClient.InstallPlans(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("unable to get InstallPlans: %w", err)
	}

	var matchingPlan v1alpha1.InstallPlan

	for _, plan := range installPlans.Items {
		for _, csv := range plan.Spec.ClusterServiceVersionNames {
			if strings.Contains(csv, csvName) {
				matchingPlan = plan

				break
			}
		}
	}

	if matchingPlan.Name == "" {
		return nil, fmt.Errorf("failed to detect InstallPlan")
	}

	return &matchingPlan, nil
}

func ApproveInstallPlan(namespace string, plan *v1alpha1.InstallPlan) error {
	plan.Spec.Approved = true

	return updateInstallPlan(namespace, plan)
}

func DeployRHCertifiedOperatorSource(ocpVersion string) error {
	err := globalhelper.APIClient.Create(context.TODO(),
		&v1alpha1.CatalogSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      tsparams.CertifiedOperatorGroup,
				Namespace: tsparams.OperatorSourceNamespace,
			},
			Spec: v1alpha1.CatalogSourceSpec{
				SourceType:  "grpc",
				Image:       "registry.redhat.io/redhat/certified-operator-index:v" + ocpVersion,
				DisplayName: "redhat-certified",
				Publisher:   "Redhat",
				Secrets:     []string{"redhat-registry-secret", "redhat-connect-registry-secret"},
				UpdateStrategy: &v1alpha1.UpdateStrategy{
					RegistryPoll: &v1alpha1.RegistryPoll{
						Interval: &metav1.Duration{
							Duration: 30 * time.Minute}},
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("can not deploy catalog source %w", err)
	}

	return nil
}

func DisableCatalogSource(name string) error {
	return setCatalogSource(true, name)
}

func EnableCatalogSource(name string) error {
	return setCatalogSource(false, name)
}

func IsCatalogSourceEnabled(name, namespace, displayName string) (bool, error) {
	source, err := globalhelper.APIClient.CatalogSources(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false, nil
	}

	return source.Spec.DisplayName == displayName, nil
}

func DeleteCatalogSource(name, namespace, displayName string) error {
	source, err := globalhelper.APIClient.CatalogSources(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}

		return err
	}

	if source.Spec.DisplayName == displayName {
		return globalhelper.APIClient.Delete(context.TODO(), source)
	}

	return nil
}

func setCatalogSource(disable bool, name string) error {
	_, err := globalhelper.APIClient.OcpClientInterface.OperatorHubs().Patch(context.TODO(),
		"cluster",
		types.MergePatchType,
		[]byte("{\"spec\":{\"sources\":[{\"disabled\": "+strconv.FormatBool(disable)+",\"name\": \""+name+"\"}]}}"),
		metav1.PatchOptions{},
	)

	if err != nil {
		return fmt.Errorf("unable to alter catalog source: %w", err)
	}

	return nil
}

func updateInstallPlan(namespace string, plan *v1alpha1.InstallPlan) error {
	_, err := globalhelper.APIClient.InstallPlans(namespace).Update(
		context.TODO(), plan, metav1.UpdateOptions{},
	)

	if err != nil {
		return fmt.Errorf("failed to update InstallPlan: %w", err)
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
				csv.Status.Phase != "Replacing" &&
				csv.Status.Phase != "Unknown"
		}

		return false
	}, 5*time.Minute, tsparams.PollingInterval).Should(Equal(true),
		//}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
		csvPrefix+" is not ready.")

	return err
}

func ApproveInstallPlanWhenReady(csvName, namespace string) {
	Eventually(func() bool {
		installPlan, err := GetInstallPlanByCSV(namespace, csvName)
		if err != nil {
			return false
		}

		if installPlan.Spec.Approval == v1alpha1.ApprovalAutomatic {
			return true
		}

		if installPlan.Status.Phase == v1alpha1.InstallPlanPhaseRequiresApproval {
			err = ApproveInstallPlan(tsparams.TestCertificationNameSpace,
				installPlan)

			return err == nil
		}

		return false
	}, 5*time.Minute, tsparams.PollingInterval).Should(Equal(true),
		csvName+" install plan is not ready.")
}

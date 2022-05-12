package affiliatedcerthelper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	goclient "sigs.k8s.io/controller-runtime/pkg/client"
)

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

func DeployOperator(namespace string, subscription *v1alpha1.Subscription) error {
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

func ApproveInstallPlan(namespace string, csvPrefix string) error {
	plan, err := GetInstallPlanByCSV(namespace, csvPrefix)

	if err != nil {
		return err
	}

	plan.Spec.Approved = true

	return updateInstallPlan(namespace, plan)
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

// IsOperatorInstalled validates if the given operator is deployed on the given cluster.
func IsOperatorInstalled(namespace string, csvPrefix string) error {
	glog.V(5).Info(fmt.Sprintf("Validate that operator namespace: %s exists", namespace))

	namespaceExists, err := namespaces.Exists(namespace, globalhelper.APIClient)
	if !namespaceExists && err == nil {
		return fmt.Errorf("operator namespace %s doesn't exist", namespace)
	}

	glog.V(5).Info(fmt.Sprintf("Validate that operator's csv %s exists", csvPrefix))
	_, err = getCsvByPrefix(csvPrefix, namespace)

	if err != nil {
		return fmt.Errorf("%s operator's CSV is not installed", csvPrefix)
	}

	return nil
}

func DeployRHCertifiedOperatorSource(ocpVersion string) error {
	err := globalhelper.APIClient.Create(context.TODO(),
		&v1alpha1.CatalogSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "certified-operators",
				Namespace: "openshift-marketplace",
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

func setCatalogSource(disable bool, name string) error {
	_, err := globalhelper.APIClient.OperatorHubs().Patch(context.TODO(),
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

func DisableCatalogSource(name string) error {
	return setCatalogSource(true, name)
}

func EnableCatalogSource(name string) error {
	return setCatalogSource(false, name)
}

func IsCatalogSourceEnabled(name, namespace string) bool {
	_, err := globalhelper.APIClient.CatalogSources(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	return err == nil
}

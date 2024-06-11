package globalhelper

import (
	"context"
	"fmt"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	v1alpha1typed "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CatalogSourceNamespace = "openshift-marketplace"
)

func ValidateCatalogSources() error {
	return validateCatalogSources(GetAPIClient().OperatorsV1alpha1Interface)
}

func validateCatalogSources(opclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	validCatalogSources := []string{"certified-operators", "community-operators"}

	catalogSources, err := opclient.CatalogSources(CatalogSourceNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(catalogSources.Items) == 0 {
		return fmt.Errorf("no catalog sources found")
	}

	var foundCatalogSources []string
	for _, catalogSource := range catalogSources.Items {
		foundCatalogSources = append(foundCatalogSources, catalogSource.Name)
	}

	for _, validCatalogSource := range validCatalogSources {
		if !contains(foundCatalogSources, validCatalogSource) {
			return fmt.Errorf("catalog source %s not found", validCatalogSource)
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func createCatalogSource(name, url string) error {
	return createCatalogSourceWithClient(GetAPIClient().OperatorsV1alpha1Interface, name, url)
}

func createCatalogSourceWithClient(opclient v1alpha1typed.OperatorsV1alpha1Interface, name, url string) error {
	_, err := opclient.CatalogSources(CatalogSourceNamespace).Create(context.Background(), &v1alpha1.CatalogSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "openshift-marketplace",
		},
		Spec: v1alpha1.CatalogSourceSpec{
			SourceType:  "grpc",
			Image:       url,
			Publisher:   "Red Hat",
			DisplayName: name,
		},
	}, metav1.CreateOptions{})

	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func CreateCommunityOperatorsCatalogSource() error {
	communityOperatorIndex414 := "registry.redhat.io/redhat/community-operator-index:v4.14"
	communityOperatorIndex415 := "registry.redhat.io/redhat/community-operator-index:v4.15"
	communityOperatorIndex416 := "registry.redhat.io/redhat/community-operator-index:v4.16"

	// determine which index to use based on ocp version
	ocpVersion, err := GetClusterVersion()
	if err != nil {
		return err
	}

	// strip off the patch version (e.g. 4.14.0 -> 4.14)
	majorMinor := ocpVersion[:4]

	// Note: Update this when new OCP versions are released and new community operator indexes are available
	switch majorMinor {
	case "4.14":
		return createCatalogSource("community-operators", communityOperatorIndex414)
	case "4.15":
		return createCatalogSource("community-operators", communityOperatorIndex415)
	case "4.16":
		return createCatalogSource("community-operators", communityOperatorIndex416)
	default:
		return fmt.Errorf("unsupported ocp version %s", ocpVersion)
	}
}

func GetClusterVersion() (string, error) {
	clusterVersion, err := GetAPIClient().ClusterVersionInterface.
		Get(context.Background(), "version", metav1.GetOptions{})

	if err != nil {
		return "", err
	}

	return clusterVersion.Status.Desired.Version, nil
}

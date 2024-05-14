package globalhelper

import (
	"context"
	"fmt"

	v1alpha1typed "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
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

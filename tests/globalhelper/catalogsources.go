package globalhelper

import (
	"context"
	"fmt"
	"slices"
	"time"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiClusterVersion "github.com/openshift-kni/eco-goinfra/pkg/clusterversion"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	v1alpha1typed "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klog "k8s.io/klog/v2"
)

const (
	CatalogSourceNamespace = "openshift-marketplace"
)

func ValidateCatalogSources() error {
	return validateCatalogSources(GetAPIClient().OperatorsV1alpha1Interface)
}

func validateCatalogSources(opclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	requiredCatalogSources := []string{"certified-operators", "community-operators"}

	const (
		timeout  = 5 * time.Minute
		interval = 10 * time.Second
	)

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		catalogSources, err := opclient.CatalogSources(
			CatalogSourceNamespace).List(context.TODO(),
			metav1.ListOptions{})
		if err != nil {
			return err
		}

		allReady := true

		for _, name := range requiredCatalogSources {
			idx := slices.IndexFunc(catalogSources.Items, func(cs v1alpha1.CatalogSource) bool {
				return cs.Name == name
			})

			if idx == -1 {
				klog.Infof("Catalog source %s not found yet, waiting...", name)
				allReady = false

				break
			}

			cs := catalogSources.Items[idx]
			if cs.Status.GRPCConnectionState == nil || cs.Status.GRPCConnectionState.LastObservedState != "READY" {
				state := "nil"
				if cs.Status.GRPCConnectionState != nil {
					state = cs.Status.GRPCConnectionState.LastObservedState
				}

				klog.Infof("Catalog source %s exists but is not READY (state: %s), waiting...", name, state)
				allReady = false

				break
			}

			klog.Infof("Catalog source %s is READY", name)
		}

		if allReady {
			return nil
		}

		time.Sleep(interval)
	}

	return fmt.Errorf("timed out after %s waiting for catalog sources %v to be READY", timeout, requiredCatalogSources)
}

func deleteCatalogSourceByName(name string) error {
	err := GetAPIClient().OperatorsV1alpha1Interface.CatalogSources(CatalogSourceNamespace).Delete(
		context.TODO(), name, metav1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete catalog source %s: %w", name, err)
	}

	return nil
}

func CreateAndValidateCatalogSources(includeRedHat bool) error {
	const maxAttempts = 3

	sources := map[string]func() error{
		"community-operators": CreateCommunityOperatorsCatalogSource,
		"certified-operators": func() error { return DeployRHCertifiedOperatorSource("") },
	}

	if includeRedHat {
		sources["redhat-operators"] = func() error { return DeployRHOperatorSource("") }
	}

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		for name, create := range sources {
			if err := create(); err != nil {
				return fmt.Errorf("failed to create catalog source %s: %w", name, err)
			}
		}

		lastErr = ValidateCatalogSources()
		if lastErr == nil {
			return nil
		}

		klog.Infof("Catalog source validation failed (attempt %d/%d): %v", attempt, maxAttempts, lastErr)

		if attempt == maxAttempts {
			break
		}

		for name := range sources {
			if delErr := deleteCatalogSourceByName(name); delErr != nil {
				klog.Infof("Warning: failed to delete catalog source %s: %v", name, delErr)
			}
		}
	}

	return fmt.Errorf("catalog sources not ready after %d attempts: %w", maxAttempts, lastErr)
}

func createCatalogSource(name, url string) error {
	return createCatalogSourceWithClient(GetAPIClient().OperatorsV1alpha1Interface, name, url)
}

func createCatalogSourceWithClient(opclient v1alpha1typed.OperatorsV1alpha1Interface, name, url string) error {
	_, err := opclient.CatalogSources(CatalogSourceNamespace).Create(context.TODO(), &v1alpha1.CatalogSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: CatalogSourceNamespace,
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
	// determine which index to use based on ocp version
	ocpVersion, err := GetClusterVersion()
	if err != nil {
		return err
	}

	return createCatalogSource("community-operators",
		"registry.redhat.io/redhat/community-operator-index:v"+ocpVersion[:4])
}

func GetClusterVersion() (string, error) {
	client := egiClients.New("")

	builder, err := egiClusterVersion.Pull(client)

	if err != nil {
		return "", err
	}

	return builder.Object.Status.Desired.Version, nil
}

// GetClusterVersionOrDefault returns the OCP cluster version, or a default value if unavailable.
// This is useful for tests that need the version but can fall back to a reasonable default
// (e.g., for kind clusters or when the version cannot be determined).
// The default value is "4.14" which uses the standard operator set.
func GetClusterVersionOrDefault() string {
	version, err := GetClusterVersion()
	if err != nil {
		// Default to 4.14 which uses the standard operator set
		return "4.14"
	}

	return version
}

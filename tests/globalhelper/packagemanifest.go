package globalhelper

import (
	"fmt"
	"strings"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiOLM "github.com/openshift-kni/eco-goinfra/pkg/olm"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// QueryPackageManifest queries the package manifest for the operator.
func QueryPackageManifestForVersion(operatorName, operatorNamespace, channel string) (string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", err
	}

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			fmt.Printf("Comparing %s with %s\n", item.Object.Status.DefaultChannel, channel)

			// skip if the default channel is not the one we are looking for
			if item.Object.Status.DefaultChannel != channel {
				continue
			}

			for _, chanObj := range item.Object.Status.Channels {
				fmt.Printf("Comparing name %s with channel %s\n", chanObj.Name, channel)

				// skip if the name of the channel is not the one we are looking for
				if chanObj.Name != channel {
					continue
				}

				return chanObj.CurrentCSVDesc.Version.String(), nil
			}
		}
	}

	return "not found", nil
}

func QueryPackageManifestForDefaultChannel(operatorName, operatorNamespace string) (string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", err
	}

	// Example of how to get the default channel
	// oc get packagemanifest cluster-logging -o json | jq .status.defaultChannel

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			return item.Object.Status.DefaultChannel, nil
		}
	}

	return "not found", nil
}

// QueryPackageManifestForOperatorName searches for packages containing the specified search string
// and returns the first matching package name. This is useful for finding operators whose
// package names may vary between OCP versions (e.g., "cloudbees-ci" vs "cloudbees-ci-rhmp").
func QueryPackageManifestForOperatorName(searchString, operatorNamespace string) (string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", err
	}

	for _, item := range pkgManifest {
		if strings.Contains(item.Object.GetName(), searchString) {
			fmt.Printf("Found package: %s matching search string: %s\n", item.Object.GetName(), searchString)

			return item.Object.GetName(), nil
		}
	}

	return "not found", nil
}

// QueryPackageManifestForOperatorNameAndCatalogSource searches for packages containing the specified search string
// and returns both the package name and the catalog source where it was found. This is useful for finding operators whose
// package names and catalog sources may vary between OCP versions (e.g., "ovms-operator" in "certified-operators"
// vs "ovms-operator-rhmp" in "redhat-marketplace").
func QueryPackageManifestForOperatorNameAndCatalogSource(searchString, operatorNamespace string) (string, string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", "", err
	}

	for _, item := range pkgManifest {
		if strings.Contains(item.Object.GetName(), searchString) {
			packageName := item.Object.GetName()
			catalogSource := item.Object.Status.CatalogSource
			fmt.Printf("Found package: %s in catalog source: %s matching search string: %s\n", packageName, catalogSource, searchString)

			return packageName, catalogSource, nil
		}
	}

	return "not found", "not found", nil
}

// QueryPackageManifestForAvailableChannelAndVersion searches for an operator and returns the first available channel
// that has versions, along with a version from that channel. This is more robust than requiring a specific channel.
func QueryPackageManifestForAvailableChannelAndVersion(operatorName, operatorNamespace string) (string, string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", "", err
	}

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			// Try the default channel first
			defaultChannel := item.Object.Status.DefaultChannel
			fmt.Printf("Default channel for %s: %s\n", operatorName, defaultChannel)

			for _, chanObj := range item.Object.Status.Channels {
				fmt.Printf("Checking channel %s for %s\n", chanObj.Name, operatorName)

				// Check if this channel has a version available
				if chanObj.CurrentCSVDesc.Version.String() != "" {
					channelName := chanObj.Name
					version := chanObj.CurrentCSVDesc.Version.String()
					fmt.Printf("Found available channel: %s with version: %s for %s\n", channelName, version, operatorName)

					return channelName, version, nil
				}
			}
		}
	}

	return "not found", "not found", nil
}

// QueryPackageManifestForAvailableChannelVersionAndCSV searches for an operator and returns the first available channel
// that has versions, along with the version and the actual CSV name from that channel. This handles cases where
// the CSV name doesn't match the package name pattern (e.g., package "ovms-operator-rhmp" has CSV "openvino-operator.v1.2.0").
func QueryPackageManifestForAvailableChannelVersionAndCSV(operatorName, operatorNamespace string) (string, string, string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", "", "", err
	}

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			// Try the default channel first
			defaultChannel := item.Object.Status.DefaultChannel
			fmt.Printf("Default channel for %s: %s\n", operatorName, defaultChannel)

			for _, chanObj := range item.Object.Status.Channels {
				fmt.Printf("Checking channel %s for %s\n", chanObj.Name, operatorName)

				// Check if this channel has a version available
				if chanObj.CurrentCSVDesc.Version.String() != "" && chanObj.CurrentCSV != "" {
					channelName := chanObj.Name
					version := chanObj.CurrentCSVDesc.Version.String()
					csvName := chanObj.CurrentCSV
					fmt.Printf("Found available channel: %s with version: %s and CSV: %s for %s\n", channelName, version, csvName, operatorName)

					return channelName, version, csvName, nil
				}
			}
		}
	}

	return "not found", "not found", "not found", nil
}

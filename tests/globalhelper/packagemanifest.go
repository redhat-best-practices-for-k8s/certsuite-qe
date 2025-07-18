package globalhelper

import (
	"fmt"
	"sort"
	"strings"

	. "github.com/onsi/ginkgo/v2"
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

// QueryPackageManifestForOperatorNameAndCatalogSourceWithPreference searches for packages containing the specified search string
// and returns the package name and catalog source, giving preference to specified catalog sources.
// This avoids non-deterministic behavior when multiple packages match the search string.
func QueryPackageManifestForOperatorNameAndCatalogSourceWithPreference(searchString, operatorNamespace string, preferredCatalogSources []string) (string, string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", "", err
	}

	// First, collect all matches
	var matches []struct {
		packageName   string
		catalogSource string
		priority      int // lower number = higher priority
	}

	for _, item := range pkgManifest {
		if strings.Contains(item.Object.GetName(), searchString) {
			packageName := item.Object.GetName()
			catalogSource := item.Object.Status.CatalogSource

			// Assign priority based on preferred catalog sources
			priority := 1000 // default low priority
			for i, preferred := range preferredCatalogSources {
				if catalogSource == preferred {
					priority = i // higher preference = lower number
					break
				}
			}

			matches = append(matches, struct {
				packageName   string
				catalogSource string
				priority      int
			}{packageName, catalogSource, priority})

			fmt.Printf("Found package: %s in catalog source: %s (priority: %d) matching search string: %s\n",
				packageName, catalogSource, priority, searchString)
		}
	}

	if len(matches) == 0 {
		return "not found", "not found", nil
	}

	// Sort by priority (lowest number first), then by package name for deterministic behavior
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].priority != matches[j].priority {
			return matches[i].priority < matches[j].priority
		}
		return matches[i].packageName < matches[j].packageName
	})

	selectedMatch := matches[0]
	fmt.Printf("Selected package: %s from catalog source: %s (priority: %d)\n",
		selectedMatch.packageName, selectedMatch.catalogSource, selectedMatch.priority)

	if len(matches) > 1 {
		fmt.Printf("Note: Found %d matches, selected based on catalog source preference\n", len(matches))
	}

	return selectedMatch.packageName, selectedMatch.catalogSource, nil
}

// DiagnosePackageManifestMatches lists all packages that match a search string along with their catalog sources.
// This is useful for debugging when multiple packages exist and understanding catalog source availability.
func DiagnosePackageManifestMatches(searchString, operatorNamespace string) error {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return fmt.Errorf("error listing package manifests: %w", err)
	}

	fmt.Printf("\n=== Package Manifest Diagnosis for search string: %s ===\n", searchString)

	var matches []struct {
		packageName   string
		catalogSource string
	}

	for _, item := range pkgManifest {
		if strings.Contains(item.Object.GetName(), searchString) {
			packageName := item.Object.GetName()
			catalogSource := item.Object.Status.CatalogSource
			matches = append(matches, struct {
				packageName   string
				catalogSource string
			}{packageName, catalogSource})
		}
	}

	if len(matches) == 0 {
		fmt.Printf("No packages found matching '%s'\n", searchString)
		return nil
	}

	fmt.Printf("Found %d packages matching '%s':\n", len(matches), searchString)
	for i, match := range matches {
		fmt.Printf("  %d. Package: %-20s | Catalog Source: %s\n", i+1, match.packageName, match.catalogSource)
	}

	if len(matches) > 1 {
		fmt.Printf("\nWarning: Multiple packages found. This could cause non-deterministic behavior.\n")
		fmt.Printf("Consider using QueryPackageManifestForOperatorNameAndCatalogSourceWithPreference() with explicit preferences.\n")
	}

	fmt.Printf("=== End Diagnosis ===\n\n")

	return nil
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

// CheckOperatorExistsOrFail checks if an operator exists in cluster packagemanifests.
// If not found, it logs the issue and fails the test using Ginkgo's Fail().
// If the operator exists, it returns the package name and catalog source for further use.
func CheckOperatorExistsOrFail(searchString, operatorNamespace string) (string, string) {
	operatorName, catalogSource, err := QueryPackageManifestForOperatorNameAndCatalogSource(searchString, operatorNamespace)
	if err != nil {
		Fail(fmt.Sprintf("Error querying package manifest for %s operator: %s", searchString, err.Error()))
	}

	if operatorName == "not found" || catalogSource == "not found" {
		Fail(fmt.Sprintf("Operator %s not found in cluster packagemanifests. This indicates a problem with catalog sources.", searchString))
	}

	fmt.Printf("Operator %s found as package %s in catalog source %s\n", searchString, operatorName, catalogSource)

	return operatorName, catalogSource
}

// CheckOperatorChannelAndVersionOrFail checks if an operator has available channel and version.
// If not found, it logs the issue and fails the test using Ginkgo's Fail().
// If found, it returns the channel, version, and CSV name for further use.
func CheckOperatorChannelAndVersionOrFail(operatorName, operatorNamespace string) (string, string, string) {
	channel, version, csvName, err := QueryPackageManifestForAvailableChannelVersionAndCSV(operatorName, operatorNamespace)
	if err != nil {
		Fail(fmt.Sprintf("Error querying package manifest for %s operator channel and version: %s", operatorName, err.Error()))
	}

	if channel == "not found" || version == "not found" || csvName == "not found" {
		Fail(fmt.Sprintf("Operator %s channel, version, or CSV not found in packagemanifests. "+
			"This indicates a problem with catalog sources.", operatorName))
	}

	fmt.Printf("Operator %s has available channel %s, version %s, CSV %s\n", operatorName, channel, version, csvName)

	return channel, version, csvName
}

// QueryPackageManifestForOperatorFromSpecificCatalogSource searches for an operator package from a specific catalog source only.
// This ensures deterministic behavior and fails clearly if the operator is not available from the expected catalog source.
func QueryPackageManifestForOperatorFromSpecificCatalogSource(searchString, operatorNamespace, requiredCatalogSource string) (string, string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", "", err
	}

	var matches []struct {
		packageName   string
		catalogSource string
	}

	// Collect all matches for diagnostics
	for _, item := range pkgManifest {
		if strings.Contains(item.Object.GetName(), searchString) {
			packageName := item.Object.GetName()
			catalogSource := item.Object.Status.CatalogSource
			matches = append(matches, struct {
				packageName   string
				catalogSource string
			}{packageName, catalogSource})
		}
	}

	// Show all available matches for debugging
	if len(matches) > 0 {
		fmt.Printf("Found %d package(s) matching '%s':\n", len(matches), searchString)
		for i, match := range matches {
			marker := ""
			if match.catalogSource == requiredCatalogSource {
				marker = " ← TARGET"
			}
			fmt.Printf("  %d. Package: %-20s | Catalog Source: %s%s\n", i+1, match.packageName, match.catalogSource, marker)
		}
	}

	// Look for a match from the required catalog source
	for _, match := range matches {
		if match.catalogSource == requiredCatalogSource {
			fmt.Printf("✓ Found required package: %s from catalog source: %s\n", match.packageName, match.catalogSource)
			return match.packageName, match.catalogSource, nil
		}
	}

	// If we get here, the required catalog source doesn't have the operator
	if len(matches) == 0 {
		return "", "", fmt.Errorf("no packages found matching '%s' in any catalog source", searchString)
	}

	availableCatalogs := make([]string, 0, len(matches))
	for _, match := range matches {
		found := false
		for _, existing := range availableCatalogs {
			if existing == match.catalogSource {
				found = true
				break
			}
		}
		if !found {
			availableCatalogs = append(availableCatalogs, match.catalogSource)
		}
	}

	return "", "", fmt.Errorf("operator '%s' not found in required catalog source '%s'. Available in: %v",
		searchString, requiredCatalogSource, availableCatalogs)
}

package globalhelper

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	egiOLM "github.com/openshift-kni/eco-goinfra/pkg/olm"
	klog "k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// QueryPackageManifestForVersion returns the version string for the given operator and channel,
// or "not found" if no matching package/channel exists.
func QueryPackageManifestForVersion(operatorName, operatorNamespace, channel string) (string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(GetEcoGoinfraClient(), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", err
	}

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			klog.InfoS("Comparing default channel with requested channel",
				"defaultChannel", item.Object.Status.DefaultChannel, "requestedChannel", channel)

			if item.Object.Status.DefaultChannel != channel {
				continue
			}

			for _, chanObj := range item.Object.Status.Channels {
				klog.InfoS("Comparing channel name with requested channel",
					"channelName", chanObj.Name, "requestedChannel", channel)

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
	pkgManifest, err := egiOLM.ListPackageManifest(GetEcoGoinfraClient(), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", err
	}

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
	pkgManifest, err := egiOLM.ListPackageManifest(GetEcoGoinfraClient(), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", err
	}

	for _, item := range pkgManifest {
		if strings.Contains(item.Object.GetName(), searchString) {
			klog.InfoS("Found matching package", "package", item.Object.GetName(), "searchString", searchString)

			return item.Object.GetName(), nil
		}
	}

	return "not found", nil
}

// QueryPackageManifestForOperatorNameAndCatalogSource searches for packages whose name starts with
// the specified search string and returns both the package name and catalog source. This handles
// operators whose names and catalog sources vary between OCP versions.
func QueryPackageManifestForOperatorNameAndCatalogSource(searchString, operatorNamespace string) (string, string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(GetEcoGoinfraClient(), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", "", err
	}

	for _, item := range pkgManifest {
		if strings.HasPrefix(item.Object.GetName(), searchString) {
			packageName := item.Object.GetName()
			catalogSource := item.Object.Status.CatalogSource
			klog.InfoS("Found package in catalog source",
				"package", packageName, "catalogSource", catalogSource, "searchString", searchString)

			return packageName, catalogSource, nil
		}
	}

	return "not found", "not found", nil
}

// QueryPackageManifestForAvailableChannelAndVersion searches for an operator and returns the first available channel
// that has versions, along with a version from that channel. This is more robust than requiring a specific channel.
func QueryPackageManifestForAvailableChannelAndVersion(operatorName, operatorNamespace string) (string, string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(GetEcoGoinfraClient(), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", "", err
	}

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			defaultChannel := item.Object.Status.DefaultChannel
			klog.InfoS("Found operator default channel", "operator", operatorName, "defaultChannel", defaultChannel)

			for _, chanObj := range item.Object.Status.Channels {
				klog.InfoS("Checking channel for operator", "channel", chanObj.Name, "operator", operatorName)

				if chanObj.CurrentCSVDesc.Version.String() != "" {
					channelName := chanObj.Name
					version := chanObj.CurrentCSVDesc.Version.String()
					klog.InfoS("Found available channel with version",
						"channel", channelName, "version", version, "operator", operatorName)

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
	pkgManifest, err := egiOLM.ListPackageManifest(GetEcoGoinfraClient(), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", "", "", err
	}

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			defaultChannel := item.Object.Status.DefaultChannel
			klog.InfoS("Found operator default channel", "operator", operatorName, "defaultChannel", defaultChannel)

			for _, chanObj := range item.Object.Status.Channels {
				klog.InfoS("Checking channel for operator", "channel", chanObj.Name, "operator", operatorName)

				if chanObj.CurrentCSVDesc.Version.String() != "" && chanObj.CurrentCSV != "" {
					channelName := chanObj.Name
					version := chanObj.CurrentCSVDesc.Version.String()
					csvName := chanObj.CurrentCSV
					klog.InfoS("Found available channel with version and CSV",
						"channel", channelName, "version", version, "csv", csvName, "operator", operatorName)

					return channelName, version, csvName, nil
				}
			}
		}
	}

	return "not found", "not found", "not found", nil
}

// checkOperatorExists is the shared implementation for CheckOperatorExistsOrFail and
// CheckOperatorExistsOrSkip. The failHandler is called when the operator query fails
// or the operator is not found (typically Ginkgo's Fail or Skip).
func checkOperatorExists(searchString, operatorNamespace string, failHandler func(string, ...int)) (string, string) {
	operatorName, catalogSource, err := QueryPackageManifestForOperatorNameAndCatalogSource(searchString, operatorNamespace)
	if err != nil {
		failHandler(fmt.Sprintf("Error querying package manifest for %s operator: %s", searchString, err.Error()))
	}

	if operatorName == "not found" || catalogSource == "not found" {
		failHandler(fmt.Sprintf("Operator %s not found in cluster packagemanifests", searchString))
	}

	klog.InfoS("Operator found in cluster packagemanifests",
		"searchString", searchString, "package", operatorName, "catalogSource", catalogSource)

	return operatorName, catalogSource
}

// CheckOperatorExistsOrFail checks if an operator exists in cluster packagemanifests.
// If not found, it fails the test using Ginkgo's Fail().
func CheckOperatorExistsOrFail(searchString, operatorNamespace string) (string, string) {
	return checkOperatorExists(searchString, operatorNamespace, Fail)
}

// CheckOperatorExistsOrSkip checks if an operator exists in cluster packagemanifests.
// If not found, it skips the test using Ginkgo's Skip(). Use this for operators that
// may not be available in all OCP versions or catalog configurations.
func CheckOperatorExistsOrSkip(searchString, operatorNamespace string) (string, string) {
	return checkOperatorExists(searchString, operatorNamespace, Skip)
}

// checkOperatorChannelAndVersion is the shared implementation for
// CheckOperatorChannelAndVersionOrFail and CheckOperatorChannelAndVersionOrSkip.
func checkOperatorChannelAndVersion(operatorName, operatorNamespace string,
	failHandler func(string, ...int)) (string, string, string) {
	channel, version, csvName, err := QueryPackageManifestForAvailableChannelVersionAndCSV(operatorName, operatorNamespace)
	if err != nil {
		failHandler(fmt.Sprintf("Error querying package manifest for %s operator channel and version: %s",
			operatorName, err.Error()))
	}

	if channel == "not found" || version == "not found" || csvName == "not found" {
		failHandler(fmt.Sprintf("Operator %s channel, version, or CSV not found in packagemanifests", operatorName))
	}

	klog.InfoS("Operator channel and version found",
		"operator", operatorName, "channel", channel, "version", version, "csv", csvName)

	return channel, version, csvName
}

// CheckOperatorChannelAndVersionOrFail checks if an operator has an available channel and version.
// If not found, it fails the test using Ginkgo's Fail().
func CheckOperatorChannelAndVersionOrFail(operatorName, operatorNamespace string) (string, string, string) {
	return checkOperatorChannelAndVersion(operatorName, operatorNamespace, Fail)
}

// CheckOperatorChannelAndVersionOrSkip checks if an operator has an available channel and version.
// If not found, it skips the test using Ginkgo's Skip(). Use this for operators that
// may have catalog issues in certain OCP versions.
func CheckOperatorChannelAndVersionOrSkip(operatorName, operatorNamespace string) (string, string, string) {
	return checkOperatorChannelAndVersion(operatorName, operatorNamespace, Skip)
}

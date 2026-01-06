// Package operatorversions provides a centralized mapping of OCP versions to operators
// that should be used in QE tests. This addresses the problem where certain operators
// are not available in all OCP versions (e.g., cockroachdb-certified is not in 4.20).
//
// Usage:
//
//	config := operatorversions.GetOperatorConfig("4.20")
//	certifiedOp := config.CertifiedOperator
//	// Use certifiedOp.PackageName, certifiedOp.CatalogSource, etc.
package operatorversions

import (
	"fmt"
	"strings"
)

// CatalogSource constants for operator catalogs.
const (
	CatalogCertifiedOperators = "certified-operators"
	CatalogCommunityOperators = "community-operators"
	CatalogRedHatOperators    = "redhat-operators"
	CatalogSourceNamespace    = "openshift-marketplace"
)

// OperatorInfo contains the configuration for a specific operator.
type OperatorInfo struct {
	// PackageName is the name of the operator package in the catalog (e.g., "cockroachdb-certified")
	PackageName string

	// CatalogSource is the catalog containing this operator (e.g., "certified-operators")
	CatalogSource string

	// CSVPrefix is the prefix used in the ClusterServiceVersion name (e.g., "cockroach-operator")
	// This is used when labeling or waiting for the operator to be ready
	CSVPrefix string

	// Description provides context about what this operator is used for in tests
	Description string
}

// OCPOperatorConfig contains all operators configured for a specific OCP version.
type OCPOperatorConfig struct {
	// OCPVersion is the OCP version this config applies to (e.g., "4.20")
	OCPVersion string

	// CertifiedOperator is the certified operator to use for tests requiring a certified operator
	// This is the operator from certified-operators catalog
	CertifiedOperator OperatorInfo

	// CommunityOperator is the community operator to use for tests
	// This is typically grafana-operator from community-operators catalog
	CommunityOperator OperatorInfo

	// LightweightOperator is a lightweight operator used for various operator tests
	// This is prometheus-exporter-operator from community-operators catalog
	LightweightOperator OperatorInfo

	// ClusterLoggingOperator is the cluster-logging operator for cluster-wide tests
	// This is from redhat-operators catalog
	ClusterLoggingOperator OperatorInfo

	// UncertifiedOperator is an uncertified operator used for negative tests
	// This is typically cockroachdb (not cockroachdb-certified) from community-operators
	UncertifiedOperator OperatorInfo
}

// DefaultOperatorConfig is the configuration used for OCP versions 4.14 through 4.19
// where all our standard operators are available.
var DefaultOperatorConfig = OCPOperatorConfig{
	OCPVersion: "default",
	CertifiedOperator: OperatorInfo{
		PackageName:   "cockroachdb-certified",
		CatalogSource: CatalogCertifiedOperators,
		CSVPrefix:     "cockroach-operator",
		Description:   "Certified CockroachDB operator for affiliated certification tests",
	},
	CommunityOperator: OperatorInfo{
		PackageName:   "grafana-operator",
		CatalogSource: CatalogCommunityOperators,
		CSVPrefix:     "grafana-operator",
		Description:   "Grafana operator for community operator tests",
	},
	LightweightOperator: OperatorInfo{
		// prometheus-exporter-operator is available in all OCP versions (4.14-4.20+)
		// Note: postgresql was previously used but does NOT exist in any OCP catalog
		// (the catalog check was producing false positives by matching relatedImages)
		PackageName:   "prometheus-exporter-operator",
		CatalogSource: CatalogCommunityOperators,
		CSVPrefix:     "prometheus-exporter-operator",
		Description:   "Prometheus Exporter operator as lightweight operator for various tests",
	},
	ClusterLoggingOperator: OperatorInfo{
		PackageName:   "cluster-logging",
		CatalogSource: CatalogRedHatOperators,
		CSVPrefix:     "cluster-logging",
		Description:   "Cluster logging operator for cluster-wide operator tests",
	},
	UncertifiedOperator: OperatorInfo{
		PackageName:   "cockroachdb",
		CatalogSource: CatalogCommunityOperators,
		CSVPrefix:     "cockroachdb",
		Description:   "Uncertified CockroachDB operator for negative certification tests",
	},
}

// OCP420OperatorConfig is the configuration for OCP 4.20+
// Several operators are NOT available in 4.20 catalogs:
//   - cockroachdb-certified: Missing from certified-operators
//
// Reference: https://github.com/redhat-best-practices-for-k8s/certsuite-qe/issues/1283
// The issue tracks which operators are missing from specific OCP version catalogs.
//
// Alternatives used for OCP 4.20:
//   - mongodb-enterprise: Replaces cockroachdb-certified (certified operator)
var OCP420OperatorConfig = OCPOperatorConfig{
	OCPVersion: "4.20",
	CertifiedOperator: OperatorInfo{
		// mongodb-enterprise is typically available in certified-operators catalog
		// It has a simpler deployment model compared to some other certified operators
		PackageName:   "mongodb-enterprise",
		CatalogSource: CatalogCertifiedOperators,
		CSVPrefix:     "mongodb-enterprise",
		Description:   "MongoDB Enterprise operator (replacement for cockroachdb-certified in 4.20)",
	},
	CommunityOperator: OperatorInfo{
		PackageName:   "grafana-operator",
		CatalogSource: CatalogCommunityOperators,
		CSVPrefix:     "grafana-operator",
		Description:   "Grafana operator for community operator tests",
	},
	LightweightOperator: OperatorInfo{
		// prometheus-exporter-operator is available in all OCP versions (4.14-4.20+)
		PackageName:   "prometheus-exporter-operator",
		CatalogSource: CatalogCommunityOperators,
		CSVPrefix:     "prometheus-exporter-operator",
		Description:   "Prometheus Exporter operator as lightweight operator for various tests",
	},
	ClusterLoggingOperator: OperatorInfo{
		PackageName:   "cluster-logging",
		CatalogSource: CatalogRedHatOperators,
		CSVPrefix:     "cluster-logging",
		Description:   "Cluster logging operator for cluster-wide operator tests",
	},
	UncertifiedOperator: OperatorInfo{
		PackageName:   "cockroachdb",
		CatalogSource: CatalogCommunityOperators,
		CSVPrefix:     "cockroachdb",
		Description:   "Uncertified CockroachDB operator for negative certification tests",
	},
}

// operatorConfigMap maps OCP version prefixes to their operator configurations.
var operatorConfigMap = map[string]*OCPOperatorConfig{
	"4.14": &DefaultOperatorConfig,
	"4.15": &DefaultOperatorConfig,
	"4.16": &DefaultOperatorConfig,
	"4.17": &DefaultOperatorConfig,
	"4.18": &DefaultOperatorConfig,
	"4.19": &DefaultOperatorConfig,
	"4.20": &OCP420OperatorConfig,
}

// GetOperatorConfig returns the operator configuration for the given OCP version.
// The version can be a full version string (e.g., "4.20.0-0.nightly-2024-12-16")
// or a short version (e.g., "4.20"). If no specific config exists for the version,
// the default config is returned.
func GetOperatorConfig(ocpVersion string) *OCPOperatorConfig {
	// Extract major.minor version (e.g., "4.20" from "4.20.0-0.nightly-2024-12-16")
	shortVersion := extractShortVersion(ocpVersion)

	if config, exists := operatorConfigMap[shortVersion]; exists {
		return config
	}

	// Return default config for unknown versions
	return &DefaultOperatorConfig
}

// GetCertifiedOperator returns the certified operator info for the given OCP version.
func GetCertifiedOperator(ocpVersion string) OperatorInfo {
	return GetOperatorConfig(ocpVersion).CertifiedOperator
}

// GetCommunityOperator returns the community operator info for the given OCP version.
func GetCommunityOperator(ocpVersion string) OperatorInfo {
	return GetOperatorConfig(ocpVersion).CommunityOperator
}

// GetLightweightOperator returns the lightweight operator info for the given OCP version.
func GetLightweightOperator(ocpVersion string) OperatorInfo {
	return GetOperatorConfig(ocpVersion).LightweightOperator
}

// GetUncertifiedOperator returns the uncertified operator info for the given OCP version.
func GetUncertifiedOperator(ocpVersion string) OperatorInfo {
	return GetOperatorConfig(ocpVersion).UncertifiedOperator
}

// IsVersion420OrLater checks if the given OCP version is 4.20 or later.
func IsVersion420OrLater(ocpVersion string) bool {
	shortVersion := extractShortVersion(ocpVersion)

	// Parse version components
	var major, minor int

	_, err := fmt.Sscanf(shortVersion, "%d.%d", &major, &minor)
	if err != nil {
		return false
	}

	return major > 4 || (major == 4 && minor >= 20)
}

// extractShortVersion extracts the major.minor version from a full version string.
// Example: "4.20.0-0.nightly-2024-12-16" -> "4.20".
func extractShortVersion(fullVersion string) string {
	// Remove any prefix spaces
	fullVersion = strings.TrimSpace(fullVersion)

	// Split by dots and take first two parts
	parts := strings.Split(fullVersion, ".")
	if len(parts) >= 2 {
		// Handle cases like "4.20.0" or "4.20.0-0.nightly"
		// The minor version might contain a dash for pre-release versions
		minor := strings.Split(parts[1], "-")[0]

		return fmt.Sprintf("%s.%s", parts[0], minor)
	}

	return fullVersion
}

// ListSupportedVersions returns a list of OCP versions that have specific configurations.
func ListSupportedVersions() []string {
	versions := make([]string, 0, len(operatorConfigMap))
	for version := range operatorConfigMap {
		versions = append(versions, version)
	}

	return versions
}

// String returns a human-readable representation of the OperatorInfo.
func (o OperatorInfo) String() string {
	return fmt.Sprintf("%s (catalog: %s, csv-prefix: %s)", o.PackageName, o.CatalogSource, o.CSVPrefix)
}

// String returns a human-readable representation of the OCPOperatorConfig.
func (c OCPOperatorConfig) String() string {
	return fmt.Sprintf("OCPOperatorConfig[%s]: certified=%s, community=%s, lightweight=%s",
		c.OCPVersion,
		c.CertifiedOperator.PackageName,
		c.CommunityOperator.PackageName,
		c.LightweightOperator.PackageName)
}

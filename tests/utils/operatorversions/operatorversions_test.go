package operatorversions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractShortVersion(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "full version with nightly suffix",
			input:    "4.20.0-0.nightly-2024-12-16",
			expected: "4.20",
		},
		{
			name:     "simple major.minor",
			input:    "4.19",
			expected: "4.19",
		},
		{
			name:     "major.minor.patch",
			input:    "4.18.5",
			expected: "4.18",
		},
		{
			name:     "rc version",
			input:    "4.17.0-rc.1",
			expected: "4.17",
		},
		{
			name:     "version with extra spaces",
			input:    "  4.14.2  ",
			expected: "4.14",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := extractShortVersion(testCase.input)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestGetOperatorConfig(t *testing.T) {
	testCases := []struct {
		name                   string
		ocpVersion             string
		expectedCertifiedPkg   string
		expectedCommunityPkg   string
		expectedLightweightPkg string
		expectedUncertifiedPkg string
	}{
		{
			name:                   "OCP 4.14 uses default config",
			ocpVersion:             "4.14",
			expectedCertifiedPkg:   "cockroachdb-certified",
			expectedCommunityPkg:   "grafana-operator",
			expectedLightweightPkg: "prometheus-exporter-operator",
			expectedUncertifiedPkg: "cockroachdb",
		},
		{
			name:                   "OCP 4.19 uses default config",
			ocpVersion:             "4.19.3",
			expectedCertifiedPkg:   "cockroachdb-certified",
			expectedCommunityPkg:   "grafana-operator",
			expectedLightweightPkg: "prometheus-exporter-operator",
			expectedUncertifiedPkg: "cockroachdb",
		},
		{
			name:                   "OCP 4.20 uses 4.20-specific config",
			ocpVersion:             "4.20.0-0.nightly-2024-12-16",
			expectedCertifiedPkg:   "mongodb-enterprise",
			expectedCommunityPkg:   "grafana-operator",
			expectedLightweightPkg: "prometheus-exporter-operator",
			expectedUncertifiedPkg: "cockroachdb",
		},
		{
			name:                   "OCP 4.21 uses 4.21-specific config",
			ocpVersion:             "4.21.0",
			expectedCertifiedPkg:   "mongodb-enterprise",
			expectedCommunityPkg:   "grafana-operator",
			expectedLightweightPkg: "prometheus-exporter-operator",
			expectedUncertifiedPkg: "cockroachdb",
		},
		{
			name:                   "Unknown future version falls back to latest 4.20+ config",
			ocpVersion:             "4.99",
			expectedCertifiedPkg:   "mongodb-enterprise",
			expectedCommunityPkg:   "grafana-operator",
			expectedLightweightPkg: "prometheus-exporter-operator",
			expectedUncertifiedPkg: "cockroachdb",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := GetOperatorConfig(testCase.ocpVersion)

			assert.Equal(t, testCase.expectedCertifiedPkg, config.CertifiedOperator.PackageName)
			assert.Equal(t, testCase.expectedCommunityPkg, config.CommunityOperator.PackageName)
			assert.Equal(t, testCase.expectedLightweightPkg, config.LightweightOperator.PackageName)
			assert.Equal(t, testCase.expectedUncertifiedPkg, config.UncertifiedOperator.PackageName)
		})
	}
}

func TestGetCertifiedOperator(t *testing.T) {
	// 4.19 should use cockroachdb-certified
	certifiedOp := GetCertifiedOperator("4.19")
	assert.Equal(t, "cockroachdb-certified", certifiedOp.PackageName)
	assert.Equal(t, CatalogCertifiedOperators, certifiedOp.CatalogSource)

	// 4.20 should use alternative certified operator (mongodb-enterprise)
	certifiedOp = GetCertifiedOperator("4.20")
	assert.Equal(t, "mongodb-enterprise", certifiedOp.PackageName)
	assert.Equal(t, CatalogCertifiedOperators, certifiedOp.CatalogSource)

	// 4.21 should use mongodb-enterprise
	certifiedOp = GetCertifiedOperator("4.21.0")
	assert.Equal(t, "mongodb-enterprise", certifiedOp.PackageName)
	assert.Equal(t, CatalogCertifiedOperators, certifiedOp.CatalogSource)
}

func TestGetCommunityOperator(t *testing.T) {
	// Should be same across all versions
	communityOp := GetCommunityOperator("4.14")
	assert.Equal(t, "grafana-operator", communityOp.PackageName)

	communityOp = GetCommunityOperator("4.20")
	assert.Equal(t, "grafana-operator", communityOp.PackageName)
}

func TestGetLightweightOperator(t *testing.T) {
	// All versions should use prometheus-exporter-operator (available in all OCP versions)
	lightweightOp := GetLightweightOperator("4.14")
	assert.Equal(t, "prometheus-exporter-operator", lightweightOp.PackageName)

	lightweightOp = GetLightweightOperator("4.17")
	assert.Equal(t, "prometheus-exporter-operator", lightweightOp.PackageName)

	lightweightOp = GetLightweightOperator("4.20")
	assert.Equal(t, "prometheus-exporter-operator", lightweightOp.PackageName)
}

func TestIsVersion420OrLater(t *testing.T) {
	testCases := []struct {
		name     string
		version  string
		expected bool
	}{
		{"4.14 is not 4.20 or later", "4.14", false},
		{"4.19 is not 4.20 or later", "4.19.5", false},
		{"4.20 is 4.20 or later", "4.20", true},
		{"4.20 nightly is 4.20 or later", "4.20.0-0.nightly", true},
		{"4.21 is 4.20 or later", "4.21", true},
		{"5.0 is 4.20 or later", "5.0", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := IsVersion420OrLater(testCase.version)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestListSupportedVersions(t *testing.T) {
	versions := ListSupportedVersions()
	assert.Contains(t, versions, "4.14")
	assert.Contains(t, versions, "4.20")
	assert.Contains(t, versions, "4.21")
	assert.True(t, len(versions) >= 7, "Should have at least 7 supported versions")
}

func TestOperatorInfoString(t *testing.T) {
	operatorInfo := OperatorInfo{
		PackageName:   "test-operator",
		CatalogSource: "test-catalog",
		CSVPrefix:     "test-prefix",
	}

	str := operatorInfo.String()
	assert.Contains(t, str, "test-operator")
	assert.Contains(t, str, "test-catalog")
	assert.Contains(t, str, "test-prefix")
}

func TestOCPOperatorConfigString(t *testing.T) {
	config := GetOperatorConfig("4.19")
	str := config.String()
	assert.Contains(t, str, "cockroachdb-certified")
	assert.Contains(t, str, "grafana-operator")
}

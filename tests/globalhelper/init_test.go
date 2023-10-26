package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOriginalTNFPaths(t *testing.T) {
	t.Setenv("TNF_CONFIG_DIR", "/tmp/tnf-config")
	t.Setenv("TNF_REPORT_DIR", "/tmp/tnf-report")

	reportDir, tnfConfigDir := GetOriginalTNFPaths()
	assert.Equal(t, "/tmp/tnf-report", reportDir)
	assert.Equal(t, "/tmp/tnf-config", tnfConfigDir)
}

func TestOverrideDirectories(t *testing.T) {
	t.Setenv("TNF_CONFIG_DIR", "/tmp/tnf-config")
	t.Setenv("TNF_REPORT_DIR", "/tmp/tnf-report")

	OverrideDirectories("abc123")

	reportDir, tnfConfigDir := GetOriginalTNFPaths()
	assert.Equal(t, "/tmp/tnf-report/abc123", reportDir)
	assert.Equal(t, "/tmp/tnf-config/abc123", tnfConfigDir)
}

func TestRestoreOriginalTNFPaths(t *testing.T) {
	RestoreOriginalTNFPaths("/tmp/tnf-report", "/tmp/tnf-config")
	assert.Equal(t, "/tmp/tnf-report", GetConfiguration().General.TnfReportDir)
	assert.Equal(t, "/tmp/tnf-config", GetConfiguration().General.TnfConfigDir)
}

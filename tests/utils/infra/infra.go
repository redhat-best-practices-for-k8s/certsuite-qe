package infra

import (
	"os"
	"strings"
)

// ShouldEnableInfrastructureTolerations checks if infrastructure tolerations should be enabled
// based on environment configuration. Returns true by default if ENABLE_INFRASTRUCTURE_TOLERATIONS
// is not set or is set to "true".
func ShouldEnableInfrastructureTolerations() bool {
	enabled := os.Getenv("ENABLE_INFRASTRUCTURE_TOLERATIONS")
	if enabled == "" {
		return true
	}

	return strings.ToLower(enabled) == "true"
}

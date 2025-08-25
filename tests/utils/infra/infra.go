package infra

import (
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

// IsConstrainedEnvironment detects if we're running in a resource-constrained environment
// like CRC, OpenShift Local, or CI with limited resources.
func IsConstrainedEnvironment() bool {
	constrainedEnv := os.Getenv("CONSTRAINED_ENVIRONMENT")
	if constrainedEnv != "" {
		return strings.ToLower(constrainedEnv) == "true"
	}

	// Auto-detect common constrained environments
	ciEnv := os.Getenv("CI")
	crctEnv := os.Getenv("CRC")
	githubActions := os.Getenv("GITHUB_ACTIONS")

	return ciEnv == "true" || crctEnv == "true" || githubActions == "true"
}

// GetDefaultResourceRequirements returns resource requirements based on environment.
// For constrained environments (CRC, CI), returns very conservative limits.
// For full clusters, returns more generous limits.
func GetDefaultResourceRequirements() corev1.ResourceRequirements {
	if IsConstrainedEnvironment() {
		return GetConstrainedResourceRequirements()
	}
	return GetStandardResourceRequirements()
}

// GetConstrainedResourceRequirements returns very conservative resource requirements
// suitable for CRC/OpenShift Local environments with limited resources (16GB RAM, 4 vCPU).
func GetConstrainedResourceRequirements() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("256Mi"),
			corev1.ResourceCPU:    resource.MustParse("500m"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("128Mi"),
			corev1.ResourceCPU:    resource.MustParse("250m"),
		},
	}
}

// GetConstrainedResourceRequirementsForDatabase returns resource requirements
// specifically tuned for database containers in constrained environments.
func GetConstrainedResourceRequirementsForDatabase() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("512Mi"),
			corev1.ResourceCPU:    resource.MustParse("750m"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("256Mi"),
			corev1.ResourceCPU:    resource.MustParse("500m"),
		},
	}
}

// GetStandardResourceRequirements returns resource requirements for full clusters.
func GetStandardResourceRequirements() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("1Gi"),
			corev1.ResourceCPU:    resource.MustParse("1000m"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("512Mi"),
			corev1.ResourceCPU:    resource.MustParse("500m"),
		},
	}
}

// GetMinimalResourceRequirements returns minimal resource requirements for lightweight containers.
func GetMinimalResourceRequirements() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("128Mi"),
			corev1.ResourceCPU:    resource.MustParse("250m"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("64Mi"),
			corev1.ResourceCPU:    resource.MustParse("100m"),
		},
	}
}

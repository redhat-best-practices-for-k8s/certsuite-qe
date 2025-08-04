package infra

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldEnableInfrastructureTolerations(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		envSet      bool
		expected    bool
		description string
	}{
		{
			name:        "env_not_set",
			envSet:      false,
			expected:    true,
			description: "should default to true when environment variable is not set",
		},
		{
			name:        "env_empty_string",
			envValue:    "",
			envSet:      true,
			expected:    true,
			description: "should default to true when environment variable is empty",
		},
		{
			name:        "env_true_lowercase",
			envValue:    "true",
			envSet:      true,
			expected:    true,
			description: "should return true when set to 'true'",
		},
		{
			name:        "env_true_uppercase",
			envValue:    "TRUE",
			envSet:      true,
			expected:    true,
			description: "should return true when set to 'TRUE'",
		},
		{
			name:        "env_true_mixedcase",
			envValue:    "True",
			envSet:      true,
			expected:    true,
			description: "should return true when set to 'True'",
		},
		{
			name:        "env_false_lowercase",
			envValue:    "false",
			envSet:      true,
			expected:    false,
			description: "should return false when set to 'false'",
		},
		{
			name:        "env_false_uppercase",
			envValue:    "FALSE",
			envSet:      true,
			expected:    false,
			description: "should return false when set to 'FALSE'",
		},
		{
			name:        "env_random_value",
			envValue:    "random",
			envSet:      true,
			expected:    false,
			description: "should return false when set to any non-'true' value",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Setup environment
			const envKey = "ENABLE_INFRASTRUCTURE_TOLERATIONS"
			oldValue, wasSet := os.LookupEnv(envKey)

			defer func() {
				if wasSet {
					t.Setenv(envKey, oldValue)
				} else {
					os.Unsetenv(envKey)
				}
			}()

			if testCase.envSet {
				t.Setenv(envKey, testCase.envValue)
			} else {
				os.Unsetenv(envKey)
			}

			// Test the function
			result := ShouldEnableInfrastructureTolerations()
			assert.Equal(t, testCase.expected, result, testCase.description)
		})
	}
}

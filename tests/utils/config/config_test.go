package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugCertsuite(t *testing.T) {
	testCases := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.value, func(t *testing.T) {
			t.Setenv("DEBUG_TNF", testCase.value)

			c, err := NewConfig()
			assert.Nil(t, err)
			result, err := c.DebugCertsuite()
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

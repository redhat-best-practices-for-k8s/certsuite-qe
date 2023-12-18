package subscription

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefineSubscription(t *testing.T) {
	testCases := []struct {
		testName      string
		testNamespace string
	}{
		{
			testName:      "test1",
			testNamespace: "testNamespace1",
		},
		{
			testName:      "test2",
			testNamespace: "testNamespace2",
		},
	}

	for _, tc := range testCases {
		subscription := DefineSubscription(tc.testName, tc.testNamespace)
		assert.NotNil(t, subscription)
		assert.Equal(t, tc.testName, subscription.Name)
		assert.Equal(t, tc.testNamespace, subscription.Namespace)
	}
}

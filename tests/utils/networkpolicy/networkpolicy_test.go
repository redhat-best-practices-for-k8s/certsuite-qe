package networkpolicy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefineDenyAllNetworkPolicy(t *testing.T) {
	testCases := []struct {
		name        string
		ns          string
		policyTypes []string
		labels      map[string]string
	}{
		{"test", "default", []string{"Ingress", "Egress"}, map[string]string{"app": "test"}},
		{"test", "default", []string{"Ingress"}, map[string]string{"app": "test"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			policyTypes := DefinePolicyTypes(testCase.policyTypes)
			networkPolicy := DefineDenyAllNetworkPolicy(testCase.name, testCase.ns, policyTypes, testCase.labels)
			assert.NotNil(t, networkPolicy)
			assert.Equal(t, testCase.name, networkPolicy.Name)
			assert.Equal(t, testCase.ns, networkPolicy.Namespace)
			assert.Equal(t, policyTypes, networkPolicy.Spec.PolicyTypes)
			assert.Equal(t, testCase.labels, networkPolicy.Spec.PodSelector.MatchLabels)
		})
	}
}

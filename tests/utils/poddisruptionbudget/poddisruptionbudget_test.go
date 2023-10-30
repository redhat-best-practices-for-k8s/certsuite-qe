package poddisruptionbudget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestDefinePodDisruptionBudgetMinAvailable(t *testing.T) {
	testCases := []struct {
		testMinAvailable  int
		testLabelSelector map[string]string
	}{
		{
			testMinAvailable:  1,
			testLabelSelector: map[string]string{"app": "test"},
		},
		{
			testMinAvailable:  2,
			testLabelSelector: map[string]string{"app": "test"},
		},
	}

	for _, testCase := range testCases {
		pdb := DefinePodDisruptionBudgetMinAvailable("test", "test", intstr.FromInt(testCase.testMinAvailable), testCase.testLabelSelector)

		assert.Equal(t, *pdb.Spec.MinAvailable, intstr.FromInt(testCase.testMinAvailable))
		assert.Equal(t, pdb.Spec.Selector.MatchLabels, testCase.testLabelSelector)
	}
}

func TestDefinePodDisruptionBudgetMaxUnAvailable(t *testing.T) {
	testCases := []struct {
		testMaxUnAvailable int
		testLabelSelector  map[string]string
	}{
		{
			testMaxUnAvailable: 1,
			testLabelSelector:  map[string]string{"app": "test"},
		},
		{
			testMaxUnAvailable: 2,
			testLabelSelector:  map[string]string{"app": "test"},
		},
	}

	for _, testCase := range testCases {
		pdb := DefinePodDisruptionBudgetMaxUnAvailable("test", "test",
			intstr.FromInt(testCase.testMaxUnAvailable), testCase.testLabelSelector)

		assert.Equal(t, *pdb.Spec.MaxUnavailable, intstr.FromInt(testCase.testMaxUnAvailable))
		assert.Equal(t, pdb.Spec.Selector.MatchLabels, testCase.testLabelSelector)
	}
}

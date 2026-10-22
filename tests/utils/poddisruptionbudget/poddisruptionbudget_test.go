package poddisruptionbudget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestDefinePDBMinAvailableWithMatchLabelsAndExpressions(t *testing.T) {
	labels := map[string]string{"app": "test"}

	expressions := []metav1.LabelSelectorRequirement{{
		Key:      "tier",
		Operator: metav1.LabelSelectorOpIn,
		Values:   []string{"frontend", "backend"},
	}}

	pdb := DefinePDBMinAvailableWithMatchLabelsAndExpressions("test-pdb", "test-ns",
		intstr.FromString("50%"), labels, expressions)

	assert.Equal(t, "test-pdb", pdb.Name)
	assert.Equal(t, "test-ns", pdb.Namespace)
	require.NotNil(t, pdb.Spec.MinAvailable)
	assert.Equal(t, intstr.FromString("50%"), *pdb.Spec.MinAvailable)

	require.NotNil(t, pdb.Spec.Selector)
	assert.Equal(t, labels, pdb.Spec.Selector.MatchLabels)
	assert.Equal(t, expressions, pdb.Spec.Selector.MatchExpressions)
	assert.Nil(t, pdb.Spec.MaxUnavailable)
}

// Each builder must leave the opposing availability field unset: the Kubernetes API
// rejects a PodDisruptionBudget that sets both minAvailable and maxUnavailable.
func TestPodDisruptionBudgetBuildersSetOneAvailabilityField(t *testing.T) {
	labels := map[string]string{"app": "test"}

	t.Run("min available leaves max unavailable unset", func(t *testing.T) {
		pdb := DefinePodDisruptionBudgetMinAvailable("min-pdb", "test-ns", intstr.FromInt(1), labels)
		assert.Equal(t, "min-pdb", pdb.Name)
		assert.Equal(t, "test-ns", pdb.Namespace)
		assert.Nil(t, pdb.Spec.MaxUnavailable)
	})

	t.Run("max unavailable leaves min available unset", func(t *testing.T) {
		pdb := DefinePodDisruptionBudgetMaxUnAvailable("max-pdb", "test-ns", intstr.FromInt(1), labels)
		assert.Equal(t, "max-pdb", pdb.Name)
		assert.Equal(t, "test-ns", pdb.Namespace)
		assert.Nil(t, pdb.Spec.MinAvailable)
	})
}

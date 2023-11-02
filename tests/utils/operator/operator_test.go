package utils

import (
	"testing"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestDefineOperatorGroup(t *testing.T) {
	og := DefineOperatorGroup("test", "default", []string{"testNamespace"})
	assert.NotNil(t, og)
	assert.Equal(t, "test", og.Name)
	assert.Equal(t, "default", og.Namespace)
	assert.Equal(t, []string{"testNamespace"}, og.Spec.TargetNamespaces)
}

func TestDefineSubscription(t *testing.T) {
	testSubscription := DefineSubscription("testName", "testNamespace",
		"testChannel", "testOperator", "testSource", "testSourceNamespace", "testCSV",
		v1alpha1.ApprovalAutomatic)
	assert.NotNil(t, testSubscription)
	assert.Equal(t, "testName", testSubscription.Name)
	assert.Equal(t, "testNamespace", testSubscription.Namespace)
	assert.Equal(t, "testChannel", testSubscription.Spec.Channel)
	assert.Equal(t, "testOperator", testSubscription.Spec.Package)
	assert.Equal(t, "testSource", testSubscription.Spec.CatalogSource)
	assert.Equal(t, "testSourceNamespace", testSubscription.Spec.CatalogSourceNamespace)
	assert.Equal(t, "testCSV", testSubscription.Spec.StartingCSV)
	assert.Equal(t, v1alpha1.ApprovalAutomatic, testSubscription.Spec.InstallPlanApproval)
}

func TestDefineSubscriptionWithNodeSelector(t *testing.T) {
	testSubscription := DefineSubscriptionWithNodeSelector("testName", "testNamespace",
		"testChannel", "testOperator", "testSource", "testSourceNamespace", "testCSV",
		v1alpha1.ApprovalAutomatic, map[string]string{"node-role.kubernetes.io/worker": ""})
	assert.NotNil(t, testSubscription)
	assert.Equal(t, "testName", testSubscription.Name)
	assert.Equal(t, "testNamespace", testSubscription.Namespace)
	assert.Equal(t, "testChannel", testSubscription.Spec.Channel)
	assert.Equal(t, "testOperator", testSubscription.Spec.Package)
	assert.Equal(t, "testSource", testSubscription.Spec.CatalogSource)
	assert.Equal(t, "testSourceNamespace", testSubscription.Spec.CatalogSourceNamespace)
	assert.Equal(t, "testCSV", testSubscription.Spec.StartingCSV)
	assert.Equal(t, v1alpha1.ApprovalAutomatic, testSubscription.Spec.InstallPlanApproval)
	assert.Equal(t, map[string]string{"node-role.kubernetes.io/worker": ""}, testSubscription.Spec.Config.NodeSelector)
}

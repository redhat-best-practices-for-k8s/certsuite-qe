package resourcequota

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestDefineResourceQuota(t *testing.T) {
	testResourceQuota := DefineResourceQuota("testResourceQuota", "testNamespace", "100m", "100Mi", "150m", "150Mi")
	assert.NotNil(t, testResourceQuota)
	assert.Equal(t, "testResourceQuota", testResourceQuota.Name)
	assert.Equal(t, "testNamespace", testResourceQuota.Namespace)
	assert.Equal(t, resource.MustParse("100m"), testResourceQuota.Spec.Hard["requests.cpu"])
	assert.Equal(t, resource.MustParse("100Mi"), testResourceQuota.Spec.Hard["requests.memory"])
	assert.Equal(t, resource.MustParse("150m"), testResourceQuota.Spec.Hard["limits.cpu"])
	assert.Equal(t, resource.MustParse("150Mi"), testResourceQuota.Spec.Hard["limits.memory"])
}

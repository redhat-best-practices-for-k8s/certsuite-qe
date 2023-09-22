package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
)

func TestDefinePodUnderTestLabels(t *testing.T) {
	testConfig := globalparameters.TnfConfig{}
	assert.Empty(t, testConfig.PodsUnderTestLabels)

	testPodLabels := []string{"test1:value", "test2:value"}
	assert.Nil(t, definePodUnderTestLabels(&testConfig, testPodLabels))
	assert.Equal(t, testPodLabels, testConfig.PodsUnderTestLabels)
}

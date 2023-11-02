package runtimeclass

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefineRunTimeClass(t *testing.T) {
	testRTC := DefineRunTimeClass("testRTC")
	assert.NotNil(t, testRTC)
	assert.Contains(t, testRTC.Name, "testRTC-")
	assert.Equal(t, "runc", testRTC.Handler)
}

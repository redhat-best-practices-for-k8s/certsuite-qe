package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			t.Setenv("DEBUG_CERTSUITE", testCase.value)

			c, err := NewConfig()
			require.NoError(t, err)
			require.NotNil(t, c)

			result, err := c.DebugCertsuite()
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestNewConfigDockerConfigDirDefaultsToTempDir(t *testing.T) {
	t.Setenv("DOCKER_CONFIG_DIR", "")

	conf, err := NewConfig()
	require.NoError(t, err)
	require.NotNil(t, conf)

	assert.Equal(t, filepath.Join(os.TempDir(), defaultDockerConfigDirName), conf.General.DockerConfigDir)
	assert.FileExists(t, filepath.Join(conf.General.DockerConfigDir, "config"))
}

func TestNewConfigDockerConfigDirFromEnv(t *testing.T) {
	dockerConfigDir := filepath.Join(t.TempDir(), "docker")
	t.Setenv("DOCKER_CONFIG_DIR", dockerConfigDir)

	conf, err := NewConfig()
	require.NoError(t, err)
	require.NotNil(t, conf)

	assert.Equal(t, dockerConfigDir, conf.General.DockerConfigDir)
	assert.FileExists(t, filepath.Join(dockerConfigDir, "config"))
}

func TestDebugCertsuiteNilConfig(t *testing.T) {
	var c *Config

	result, err := c.DebugCertsuite()
	assert.Error(t, err)
	assert.False(t, result)
}

func TestMakeDockerConfigKeepsExistingConfig(t *testing.T) {
	dockerConfigDir := t.TempDir()
	configPath := filepath.Join(dockerConfigDir, "config")
	sentinel := `{"auths":{"registry.example.com":{"auth":"c2VjcmV0"}}}`

	require.NoError(t, os.WriteFile(configPath, []byte(sentinel), 0o600))

	conf := &Config{}
	conf.General.DockerConfigDir = dockerConfigDir

	require.NoError(t, conf.makeDockerConfig())

	contents, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Equal(t, sentinel, string(contents), "an existing docker config must never be rewritten")
}

func TestMakeDockerConfigCreatesMissingConfig(t *testing.T) {
	conf := &Config{}
	conf.General.DockerConfigDir = filepath.Join(t.TempDir(), "docker")

	require.NoError(t, conf.makeDockerConfig())

	contents, err := os.ReadFile(filepath.Join(conf.General.DockerConfigDir, "config"))
	require.NoError(t, err)
	assert.JSONEq(t, `{"auths": {}}`, string(contents))
}

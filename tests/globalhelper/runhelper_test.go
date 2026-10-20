//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithRetry(t *testing.T) {
	success := writeExecutable(t, "#!/bin/sh\nexit 0\n")
	assert.NoError(t, executeWithRetry(success, nil, "tc", nil, nil, nil))

	failure := writeExecutable(t, "#!/bin/sh\nexit 7\n")
	err := executeWithRetry(failure, nil, "tc", nil, nil, []string{"PATH=/bin:/usr/bin"})
	assert.Error(t, err)

	oldTimeout := commandTimeout
	commandTimeout = 10 * time.Millisecond

	t.Cleanup(func() { commandTimeout = oldTimeout })
	sleeper := writeExecutable(t, "#!/bin/sh\nsleep 1\n")
	err = executeWithRetry(sleeper, nil, "tc", nil, nil, nil)
	assert.Error(t, err)
}

func TestLaunchTestsViaBinary(t *testing.T) {
	oldConfig := conf

	t.Cleanup(func() { conf = oldConfig })
	root := t.TempDir()
	binary := writeExecutable(t, "#!/bin/sh\nexit 0\n")
	conf = &config.Config{}
	conf.General.CertsuiteRepoPath = root
	conf.General.CertsuiteEntryPointBinary = filepath.Base(binary)
	conf.General.ReportDirAbsPath = filepath.Join(root, "reports")
	conf.General.DebugCertsuite = "false"
	conf.General.UseBinary = "true"
	reportDir, configDir := filepath.Join(root, "results"), filepath.Join(root, "config")
	assert.Error(t, launchTestsViaBinary("networking-case", "networking_case", reportDir, configDir)) // absent binary
	conf.General.CertsuiteRepoPath = root

	assert.NoError(t, os.MkdirAll(reportDir, 0o755))
	assert.NoError(t, os.MkdirAll(configDir, 0o755))
	assert.NoError(t, os.Rename(binary, filepath.Join(root, filepath.Base(binary))))
	assert.NoError(t, launchTestsViaBinary("networking-case", "networking_case", reportDir, configDir))
	assert.NoError(t, LaunchTests("networking-case", "networking_case", reportDir, configDir))
	assert.Equal(t, globalparameters.NetworkSuiteName, mustSuiteName(t, "networking-case"))
	_, err := getTestSuiteName("unknown")
	assert.Error(t, err)

	conf.General.DebugCertsuite = "true"

	assert.Error(t, launchTestsViaBinary("unknown-case", "unknown", reportDir, configDir))
	assert.NoError(t, launchTestsViaBinary("networking-case", "networking_debug", reportDir, configDir))
	logPath := filepath.Join(conf.General.ReportDirAbsPath, "Debug", globalparameters.NetworkSuiteName, "networking_debug", "networking_debug.log")
	logData, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	assert.Contains(t, string(logData), "Running test: networking_debug")

	oldDebug, oldWrite, oldSet := debugCertsuite, writeDebugLog, setEnv

	t.Cleanup(func() { debugCertsuite, writeDebugLog, setEnv = oldDebug, oldWrite, oldSet })
	debugCertsuite = func() (bool, error) { return false, errors.New("debug env failed") }

	assert.Error(t, launchTestsViaBinary("networking-case", "case", reportDir, configDir))
	debugCertsuite = func() (bool, error) { return true, nil }
	writeDebugLog = func(io.Writer, string, ...any) (int, error) { return 0, errors.New("write failed") }

	assert.Error(t, launchTestsViaBinary("networking-case", "case", reportDir, configDir))
	writeDebugLog = oldWrite
	setEnv = func(key, value string) error {
		if key == "CERTSUITE_OMIT_ARTIFACTS_ZIP_FILE" {
			return errors.New("set omit env failed")
		}

		t.Setenv(key, value)

		return nil
	}

	assert.Error(t, launchTestsViaBinary("networking-case", "case", reportDir, configDir))
	setEnv = func(key, value string) error {
		if key == "CERTSUITE_ENABLE_DATA_COLLECTION" {
			return errors.New("set env failed")
		}

		t.Setenv(key, value)

		return nil
	}

	assert.Error(t, launchTestsViaBinary("networking-case", "case", reportDir, configDir))

	t.Setenv("ENABLE_CERTSUITE_MEMORY_LIMIT", "true")
	t.Setenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB", "256")
	setEnv = func(key, value string) error {
		t.Setenv(key, value)

		return nil
	}
	conf.General.DebugCertsuite = "false"

	assert.NoError(t, launchTestsViaBinary("networking-case", "memory_limit", reportDir, configDir))
	assert.NoError(t, os.WriteFile(filepath.Join(root, filepath.Base(binary)), []byte("#!/bin/sh\nexit 7\n"), 0o700))
	assert.Error(t, launchTestsViaBinary("networking-case", "binary_failure", reportDir, configDir))
	conf.General.DebugCertsuite = "true"

	assert.Error(t, launchTestsViaImage("unknown-case", "unknown", reportDir, configDir))
}
func TestLaunchTestsViaImage(t *testing.T) {
	oldConfig := conf

	t.Cleanup(func() { conf = oldConfig })
	root := t.TempDir()
	engine := writeExecutable(t, "#!/bin/sh\nexit 0\n")
	conf = &config.Config{}
	conf.General.ContainerEngine = engine
	conf.General.CertsuiteImage = "example/certsuite"
	conf.General.CertsuiteImageTag = "latest"
	conf.General.DockerConfigDir = filepath.Join(root, "docker")
	conf.General.ReportDirAbsPath = filepath.Join(root, "reports")
	conf.General.DebugCertsuite = "false"

	assert.NoError(t, launchTestsViaImage("networking-case", "networking_case", filepath.Join(root, "results"), filepath.Join(root, "config")))
	conf.General.UseBinary = "false"

	assert.NoError(t, LaunchTests("networking-case", "networking_case", filepath.Join(root, "results"), filepath.Join(root, "config")))

	oldDebug, oldWrite := debugCertsuite, writeDebugLog

	t.Cleanup(func() { debugCertsuite, writeDebugLog = oldDebug, oldWrite })
	debugCertsuite = func() (bool, error) { return false, errors.New("debug env failed") }

	assert.Error(t, launchTestsViaImage("networking-case", "case", filepath.Join(root, "results"), filepath.Join(root, "config"))) //nolint:lll
	debugCertsuite = func() (bool, error) { return true, nil }
	writeDebugLog = func(io.Writer, string, ...any) (int, error) { return 0, errors.New("write failed") } //nolint:lll

	assert.Error(t, launchTestsViaImage("networking-case", "case", filepath.Join(root, "results"), filepath.Join(root, "config")))
	writeDebugLog = oldWrite
	conf.General.DebugCertsuite = "true"

	assert.NoError(t, launchTestsViaImage("networking-case", "networking_debug", filepath.Join(root, "results"), filepath.Join(root, "config")))

	failedEngine := writeExecutable(t, "#!/bin/sh\necho failed >&2\nexit 3\n")
	conf.General.ContainerEngine = failedEngine
	conf.General.DebugCertsuite = "false"
	err := launchTestsViaImage("networking-case", "case", filepath.Join(root, "results"), filepath.Join(root, "config"))
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "failed to run tc"))

	t.Setenv("ENABLE_CERTSUITE_MEMORY_LIMIT", "true")
	t.Setenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB", "bad")
	conf.General.ContainerEngine = engine

	assert.NoError(t, launchTestsViaImage("networking-case", "memory_default", filepath.Join(root, "results"), filepath.Join(root, "config")))
	t.Setenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB", "256")
	assert.NoError(t, launchTestsViaImage("networking-case", "memory_configured", filepath.Join(root, "results"), filepath.Join(root, "config")))
}

func writeExecutable(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "fake-executable")
	if err := os.WriteFile(path, []byte(contents), 0o700); err != nil {
		t.Fatal(err)
	}

	return path
}

func mustSuiteName(t *testing.T, testCase string) string {
	t.Helper()

	name, err := getTestSuiteName(testCase)
	if err != nil {
		t.Fatal(err)
	}

	return name
}

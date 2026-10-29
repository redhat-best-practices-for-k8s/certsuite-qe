//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"errors"
	"fmt"
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

func TestTailWriter(t *testing.T) {
	tail := &tailWriter{limit: 8}
	assert.Equal(t, 0, len(tail.buf))
	assert.Empty(t, tail.failureDetails())

	_, err := tail.Write([]byte("1234"))
	assert.NoError(t, err)

	_, err = tail.Write([]byte("5678"))
	assert.NoError(t, err)
	assert.Equal(t, 8, len(tail.buf))
	assert.Equal(t, "12345678", string(tail.buf))

	_, err = tail.Write([]byte("90abcdef"))
	assert.NoError(t, err)
	assert.Equal(t, 8, len(tail.buf))
	assert.Equal(t, "90abcdef", string(tail.buf))

	_, err = tail.Write([]byte("0123456789ABCDEF"))
	assert.NoError(t, err)
	assert.Equal(t, 8, len(tail.buf))
	assert.Equal(t, "89ABCDEF", string(tail.buf))

	tail = &tailWriter{limit: 8}
	_, err = tail.Write([]byte("abcdef"))
	assert.NoError(t, err)

	_, err = tail.Write([]byte("ghi"))
	assert.NoError(t, err)
	assert.Equal(t, 8, len(tail.buf))
	assert.Equal(t, "bcdefghi", string(tail.buf))

	var builder strings.Builder

	for i := range 100 {
		fmt.Fprintf(&builder, "line-%03d\n", i)
	}

	tail = &tailWriter{limit: capturedOutputLimitBytes}

	_, err = tail.Write([]byte(builder.String()))
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(tail.buf), capturedOutputLimitBytes)

	details := tail.failureDetails()
	assert.Contains(t, details, "--- last 50 lines of certsuite output ---")
	assert.Contains(t, details, "line-050")
	assert.Contains(t, details, "line-099")
	assert.NotContains(t, details, "line-049")
	assert.Contains(t, details, "DEBUG_CERTSUITE=true CERTSUITE_LOG_LEVEL=debug")
}

func TestFailureOutputCapturedViaBinary(t *testing.T) {
	oldConfig := conf

	t.Cleanup(func() { conf = oldConfig })
	root := t.TempDir()
	binary := writeExecutable(t, "#!/bin/sh\necho stdout-marker\necho stderr-marker >&2\nexit 7\n")
	conf = &config.Config{}
	conf.General.CertsuiteRepoPath = root
	conf.General.CertsuiteEntryPointBinary = filepath.Base(binary)
	conf.General.ReportDirAbsPath = filepath.Join(root, "reports")
	conf.General.DebugCertsuite = "false"
	conf.General.UseBinary = "true"
	reportDir, configDir := filepath.Join(root, "results"), filepath.Join(root, "config")
	assert.NoError(t, os.MkdirAll(reportDir, 0o755))
	assert.NoError(t, os.MkdirAll(configDir, 0o755))
	assert.NoError(t, os.Rename(binary, filepath.Join(root, filepath.Base(binary))))

	err := launchTestsViaBinary("networking-case", "networking_failure", reportDir, configDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to run tc")
	assert.Contains(t, err.Error(), "stdout-marker")
	assert.Contains(t, err.Error(), "stderr-marker")
	assert.Contains(t, err.Error(), "DEBUG_CERTSUITE=true CERTSUITE_LOG_LEVEL=debug")
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

	assert.NoError(t, os.WriteFile(filepath.Join(root, filepath.Base(binary)), []byte("#!/bin/sh\necho debug-tee-marker\n"), 0o700))
	assert.NoError(t, launchTestsViaBinary("networking-case", "networking_tee", reportDir, configDir))

	teeLogPath := filepath.Join(conf.General.ReportDirAbsPath, "Debug", globalparameters.NetworkSuiteName, "networking_tee", "networking_tee.log")
	teeLogData, err := os.ReadFile(teeLogPath)
	assert.NoError(t, err)
	assert.Contains(t, string(teeLogData), "debug-tee-marker")

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

	failedEngine := writeExecutable(t, "#!/bin/sh\necho engine-stderr-marker >&2\nexit 3\n")
	conf.General.ContainerEngine = failedEngine
	conf.General.DebugCertsuite = "false"
	err := launchTestsViaImage("networking-case", "case", filepath.Join(root, "results"), filepath.Join(root, "config"))
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "failed to run tc"))
	assert.Contains(t, err.Error(), "engine-stderr-marker")
	assert.Contains(t, err.Error(), "DEBUG_CERTSUITE=true CERTSUITE_LOG_LEVEL=debug")

	t.Setenv("ENABLE_CERTSUITE_MEMORY_LIMIT", "true")
	t.Setenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB", "bad")
	conf.General.ContainerEngine = engine

	assert.NoError(t, launchTestsViaImage("networking-case", "memory_default", filepath.Join(root, "results"), filepath.Join(root, "config")))
	t.Setenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB", "256")
	assert.NoError(t, launchTestsViaImage("networking-case", "memory_configured", filepath.Join(root, "results"), filepath.Join(root, "config")))
}

func TestCertsuiteMemoryLimitMB(t *testing.T) {
	tests := []struct {
		name     string
		enabled  string
		override string
		wantMB   int
		wantOK   bool
	}{
		{name: "unset", wantMB: 0, wantOK: false},
		{name: "disabled", enabled: "false", wantMB: 0, wantOK: false},
		{name: "enabled with default", enabled: "true", wantMB: MemoryLimitMB, wantOK: true},
		{name: "enabled with override", enabled: "true", override: "256", wantMB: 256, wantOK: true},
		{name: "non-numeric override", enabled: "true", override: "bad", wantMB: MemoryLimitMB, wantOK: true},
		{name: "zero override", enabled: "true", override: "0", wantMB: MemoryLimitMB, wantOK: true},
		{name: "negative override", enabled: "true", override: "-5", wantMB: MemoryLimitMB, wantOK: true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("ENABLE_CERTSUITE_MEMORY_LIMIT", testCase.enabled)
			t.Setenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB", testCase.override)

			memLimitMB, enabled := certsuiteMemoryLimitMB()

			assert.Equal(t, testCase.wantMB, memLimitMB)
			assert.Equal(t, testCase.wantOK, enabled)
		})
	}
}

func TestCertsuiteMemoryLimitArgs(t *testing.T) {
	t.Setenv("ENABLE_CERTSUITE_MEMORY_LIMIT", "false")
	t.Setenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB", "256")
	assert.Empty(t, certsuiteMemoryLimitArgs())

	t.Setenv("ENABLE_CERTSUITE_MEMORY_LIMIT", "true")
	assert.Equal(t, []string{"-e", "GOMEMLIMIT=256MiB"}, certsuiteMemoryLimitArgs())
}

func TestLaunchTestsViaImagePassesMemoryLimit(t *testing.T) {
	oldConfig := conf

	t.Cleanup(func() { conf = oldConfig })

	argsFile := filepath.Join(t.TempDir(), "args")
	engine := writeExecutable(t, "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$FAKE_ARGS_FILE\"\n")

	t.Setenv("FAKE_ARGS_FILE", argsFile)

	conf = &config.Config{}
	conf.General.ContainerEngine = engine
	conf.General.CertsuiteImage = "example/certsuite"
	conf.General.CertsuiteImageTag = "latest"
	conf.General.DockerConfigDir = filepath.Join(t.TempDir(), "docker")
	conf.General.ReportDirAbsPath = t.TempDir()
	conf.General.DebugCertsuite = "false"

	runImage := func() string {
		t.Helper()

		assert.NoError(t, launchTestsViaImage("networking-case", "case", t.TempDir(), t.TempDir()))

		contents, err := os.ReadFile(argsFile)
		assert.NoError(t, err)

		return string(contents)
	}

	t.Setenv("ENABLE_CERTSUITE_MEMORY_LIMIT", "true")
	t.Setenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB", "256")

	enabledArgs := runImage()

	assert.Contains(t, enabledArgs, "-e\nGOMEMLIMIT=256MiB\n")

	gomemlimitIdx := strings.Index(enabledArgs, "GOMEMLIMIT=256MiB")
	imageIdx := strings.Index(enabledArgs, "example/certsuite:latest")

	assert.GreaterOrEqual(t, gomemlimitIdx, 0)
	assert.Greater(t, imageIdx, gomemlimitIdx, "-e flags must precede the image reference")

	t.Setenv("ENABLE_CERTSUITE_MEMORY_LIMIT", "false")

	assert.NotContains(t, runImage(), "GOMEMLIMIT")
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

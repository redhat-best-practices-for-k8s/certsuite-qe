package globalhelper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	klog "k8s.io/klog/v2"
)

const (
	// MaxRetries is the maximum number of retries for test execution.
	MaxRetries = 2
	// TestTimeout is the timeout duration for test execution.
	TestTimeout = 30 * time.Minute
	// MemoryLimitMB is the default soft cap for Go-managed memory when enabled.
	MemoryLimitMB = 50
	// capturedOutputLimitBytes bounds the certsuite output retained in memory for failure reporting.
	capturedOutputLimitBytes = 64 * 1024
	// failureOutputTailLines is how many trailing output lines a failure error includes.
	failureOutputTailLines = 50
)

var (
	commandTimeout = TestTimeout
	setEnv         = os.Setenv
	debugCertsuite = func() (bool, error) { return GetConfiguration().DebugCertsuite() }
	writeDebugLog  = fmt.Fprintf
)

// certsuiteMemoryLimitMB returns the Go soft memory limit to apply and whether
// the limit is enabled. Invalid or non-positive overrides fall back to
// MemoryLimitMB.
func certsuiteMemoryLimitMB() (int, bool) {
	if os.Getenv("ENABLE_CERTSUITE_MEMORY_LIMIT") != "true" {
		return 0, false
	}

	memLimitMB := MemoryLimitMB

	if v := os.Getenv("CERTSUITE_MEMORY_SOFT_LIMIT_MB"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			memLimitMB = parsed
		}
	}

	return memLimitMB, true
}

// certsuiteMemoryLimitArgs returns the container engine flags that set
// GOMEMLIMIT inside the certsuite container, or nil when the limit is disabled.
func certsuiteMemoryLimitArgs() []string {
	memLimitMB, enabled := certsuiteMemoryLimitMB()
	if !enabled {
		return nil
	}

	return []string{"-e", fmt.Sprintf("GOMEMLIMIT=%dMiB", memLimitMB)}
}

// tailWriter keeps only the trailing limit bytes of everything written to it, so
// certsuite output can be reported on failure without unbounded memory growth.
type tailWriter struct {
	limit int
	buf   []byte
}

func (w *tailWriter) Write(data []byte) (int, error) {
	written := len(data)

	if written >= w.limit {
		w.buf = append(w.buf[:0], data[written-w.limit:]...)

		return written, nil
	}

	w.buf = append(w.buf, data...)

	if overflow := len(w.buf) - w.limit; overflow > 0 {
		w.buf = append(w.buf[:0], w.buf[overflow:]...)
	}

	return written, nil
}

// newOutputWriter returns a writer for certsuite output that always feeds the bounded
// tail retained for failure reporting, and the debug log when one is available.
func newOutputWriter(outfile *os.File) (io.Writer, *tailWriter) {
	tail := &tailWriter{limit: capturedOutputLimitBytes}

	if outfile == nil {
		return tail, tail
	}

	return io.MultiWriter(outfile, tail), tail
}

// failureDetails renders the last failureOutputTailLines lines of captured output
// plus a hint on how to get the full logs.
func (w *tailWriter) failureDetails() string {
	if len(w.buf) == 0 {
		return ""
	}

	lines := strings.Split(strings.TrimRight(string(w.buf), "\n"), "\n")
	if len(lines) > failureOutputTailLines {
		lines = lines[len(lines)-failureOutputTailLines:]
	}

	return fmt.Sprintf("\n--- last %d lines of certsuite output ---\n%s\n"+
		"(hint: set DEBUG_CERTSUITE=true CERTSUITE_LOG_LEVEL=debug for full logs)",
		failureOutputTailLines, strings.Join(lines, "\n"))
}

// executeWithRetry executes a command with timeout and retry logic.
func executeWithRetry(cmdPath string, args []string, testCaseName string, stdout, stderr io.Writer, env []string) error {
	var err error

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		// Create a context with timeout for each attempt
		ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)

		// Create a new command for each attempt
		cmd := exec.CommandContext(ctx, cmdPath, args...)
		if env != nil {
			cmd.Env = env
		}

		if stdout != nil {
			cmd.Stdout = stdout
		}

		if stderr != nil {
			cmd.Stderr = stderr
		}

		klog.V(5).Infof("Attempt %d/%d: Running test: %s", attempt, MaxRetries, testCaseName)

		err = cmd.Run()

		// If the command succeeded, return success
		if err == nil {
			cancel()

			return nil
		}

		// Check if it was a timeout error
		if ctx.Err() == context.DeadlineExceeded {
			klog.V(5).Infof("Attempt %d/%d timed out after %v for test: %s",
				attempt, MaxRetries, TestTimeout, testCaseName)
			cancel()

			// If this wasn't the last attempt, continue retrying
			if attempt < MaxRetries {
				klog.V(5).Infof("Retrying test: %s", testCaseName)

				continue
			}
		} else {
			// If it's not a timeout error, don't retry
			cancel()

			break
		}

		cancel()
	}

	return err
}

func launchTestsViaBinary(testCaseName string, tcNameForReport string, reportDir string, configDir string) error {
	// check that the binary exists and is executable in the certsuite repo path
	_, err := os.Stat(fmt.Sprintf("%s/%s", GetConfiguration().General.CertsuiteRepoPath,
		GetConfiguration().General.CertsuiteEntryPointBinary))
	if err != nil {
		klog.V(5).Infof("binary does not exist: %s. "+
			"Please run `make build-certsuite-tool` in the certsuite repo.", err)

		return fmt.Errorf("binary does not exist: %w", err)
	}

	// disable the zip file creation
	err = setEnv("CERTSUITE_OMIT_ARTIFACTS_ZIP_FILE", "true")
	if err != nil {
		return fmt.Errorf("failed to set CERTSUITE_OMIT_ARTIFACTS_ZIP_FILE: %w", err)
	}

	// enable the collector
	err = setEnv("CERTSUITE_ENABLE_DATA_COLLECTION", "false")
	if err != nil {
		return fmt.Errorf("failed to set CERTSUITE_ENABLE_DATA_COLLECTION: %w", err)
	}

	// populate the arguments for the binary
	testArgs := []string{
		"run",
		"--config-file", configDir + "/" + globalparameters.DefaultCertsuiteConfigFileName,
		"--output-dir", reportDir,
		"--label-filter", testCaseName,
		"--sanitize-claim", "true",
		"--cleanup-probe", "false",
	}

	cmdPath := fmt.Sprintf("%s/%s", GetConfiguration().General.CertsuiteRepoPath,
		GetConfiguration().General.CertsuiteEntryPointBinary)

	klog.Infof("cmd: %s %s", cmdPath, strings.Join(testArgs, " "))

	debugEnabled, err := debugCertsuite()
	if err != nil {
		return fmt.Errorf("failed to set env var CERTSUITE_LOG_LEVEL: %w", err)
	}

	var outfile *os.File

	if debugEnabled {
		suiteName, suiteErr := getTestSuiteName(testCaseName)
		if suiteErr != nil {
			return fmt.Errorf("failed to create debug log file: %w", suiteErr)
		}

		outfile = GetConfiguration().CreateLogFile(suiteName, tcNameForReport)

		defer outfile.Close()

		_, err = writeDebugLog(outfile, "Running test: %s\n", tcNameForReport)
		if err != nil {
			return fmt.Errorf("failed to write to debug file: %w", err)
		}
	}

	// before running, construct env with optional GOMEMLIMIT
	env := os.Environ()

	if memLimitMB, enabled := certsuiteMemoryLimitMB(); enabled {
		env = append(env, fmt.Sprintf("GOMEMLIMIT=%dMiB", memLimitMB))
	}

	output, tail := newOutputWriter(outfile)

	err = executeWithRetry(cmdPath, testArgs, testCaseName, output, output, env)
	if err != nil {
		err = fmt.Errorf("failed to run tc: %s, err: %w, cmd: %s %s%s",
			testCaseName, err, cmdPath, strings.Join(testArgs, " "), tail.failureDetails())
	}

	CopyClaimFileToTcFolder(testCaseName, tcNameForReport, reportDir)

	return err
}

func launchTestsViaImage(testCaseName string, tcNameForReport string, reportDir string, configDir string) error {
	// use the container to run the tests
	// Note: Unlike launchTestsViaBinary, this function does not use executeWithRetry to avoid
	// abandoned containers. When executeWithRetry times out, the container process may continue
	// running even after the parent process gives up, leading to resource leaks.
	containerEngine := GetConfiguration().General.ContainerEngine
	klog.V(5).Infof("Selected Container engine:%s", containerEngine)

	certsuiteCmdArgs := []string{
		"run",
		"--rm",
		"--network", "host",
		"-v", fmt.Sprintf("%s:%s", os.Getenv("KUBECONFIG"), "/usr/certsuite/kubeconfig/config:Z"),
		"-v", fmt.Sprintf("%s:%s", GetConfiguration().General.DockerConfigDir+"/config", "/usr/certsuite/dockerconfig/config:Z"),
		"-v", fmt.Sprintf("%s:%s", configDir, "/usr/certsuite/config:Z"),
		"-v", fmt.Sprintf("%s:%s", reportDir, "/usr/certsuite/results:Z"),
	}

	// container engine flags (such as -e GOMEMLIMIT) must precede the image reference
	certsuiteCmdArgs = append(certsuiteCmdArgs, certsuiteMemoryLimitArgs()...)
	certsuiteCmdArgs = append(certsuiteCmdArgs,
		fmt.Sprintf("%s:%s", GetConfiguration().General.CertsuiteImage, GetConfiguration().General.CertsuiteImageTag),
		"certsuite",
		"run",
		"--kubeconfig", "/usr/certsuite/kubeconfig/config",
		"--preflight-dockerconfig", "/usr/certsuite/dockerconfig/config",
		"--config-file", "/usr/certsuite/config/certsuite_config.yml",
		"--output-dir", "/usr/certsuite/results",
		"--omit-artifacts-zip-file", "true",
		"--enable-data-collection", "false",
		"--sanitize-claim", "true",
		"--cleanup-probe", "false",
		"--label-filter", testCaseName,
	)

	// print the command
	klog.V(5).Infof("Running command: %s %s", containerEngine, strings.Join(certsuiteCmdArgs, " "))

	cmd := exec.CommandContext(context.TODO(), containerEngine, certsuiteCmdArgs...)

	debugEnabled, err := debugCertsuite()
	if err != nil {
		return fmt.Errorf("failed to set env var CERTSUITE_LOG_LEVEL: %w", err)
	}

	var outfile *os.File

	if debugEnabled {
		suiteName, suiteErr := getTestSuiteName(testCaseName)
		if suiteErr != nil {
			return fmt.Errorf("failed to create debug log file: %w", suiteErr)
		}

		outfile = GetConfiguration().CreateLogFile(suiteName, tcNameForReport)

		defer outfile.Close()

		_, err = writeDebugLog(outfile, "Running test: %s\n", tcNameForReport)
		if err != nil {
			return fmt.Errorf("failed to write to debug file: %w", err)
		}
	}

	output, tail := newOutputWriter(outfile)

	cmd.Stdout = output
	cmd.Stderr = output

	err = cmd.Run()
	if err != nil {
		errStr := fmt.Sprintf("failed to run tc: %s, err: %v, cmd: %s",
			testCaseName, err, cmd.String())
		if outfile != nil {
			errStr += ", outFile=" + outfile.Name()
		}

		return errors.New(errStr + tail.failureDetails())
	}

	CopyClaimFileToTcFolder(testCaseName, tcNameForReport, reportDir)

	return nil
}

// LaunchTests stats tests based on given parameters.
func LaunchTests(testCaseName string, tcNameForReport string, reportDir string, configDir string) error {
	// check if the `USE_BINARY` flag is set, if so, run the binary version of the tests
	if GetConfiguration().General.UseBinary == "true" {
		return launchTestsViaBinary(testCaseName, tcNameForReport, reportDir, configDir)
	}

	return launchTestsViaImage(testCaseName, tcNameForReport, reportDir, configDir)
}

// suiteNames lists every known suite name for getTestSuiteName lookups.
// Add new suites here instead of extending an if/else chain.
var suiteNames = []string{
	globalparameters.NetworkSuiteName,
	globalparameters.AffiliatedCertificationSuiteName,
	globalparameters.LifecycleSuiteName,
	globalparameters.PlatformAlterationSuiteName,
	globalparameters.ObservabilitySuiteName,
	globalparameters.AccessControlSuiteName,
	globalparameters.PerformanceSuiteName,
	globalparameters.ManageabilitySuiteName,
	globalparameters.OperatorSuiteName,
	globalparameters.PreflightSuiteName,
}

func getTestSuiteName(testCaseName string) (string, error) {
	for _, name := range suiteNames {
		if strings.Contains(testCaseName, name) {
			return name, nil
		}
	}

	return "", fmt.Errorf("unable to retrieve test suite name from test case name %s", testCaseName)
}

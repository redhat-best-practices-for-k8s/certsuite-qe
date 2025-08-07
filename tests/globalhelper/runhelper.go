package globalhelper

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	"github.com/golang/glog"
)

const (
	// MaxRetries is the maximum number of retries for test execution.
	MaxRetries = 3
	// TestTimeout is the timeout duration for test execution.
	TestTimeout = 15 * time.Minute
)

// executeWithRetry executes a command with timeout and retry logic.
func executeWithRetry(cmdPath string, args []string, testCaseName string, stdout, stderr *os.File) error {
	var err error

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		// Create a context with timeout for each attempt
		ctx, cancel := context.WithTimeout(context.Background(), TestTimeout)

		// Create a new command for each attempt
		cmd := exec.CommandContext(ctx, cmdPath, args...)
		if stdout != nil {
			cmd.Stdout = stdout
		}

		if stderr != nil {
			cmd.Stderr = stderr
		}

		glog.V(5).Info(fmt.Sprintf("Attempt %d/%d: Running test: %s", attempt, MaxRetries, testCaseName))

		err = cmd.Run()

		// If the command succeeded, return success
		if err == nil {
			cancel()

			return nil
		}

		// Check if it was a timeout error
		if ctx.Err() == context.DeadlineExceeded {
			glog.V(5).Info(fmt.Sprintf("Attempt %d/%d timed out after %v for test: %s",
				attempt, MaxRetries, TestTimeout, testCaseName))
			cancel()

			// If this wasn't the last attempt, continue retrying
			if attempt < MaxRetries {
				glog.V(5).Info(fmt.Sprintf("Retrying test: %s", testCaseName))

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
		glog.V(5).Info(fmt.Sprintf("binary does not exist: %s. "+
			"Please run `make build-certsuite-tool` in the certsuite repo.", err))

		return fmt.Errorf("binary does not exist: %w", err)
	}

	// disable the zip file creation
	err = os.Setenv("CERTSUITE_OMIT_ARTIFACTS_ZIP_FILE", "true")
	if err != nil {
		return fmt.Errorf("failed to set CERTSUITE_OMIT_ARTIFACTS_ZIP_FILE: %w", err)
	}

	// enable the collector
	err = os.Setenv("CERTSUITE_ENABLE_DATA_COLLECTION", "true")
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
	}

	cmdPath := fmt.Sprintf("%s/%s", GetConfiguration().General.CertsuiteRepoPath,
		GetConfiguration().General.CertsuiteEntryPointBinary)

	fmt.Printf("cmd: %s %s\n", cmdPath, strings.Join(testArgs, " "))

	debugCertsuite, err := GetConfiguration().DebugCertsuite()
	if err != nil {
		return fmt.Errorf("failed to set env var CERTSUITE_LOG_LEVEL: %w", err)
	}

	var outfile *os.File
	if debugCertsuite {
		outfile = GetConfiguration().CreateLogFile(getTestSuiteName(testCaseName), tcNameForReport)

		defer outfile.Close()

		_, err = fmt.Fprintf(outfile, "Running test: %s\n", tcNameForReport)
		if err != nil {
			return fmt.Errorf("failed to write to debug file: %w", err)
		}
	}

	err = executeWithRetry(cmdPath, testArgs, testCaseName, outfile, outfile)
	if err != nil {
		err = fmt.Errorf("failed to run tc: %s, err: %w, cmd: %s %s",
			testCaseName, err, cmdPath, strings.Join(testArgs, " "))
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
	glog.V(5).Info(fmt.Sprintf("Selected Container engine:%s", containerEngine))

	certsuiteCmdArgs := []string{
		"run",
		"--rm",
		"--network", "host",
		"-v", fmt.Sprintf("%s:%s", os.Getenv("KUBECONFIG"), "/usr/certsuite/kubeconfig/config:Z"),
		"-v", fmt.Sprintf("%s:%s", GetConfiguration().General.DockerConfigDir+"/config", "/usr/certsuite/dockerconfig/config:Z"),
		"-v", fmt.Sprintf("%s:%s", configDir, "/usr/certsuite/config:Z"),
		"-v", fmt.Sprintf("%s:%s", reportDir, "/usr/certsuite/results:Z"),
		fmt.Sprintf("%s:%s", GetConfiguration().General.CertsuiteImage, GetConfiguration().General.CertsuiteImageTag),
		"certsuite",
		"run",
		"--kubeconfig", "/usr/certsuite/kubeconfig/config",
		"--preflight-dockerconfig", "/usr/certsuite/dockerconfig/config",
		"--config-file", "/usr/certsuite/config/certsuite_config.yml",
		"--output-dir", "/usr/certsuite/results",
		"--omit-artifacts-zip-file", "true",
		"--enable-data-collection", "true",
		"--sanitize-claim", "true",
		"--label-filter", testCaseName,
	}

	// print the command
	glog.V(5).Info(fmt.Sprintf("Running command: %s %s", containerEngine, strings.Join(certsuiteCmdArgs, " ")))

	cmd := exec.CommandContext(context.TODO(), containerEngine, certsuiteCmdArgs...)

	debugCertsuite, err := GetConfiguration().DebugCertsuite()
	if err != nil {
		return fmt.Errorf("failed to set env var CERTSUITE_LOG_LEVEL: %w", err)
	}

	if debugCertsuite {
		outfile := GetConfiguration().CreateLogFile(getTestSuiteName(testCaseName), tcNameForReport)

		defer outfile.Close()

		_, err = fmt.Fprintf(outfile, "Running test: %s\n", tcNameForReport)
		if err != nil {
			return fmt.Errorf("failed to write to debug file: %w", err)
		}

		cmd.Stdout = outfile
		cmd.Stderr = outfile
	}

	err = cmd.Run()
	if err != nil {
		errStr := fmt.Sprintf("failed to run tc: %s, err: %v, cmd: %s",
			testCaseName, err, cmd.String())
		if debugCertsuite && cmd.Stdout != nil {
			if outfile, ok := cmd.Stdout.(*os.File); ok {
				errStr += ", outFile=" + outfile.Name()
			}
		}

		return errors.New(errStr)
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

func getTestSuiteName(testCaseName string) string {
	if strings.Contains(testCaseName, globalparameters.NetworkSuiteName) {
		return globalparameters.NetworkSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.AffiliatedCertificationSuiteName) {
		return globalparameters.AffiliatedCertificationSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.LifecycleSuiteName) {
		return globalparameters.LifecycleSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.PlatformAlterationSuiteName) {
		return globalparameters.PlatformAlterationSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.ObservabilitySuiteName) {
		return globalparameters.ObservabilitySuiteName
	}

	if strings.Contains(testCaseName, globalparameters.AccessControlSuiteName) {
		return globalparameters.AccessControlSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.PerformanceSuiteName) {
		return globalparameters.PerformanceSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.ManageabilitySuiteName) {
		return globalparameters.ManageabilitySuiteName
	}

	if strings.Contains(testCaseName, globalparameters.OperatorSuiteName) {
		return globalparameters.OperatorSuiteName
	}

	panic(fmt.Sprintf("unable to retrieve test suite name from test case name %s", testCaseName))
}

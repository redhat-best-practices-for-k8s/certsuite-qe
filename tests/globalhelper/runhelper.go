package globalhelper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/container"
)

// LaunchTests stats tests based on given parameters.
func LaunchTests(testCaseName string, tcNameForReport string) error {
	containerEngine, err := container.SelectEngine()
	if err != nil {
		return fmt.Errorf("failed to select engine: %w", err)
	}

	err = os.Setenv("TNF_CONTAINER_CLIENT", containerEngine)
	if err != nil {
		return fmt.Errorf("failed to set TNF_CONTAINER_CLIENT: %w", err)
	}

	glog.V(5).Info(fmt.Sprintf("container engine set to %s", containerEngine))
	testArgs := []string{
		"-k", os.Getenv("KUBECONFIG"),
		"-c", GetConfiguration().General.DockerConfigDir + "/config",
		"-t", GetConfiguration().General.TnfConfigDir,
		"-o", GetConfiguration().General.TnfReportDir,
		"-i", fmt.Sprintf("%s:%s", GetConfiguration().General.TnfImage, GetConfiguration().General.TnfImageTag),
		"-l", testCaseName,
	}

	cmd := exec.Command(fmt.Sprintf("./%s", GetConfiguration().General.TnfEntryPointScript))
	cmd.Args = append(cmd.Args, testArgs...)
	cmd.Dir = GetConfiguration().General.TnfRepoPath

	debugTnf, err := GetConfiguration().DebugTnf()
	if err != nil {
		return fmt.Errorf("failed to set env var TNF_LOG_LEVEL: %w", err)
	}

	if debugTnf {
		outfile := GetConfiguration().CreateLogFile(getTestSuiteName(testCaseName), tcNameForReport)

		defer outfile.Close()

		_, err = outfile.WriteString(fmt.Sprintf("Running test: %s\n", tcNameForReport))
		if err != nil {
			return fmt.Errorf("failed to write to debug file: %w", err)
		}

		cmd.Stdout = outfile
		cmd.Stderr = outfile
	}

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run tc: %s, err: %w, cmd: %s",
			testCaseName, err, cmd.String())
	}

	CopyClaimFileToTcFolder(testCaseName, tcNameForReport)

	return err
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

	panic(fmt.Sprintf("unable to retrieve test suite name from test case name %s", testCaseName))
}

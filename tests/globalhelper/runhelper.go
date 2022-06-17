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
		return err
	}

	err = os.Setenv("TNF_CONTAINER_CLIENT", containerEngine)

	if err != nil {
		return err
	}

	glog.V(5).Info(fmt.Sprintf("container engine set to %s", containerEngine))
	testArgs := []string{
		"-k", os.Getenv("KUBECONFIG"),
		"-t", Configuration.General.TnfConfigDir,
		"-o", Configuration.General.TnfReportDir,
		"-i", fmt.Sprintf("%s:%s", Configuration.General.TnfImage, Configuration.General.TnfImageTag),
	}

	if len(testCaseName) > 0 {
		testArgs = append(testArgs, "-f")
		testArgs = append(testArgs, testCaseName)
		glog.V(5).Info(fmt.Sprintf("add test suite %s", testCaseName))
	} else {
		panic("No test suite name provided.")
	}

	cmd := exec.Command(fmt.Sprintf("./%s", Configuration.General.TnfEntryPointScript))
	cmd.Args = append(cmd.Args, testArgs...)
	cmd.Dir = Configuration.General.TnfRepoPath

	debugTnf, err := Configuration.DebugTnf()

	if err != nil {
		return err
	}

	if debugTnf {
		outfile := Configuration.CreateLogFile(getTestSuteName(testCaseName), tcNameForReport)

		defer outfile.Close()
		_, err = outfile.WriteString(fmt.Sprintf("Running test: %s\n", tcNameForReport))

		if err != nil {
			return err
		}

		cmd.Stdout = outfile
		cmd.Stderr = outfile
	}

	return cmd.Run()
}

func getTestSuteName(testCaseName string) string {
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

	panic(fmt.Sprintf("can't retrieve test suite name from test case name %s", testCaseName))
}

package globalhelper

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/container"
)

func LaunchTests(testSuite string, tcNameForReport string, skipRegEx string) error {
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

	if skipRegEx != "" {
		testArgs = append(testArgs, []string{"-s", skipRegEx}...)
		glog.V(5).Info(fmt.Sprintf("set skip regex to %s", skipRegEx))
	}

	if len(testSuite) > 0 {
		testArgs = append(testArgs, "-f")
		testArgs = append(testArgs, testSuite)
		glog.V(5).Info(fmt.Sprintf("add test suite %s", testSuite))
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
		outfile := Configuration.CreateLogFile(testSuite, tcNameForReport)

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

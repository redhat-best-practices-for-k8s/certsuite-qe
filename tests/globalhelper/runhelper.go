package globalhelper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golang/glog"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/container"
)

func LaunchTests(testSuites []string, tcNameForFolder string, tcNameForReport string, skipRegEx string) error {
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

	if len(testSuites) > 0 {
		testArgs = append(testArgs, "-f")
		for _, testSuite := range testSuites {
			testArgs = append(testArgs, testSuite)
			glog.V(5).Info(fmt.Sprintf("add test suite %s", testSuite))
		}
	}

	cmd := exec.Command(fmt.Sprintf("./%s", Configuration.General.TnfEntryPointScript))
	cmd.Args = append(cmd.Args, testArgs...)
	cmd.Dir = Configuration.General.TnfRepoPath

	if os.Getenv("DEBUG_TNF") == "true" {
		os.Setenv("LOG_LEVEL", "trace")

		folderPath := filepath.Join(Configuration.General.ReportDirAbsPath, "Debug", tcNameForFolder)

		_, err := VerifyFolderExists(folderPath, 0755)
		if err != nil {
			panic(err)
		}

		tcFile := filepath.Join(folderPath, tcNameForReport+".log")
		outfile, err := os.OpenFile(tcFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)

		if err != nil {
			panic(err)
		}

		defer outfile.Close()
		_, err = outfile.WriteString(fmt.Sprintf("Running test: %s\n", tcNameForReport))

		if err != nil {
			panic(err)
		}

		cmd.Stdout = outfile
		cmd.Stderr = outfile
	}

	return cmd.Run()
}

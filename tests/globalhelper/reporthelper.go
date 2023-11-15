package globalhelper

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/glog"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/test-network-function-claim/pkg/claim"
)

// OpenClaimReport opens claim.json file and returns struct.
func OpenClaimReport() (*claim.Root, error) {
	dataClaim, err := os.Open(path.Join(GetConfiguration().General.TnfReportDir, globalparameters.DefaultClaimFileName))
	if err != nil {
		return nil, fmt.Errorf("error opening %s report file: %w", globalparameters.DefaultClaimFileName, err)
	}

	byteValueClaim, err := io.ReadAll(dataClaim)

	if err != nil {
		return nil, fmt.Errorf("error reading %s report file: %w", globalparameters.DefaultClaimFileName, err)
	}

	var claimRootReport claim.Root
	err = json.Unmarshal(byteValueClaim, &claimRootReport)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling %s report file: %w", globalparameters.DefaultClaimFileName, err)
	}

	return &claimRootReport, nil
}

// IsTestCasePassedInClaimReport tests if test case is passed as expected in claim.json file.
func IsTestCasePassedInClaimReport(testCaseName string, claimReport claim.Root) (bool, error) {
	return isTestCaseInExpectedStatusInClaimReport(testCaseName, claimReport, globalparameters.TestCasePassed)
}

// IsTestCaseFailedInClaimReport test if test case is failed as expected in claim.json file.
func IsTestCaseFailedInClaimReport(testCaseName string, claimReport claim.Root) (bool, error) {
	return isTestCaseInExpectedStatusInClaimReport(testCaseName, claimReport, globalparameters.TestCaseFailed)
}

// IsTestCaseSkippedInClaimReport test if test case is failed as expected in claim.json file.
func IsTestCaseSkippedInClaimReport(testCaseName string, claimReport claim.Root) (bool, error) {
	return isTestCaseInExpectedStatusInClaimReport(testCaseName, claimReport, globalparameters.TestCaseSkipped)
}

// RemoveContentsFromReportDir removes all files from report dir.
func RemoveContentsFromReportDir() error {
	glog.V(5).Info(fmt.Sprintf("removing all files from %s directory", GetConfiguration().General.TnfReportDir))

	tnfReportDir, err := os.Open(GetConfiguration().General.TnfReportDir)
	if err != nil {
		return fmt.Errorf("failed to open report directory: %w", err)
	}

	defer tnfReportDir.Close()

	names, err := tnfReportDir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(GetConfiguration().General.TnfReportDir, name))
		if err != nil {
			return fmt.Errorf("failed to remove content from report directory: %w", err)
		}

		glog.V(5).Info(fmt.Sprintf("file %s removed from %s directory",
			name,
			GetConfiguration().General.TnfReportDir))
	}

	// Delete the report directory
	err = os.Remove(GetConfiguration().General.TnfReportDir)
	if err != nil {
		return fmt.Errorf("failed to remove report directory: %w", err)
	}

	return nil
}

func RemoveContentsFromConfigDir() error {
	tnfConfigDir, err := os.Open(GetConfiguration().General.TnfConfigDir)
	if err != nil {
		return fmt.Errorf("failed to open config directory: %w", err)
	}

	defer tnfConfigDir.Close()

	names, err := tnfConfigDir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(GetConfiguration().General.TnfConfigDir, name))
		if err != nil {
			return fmt.Errorf("failed to remove content from config directory: %w", err)
		}

		glog.V(5).Info(fmt.Sprintf("file %s removed from %s directory",
			name,
			GetConfiguration().General.TnfConfigDir))
	}

	// Delete the config directory
	err = os.Remove(GetConfiguration().General.TnfConfigDir)
	if err != nil {
		return fmt.Errorf("failed to remove config directory: %w", err)
	}

	return nil
}

// ConvertSpecNameToFileName converts given spec name to file name.
func ConvertSpecNameToFileName(specName string) string {
	formatString := specName
	for _, symbol := range []string{" ", ", ", "-", "/"} {
		formatString = strings.ReplaceAll(formatString, symbol, "_")
	}

	return strings.ToLower(removeCharactersFromString(formatString, []string{","}))
}

func isTestCaseInExpectedStatusInClaimReport(
	testCaseName string,
	claimReport claim.Root,
	expectedStatus string) (bool, error) {
	for testCaseClaim := range claimReport.Claim.Results {
		if formatTestCaseName(testCaseClaim) == formatTestCaseName(testCaseName) {
			var testCaseResult []*claim.Result

			encodedTestResult, err := json.Marshal(claimReport.Claim.Results[testCaseClaim])

			if err != nil {
				return false, err
			}

			err = json.Unmarshal(encodedTestResult, &testCaseResult)

			if err != nil {
				return false, err
			}

			if testCaseResult[0].State == expectedStatus {
				glog.V(5).Info("claim report test case status passed")

				return true, nil
			}

			return false, fmt.Errorf("invalid test status %s instead expected %s",
				testCaseResult[0].State,
				expectedStatus)
		}
	}

	return false, fmt.Errorf("test case is not found in the claim report")
}

func removeCharactersFromString(stringToFormat string, charactersToRemove []string) string {
	var formattedString string

	for index, element := range charactersToRemove {
		if index == 0 {
			formattedString = strings.ReplaceAll(stringToFormat, element, "")
		}

		formattedString = strings.ReplaceAll(formattedString, element, "")
	}

	return formattedString
}

func formatTestCaseName(tcName string) string {
	return removeCharactersFromString(tcName, []string{"-", "_", " ", "online,"})
}

func CopyClaimFileToTcFolder(tcName, formattedTcName string) {
	srcClaim := path.Join(GetConfiguration().General.TnfReportDir, globalparameters.DefaultClaimFileName)
	dstDir := path.Join(GetConfiguration().General.ReportDirAbsPath, "Debug", getTestSuiteName(tcName), formattedTcName)
	dstClaim := path.Join(dstDir, globalparameters.DefaultClaimFileName)

	_, err := os.Stat(srcClaim)
	if err != nil {
		glog.Error("file does not exist ", srcClaim)
	}

	// create destination folder
	err = os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		glog.Error("could not create dest directory= %s, err=%s", dstDir, err)
	}

	err = CopyFiles(srcClaim, dstClaim)
	if err != nil {
		glog.Fatalf("failed to copy %s to %s", srcClaim, dstClaim)
	}
}

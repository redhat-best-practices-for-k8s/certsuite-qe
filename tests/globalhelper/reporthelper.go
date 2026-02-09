package globalhelper

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	klog "k8s.io/klog/v2"

	"github.com/redhat-best-practices-for-k8s/certsuite-claim/pkg/claim"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
)

// OpenClaimReport opens claim.json file and returns struct.
func OpenClaimReport(reportDir string) (*claim.Root, error) {
	dataClaim, err := os.Open(path.Join(reportDir, globalparameters.DefaultClaimFileName))
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
// Returns nil if the directory is empty or doesn't exist (graceful cleanup).
func RemoveContentsFromReportDir(reportDir string) error {
	// Handle empty path gracefully (e.g., when test skips before directory creation).
	if reportDir == "" {
		klog.V(5).Info("report directory path is empty, skipping cleanup")

		return nil
	}

	// Check if directory exists before attempting cleanup.
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		klog.V(5).Info(fmt.Sprintf("report directory %s does not exist, skipping cleanup", reportDir))

		return nil
	}

	klog.V(5).Info(fmt.Sprintf("removing all files from %s directory", reportDir))

	certsuiteReportDir, err := os.Open(reportDir)
	if err != nil {
		return fmt.Errorf("failed to open report directory: %w", err)
	}

	defer certsuiteReportDir.Close()

	names, err := certsuiteReportDir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(reportDir, name))
		if err != nil {
			return fmt.Errorf("failed to remove content from report directory: %w", err)
		}

		klog.V(5).Info(fmt.Sprintf("file %s removed from %s directory",
			name,
			reportDir))
	}

	// Delete the report directory
	err = os.Remove(reportDir)
	if err != nil {
		return fmt.Errorf("failed to remove report directory: %w", err)
	}

	return nil
}

// RemoveContentsFromConfigDir removes all files from config dir.
// Returns nil if the directory is empty or doesn't exist (graceful cleanup).
func RemoveContentsFromConfigDir(configDir string) error {
	// Handle empty path gracefully (e.g., when test skips before directory creation).
	if configDir == "" {
		klog.V(5).Info("config directory path is empty, skipping cleanup")

		return nil
	}

	// Check if directory exists before attempting cleanup.
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		klog.V(5).Info(fmt.Sprintf("config directory %s does not exist, skipping cleanup", configDir))

		return nil
	}

	certsuiteConfigDir, err := os.Open(configDir)
	if err != nil {
		return fmt.Errorf("failed to open config directory: %w", err)
	}

	defer certsuiteConfigDir.Close()

	names, err := certsuiteConfigDir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(configDir, name))
		if err != nil {
			return fmt.Errorf("failed to remove content from config directory: %w", err)
		}

		klog.V(5).Info(fmt.Sprintf("file %s removed from %s directory",
			name,
			configDir))
	}

	// Delete the config directory
	err = os.Remove(configDir)
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
			var testCaseResult *claim.Result

			encodedTestResult, err := json.Marshal(claimReport.Claim.Results[testCaseClaim])

			if err != nil {
				return false, err
			}

			err = json.Unmarshal(encodedTestResult, &testCaseResult)

			if err != nil {
				return false, err
			}

			if testCaseResult.State == expectedStatus {
				klog.V(5).Info("claim report test case status passed")

				return true, nil
			}

			return false, fmt.Errorf("invalid test status %s instead expected %s",
				testCaseResult.State,
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

func CopyClaimFileToTcFolder(tcName, formattedTcName, reportDir string) {
	srcClaim := path.Join(reportDir, globalparameters.DefaultClaimFileName)
	dstDir := path.Join(GetConfiguration().General.ReportDirAbsPath, "Debug", getTestSuiteName(tcName), formattedTcName)
	dstClaim := path.Join(dstDir, globalparameters.DefaultClaimFileName)

	_, err := os.Stat(srcClaim)
	if err != nil {
		klog.Error("file does not exist ", srcClaim)
	}

	// create destination folder
	err = os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		klog.ErrorS(err, "could not create dest directory", "dir", dstDir)
	}

	err = CopyFiles(srcClaim, dstClaim)
	if err != nil {
		klog.ErrorS(err, "failed to copy claim file", "src", srcClaim, "dst", dstClaim)
	}
}

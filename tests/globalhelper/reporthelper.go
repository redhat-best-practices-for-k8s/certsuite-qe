package globalhelper

import (
	"encoding/json"
	"encoding/xml"
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
	dataClaim, err := os.Open(path.Join(Configuration.General.TnfReportDir, globalparameters.DefaultClaimFileName))
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
		return nil, fmt.Errorf("error unmarshaling %s report file: %w", globalparameters.DefaultClaimFileName, err)
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

// OpenJunitTestReport returns junit struct.
func OpenJunitTestReport() (*globalparameters.JUnitTestSuites, error) {
	junitReportFile, err := os.Open(
		path.Join(Configuration.General.TnfReportDir, globalparameters.DefaultJunitReportName),
	)
	if err != nil {
		return nil, fmt.Errorf("error opening %s report file: %w", globalparameters.DefaultJunitReportName, err)
	}

	junitReportByte, err := io.ReadAll(junitReportFile)

	if err != nil {
		return nil, fmt.Errorf("error reading %s report file: %w", globalparameters.DefaultJunitReportName, err)
	}

	var junitReport globalparameters.JUnitTestSuites
	err = xml.Unmarshal(junitReportByte, &junitReport)

	if err != nil {
		return nil, err
	}

	return &junitReport, nil
}

// IsTestCasePassedInJunitReport tests if test case is passed as expected in junit report file.
func IsTestCasePassedInJunitReport(report *globalparameters.JUnitTestSuites, testCaseName string) bool {
	return isTestCaseInRequiredStatusInJunitReport(report, testCaseName, globalparameters.TestCasePassed)
}

// IsTestCaseFailedInJunitReport tests if test case is failed as expected in junit report file.
func IsTestCaseFailedInJunitReport(report *globalparameters.JUnitTestSuites, testCaseName string) bool {
	return isTestCaseInRequiredStatusInJunitReport(report, testCaseName, globalparameters.TestCaseFailed)
}

// IsTestCaseSkippedInJunitReport tests if test case is skipped as expected in junit report file.
func IsTestCaseSkippedInJunitReport(report *globalparameters.JUnitTestSuites, testCaseName string) bool {
	return isTestCaseInRequiredStatusInJunitReport(report, testCaseName, globalparameters.TestCaseSkipped)
}

// RemoveContentsFromReportDir removes all files from report dir.
func RemoveContentsFromReportDir() error {
	tnfReportDir, err := os.Open(Configuration.General.TnfReportDir)
	if err != nil {
		return err
	}
	defer tnfReportDir.Close()
	names, err := tnfReportDir.Readdirnames(-1)

	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(Configuration.General.TnfReportDir, name))
		if err != nil {
			return err
		}

		glog.V(5).Info(fmt.Sprintf("file %s removed from %s directory",
			name,
			Configuration.General.TnfReportDir))
	}

	return nil
}

func isTestCaseInRequiredStatusInJunitReport(
	report *globalparameters.JUnitTestSuites,
	testCaseName string,
	status string) bool {
	for _, tc := range report.Suites[0].Testcases {
		if strings.Contains(formatTestCaseName(tc.Name), formatTestCaseName(testCaseName)) {
			glog.V(5).Info(fmt.Sprintf("test case status %s", tc.Status))

			return tc.Status == status
		}
	}

	return false
}

func isTestCaseInExpectedStatusInClaimReport(
	testCaseName string,
	claimReport claim.Root,
	expectedStatus string) (bool, error) {
	for testCaseClaim := range claimReport.Claim.Results {
		if formatTestCaseName(testCaseName) == formatTestCaseName(testCaseClaim) {
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
	return removeCharactersFromString(tcName, []string{"-", "_", " "})
}

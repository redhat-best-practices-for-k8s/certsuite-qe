package globalhelper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
	"github.com/stretchr/testify/assert"
)

func writeClaimReportForTest(t *testing.T, reportDir, testCaseName, state string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(reportDir, 0o755))
	claimJSON := fmt.Sprintf(`{
		"claim": {
			"configurations": {},
			"metadata": {"endTime": "", "startTime": ""},
			"nodes": {},
			"versions": {"certSuite": "", "claimFormat": ""},
			"results": {
				%q: {
					"capturedTestOutput": "",
					"catalogInfo": {"bestPracticeReference": "", "description": "", "exceptionProcess": "", "remediation": ""},
					"categoryClassification": {},
					"checkDetails": "",
					"duration": 0,
					"failureLineContent": "",
					"failureLocation": "",
					"skipReason": "",
					"startTime": "",
					"state": %q,
					"testID": {"id": "", "suite": ""}
				}
			}
		}
	}`, testCaseName, state)
	claimPath := filepath.Join(reportDir, globalparameters.DefaultClaimFileName)
	require.NoError(t, os.WriteFile(claimPath, []byte(claimJSON), 0o600))
}

func TestOpenClaimReport(t *testing.T) {
	t.Run("valid report", func(t *testing.T) {
		reportDir := t.TempDir()
		writeClaimReportForTest(t, reportDir, "test_case", globalparameters.TestCasePassed)

		report, err := OpenClaimReport(reportDir)
		if !assert.NoError(t, err) || !assert.NotNil(t, report) {
			return
		}

		if !assert.NotNil(t, report.Claim) {
			return
		}

		assert.Equal(t, globalparameters.TestCasePassed, report.Claim.Results["test_case"].State)
	})

	t.Run("missing report", func(t *testing.T) {
		_, err := OpenClaimReport(t.TempDir())
		if !assert.Error(t, err) {
			return
		}

		assert.Contains(t, err.Error(), globalparameters.DefaultClaimFileName)
	})

	t.Run("malformed report", func(t *testing.T) {
		reportDir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(reportDir, globalparameters.DefaultClaimFileName), []byte("{"), 0o600))

		_, err := OpenClaimReport(reportDir)
		if !assert.Error(t, err) {
			return
		}

		assert.Contains(t, err.Error(), "unmarshalling")
	})
}

func TestValidateIfReportsAreValid(t *testing.T) {
	testCases := []struct {
		name           string
		reportState    string
		expectedStatus string
		wantErr        bool
	}{
		{name: "passed", reportState: globalparameters.TestCasePassed, expectedStatus: globalparameters.TestCasePassed},
		{name: "failed", reportState: globalparameters.TestCaseFailed, expectedStatus: globalparameters.TestCaseFailed},
		{name: "skipped", reportState: globalparameters.TestCaseSkipped, expectedStatus: globalparameters.TestCaseSkipped},
		{
			name:           "status mismatch",
			reportState:    globalparameters.TestCasePassed,
			expectedStatus: globalparameters.TestCaseFailed,
			wantErr:        true,
		},
		{name: "invalid expected status", reportState: globalparameters.TestCasePassed, expectedStatus: "invalid", wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reportDir := t.TempDir()
			writeClaimReportForTest(t, reportDir, "test_case", testCase.reportState)

			err := ValidateIfReportsAreValid("test-case", testCase.expectedStatus, reportDir)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("test case missing from report", func(t *testing.T) {
		reportDir := t.TempDir()
		writeClaimReportForTest(t, reportDir, "different_test", globalparameters.TestCasePassed)

		err := ValidateIfReportsAreValid("test_case", globalparameters.TestCasePassed, reportDir)
		if !assert.Error(t, err) {
			return
		}

		assert.Contains(t, err.Error(), "test case is not found")
	})

	t.Run("missing report file", func(t *testing.T) {
		assert.Error(t, ValidateIfReportsAreValid("test", globalparameters.TestCasePassed, t.TempDir()))
	})

	t.Run("expected failed but report passed", func(t *testing.T) {
		reportDir := t.TempDir()
		writeClaimReportForTest(t, reportDir, "test", globalparameters.TestCasePassed)
		assert.Error(t, ValidateIfReportsAreValid("test", globalparameters.TestCaseFailed, reportDir))
	})
}

func TestValidateIfReportsAreValidWithAcceptedStatuses(t *testing.T) {
	testCases := []struct {
		name             string
		reportState      string
		acceptedStatuses []string
		wantErr          bool
	}{
		{
			name:             "first status accepted",
			reportState:      globalparameters.TestCasePassed,
			acceptedStatuses: []string{globalparameters.TestCasePassed, globalparameters.TestCaseSkipped},
		},
		{
			name:             "later status accepted",
			reportState:      globalparameters.TestCaseSkipped,
			acceptedStatuses: []string{"invalid", globalparameters.TestCasePassed, globalparameters.TestCaseSkipped},
		},
		{
			name:             "no accepted status matches",
			reportState:      globalparameters.TestCaseFailed,
			acceptedStatuses: []string{globalparameters.TestCasePassed},
			wantErr:          true,
		},
		{
			name:             "invalid statuses skipped",
			reportState:      globalparameters.TestCasePassed,
			acceptedStatuses: []string{"invalid"},
			wantErr:          true,
		},
		{name: "empty accepted statuses", reportState: globalparameters.TestCasePassed, wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reportDir := t.TempDir()
			writeClaimReportForTest(t, reportDir, "test_case", testCase.reportState)

			err := ValidateIfReportsAreValidWithAcceptedStatuses("test_case", testCase.acceptedStatuses, reportDir)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("accept failed result", func(t *testing.T) {
		reportDir := t.TempDir()
		writeClaimReportForTest(t, reportDir, "test_case", globalparameters.TestCaseFailed)
		assert.NoError(t, ValidateIfReportsAreValidWithAcceptedStatuses("test_case", []string{globalparameters.TestCaseFailed}, reportDir))
	})

	t.Run("missing report file", func(t *testing.T) {
		assert.Error(t, ValidateIfReportsAreValidWithAcceptedStatuses("test", []string{globalparameters.TestCasePassed}, t.TempDir()))
	})
}

func TestRemoveContentsFromReportDir(t *testing.T) {
	t.Run("empty and missing paths are safe", func(t *testing.T) {
		assert.NoError(t, RemoveContentsFromReportDir(""))
		assert.NoError(t, RemoveContentsFromReportDir(filepath.Join(t.TempDir(), "missing")))
	})

	t.Run("removes files and nested directories", func(t *testing.T) {
		reportDir := filepath.Join(t.TempDir(), "report")
		nestedDir := filepath.Join(reportDir, "nested")
		require.NoError(t, os.MkdirAll(nestedDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(nestedDir, "claim.json"), []byte("{}"), 0o600))

		require.NoError(t, RemoveContentsFromReportDir(reportDir))
		_, err := os.Stat(reportDir)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("reports directory read and filesystem errors", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "file")
		require.NoError(t, os.WriteFile(filePath, []byte("x"), 0o600))
		assert.Error(t, removeContentsFromDir(filePath, "report"))

		oldOpen, oldRemoveContent, oldRemoveDir := openReportDirectory, removeReportContent, removeReportDir

		t.Cleanup(func() {
			openReportDirectory, removeReportContent, removeReportDir = oldOpen, oldRemoveContent, oldRemoveDir
		})
		openReportDirectory = func(string) (*os.File, error) { return nil, fmt.Errorf("open failed") }
		existingDir := t.TempDir()
		assert.Error(t, removeContentsFromDir(existingDir, "report"))

		openReportDirectory = os.Open
		contentDir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(contentDir, "file"), []byte("x"), 0o600))
		removeReportContent = func(string) error { return fmt.Errorf("remove content failed") }

		assert.Error(t, removeContentsFromDir(contentDir, "report"))

		removeReportContent = os.RemoveAll
		removeReportDir = func(string) error { return fmt.Errorf("remove dir failed") }

		assert.Error(t, removeContentsFromDir(contentDir, "report"))
	})
}

func TestGetTestCaseResultMarshalErrors(t *testing.T) {
	reportDir := t.TempDir()
	writeClaimReportForTest(t, reportDir, "test", globalparameters.TestCasePassed)
	report, err := OpenClaimReport(reportDir)
	require.NoError(t, err)
	oldMarshal, oldUnmarshal := marshalClaimResult, unmarshalClaimResult

	t.Cleanup(func() { marshalClaimResult, unmarshalClaimResult = oldMarshal, oldUnmarshal })
	marshalClaimResult = func(any) ([]byte, error) { return nil, fmt.Errorf("marshal failed") }
	_, err = getTestCaseResult("test", *report)
	assert.Error(t, err)
	marshalClaimResult = oldMarshal
	unmarshalClaimResult = func([]byte, any) error { return fmt.Errorf("unmarshal failed") }
	_, err = getTestCaseResult("test", *report)
	assert.Error(t, err)
}

func TestRemoveContentsFromConfigDir(t *testing.T) {
	configDir := filepath.Join(t.TempDir(), "config")
	require.NoError(t, os.Mkdir(configDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.yml"), []byte("content"), 0o600))

	require.NoError(t, RemoveContentsFromConfigDir(configDir))
	_, err := os.Stat(configDir)
	assert.True(t, os.IsNotExist(err))
}

func TestConvertSpecNameToFileName(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{input: "Networking, default-network", want: "networking_default_network"},
		{input: "Lifecycle / Pod Startup", want: "lifecycle___pod_startup"},
		{input: "Already_lowercase", want: "already_lowercase"},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.want, ConvertSpecNameToFileName(testCase.input))
	}
}

func TestCopyClaimFileToTcFolderFilesystemFailures(t *testing.T) {
	oldConfig, oldMkdir, oldCopy := conf, makeReportDir, copyFileContents

	t.Cleanup(func() { conf, makeReportDir, copyFileContents = oldConfig, oldMkdir, oldCopy })
	conf = &config.Config{}
	conf.General.ReportDirAbsPath = t.TempDir()
	reportDir := t.TempDir()
	writeClaimReportForTest(t, reportDir, "networking-test", globalparameters.TestCasePassed)

	makeReportDir = func(string, os.FileMode) error { return fmt.Errorf("mkdir failed") }

	CopyClaimFileToTcFolder("networking-test", "copy-fail", reportDir)

	makeReportDir = os.MkdirAll
	copyFileContents = func(io.Writer, io.Reader) (int64, error) { return 0, fmt.Errorf("copy failed") }

	CopyClaimFileToTcFolder("networking-test", "copy-fail", reportDir)
}

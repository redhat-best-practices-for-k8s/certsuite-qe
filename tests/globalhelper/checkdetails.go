package globalhelper

import (
	"encoding/json"
	"fmt"
)

// ReportObject mirrors the certsuite's testhelper.ReportObject structure.
// Defined locally because certsuite-qe cannot import certsuite's internal packages.
type ReportObject struct {
	ObjectType         string   `json:"ObjectType"`
	ObjectFieldsKeys   []string `json:"ObjectFieldsKeys"`
	ObjectFieldsValues []string `json:"ObjectFieldsValues"`
}

// CheckDetails mirrors the certsuite's testhelper.CheckDetails structure.
type CheckDetails struct {
	CompliantObjectsOut    []*ReportObject `json:"CompliantObjectsOut"`
	NonCompliantObjectsOut []*ReportObject `json:"NonCompliantObjectsOut"`
}

// ParseCheckDetails unmarshals a CheckDetails JSON string from a claim result's CheckDetails field.
func ParseCheckDetails(checkDetailsStr string) (*CheckDetails, error) {
	if checkDetailsStr == "" {
		return nil, fmt.Errorf("checkDetails string is empty")
	}

	var details CheckDetails

	err := json.Unmarshal([]byte(checkDetailsStr), &details)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkDetails: %w", err)
	}

	return &details, nil
}

// GetTestCaseCheckDetails opens the claim report, finds the named test case, and parses its CheckDetails.
func GetTestCaseCheckDetails(tcName, reportDir string) (*CheckDetails, error) {
	claimReport, err := OpenClaimReport(reportDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open claim report: %w", err)
	}

	result, err := getTestCaseResult(tcName, *claimReport)
	if err != nil {
		return nil, fmt.Errorf("failed to get test case result for %q: %w", tcName, err)
	}

	details, err := ParseCheckDetails(result.CheckDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to parse checkDetails for %q: %w", tcName, err)
	}

	return details, nil
}

// GetReportObjectFieldValue returns the value for a given key from a ReportObject's parallel arrays.
// Returns an empty string if the key is not found.
func GetReportObjectFieldValue(obj *ReportObject, key string) string {
	for i, k := range obj.ObjectFieldsKeys {
		if k == key && i < len(obj.ObjectFieldsValues) {
			return obj.ObjectFieldsValues[i]
		}
	}

	return ""
}

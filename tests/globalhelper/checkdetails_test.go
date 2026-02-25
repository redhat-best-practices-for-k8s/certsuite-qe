package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCheckDetails(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		details, err := ParseCheckDetails("")
		assert.Nil(t, details)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
	})

	t.Run("valid JSON with both arrays", func(t *testing.T) {
		input := `{
			"CompliantObjectsOut": [
				{
					"ObjectType": "Pod",
					"ObjectFieldsKeys": ["Name", "Reason"],
					"ObjectFieldsValues": ["test-pod", "process disappeared before scheduling check"]
				}
			],
			"NonCompliantObjectsOut": [
				{
					"ObjectType": "Pod",
					"ObjectFieldsKeys": ["Name", "Reason"],
					"ObjectFieldsValues": ["bad-pod", "RT scheduling policy detected"]
				}
			]
		}`

		details, err := ParseCheckDetails(input)
		assert.NoError(t, err)
		assert.Len(t, details.CompliantObjectsOut, 1)
		assert.Len(t, details.NonCompliantObjectsOut, 1)
		assert.Equal(t, "Pod", details.CompliantObjectsOut[0].ObjectType)
		assert.Equal(t, "Pod", details.NonCompliantObjectsOut[0].ObjectType)
	})

	t.Run("valid JSON with empty arrays", func(t *testing.T) {
		input := `{"CompliantObjectsOut": [], "NonCompliantObjectsOut": []}`

		details, err := ParseCheckDetails(input)
		assert.NoError(t, err)
		assert.Empty(t, details.CompliantObjectsOut)
		assert.Empty(t, details.NonCompliantObjectsOut)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		details, err := ParseCheckDetails("{not valid json")
		assert.Nil(t, details)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal")
	})
}

func TestGetReportObjectFieldValue(t *testing.T) {
	obj := &ReportObject{
		ObjectType:         "Pod",
		ObjectFieldsKeys:   []string{"Name", "Namespace", "Reason"},
		ObjectFieldsValues: []string{"test-pod", "default", "compliant"},
	}

	t.Run("key found", func(t *testing.T) {
		assert.Equal(t, "test-pod", GetReportObjectFieldValue(obj, "Name"))
		assert.Equal(t, "default", GetReportObjectFieldValue(obj, "Namespace"))
		assert.Equal(t, "compliant", GetReportObjectFieldValue(obj, "Reason"))
	})

	t.Run("key not found", func(t *testing.T) {
		assert.Equal(t, "", GetReportObjectFieldValue(obj, "NonExistentKey"))
	})

	t.Run("mismatched array lengths", func(t *testing.T) {
		shortObj := &ReportObject{
			ObjectFieldsKeys:   []string{"Name", "Extra"},
			ObjectFieldsValues: []string{"test-pod"},
		}
		// "Name" is at index 0, which is within values range
		assert.Equal(t, "test-pod", GetReportObjectFieldValue(shortObj, "Name"))
		// "Extra" is at index 1, which is out of values range
		assert.Equal(t, "", GetReportObjectFieldValue(shortObj, "Extra"))
	})
}

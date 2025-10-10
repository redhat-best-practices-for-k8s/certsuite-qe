package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//nolint:funlen
func TestVerifyRHCOSVersions(t *testing.T) {
	testCases := []struct {
		name                   string
		nodes                  []corev1.Node
		rhcosVersionMapContent string
		expectError            bool
		errorContains          string
	}{
		{
			name: "Valid single node with version in map",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 410.84.202205031645-0 (Ootpa)",
						},
					},
				},
			},
			rhcosVersionMapContent: `410.84.202205031645-0
411.85.202206011234-0
412.86.202207011234-0`,
			expectError: false,
		},
		{
			name: "Valid multiple nodes with versions in map",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 410.84.202205031645-0 (Ootpa)",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-2",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 411.85.202206011234-0 (Plow)",
						},
					},
				},
			},
			rhcosVersionMapContent: `410.84.202205031645-0
411.85.202206011234-0
412.86.202207011234-0`,
			expectError: false,
		},
		{
			name: "Version not found in map",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 999.99.202299999999-0 (Unknown)",
						},
					},
				},
			},
			rhcosVersionMapContent: `410.84.202205031645-0
411.85.202206011234-0
412.86.202207011234-0`,
			expectError:   true,
			errorContains: "not found in rhcos_version_map",
		},
		{
			name:                   "Empty nodes list",
			nodes:                  []corev1.Node{},
			rhcosVersionMapContent: "410.84.202205031645-0",
			expectError:            true,
			errorContains:          "no nodes provided",
		},
		{
			name: "Empty rhcos_version_map content",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 410.84.202205031645-0 (Ootpa)",
						},
					},
				},
			},
			rhcosVersionMapContent: "",
			expectError:            true,
			errorContains:          "rhcos_version_map content is empty",
		},
		{
			name: "Invalid OSImage format - missing RHCOS name",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Ubuntu 20.04 LTS",
						},
					},
				},
			},
			rhcosVersionMapContent: "410.84.202205031645-0",
			expectError:            true,
			errorContains:          "does not contain Red Hat Enterprise Linux CoreOS",
		},
		{
			name: "Invalid OSImage format - no version after RHCOS name",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS",
						},
					},
				},
			},
			rhcosVersionMapContent: "410.84.202205031645-0",
			expectError:            true,
			errorContains:          "cannot extract version",
		},
		{
			name: "Version with different codename",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 410.84.202205031645-0 (Plow)",
						},
					},
				},
			},
			rhcosVersionMapContent: "410.84.202205031645-0",
			expectError:            false,
		},
		{
			name: "Version without parentheses",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 410.84.202205031645-0",
						},
					},
				},
			},
			rhcosVersionMapContent: "410.84.202205031645-0",
			expectError:            false,
		},
		{
			name: "Mixed versions - one valid, one invalid",
			nodes: []corev1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-1",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 410.84.202205031645-0 (Ootpa)",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node-2",
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							OSImage: "Red Hat Enterprise Linux CoreOS 999.99.202299999999-0 (Unknown)",
						},
					},
				},
			},
			rhcosVersionMapContent: `410.84.202205031645-0
411.85.202206011234-0`,
			expectError:   true,
			errorContains: "999.99.202299999999-0",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := verifyRHCOSVersions(testCase.nodes, testCase.rhcosVersionMapContent)

			if testCase.expectError {
				assert.NotNil(t, err, "Expected an error but got none")

				if testCase.errorContains != "" {
					assert.Contains(t, err.Error(), testCase.errorContains,
						"Error message should contain expected substring")
				}
			} else {
				assert.Nil(t, err, "Expected no error but got: %v", err)
			}
		})
	}
}

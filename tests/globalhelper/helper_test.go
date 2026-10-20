package globalhelper

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDefinePodUnderTestLabels(t *testing.T) {
	testConfig := globalparameters.CertsuiteConfig{}
	assert.Empty(t, testConfig.PodsUnderTestLabels)

	testPodLabels := []string{"test1:value", "test2:value"}
	assert.Nil(t, definePodUnderTestLabels(&testConfig, testPodLabels))
	assert.Equal(t, testPodLabels, testConfig.PodsUnderTestLabels)

	// config is nil
	assert.NotNil(t, definePodUnderTestLabels(nil, []string{}))
	assert.NotNil(t, definePodUnderTestLabels(nil, []string{"valid:value"}))
	assert.Error(t, definePodUnderTestLabels(&globalparameters.CertsuiteConfig{}, []string{"invalid"}))
}

func TestIsExpectedStatusParamValid(t *testing.T) {
	testCases := []struct {
		expectedStatus string
		expectedError  error
	}{
		{
			expectedStatus: globalparameters.TestCasePassed,
			expectedError:  nil,
		},
		{
			expectedStatus: globalparameters.TestCaseFailed,
			expectedError:  nil,
		},
		{
			expectedStatus: globalparameters.TestCaseSkipped,
			expectedError:  nil,
		},
		{
			expectedStatus: "invalid",
			expectedError:  errors.New("this is an error)"),
		},
	}

	for _, tc := range testCases {
		err := IsExpectedStatusParamValid(tc.expectedStatus)
		if tc.expectedError != nil {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestAppendContainersToDeployment(t *testing.T) {
	dep := appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	// 2 expected containers
	AppendContainersToDeployment(&dep, 2, "nginx")
	assert.Equal(t, 2, len(dep.Spec.Template.Spec.Containers))
	assert.Equal(t, "container1", dep.Spec.Template.Spec.Containers[0].Name)
	assert.Equal(t, "container2", dep.Spec.Template.Spec.Containers[1].Name)
}

func TestAppendContainersToStatefulSet(t *testing.T) {
	statefulSet := appsv1.StatefulSet{}

	AppendContainersToStatefulSet(&statefulSet, 2, "nginx")

	if !assert.Len(t, statefulSet.Spec.Template.Spec.Containers, 2) {
		return
	}

	assert.Equal(t, "container1", statefulSet.Spec.Template.Spec.Containers[0].Name)
	assert.Equal(t, "container2", statefulSet.Spec.Template.Spec.Containers[1].Name)
	assert.Equal(t, "nginx", statefulSet.Spec.Template.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/bin/bash", "-c", "sleep INF"}, statefulSet.Spec.Template.Spec.Containers[0].Command)
}

func TestDefineCertifiedContainersInfo(t *testing.T) {
	testCases := []struct {
		name      string
		input     []string
		wantCount int
		wantInfo  *globalparameters.CertifiedContainerRepoInfo
		wantErr   bool
	}{
		{name: "optional empty list", wantCount: 0},
		{
			name:      "single value adds empty info",
			input:     []string{"registry/repo"},
			wantCount: 1,
		},
		{
			name:      "four fields are trimmed",
			input:     []string{" repo ; registry ; tag ; digest "},
			wantCount: 1,
			wantInfo: &globalparameters.CertifiedContainerRepoInfo{
				Repository: "repo",
				Registry:   "registry",
				Tag:        "tag",
				Digest:     "digest",
			},
		},
		{name: "invalid field count", input: []string{"repo;registry;tag"}, wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := globalparameters.CertsuiteConfig{}

			err := defineCertifiedContainersInfo(&config, testCase.input)
			if testCase.wantErr {
				assert.Error(t, err)

				return
			}

			if !assert.NoError(t, err) {
				return
			}

			if !assert.Len(t, config.Certifiedcontainerinfo, testCase.wantCount) {
				return
			}

			if testCase.wantInfo != nil {
				assert.Equal(t, *testCase.wantInfo, config.Certifiedcontainerinfo[0])
			}
		})
	}
}

func TestDefineCertsuiteConfig(t *testing.T) {
	testCases := []struct {
		name            string
		jobID           string
		wantPartnerName string
	}{
		{name: "without job ID", wantPartnerName: PartnerName},
		{name: "with job ID", jobID: "1234", wantPartnerName: PartnerName + "_1234"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("JOB_ID", testCase.jobID)
			configDir := t.TempDir()

			err := DefineCertsuiteConfig(
				[]string{"namespace-a"},
				[]string{"pod-label:value"},
				[]string{"operator-label:value"},
				[]string{"repo;registry;tag;digest"},
				[]string{"example.io"},
				configDir,
			)
			if !assert.NoError(t, err) {
				return
			}

			configYAML, err := os.ReadFile(filepath.Join(configDir, globalparameters.DefaultCertsuiteConfigFileName))
			if !assert.NoError(t, err) {
				return
			}

			var config globalparameters.CertsuiteConfig
			if !assert.NoError(t, yaml.Unmarshal(configYAML, &config)) {
				return
			}

			assert.Equal(t, []globalparameters.TargetNameSpace{{Name: "namespace-a"}}, config.TargetNameSpaces)
			assert.Equal(t, []string{"pod-label:value"}, config.PodsUnderTestLabels)
			assert.Equal(t, []string{"operator-label:value"}, config.OperatorsUnderTestLabels)
			wantCertifiedContainerInfo := []globalparameters.CertifiedContainerRepoInfo{{
				Repository: "repo",
				Registry:   "registry",
				Tag:        "tag",
				Digest:     "digest",
			}}

			assert.Equal(t, wantCertifiedContainerInfo, config.Certifiedcontainerinfo)
			assert.Equal(t, []globalparameters.TargetCrdFilter{{NameSuffix: "example.io"}}, config.TargetCrdFilters)
			assert.Equal(t, ExecutedBy, config.ExecutedBy)
			assert.Equal(t, testCase.wantPartnerName, config.PartnerName)
		})
	}

	t.Run("empty namespace list is rejected", func(t *testing.T) {
		err := DefineCertsuiteConfig(nil, nil, nil, nil, nil, t.TempDir())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target namespaces cannot be empty")
	})

	t.Run("invalid pod labels are rejected", func(t *testing.T) {
		err := DefineCertsuiteConfig([]string{"ns"}, []string{"invalid"}, nil, nil, nil, t.TempDir())
		assert.Error(t, err)
	})

	t.Run("invalid certified container details are rejected", func(t *testing.T) {
		err := DefineCertsuiteConfig([]string{"ns"}, []string{"label:value"}, nil, []string{"repo;registry"}, nil, t.TempDir())
		assert.Error(t, err)
	})

	t.Run("reports config encoding failures", func(t *testing.T) {
		oldEncode := encodeCertsuiteConfig
		encodeCertsuiteConfig = func(*os.File, globalparameters.CertsuiteConfig) error { return errors.New("encode failed") }

		t.Cleanup(func() { encodeCertsuiteConfig = oldEncode })
		err := DefineCertsuiteConfig([]string{"ns"}, []string{"label:value"}, nil, nil, nil, t.TempDir())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encode failed")
	})
}

func TestGenerateRandomString(t *testing.T) {
	assert.Empty(t, GenerateRandomString(0))

	generated := GenerateRandomString(32)
	if !assert.Len(t, generated, 32) {
		return
	}

	for _, char := range generated {
		assert.Contains(t, "abcdefghijklmnopqrstuvwxyz", string(char))
	}
}

func TestRunShellCommand(t *testing.T) {
	t.Run("returns combined output on success", func(t *testing.T) {
		output, err := RunShellCommand("printf 'stdout'; printf 'stderr' >&2")
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, "stdoutstderr", output)
	})

	t.Run("returns output and command context on failure", func(t *testing.T) {
		output, err := RunShellCommand("printf 'failure output'; exit 7")
		if !assert.Error(t, err) {
			return
		}

		assert.Equal(t, "failure output", output)
		assert.Contains(t, err.Error(), "exit status 7")
		assert.Contains(t, err.Error(), "failure output")
	})
}

func TestIsTransientDaemonSetError(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil error", want: false},
		{name: "not schedulable", err: errors.New("pod is not schedulable"), want: true},
		{name: "timeout", err: errors.New("Timed out waiting for daemonset"), want: true},
		{name: "not running", err: errors.New("container is not running"), want: true},
		{name: "not ready", err: errors.New("daemonset not ready"), want: true},
		{name: "permanent error", err: errors.New("invalid image configuration"), want: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, IsTransientDaemonSetError(testCase.err))
		})
	}
}

func TestCopyFiles(t *testing.T) {
	t.Run("copies file contents", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "source.txt")
		dst := filepath.Join(dir, "destination.txt")

		require.NoError(t, os.WriteFile(src, []byte("contents"), 0o600))

		require.NoError(t, CopyFiles(src, dst))
		contents, err := os.ReadFile(dst)
		require.NoError(t, err)
		assert.Equal(t, "contents", string(contents))
	})

	t.Run("reports missing source", func(t *testing.T) {
		err := CopyFiles(filepath.Join(t.TempDir(), "missing"), filepath.Join(t.TempDir(), "destination"))
		if !assert.Error(t, err) {
			return
		}

		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("reports copy errors", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "source.txt")
		require.NoError(t, os.WriteFile(src, []byte("contents"), 0o600))
		oldCopy := copyFileContents
		copyFileContents = func(_ io.Writer, _ io.Reader) (int64, error) { return 0, errors.New("copy failed") }

		t.Cleanup(func() { copyFileContents = oldCopy })
		err := CopyFiles(src, filepath.Join(dir, "destination.txt"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "copy failed")
	})
}

func TestDefineCertsuiteNamespaces(t *testing.T) {
	testCases := []struct {
		testNamespaces []string
		nilConfig      bool
	}{
		{
			testNamespaces: []string{"test1"},
		},
		{
			testNamespaces: []string{"test1", "test2"},
		},
		{
			testNamespaces: []string{},
		},
		{
			testNamespaces: []string{"test1", "test2"},
			nilConfig:      true,
		},
	}

	for _, testCase := range testCases {
		if testCase.nilConfig {
			assert.NotNil(t, defineCertsuiteNamespaces(nil, testCase.testNamespaces))
		} else {
			config := globalparameters.CertsuiteConfig{}
			if len(testCase.testNamespaces) > 0 {
				assert.Nil(t, defineCertsuiteNamespaces(&config, testCase.testNamespaces))
			} else {
				assert.NotNil(t, defineCertsuiteNamespaces(&config, testCase.testNamespaces))
				assert.Equal(t, len(testCase.testNamespaces), len(config.TargetNameSpaces))
			}
		}
	}
}

func TestDefineOperatorsUnderTestLabels(t *testing.T) {
	testConfig := globalparameters.CertsuiteConfig{}
	assert.Empty(t, testConfig.OperatorsUnderTestLabels)

	testOperatorLabels := []string{"test1:value", "test2:value"}
	defineOperatorsUnderTestLabels(&testConfig, testOperatorLabels)
	assert.Equal(t, testOperatorLabels, testConfig.OperatorsUnderTestLabels)
}

func TestDefineCrdFilters(t *testing.T) {
	testCases := [][]string{
		{"test1"},
		{"test1", "test2"},
		{},
	}

	for _, testCrdFilters := range testCases {
		config := globalparameters.CertsuiteConfig{}
		defineCrdFilters(&config, testCrdFilters)
		assert.Len(t, config.TargetCrdFilters, len(testCrdFilters))

		for index, suffix := range testCrdFilters {
			assert.Equal(t, suffix, config.TargetCrdFilters[index].NameSuffix)
		}
	}
}

func TestValidateIfParamInAllowedListOfParams(t *testing.T) {
	testCases := []struct {
		testParam         string
		testAllowedParams []string
		expectedError     error
	}{
		{
			testParam:         "test1",
			testAllowedParams: []string{"test1", "test2"},
			expectedError:     nil,
		},
		{
			testParam:         "test1",
			testAllowedParams: []string{"test2", "test3"},
			expectedError:     errors.New("this is an error"),
		},
	}

	for _, testCase := range testCases {
		err := validateIfParamInAllowedListOfParams(testCase.testParam, testCase.testAllowedParams)
		if testCase.expectedError != nil {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

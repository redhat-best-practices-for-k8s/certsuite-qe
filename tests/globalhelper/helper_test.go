package globalhelper

import (
	"errors"
	"testing"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, defineOperatorsUnderTestLabels(&testConfig, testOperatorLabels))
	assert.Equal(t, testOperatorLabels, testConfig.OperatorsUnderTestLabels)

	// config is nil
	assert.NotNil(t, defineOperatorsUnderTestLabels(nil, []string{}))
}

func TestDefineCrdFilters(t *testing.T) {
	testCases := []struct {
		testCrdFilters []string
		nilConfig      bool
	}{
		{
			testCrdFilters: []string{"test1"},
		},
		{
			testCrdFilters: []string{"test1", "test2"},
		},
		{
			testCrdFilters: []string{},
		},
		{
			testCrdFilters: []string{"test1", "test2"},
			nilConfig:      true,
		},
	}

	for _, testCase := range testCases {
		if testCase.nilConfig {
			assert.NotNil(t, defineCrdFilters(nil, testCase.testCrdFilters))
		} else {
			config := globalparameters.CertsuiteConfig{}
			if len(testCase.testCrdFilters) > 0 {
				assert.Nil(t, defineCrdFilters(&config, testCase.testCrdFilters))
			} else {
				assert.Nil(t, defineCrdFilters(&config, testCase.testCrdFilters))
				assert.Equal(t, len(testCase.testCrdFilters), len(config.TargetCrdFilters))
			}
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

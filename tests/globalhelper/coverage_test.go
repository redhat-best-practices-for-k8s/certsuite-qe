//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	semver "github.com/blang/semver/v4"
	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/onsi/gomega"
	egipod "github.com/openshift-kni/eco-goinfra/pkg/pod"
	packagemanifestv1 "github.com/openshift-kni/eco-goinfra/pkg/schemes/olm/package-server/operators/v1"
	packageversion "github.com/openshift-kni/eco-goinfra/pkg/schemes/olm/version"
	configv1api "github.com/openshift/api/config/v1"
	configv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	machinev1 "github.com/openshift/client-go/machineconfiguration/clientset/versioned/typed/machineconfiguration/v1"
	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/client"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	nodev1 "k8s.io/api/node/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestCoverageInitAndAccessors(t *testing.T) {
	client, _ := coverageFixture(t)
	assert.Same(t, client, GetAPIClient().K8sClient)
	SetTestK8sAPIClient(client)
	assert.Same(t, client, GetAPIClient().K8sClient)

	old := apiclient

	UnsetTestK8sAPIClient()
	assert.Nil(t, apiclient)
	SetTestK8sAPIClient(client)
	assert.Same(t, client, GetAPIClient().K8sClient)
	apiclient = old

	conf = &config.Config{}
	assert.Same(t, conf, GetConfiguration())
}

func TestCoverageLazyClientErrorsAndEcoClient(t *testing.T) {
	oldAPIClientFactory, oldConfigFactory, oldFatal := newAPIClient, newConfig, fatalf

	t.Cleanup(func() { newAPIClient, newConfig, fatalf = oldAPIClientFactory, oldConfigFactory, oldFatal })

	apiclient = nil
	newAPIClient = func() (*client.ClientSet, error) { return nil, errors.New("client init failed") }
	fatalf = func(string, ...any) {}

	assert.Nil(t, GetAPIClient())

	conf = nil
	newConfig = func() (*config.Config, error) { return nil, errors.New("config init failed") }

	assert.Nil(t, GetConfiguration())

	oldEGIClient := egiClient

	t.Cleanup(func() { egiClient = oldEGIClient })
	egiClient = nil

	t.Setenv("KUBECONFIG", filepath.Join(t.TempDir(), "does-not-exist"))
	assert.Nil(t, GetEcoGoinfraClient())
}

func TestCoverageGenerateDirectoryFailures(t *testing.T) {
	old := conf

	t.Cleanup(func() { conf = old })
	root := t.TempDir()
	blocker := filepath.Join(root, "file")
	require.NoError(t, os.WriteFile(blocker, []byte("x"), 0o600))
	conf = &config.Config{}
	conf.General.CertsuiteReportDir = filepath.Join(blocker, "reports")

	if _, _, err := GenerateDirectories("case"); assert.Error(t, err) {
		assert.Contains(t, err.Error(), "report directory")
	}
	conf.General.CertsuiteReportDir = filepath.Join(root, "reports")
	conf.General.CertsuiteConfigDir = filepath.Join(blocker, "configs")

	if _, _, err := GenerateDirectories("case"); assert.Error(t, err) {
		assert.Contains(t, err.Error(), "config directory")
	}
}

func TestCoverageOCPChecks(t *testing.T) {
	k8s, _ := coverageFixture(t)
	worker := &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "worker-0", Labels: map[string]string{"node-role.kubernetes.io/worker": ""}}}
	_, err := k8s.CoreV1().Nodes().Create(context.Background(), worker, metav1.CreateOptions{})
	require.NoError(t, err)

	setOCPRESTClients(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		switch request.URL.Path {
		case "/apis/machineconfiguration.openshift.io/v1/machineconfigs":
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"machineconfiguration.openshift.io/v1","kind":"MachineConfigList","items":[]}`))
		case "/apis/machineconfiguration.openshift.io/v1/machineconfigpools":
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"machineconfiguration.openshift.io/v1","kind":"MachineConfigPoolList","items":[{"metadata":{"name":"worker"}}]}`))
		case "/apis/performance.openshift.io/v2/performanceprofiles":
			_, _ = responseWriter.Write([]byte(`{"items":[{"metadata":{"name":"profile"}}]}`))
		default:
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"v1","kind":"NodeList","items":[{"metadata":{"name":"worker-0","labels":{"node-role.kubernetes.io/worker":""}}}]}`))
		}
	})
	healthy, err := IsMCOHealthy()
	require.NoError(t, err)
	assert.True(t, healthy)
	assert.True(t, HasWorkerNodes())
	hasProfiles, err := HasPerformanceProfiles()
	require.NoError(t, err)
	assert.True(t, hasProfiles)
	static, err := HasStaticCPUManagerPolicy()
	require.NoError(t, err)
	assert.True(t, static)
	configured, reason, err := IsClusterConfiguredForExclusiveCPUs()
	require.NoError(t, err)
	assert.True(t, configured)
	assert.Empty(t, reason)

	setOCPRESTClients(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		if request.URL.Path == "/apis/machineconfiguration.openshift.io/v1/machineconfigpools" {
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"machineconfiguration.openshift.io/v1","kind":"MachineConfigPoolList","items":[]}`))

			return
		}

		if request.URL.Path == "/apis/performance.openshift.io/v2/performanceprofiles" {
			_, _ = responseWriter.Write([]byte(`{"items":[]}`))

			return
		}

		if request.URL.Path == "/api/v1/nodes" {
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"v1","kind":"NodeList","items":[{"metadata":{"name":"worker-0","labels":{"node-role.kubernetes.io/worker":""}}}]}`))

			return
		}
		_, _ = responseWriter.Write([]byte(`{"apiVersion":"machineconfiguration.openshift.io/v1","kind":"MachineConfigList","items":[]}`))
	})
	healthy, err = IsMCOHealthy()
	require.NoError(t, err)
	assert.False(t, healthy)
	hasProfiles, err = HasPerformanceProfiles()
	require.NoError(t, err)
	assert.False(t, hasProfiles)
	static, err = HasStaticCPUManagerPolicy()
	require.NoError(t, err)
	assert.False(t, static)
	configured, reason, err = IsClusterConfiguredForExclusiveCPUs()
	require.NoError(t, err)
	assert.False(t, configured)
	assert.Contains(t, reason, "no PerformanceProfile")

	setOCPRESTClients(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"api down","code":500}`, http.StatusInternalServerError)
	})
	_, err = IsMCOHealthy()
	assert.Error(t, err)
	_, err = HasPerformanceProfiles()
	assert.Error(t, err)
	assert.False(t, HasWorkerNodes())
	static, err = HasStaticCPUManagerPolicy()
	assert.False(t, static)
	assert.Error(t, err)
	configured, reason, err = IsClusterConfiguredForExclusiveCPUs()
	require.NoError(t, err)
	assert.False(t, configured)
	assert.Equal(t, "cluster has no worker nodes", reason)

	setOCPRESTClients(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		_, _ = responseWriter.Write([]byte(`not-json`))
	})
	_, err = HasPerformanceProfiles()
	assert.Error(t, err)

	setOCPRESTClients(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		switch request.URL.Path {
		case "/apis/machineconfiguration.openshift.io/v1/machineconfigs":
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"machineconfiguration.openshift.io/v1","kind":"MachineConfigList","items":[]}`))
		case "/apis/machineconfiguration.openshift.io/v1/machineconfigpools":
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"machineconfiguration.openshift.io/v1","kind":"MachineConfigPoolList","items":[{"metadata":{"name":"worker"}}]}`))
		case "/api/v1/nodes":
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"v1","kind":"NodeList","items":[{"metadata":{"name":"worker-0","labels":{"node-role.kubernetes.io/worker":""}}}]}`))
		default:
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
	})
	_, err = HasStaticCPUManagerPolicy()
	assert.Error(t, err)

	setOCPRESTClients(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		_, _ = responseWriter.Write([]byte(`{"apiVersion":"v1","kind":"NodeList","items":[]}`))
	})
	static, err = HasStaticCPUManagerPolicy()
	assert.False(t, static)
	assert.Error(t, err)
	configured, reason, err = IsClusterConfiguredForExclusiveCPUs()
	require.NoError(t, err)
	assert.False(t, configured)
	assert.Equal(t, "cluster has no worker nodes", reason)

	setOCPRESTClients(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		if request.URL.Path == "/api/v1/nodes" {
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"v1","kind":"NodeList","items":[{"metadata":{"name":"worker-0","labels":{"node-role.kubernetes.io/worker":""}}}]}`))

			return
		}

		http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"profile API down","code":500}`, http.StatusInternalServerError)
	})
	static, err = HasStaticCPUManagerPolicy()
	assert.False(t, static)
	assert.Error(t, err)
	configured, reason, err = IsClusterConfiguredForExclusiveCPUs()
	require.NoError(t, err)
	assert.False(t, configured)
	assert.Contains(t, reason, "cannot query PerformanceProfiles")

	setOCPRESTClients(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		if strings.Contains(request.URL.Path, "/machineconfigs") {
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"machineconfiguration.openshift.io/v1","kind":"MachineConfigList","items":[]}`))

			return
		}

		http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"pool API down","code":500}`, http.StatusInternalServerError)
	})
	_, err = IsMCOHealthy()
	assert.Error(t, err)
}

func TestCoverageCatalogSourceHelpers(t *testing.T) {
	k8s, olm := coverageFixture(t)
	_ = k8s

	setCoverageCtrlObjects(&configv1api.ClusterVersion{ObjectMeta: metav1.ObjectMeta{Name: "version"}, Status: configv1api.ClusterVersionStatus{Desired: configv1api.Release{Version: "4.18.0"}}})
	require.NoError(t, ValidateCatalogSources())
	require.NoError(t, validateCatalogSources(olm))
	require.NoError(t, createCatalogSource("coverage-source", "example/image"))
	require.NoError(t, createCatalogSourceWithClient(olm, "coverage-source", "example/image"))
	require.NoError(t, CreateCommunityOperatorsCatalogSource())
	version, err := GetClusterVersion()
	require.NoError(t, err)
	assert.Equal(t, "4.18.0", version)
	assert.Equal(t, "4.18.0", GetClusterVersionOrDefault())
	require.NoError(t, DeployRHOperatorSource("4.18"))
	require.NoError(t, DeployRHCertifiedOperatorSource("4.18"))
	require.NoError(t, DeployCustomOperatorSource("example/custom:latest"))
	require.NoError(t, CreateAndValidateCatalogSources(false))
	require.NoError(t, deleteCatalogSourceByName("coverage-source"))

	oldTimeout, oldInterval := catalogValidationTimeout, catalogValidationInterval
	catalogValidationTimeout, catalogValidationInterval = 4*time.Millisecond, time.Millisecond

	t.Cleanup(func() { catalogValidationTimeout, catalogValidationInterval = oldTimeout, oldInterval })
	unready := func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"CatalogSourceList","items":[]}`))
	}
	_, unreadyClient := coverageFixtureWithOLMHandler(t, unready)
	err = validateCatalogSources(unreadyClient)
	assert.Error(t, err)

	apiFailure := func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"catalog API down","code":500}`, http.StatusInternalServerError)
	}
	_, failedClient := coverageFixtureWithOLMHandler(t, apiFailure)
	err = validateCatalogSources(failedClient)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "catalog API down")
	_, err = GetClusterVersion()
	assert.Error(t, err)
	assert.Equal(t, "4.14", GetClusterVersionOrDefault())

	_, deleteFailClient := coverageFixtureWithOLMHandler(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodDelete {
			http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"delete failed","code":500}`, http.StatusInternalServerError)

			return
		}

		defaultOLMHTTPHandler(responseWriter, request)
	})
	assert.Error(t, deleteCatalogSourceByName("bad-delete"))
	assert.NoError(t, createCatalogSourceWithClient(deleteFailClient, "broken", "image"))

	_, unreadyStateClient := coverageFixtureWithOLMHandler(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		if strings.HasSuffix(request.URL.Path, "/catalogsources") {
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"CatalogSourceList","items":[{"metadata":{"name":"certified-operators"},"status":{"connectionState":{"lastObservedState":"CONNECTING"}}},{"metadata":{"name":"community-operators"},"status":{"connectionState":{"lastObservedState":"READY"}}}]}`))

			return
		}

		defaultOLMHTTPHandler(responseWriter, request)
	})
	err = validateCatalogSources(unreadyStateClient)
	assert.Error(t, err)

	coverageFixture(t)
	assert.Error(t, CreateCommunityOperatorsCatalogSource()) // cluster version missing

	_, createFailClient := coverageFixtureWithOLMHandler(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"create rejected","code":500}`, http.StatusInternalServerError)

			return
		}

		defaultOLMHTTPHandler(responseWriter, request)
	})
	assert.Error(t, createCatalogSourceWithClient(createFailClient, "broken", "image"))

	coverageFixture(t)
	setCoverageCtrlObjects(&configv1api.ClusterVersion{ObjectMeta: metav1.ObjectMeta{Name: "version"}, Status: configv1api.ClusterVersionStatus{Desired: configv1api.Release{Version: "4.18.0"}}})
	assert.NoError(t, CreateAndValidateCatalogSources(true))

	oldTimeout, oldInterval = catalogValidationTimeout, catalogValidationInterval
	catalogValidationTimeout, catalogValidationInterval = 4*time.Millisecond, time.Millisecond

	t.Cleanup(func() { catalogValidationTimeout, catalogValidationInterval = oldTimeout, oldInterval })
	_, retryFailClient := coverageFixtureWithOLMHandler(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/catalogsources") {
			responseWriter.Header().Set("Content-Type", "application/json")
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"CatalogSourceList","items":[]}`))

			return
		}

		if request.Method == http.MethodDelete {
			http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"delete rejected","code":500}`, http.StatusInternalServerError)

			return
		}

		defaultOLMHTTPHandler(responseWriter, request)
	})

	setCoverageCtrlObjects(&configv1api.ClusterVersion{ObjectMeta: metav1.ObjectMeta{Name: "version"}, Status: configv1api.ClusterVersionStatus{Desired: configv1api.Release{Version: "4.18.0"}}})
	err = CreateAndValidateCatalogSources(true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not ready after 3 attempts")
	_ = retryFailClient
}

func TestCoveragePackageManifestHelpers(t *testing.T) {
	coverageFixture(t)
	manifest := &packagemanifestv1.PackageManifest{
		ObjectMeta: metav1.ObjectMeta{Name: "cloudbees-ci-rhmp", Namespace: "marketplace"},
		Status: packagemanifestv1.PackageManifestStatus{
			CatalogSource:  "certified-operators",
			DefaultChannel: "stable",
			Channels: []packagemanifestv1.PackageChannel{
				{Name: "alpha", CurrentCSV: "cloudbees.v0.9.0", CurrentCSVDesc: packagemanifestv1.CSVDescription{Version: packageversion.OperatorVersion{Version: semver.Version{Major: 0, Minor: 9}}}},
				{Name: "stable", CurrentCSV: "cloudbees.v1.2.3", CurrentCSVDesc: packagemanifestv1.CSVDescription{Version: packageversion.OperatorVersion{Version: semver.Version{Major: 1, Minor: 2, Patch: 3}}}},
			},
		},
	}
	setCoverageCtrlObjects(manifest)

	version, err := QueryPackageManifestForVersion("cloudbees-ci-rhmp", "marketplace", "stable")
	require.NoError(t, err)
	assert.Equal(t, "1.2.3", version)
	version, err = QueryPackageManifestForVersion("cloudbees-ci-rhmp", "marketplace", "alpha")
	require.NoError(t, err)
	assert.Equal(t, "not found", version)
	channel, err := QueryPackageManifestForDefaultChannel("cloudbees-ci-rhmp", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "stable", channel)
	name, err := QueryPackageManifestForOperatorName("cloudbees-ci", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "cloudbees-ci-rhmp", name)
	name, catalog, err := QueryPackageManifestForOperatorNameAndCatalogSource("cloudbees-ci", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "cloudbees-ci-rhmp", name)
	assert.Equal(t, "certified-operators", catalog)
	channel, version, err = QueryPackageManifestForAvailableChannelAndVersion("cloudbees-ci-rhmp", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "alpha", channel)
	assert.Equal(t, "0.9.0", version)
	channel, version, csv, err := QueryPackageManifestForAvailableChannelVersionAndCSV("cloudbees-ci-rhmp", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "alpha", channel)
	assert.Equal(t, "0.9.0", version)
	assert.Equal(t, "cloudbees.v0.9.0", csv)
	name, catalog = CheckOperatorExistsOrFail("cloudbees-ci", "marketplace")
	assert.Equal(t, "cloudbees-ci-rhmp", name)
	assert.Equal(t, "certified-operators", catalog)
	channel, version, csv = CheckOperatorChannelAndVersionOrFail("cloudbees-ci-rhmp", "marketplace")
	assert.Equal(t, "alpha", channel)
	assert.Equal(t, "0.9.0", version)
	assert.Equal(t, "cloudbees.v0.9.0", csv)

	name, err = QueryPackageManifestForOperatorName("missing", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "not found", name)

	var failCalls int
	name, catalog = checkOperatorExists("missing", "marketplace", func(string, ...int) { failCalls++ })

	assert.Equal(t, 1, failCalls)
	assert.Equal(t, "not found", name)
	assert.Equal(t, "not found", catalog)
	channel, version, csv = checkOperatorChannelAndVersion("missing", "marketplace", func(string, ...int) { failCalls++ })

	assert.Equal(t, 2, failCalls)
	assert.Equal(t, "not found", channel)
	assert.Equal(t, "not found", version)
	assert.Equal(t, "not found", csv)
	assert.Error(t, ValidateIfReportsAreValid("present", globalparameters.TestCasePassed, func() string {
		dir := t.TempDir()
		writeClaimReportForTest(t, dir, "present", globalparameters.TestCaseFailed)

		return dir
	}()))
	reportDir := t.TempDir()
	writeClaimReportForTest(t, reportDir, "present", globalparameters.TestCasePassed)
	assert.NoError(t, ValidateIfReportsAreValidWithAcceptedStatuses("present", []string{"not-a-status", globalparameters.TestCasePassed}, reportDir))
}

func TestCoveragePackageManifestEmptyAndAPIError(t *testing.T) {
	coverageFixture(t)
	setCoverageCtrlObjects()
	version, err := QueryPackageManifestForVersion("missing", "marketplace", "stable")
	require.NoError(t, err)
	assert.Equal(t, "not found", version)
	channel, err := QueryPackageManifestForDefaultChannel("missing", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "not found", channel)
	name, err := QueryPackageManifestForOperatorName("missing", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "not found", name)
	name, catalog, err := QueryPackageManifestForOperatorNameAndCatalogSource("missing", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "not found", name)
	assert.Equal(t, "not found", catalog)
	channel, version, err = QueryPackageManifestForAvailableChannelAndVersion("missing", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "not found", channel)
	assert.Equal(t, "not found", version)
	channel, version, csv, err := QueryPackageManifestForAvailableChannelVersionAndCSV("missing", "marketplace")
	require.NoError(t, err)
	assert.Equal(t, "not found", channel)
	assert.Equal(t, "not found", version)
	assert.Equal(t, "not found", csv)

	broken := failingControllerClient{Client: egiClient.Client, listErr: errors.New("package manifest list failed")}
	egiClient.Client, apiclient.Client = broken, broken
	_, err = QueryPackageManifestForVersion("missing", "marketplace", "stable")
	require.Error(t, err)
	_, err = QueryPackageManifestForDefaultChannel("missing", "marketplace")
	require.Error(t, err)
	_, err = QueryPackageManifestForOperatorName("missing", "marketplace")
	require.Error(t, err)
	_, _, err = QueryPackageManifestForOperatorNameAndCatalogSource("missing", "marketplace")
	require.Error(t, err)
	_, _, err = QueryPackageManifestForAvailableChannelAndVersion("missing", "marketplace")
	require.Error(t, err)
	_, _, _, err = QueryPackageManifestForAvailableChannelVersionAndCSV("missing", "marketplace")
	require.Error(t, err)

	var failures int
	_, _ = checkOperatorExists("missing", "marketplace", func(string, ...int) { failures++ })
	_, _, _ = checkOperatorChannelAndVersion("missing", "marketplace", func(string, ...int) { failures++ })

	assert.Equal(t, 2, failures)
}

func TestCoverageOperatorAndSubscriptionHelpers(t *testing.T) {
	k8s, olm := coverageFixture(t)
	operatorGroup := &olmv1.OperatorGroup{ObjectMeta: metav1.ObjectMeta{Name: "og", Namespace: "operators"}, Spec: olmv1.OperatorGroupSpec{TargetNamespaces: []string{"operators"}}}
	require.NoError(t, DeployOperatorGroup("operators", operatorGroup))
	require.NoError(t, IsOperatorGroupInstalled("og", "operators"))
	assert.Error(t, IsOperatorGroupInstalled("missing", "operators"))

	sub := &operatorv1alpha1.Subscription{ObjectMeta: metav1.ObjectMeta{Name: "sub", Namespace: "operators"}}
	require.NoError(t, DeployOperator(sub))
	require.NoError(t, CreateSubscription("operators", sub))
	require.NoError(t, DeleteSubscription("operators", "sub"))
	require.NoError(t, createSubscription("operators", sub, olm))
	require.NoError(t, deleteSubscription("operators", "absent", olm))

	plan, err := GetInstallPlanByCSV("operators", "missing-csv")
	assert.Nil(t, plan)
	assert.Error(t, err)
	planObj := &operatorv1alpha1.InstallPlan{ObjectMeta: metav1.ObjectMeta{Name: "plan", Namespace: "operators"}}
	require.NoError(t, ApproveInstallPlan("operators", planObj))
	assert.True(t, planObj.Spec.Approved)
	require.NoError(t, updateInstallPlan("operators", planObj))
	require.NoError(t, CreateInstallPlan(planObj))
	require.NoError(t, CreateInstallPlan(planObj))

	setCoverageCtrlObjects(&operatorv1alpha1.CatalogSource{ObjectMeta: metav1.ObjectMeta{Name: "custom-catalog", Namespace: CatalogSourceNamespace}, Spec: operatorv1alpha1.CatalogSourceSpec{DisplayName: "Custom Index"}})
	require.NoError(t, DeleteCustomOperatorSource())
	require.NoError(t, DeployCustomOperatorSource("example/custom:latest"))
	require.NoError(t, DisableCatalogSource("community-operators"))
	require.NoError(t, EnableCatalogSource("community-operators"))
	enabled, err := IsCatalogSourceEnabled("community-operators", CatalogSourceNamespace, "")
	require.NoError(t, err)
	assert.True(t, enabled)
	require.NoError(t, DeleteCatalogSource("custom-catalog", CatalogSourceNamespace, "Custom Index"))
	require.NoError(t, DeleteCatalogSource("missing", CatalogSourceNamespace, "missing"))

	_, err = k8s.CoreV1().Namespaces().Get(context.Background(), "operators", metav1.GetOptions{})
	assert.NoError(t, err)
}

func TestCoverageNamespaceCleanupHelpers(t *testing.T) {
	_, olm := coverageFixture(t,
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "clean"}},
		&corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "claim", Namespace: "clean"}},
		&policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: "pdb", Namespace: "clean"}},
	)
	assert.NoError(t, cleanSubscriptions("clean", olm))
	assert.NoError(t, cleanCSVs("clean", olm))
	assert.NoError(t, cleanInstallPlans("clean", olm))
	assert.NoError(t, cleanPVCs("clean", egiClient))
	assert.NoError(t, cleanPodDisruptionBudget("clean", egiClient))

	assert.NoError(t, cleanNetworkAttachmentDefinition(egiClient.Client))
	setCoverageCtrlObjects(
		&nadv1.NetworkAttachmentDefinition{ObjectMeta: metav1.ObjectMeta{Name: "dummy-dhcp-network", Namespace: "clean"}},
		&nadv1.NetworkAttachmentDefinition{ObjectMeta: metav1.ObjectMeta{Name: "delete-me", Namespace: "clean"}},
	)
	assert.NoError(t, cleanNetworkAttachmentDefinition(egiClient.Client))

	var remaining nadv1.NetworkAttachmentDefinitionList
	assert.NoError(t, egiClient.List(context.Background(), &remaining))
	assert.Len(t, remaining.Items, 1)
	assert.Equal(t, "dummy-dhcp-network", remaining.Items[0].Name)
	assert.NoError(t, cleanNamespace("clean", GetAPIClient(), egiClient))

	k8s, _ := coverageFixture(t, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "delete-ns"}}, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "still-here"}})
	assert.NoError(t, DeleteNamespaceAndWait("delete-ns", time.Second))
	assert.NoError(t, DeleteNamespaceAndWait("", time.Second))
	assert.Error(t, WaitForNamespaceDeletion(k8s.CoreV1(), "still-here", time.Millisecond))
	assert.NoError(t, CleanNamespace("not-found"))
	assert.NoError(t, cleanNamespace("not-found", GetAPIClient(), egiClient))
	assert.NoError(t, cleanNamespace("delete-ns", GetAPIClient(), egiClient))

	badList := failingControllerClient{Client: egiClient.Client, listErr: errors.New("list failed")}
	assert.Error(t, cleanNetworkAttachmentDefinition(badList))
	setCoverageCtrlObjects(
		&nadv1.NetworkAttachmentDefinition{ObjectMeta: metav1.ObjectMeta{Name: "one", Namespace: "clean"}},
		&nadv1.NetworkAttachmentDefinition{ObjectMeta: metav1.ObjectMeta{Name: "two", Namespace: "clean"}},
	)
	badDelete := failingControllerClient{Client: egiClient.Client, deleteErr: errors.New("delete failed")}
	assert.Error(t, cleanNetworkAttachmentDefinition(badDelete))

	failedHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"delete rejected","code":500}`, http.StatusInternalServerError)
	}
	_, failingOLM := coverageFixtureWithOLMHandler(t, failedHandler)
	assert.Error(t, cleanSubscriptions("clean", failingOLM))
	assert.Error(t, cleanCSVs("clean", failingOLM))
	assert.Error(t, cleanInstallPlans("clean", failingOLM))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("list", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("pvc list failed") })
	assert.Error(t, cleanPVCs("clean", egiClient))
	k8s.PrependReactor("delete-collection", "poddisruptionbudgets", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("pdb cleanup failed")
	})
	assert.Error(t, cleanPodDisruptionBudget("clean", egiClient))
}

type failingControllerClient struct {
	runtimeclient.Client
	listErr   error
	deleteErr error
	getErr    error
	createErr error
	updateErr error
	getOnce   *bool
}

func (c failingControllerClient) List(ctx context.Context, list runtimeclient.ObjectList, options ...runtimeclient.ListOption) error {
	if c.listErr != nil {
		return c.listErr
	}

	return c.Client.List(ctx, list, options...)
}

func (c failingControllerClient) Delete(ctx context.Context, object runtimeclient.Object, options ...runtimeclient.DeleteOption) error {
	if c.deleteErr != nil {
		return c.deleteErr
	}

	return c.Client.Delete(ctx, object, options...)
}

func (c failingControllerClient) Get(ctx context.Context, key runtimeclient.ObjectKey, object runtimeclient.Object, options ...runtimeclient.GetOption) error {
	if c.getErr != nil && (c.getOnce == nil || *c.getOnce) {
		if c.getOnce != nil {
			*c.getOnce = false
		}

		return c.getErr
	}

	return c.Client.Get(ctx, key, object, options...)
}

func (c failingControllerClient) Create(ctx context.Context, object runtimeclient.Object, options ...runtimeclient.CreateOption) error {
	if c.createErr != nil {
		return c.createErr
	}

	return c.Client.Create(ctx, object, options...)
}

func (c failingControllerClient) Update(ctx context.Context, object runtimeclient.Object, options ...runtimeclient.UpdateOption) error {
	if c.updateErr != nil {
		return c.updateErr
	}

	return c.Client.Update(ctx, object, options...)
}

func setOCPRESTClients(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	apiConfig := &rest.Config{Host: server.URL}
	core, err := typedcorev1.NewForConfig(apiConfig)
	require.NoError(t, err)
	configClient, err := configv1.NewForConfig(apiConfig)
	require.NoError(t, err)
	machineClient, err := machinev1.NewForConfig(apiConfig)
	require.NoError(t, err)
	apiclient.CoreV1Interface = core
	apiclient.OcpClientInterface = configClient
	apiclient.MachineconfigurationV1Interface = machineClient
}

func TestCoverageNodeHelpers(t *testing.T) {
	k8s, _ := coverageFixture(t, &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node"}})
	node, err := k8s.CoreV1().Nodes().Get(context.Background(), "node", metav1.GetOptions{})
	require.NoError(t, err)
	require.NoError(t, addControlPlaneTaint(k8s.CoreV1(), node))
	assert.Equal(t, controlPlaneTaintKey, node.Spec.Taints[0].Key)
	require.NoError(t, addControlPlaneTaint(k8s.CoreV1(), node))
	node.Spec.Taints = append(node.Spec.Taints, corev1.Taint{Key: "keep"}, corev1.Taint{Key: masterTaintKey})
	require.NoError(t, removeControlPlaneTaint(k8s.CoreV1(), node))
	assert.Equal(t, []corev1.Taint{{Key: "keep"}}, node.Spec.Taints)

	huge := &corev1.Node{Status: corev1.NodeStatus{Capacity: corev1.ResourceList{
		corev1.ResourceName("hugepages-2Mi"): resource.MustParse("2Mi"),
	}}}
	assert.True(t, NodeHasHugePagesEnabled(huge, "2Mi"))
	assert.False(t, NodeHasHugePagesEnabled(huge, "1Gi"))
	assert.False(t, NodeHasHugePagesEnabled(&corev1.Node{}, "2Mi"))

	huge.Labels = map[string]string{"node-role.kubernetes.io/worker-cnf": ""}
	k8s, _ = coverageFixture(t, huge, &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "worker", Labels: map[string]string{"node-role.kubernetes.io/worker-cnf": ""}}})
	assert.True(t, NodesHaveHugePagesEnabled("2Mi"))
	assert.False(t, NodesHaveHugePagesEnabled("1Gi"))
	assert.Equal(t, 2, GetNumberOfNodes(k8s.CoreV1()))
	assert.Equal(t, 2, getNumberOfNodes(k8s.CoreV1()))

	assert.False(t, isNodeOvercommitted(&corev1.Node{}))
	assert.False(t, isNodeOvercommitted(&corev1.Node{Status: corev1.NodeStatus{Allocatable: corev1.ResourceList{corev1.ResourceMemory: resource.MustParse("0"), corev1.ResourceCPU: resource.MustParse("0")}}}))
	assert.True(t, isNodeOvercommitted(&corev1.Node{Status: corev1.NodeStatus{Allocatable: corev1.ResourceList{corev1.ResourceMemory: resource.MustParse("1Mi")}}}))
	assert.True(t, isNodeOvercommitted(&corev1.Node{Status: corev1.NodeStatus{Allocatable: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("1")}}}))

	k8s, _ = coverageFixture(t, &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "worker"}, Status: corev1.NodeStatus{Allocatable: corev1.ResourceList{corev1.ResourceMemory: resource.MustParse("1Mi")}}})
	over, err := IsClusterOvercommitted()
	require.NoError(t, err)
	assert.True(t, over)
	_ = k8s
}

func TestCoverageNamespaceResourceHelpers(t *testing.T) {
	gomega.RegisterTestingT(t)
	k8s, olm := coverageFixture(t, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns"}})
	require.NoError(t, CreateNamespace("new-ns"))
	require.NoError(t, CreatePrivilegedNamespace("priv-ns"))
	priv, err := k8s.CoreV1().Namespaces().Get(context.Background(), "priv-ns", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "privileged", priv.Labels["pod-security.kubernetes.io/enforce"])
	exists, err := NamespaceExists("ns")
	require.NoError(t, err)
	assert.True(t, exists)
	exists, err = NamespaceExists("not-here")
	require.NoError(t, err)
	assert.False(t, exists)

	require.NoError(t, createResourceQuota(egiClient, &corev1.ResourceQuota{ObjectMeta: metav1.ObjectMeta{Name: "q", Namespace: "ns"}}))
	require.NoError(t, createResourceQuota(egiClient, &corev1.ResourceQuota{ObjectMeta: metav1.ObjectMeta{Name: "q", Namespace: "ns"}}))
	require.NoError(t, createResourceQuota(egiClient, &corev1.ResourceQuota{ObjectMeta: metav1.ObjectMeta{Name: "q", Namespace: "missing"}}))
	require.NoError(t, CreateResourceQuota(&corev1.ResourceQuota{ObjectMeta: metav1.ObjectMeta{Name: "public-q", Namespace: "ns"}}))

	require.NoError(t, createStorageClass(k8s.StorageV1(), "coverage-sc", true))
	require.NoError(t, createStorageClass(k8s.StorageV1(), "coverage-sc", false))
	require.NoError(t, CreateStorageClass("public-sc", true))
	require.NoError(t, deleteStorageClass(k8s.StorageV1(), "missing-sc"))
	require.NoError(t, deleteStorageClass(k8s.StorageV1(), "coverage-sc"))
	require.NoError(t, DeleteStorageClass("public-sc"))

	require.NoError(t, createInstallPlan(installPlanForCoverage(), olm))
	require.NoError(t, createInstallPlan(installPlanForCoverage(), olm))
	require.NoError(t, deleteMutatingWebhookConfiguration(k8s.AdmissionregistrationV1(), "missing"))
	require.NoError(t, deleteValidatingWebhookConfiguration(k8s.AdmissionregistrationV1(), "missing"))
	require.NoError(t, DeleteMutatingWebhookConfiguration("missing"))
	require.NoError(t, DeleteValidatingWebhookConfiguration("missing"))
	_, err = getServicesFromNamespace("ns", k8s)
	require.NoError(t, err)
	_, err = GetServicesFromNamespace("ns")
	require.NoError(t, err)

	runtimeClass := &nodev1.RuntimeClass{ObjectMeta: metav1.ObjectMeta{Name: "coverage-rtc", UID: "uid"}}
	require.NoError(t, createRunTimeClass(k8s, runtimeClass))
	require.NoError(t, deleteRunTimeClass(k8s, runtimeClass))
	runtimeClass.Name = "public-rtc"
	require.NoError(t, CreateRunTimeClass(runtimeClass))
	require.NoError(t, DeleteRunTimeClass(runtimeClass))

	require.NoError(t, createServiceAccount(k8s.CoreV1(), "sa", "ns"))
	require.NoError(t, deleteServiceAccount(k8s.CoreV1(), "sa", "ns"))
	require.NoError(t, CreateServiceAccount("public-sa", "ns"))
	require.NoError(t, DeleteServiceAccount("public-sa", "ns"))
	require.NoError(t, createRole(k8s.RbacV1(), &rbacv1.Role{ObjectMeta: metav1.ObjectMeta{Name: "role", Namespace: "ns"}, Rules: []rbacv1.PolicyRule{{Verbs: []string{"get"}}}}))
	publicRole := DefineRole("public-role", "ns")
	RedefineRoleWithAPIGroups(publicRole, []string{"apps"})
	RedefineRoleWithResources(publicRole, []string{"deployments"})
	require.NoError(t, CreateRole(publicRole))
	require.NoError(t, DeleteRole("public-role", "ns"))
	require.NoError(t, deleteRole(k8s.RbacV1(), "role", "ns"))
	require.NoError(t, deleteRoleBinding(k8s.RbacV1(), "binding", "ns"))
	require.NoError(t, createRoleBindingWithServiceAccountSubject(k8s.RbacV1(), "binding", "role", "sa", "ns", "ns"))
	require.NoError(t, CreateRoleBindingWithServiceAccountSubject("public-binding", "public-role", "public-sa", "ns", "ns"))
	require.NoError(t, DeleteRoleBinding("public-binding", "ns"))
	require.NoError(t, createPersistentVolume(k8s.CoreV1(), &corev1.PersistentVolume{ObjectMeta: metav1.ObjectMeta{Name: "pv"}}))
	require.NoError(t, CreatePersistentVolume(&corev1.PersistentVolume{ObjectMeta: metav1.ObjectMeta{Name: "public-pv"}}))
	require.NoError(t, deletePersistentVolume(k8s.CoreV1(), "pv", time.Second))
	require.NoError(t, DeletePersistentVolume("public-pv", time.Second))
	require.NoError(t, deletePersistentVolumeClaim(k8s.CoreV1(), &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "pvc", Namespace: "ns"}}))
	require.NoError(t, DeletePersistentVolumeClaim(&corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "public-pvc", Namespace: "ns"}}))

	require.NoError(t, createNetworkPolicy(egiClient, &networkingv1.NetworkPolicy{ObjectMeta: metav1.ObjectMeta{Name: "np", Namespace: "ns"}}))
	require.NoError(t, createClusterRoleBinding(egiClient, &rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{Name: "crb"}, RoleRef: rbacv1.RoleRef{Name: "role"}, Subjects: []rbacv1.Subject{{Kind: "User", Name: "u"}}}))
	require.NoError(t, deleteClusterRoleBindingByName(egiClient, "crb"))
	require.NoError(t, deleteClusterRole(egiClient, "coverage-role"))
	publicCRB := &rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{Name: "public-crb"}, RoleRef: rbacv1.RoleRef{Name: "role"}, Subjects: []rbacv1.Subject{{Kind: "User", Name: "u"}}}
	require.NoError(t, CreateClusterRoleBinding(publicCRB))
	require.NoError(t, DeleteClusterRoleBinding(publicCRB))
	require.NoError(t, DeleteClusterRoleBindingByName("public-crb"))
	require.NoError(t, DeleteClusterRole("public-role"))
	require.NoError(t, CreateNetworkPolicy(&networkingv1.NetworkPolicy{ObjectMeta: metav1.ObjectMeta{Name: "public-np", Namespace: "ns"}}, time.Second))
}

//nolint:funlen // The cases share serial fake-client state and cover related workload helpers.
func TestCoveragePodAndWorkloadHelpers(t *testing.T) {
	gomega.RegisterTestingT(t)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "ready-pod", Namespace: "ns"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "c"}}},
		Status: corev1.PodStatus{
			Phase:             corev1.PodRunning,
			ContainerStatuses: []corev1.ContainerStatus{{Name: "c", Ready: true}},
		},
	}
	k8s, _ := coverageFixture(t, pod)
	k8s.PrependReactor("get", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, pod.DeepCopy(), nil
	})
	assert.NoError(t, CreateAndWaitUntilPodIsReady(pod, time.Second))
	assert.True(t, func() bool {
		ready, err := isPodReady("ns", "ready-pod")

		return err == nil && ready
	}())
	assert.NoError(t, CreateAndWaitUntilPodIsReady(pod, time.Second)) // already exists
	pods, err := GetListOfPodsInNamespace("ns")
	assert.NoError(t, err)
	assert.Len(t, pods.Items, 1)
	gotPod, err := GetRunningPod("ns", "ready-pod")
	assert.NoError(t, err)
	assert.Equal(t, "ready-pod", gotPod.Name)
	controller, err := GetControllerPodFromOperator("ns", "ready")
	assert.NoError(t, err)
	assert.Equal(t, pod.Name, controller.Name)
	oldExec := execPodCommand
	execPodCommand = func(_ *egipod.Builder, command []string, _ ...string) (bytes.Buffer, error) {
		var output bytes.Buffer
		output.WriteString(strings.Join(command, " "))

		return output, nil
	}
	output, err := ExecCommand(*pod, []string{"echo", "ok"})
	assert.NoError(t, err)
	assert.Equal(t, "echo ok", output.String())
	execPodCommand = func(*egipod.Builder, []string, ...string) (bytes.Buffer, error) {
		return bytes.Buffer{}, errors.New("exec failed")
	}
	_, err = ExecCommand(*pod, []string{"false"})
	assert.Error(t, err)
	execPodCommand = oldExec
	k8s, _ = coverageFixture(t)
	_, err = ExecCommand(*pod, []string{"true"})
	assert.Error(t, err)
	_, err = getRunningPod(k8s.CoreV1(), "ns", "missing")
	assert.Error(t, err)

	appended := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"old": "v"}}}
	AppendContainersToPod(appended, 2, "image")
	AppendLabelsToPod(appended, map[string]string{"old": "new", "added": "v"})
	assert.Len(t, appended.Spec.Containers, 2)
	assert.Equal(t, "new", appended.Labels["old"])
	assert.Equal(t, "v", appended.Labels["added"])

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "dep", Namespace: "ns"},
		Spec:       appsv1.DeploymentSpec{Replicas: ptr.To[int32](1)},
		Status:     appsv1.DeploymentStatus{Replicas: 1, ReadyReplicas: 1},
	}
	k8s, _ = coverageFixture(t)
	deploymentPollInterval = time.Millisecond

	t.Cleanup(func() { deploymentPollInterval = 10 * time.Second })
	assert.NoError(t, CreateAndWaitUntilDeploymentIsReady(dep, time.Second))
	assert.NoError(t, CreateAndWaitUntilDeploymentIsReady(dep, time.Second)) // already exists
	assert.True(t, func() bool {
		ready, err := IsDeploymentReady(egiClient, "ns", "dep")

		return err == nil && ready
	}())
	gotDep, err := GetRunningDeployment("ns", "dep")
	assert.NoError(t, err)
	assert.Equal(t, "dep", gotDep.Name)
	assert.NoError(t, DeleteDeployment("dep", "ns"))
	_, err = GetRunningDeployment("ns", "missing")
	assert.Error(t, err)
	_, err = k8s.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns"}}, metav1.CreateOptions{})
	assert.NoError(t, err)

	coverageFixture(t, &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node"}})
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{Name: "ds", Namespace: "ns"},
		Spec:       appsv1.DaemonSetSpec{Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "c", ImagePullPolicy: corev1.PullNever}}}}},
		Status:     appsv1.DaemonSetStatus{DesiredNumberScheduled: 1, NumberReady: 1, NumberAvailable: 1, CurrentNumberScheduled: 1},
	}
	assert.NoError(t, CreateAndWaitUntilDaemonSetIsReady(daemonSet, time.Second))
	assert.NoError(t, CreateAndWaitUntilDaemonSetIsReady(daemonSet, time.Second))
	assert.Equal(t, corev1.PullNever, func() corev1.PullPolicy {
		p, _ := GetDaemonSetPullPolicy(daemonSet)

		return p
	}())
	assert.Equal(t, "ds", func() string {
		got, _ := GetRunningDaemonset(daemonSet)

		return got.Name
	}())
	assert.False(t, IsTransientDaemonSetError(nil))
	assert.True(t, IsTransientDaemonSetError(errors.New("pod is not schedulable")))

	replicas := ptr.To[int32](1)
	replicaSet := &appsv1.ReplicaSet{ObjectMeta: metav1.ObjectMeta{Name: "rs", Namespace: "ns"}, Spec: appsv1.ReplicaSetSpec{Replicas: replicas}, Status: appsv1.ReplicaSetStatus{AvailableReplicas: 1}}
	k8s, _ = coverageFixture(t, replicaSet)
	k8s.PrependReactor("get", "replicasets", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, replicaSet.DeepCopy(), nil
	})
	assert.True(t, func() bool {
		ready, err := isReplicaSetReady("ns", "replicaSet")

		return err == nil && ready
	}())
	assert.NoError(t, CreateAndWaitUntilReplicaSetIsReady(replicaSet, time.Second))
	k8s.PrependReactor("get", "replicasets", func(k8stesting.Action) (bool, runtime.Object, error) {
		mismatch := replicaSet.DeepCopy()
		mismatch.Spec.Replicas = ptr.To[int32](2)

		return true, mismatch, nil
	})
	ready, err := isReplicaSetReady("ns", "replicaSet")
	assert.NoError(t, err)
	assert.False(t, ready)

	sts := &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "sts", Namespace: "ns"}, Spec: appsv1.StatefulSetSpec{Replicas: replicas}, Status: appsv1.StatefulSetStatus{Replicas: 1, ReadyReplicas: 1}}

	coverageFixture(t)
	assert.NoError(t, createAndWaitUntilStatefulSetIsReady(egiClient, sts, time.Second))
	assert.NoError(t, createAndWaitUntilStatefulSetIsReady(egiClient, sts, time.Second))
	assert.NoError(t, CreateAndWaitUntilStatefulSetIsReady(sts, time.Second))
	assert.Equal(t, "sts", func() string {
		got, _ := GetRunningStatefulSet("ns", "sts")

		return got.Name
	}())

	pdb := &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: "pdb", Namespace: "ns", UID: "uid"}}
	k8s, _ = coverageFixture(t)
	assert.NoError(t, CreatePodDisruptionBudget(pdb, time.Second))
	created, err := IsPodDisruptionBudgetCreated("pdb", "ns")
	assert.NoError(t, err)
	assert.True(t, created)
	_ = k8s

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: "bound", Namespace: "ns"},
		Spec:       corev1.PersistentVolumeClaimSpec{VolumeName: "pv"},
		Status:     corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimBound},
	}
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) { return true, pvc.DeepCopy(), nil })
	assert.True(t, func() bool {
		bound, err := isPvcBound(k8s.CoreV1(), "bound", "ns", "pv")

		return err == nil && bound
	}())
	assert.False(t, func() bool {
		bound, err := isPvcBound(k8s.CoreV1(), "bound", "ns", "other")

		return err == nil && bound
	}())
	assert.NoError(t, CreateAndWaitUntilPVCIsBound(pvc, time.Second, "pv"))
	assert.NoError(t, DeletePersistentVolumeClaim(pvc))
	k8s, _ = coverageFixture(t)
	assert.Error(t, func() error {
		_, err := isPvcBound(k8s.CoreV1(), "missing", "ns", "pv")

		return err
	}())
}

func TestCoveragePodFailureBranches(t *testing.T) {
	gomega.RegisterTestingT(t)
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pending", Namespace: "ns"}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "c"}}}}
	k8s, _ := coverageFixture(t)
	k8s.PrependReactor("create", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("pod create failed")
	})
	assert.Error(t, CreateAndWaitUntilPodIsReady(pod, time.Millisecond))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		unschedulable := pod.DeepCopy()
		unschedulable.Status.Conditions = []corev1.PodCondition{{Type: corev1.PodScheduled, Status: corev1.ConditionFalse, Reason: "Unschedulable"}}

		return true, unschedulable, nil
	})
	assert.Error(t, CreateAndWaitUntilPodIsReady(pod, time.Second))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "pods", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("pod get failed") })
	_, err := GetListOfPodsInNamespace("ns")
	assert.NoError(t, err)
	_, err = isPodReady("ns", "missing")
	assert.Error(t, err)
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("list", "pods", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("pod list failed") })
	_, err = GetListOfPodsInNamespace("ns")
	assert.Error(t, err)

	pending := pod.DeepCopy()
	k8s, _ = coverageFixture(t, pending)
	assert.False(t, func() bool {
		ready, err := isPodReady("ns", "pending")

		return err == nil && ready
	}())
	pending.Status.ContainerStatuses = []corev1.ContainerStatus{{Name: "c", Ready: false}}

	k8s.PrependReactor("get", "pods", func(k8stesting.Action) (bool, runtime.Object, error) { return true, pending.DeepCopy(), nil })
	assert.False(t, func() bool {
		ready, err := isPodReady("ns", "pending")

		return err == nil && ready
	}())

	oldTimeout, oldInterval := podDiscoveryTimeout, podDiscoveryInterval
	podDiscoveryTimeout, podDiscoveryInterval = 3*time.Millisecond, time.Millisecond

	t.Cleanup(func() { podDiscoveryTimeout, podDiscoveryInterval = oldTimeout, oldInterval })
	coverageFixture(t)
	_, err = GetControllerPodFromOperator("ns", "operator") // empty list times out
	assert.Error(t, err)
	coverageFixture(t, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: "ns"}})
	_, err = GetControllerPodFromOperator("ns", "operator") // no matching pod
	assert.Error(t, err)
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("list", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("list pods failed")
	})
	_, err = GetControllerPodFromOperator("ns", "operator")
	assert.Error(t, err)
}

func TestCoverageDirectoriesAndFileHelpers(t *testing.T) {
	root := t.TempDir()
	conf = &config.Config{}
	conf.General.CertsuiteReportDir = filepath.Join(root, "reports")
	conf.General.CertsuiteConfigDir = filepath.Join(root, "configs")
	report, configDir, err := GenerateDirectories("case")
	require.NoError(t, err)
	assert.DirExists(t, report)
	assert.DirExists(t, configDir)

	source := filepath.Join(root, "src")
	require.NoError(t, os.WriteFile(source, []byte("data"), 0o600))
	destination := filepath.Join(root, "dst")
	require.NoError(t, CopyFiles(source, destination))
	data, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, "data", string(data))
	assert.Error(t, CopyFiles(filepath.Join(root, "missing"), destination))
	assert.Error(t, CopyFiles(source, filepath.Join(root, "absent", "file")))
	assert.Error(t, DefineCertsuiteConfig(nil, nil, nil, nil, nil, filepath.Join(root, "no-dir")))

	reportDir := t.TempDir()
	writeClaimReportForTest(t, reportDir, "networking-test", globalparameters.TestCasePassed)
	conf.General.ReportDirAbsPath = root

	CopyClaimFileToTcFolder("networking-test", "networking_test", reportDir)
	assert.FileExists(t, filepath.Join(root, "Debug", globalparameters.NetworkSuiteName, "networking_test", globalparameters.DefaultClaimFileName))
	CopyClaimFileToTcFolder("unknown-test", "unknown", reportDir)
	CopyClaimFileToTcFolder("networking-test", "missing", filepath.Join(root, "missing-report"))
}

func TestIsCRCSingleNode(t *testing.T) {
	crcLabels := map[string]string{
		"node-role.kubernetes.io/master":        "",
		"node-role.kubernetes.io/control-plane": "",
	}

	testCases := []struct {
		name  string
		nodes []corev1.Node
		want  bool
	}{
		{name: "no nodes", nodes: nil, want: false},
		{
			name:  "single CRC node",
			nodes: []corev1.Node{{ObjectMeta: metav1.ObjectMeta{Name: "crc-abcdef-master-0"}}},
			want:  true,
		},
		{
			name:  "single CRC node with role labels",
			nodes: []corev1.Node{{ObjectMeta: metav1.ObjectMeta{Name: "crc-abcdef-master-0", Labels: crcLabels}}},
			want:  true,
		},
		{
			name:  "single node with role labels but no CRC name",
			nodes: []corev1.Node{{ObjectMeta: metav1.ObjectMeta{Name: "sno-master-0", Labels: crcLabels}}},
			want:  false,
		},
		{
			name:  "single CRC-named worker node",
			nodes: []corev1.Node{{ObjectMeta: metav1.ObjectMeta{Name: "crc-abcdef-worker-0"}}},
			want:  false,
		},
		{
			name: "multiple nodes including a CRC-named one",
			nodes: []corev1.Node{
				{ObjectMeta: metav1.ObjectMeta{Name: "crc-abcdef-master-0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "crc-abcdef-worker-0"}},
			},
			want: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, isCRCSingleNode(testCase.nodes))
		})
	}
}

func TestCoverageClusterDetectionAndSCCHelper(t *testing.T) {
	coverageFixture(t)

	pathDir := t.TempDir()
	writeOC := func(contents string) {
		t.Helper()
		require.NoError(t, os.WriteFile(filepath.Join(pathDir, "oc"), []byte("#!/bin/sh\n"+contents), 0o700))
		t.Setenv("PATH", pathDir)
	}
	writeOC("exit 0\n")
	assert.True(t, IsKindCluster())
	assert.False(t, IsCRCCluster())
	writeOC("echo 'https://api.crc.testing:6443'\nexit 0\n")
	assert.True(t, IsCRCCluster())
	writeOC("echo 'https://apps-crc.testing'\nexit 0\n")
	assert.True(t, IsCRCCluster())
	writeOC("echo 'https://api.example.test'\nexit 0\n")
	coverageFixture(t, &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "crc-abcdef-master-0"}})
	assert.True(t, IsCRCCluster())
	// A single node carrying both control-plane role labels is a Single Node
	// OpenShift cluster, not a CRC cluster, and must not take the CRC branch.
	coverageFixture(t, &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "single", Labels: map[string]string{"node-role.kubernetes.io/master": "", "node-role.kubernetes.io/control-plane": ""}}})
	assert.False(t, IsCRCCluster())
	coverageFixture(t, &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "worker"}})
	assert.False(t, IsCRCCluster())
	k8s, _ := coverageFixture(t)
	k8s.PrependReactor("list", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("node list failed")
	})
	assert.False(t, IsCRCCluster())
	writeOC("exit 1\n")
	assert.False(t, IsKindCluster())
	assert.False(t, IsCRCCluster())
	_ = k8s

	k8s, _ = coverageFixture(t)
	assert.NoError(t, AllowAuthenticatedUsersRunPrivilegedContainers())
	assert.NoError(t, AllowAuthenticatedUsersRunPrivilegedContainers())
	_, err := k8s.RbacV1().ClusterRoleBindings().Get(context.Background(), "system:openshift:scc:privileged", metav1.GetOptions{})
	assert.NoError(t, err)

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "clusterrolebindings", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create denied")
	})
	assert.Error(t, AllowAuthenticatedUsersRunPrivilegedContainers())

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "clusterrolebindings", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("get denied")
	})
	assert.NoError(t, AllowAuthenticatedUsersRunPrivilegedContainers())

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "clusterrolebindings", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, k8serrors.NewAlreadyExists(rbacv1.Resource("clusterrolebindings"), "system:openshift:scc:privileged")
	})
	assert.NoError(t, AllowAuthenticatedUsersRunPrivilegedContainers())
}

func TestCoverageCheckDetailsFromReport(t *testing.T) {
	reportDir := t.TempDir()
	writeClaimReportForTest(t, reportDir, "test-case", globalparameters.TestCasePassed)
	claimPath := filepath.Join(reportDir, globalparameters.DefaultClaimFileName)
	content, err := os.ReadFile(claimPath)
	require.NoError(t, err)
	content = []byte(strings.Replace(string(content), `"checkDetails": ""`, `"checkDetails": "{\"CompliantObjectsOut\":[],\"NonCompliantObjectsOut\":[]}"`, 1))
	require.NoError(t, os.WriteFile(claimPath, content, 0o600))
	details, err := GetTestCaseCheckDetails("test-case", reportDir)
	require.NoError(t, err)
	assert.NotNil(t, details)
	_, err = GetTestCaseCheckDetails("missing", reportDir)
	assert.Error(t, err)
	_, err = GetTestCaseCheckDetails("test-case", t.TempDir())
	assert.Error(t, err)

	content = []byte(strings.Replace(string(content), `"checkDetails": "{\"CompliantObjectsOut\":[],\"NonCompliantObjectsOut\":[]}"`, `"checkDetails": "{"`, 1))
	require.NoError(t, os.WriteFile(claimPath, content, 0o600))
	_, err = GetTestCaseCheckDetails("test-case", reportDir)
	assert.Error(t, err)
	assert.Empty(t, GetReportObjectFieldValue(nil, "anything"))
}

func TestCoverageFailureBranches(t *testing.T) {
	k8s, _ := coverageFixture(t)
	k8s.PrependReactor("create", "namespaces", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("denied") })
	assert.Error(t, createNamespace("denied", egiClient))
	assert.Error(t, createPrivilegedNamespace("denied", egiClient))

	k8s, _ = coverageFixture(t, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns"}})
	k8s.PrependReactor("create", "resourcequotas", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("quota denied") })
	assert.Error(t, createResourceQuota(egiClient, &corev1.ResourceQuota{ObjectMeta: metav1.ObjectMeta{Name: "q", Namespace: "ns"}}))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("list", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("nodes denied") })
	assert.False(t, NodesHaveHugePagesEnabled("2Mi"))
	assert.Equal(t, 0, GetNumberOfNodes(k8s.CoreV1()))
	_, err := IsClusterOvercommitted()
	assert.Error(t, err)

	k8s, _ = coverageFixture(t, &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node"}})
	k8s.PrependReactor("update", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("update denied") })
	node, _ := k8s.CoreV1().Nodes().Get(context.Background(), "node", metav1.GetOptions{})
	assert.Error(t, addControlPlaneTaint(k8s.CoreV1(), node))
	assert.Error(t, removeControlPlaneTaint(k8s.CoreV1(), node))

	coverageFixture(t)
	over, err := IsClusterOvercommitted()
	assert.NoError(t, err)
	assert.False(t, over)
}

func installPlanForCoverage() *operatorv1alpha1.InstallPlan {
	return &operatorv1alpha1.InstallPlan{ObjectMeta: metav1.ObjectMeta{Name: "ip", Namespace: "ns"}}
}

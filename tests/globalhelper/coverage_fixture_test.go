//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	packagemanifestv1 "github.com/openshift-kni/eco-goinfra/pkg/schemes/olm/package-server/operators/v1"
	configv1api "github.com/openshift/api/config/v1"
	configv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	operatorv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	operatorclient "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/client"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// coverageFixture installs serial, fake clients for tests that exercise the
// package-level client accessors. Restore globals after each test because the
// helpers intentionally cache their clients.
//
//nolint:ireturn // Typed client interface lets fixtures switch between fake and HTTP-backed clients.
func coverageFixture(t *testing.T, objects ...runtime.Object) (*k8sfake.Clientset, operatorclient.OperatorsV1alpha1Interface) {
	t.Helper()

	return coverageFixtureWithOLMHandler(t, defaultOLMHTTPHandler, objects...)
}

//nolint:ireturn // Typed client interface lets fixtures switch between fake and HTTP-backed clients.
func coverageFixtureWithOLMHandler(t *testing.T, handler http.HandlerFunc, objects ...runtime.Object) (*k8sfake.Clientset, operatorclient.OperatorsV1alpha1Interface) {
	t.Helper()
	oldAPIClient, oldConfig, oldEGIClient := apiclient, conf, egiClient

	t.Cleanup(func() {
		apiclient, conf, egiClient = oldAPIClient, oldConfig, oldEGIClient
	})

	k8sClient := k8sfake.NewClientset(objects...)
	operatorServer := httptest.NewServer(handler)
	t.Cleanup(operatorServer.Close)

	operatorClient, err := operatorclient.NewForConfig(&rest.Config{Host: operatorServer.URL})
	if err != nil {
		t.Fatal(err)
	}

	configClient, err := configv1.NewForConfig(&rest.Config{Host: operatorServer.URL})
	if err != nil {
		t.Fatal(err)
	}

	api := &client.ClientSet{
		CoreV1Interface:            k8sClient.CoreV1(),
		AppsV1Interface:            k8sClient.AppsV1(),
		RbacV1Interface:            k8sClient.RbacV1(),
		NodeV1Interface:            k8sClient.NodeV1(),
		K8sClient:                  k8sClient,
		OperatorsV1alpha1Interface: operatorClient,
		PolicyV1Interface:          k8sClient.PolicyV1(),
		OcpClientInterface:         configClient,
	}

	egiClient = egiClients.GetTestClients(egiClients.TestClientParams{K8sMockObjects: objects})
	egiClient.K8sClient = k8sClient
	egiClient.CoreV1Interface = k8sClient.CoreV1()
	egiClient.AppsV1Interface = k8sClient.AppsV1()
	egiClient.RbacV1Interface = k8sClient.RbacV1()
	egiClient.NetworkingV1Interface = k8sClient.NetworkingV1()
	egiClient.StorageV1Interface = k8sClient.StorageV1()
	egiClient.PolicyV1Interface = k8sClient.PolicyV1()
	egiClient.ConfigV1Interface = configClient
	ctrlScheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(ctrlScheme)
	_ = operatorv1.AddToScheme(ctrlScheme)
	_ = operatorv1alpha1.AddToScheme(ctrlScheme)
	_ = nadv1.AddToScheme(ctrlScheme)
	_ = configv1api.Install(ctrlScheme)
	_ = packagemanifestv1.AddToScheme(ctrlScheme)

	ctrlObjects := make([]runtimeclient.Object, 0, len(objects))
	for _, object := range objects {
		if clientObject, ok := object.(runtimeclient.Object); ok {
			ctrlObjects = append(ctrlObjects, clientObject)
		}
	}
	ctrlClient := fake.NewClientBuilder().WithScheme(ctrlScheme).WithObjects(ctrlObjects...).Build()
	egiClient.Client = ctrlClient
	api.Client = ctrlClient

	apiclient = api
	conf = nil

	return k8sClient, operatorClient
}

func setCoverageCtrlObjects(objects ...runtimeclient.Object) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = operatorv1.AddToScheme(scheme)
	_ = operatorv1alpha1.AddToScheme(scheme)
	_ = nadv1.AddToScheme(scheme)
	_ = configv1api.Install(scheme)
	_ = packagemanifestv1.AddToScheme(scheme)
	ctrlClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()
	egiClient.Client = ctrlClient
	apiclient.Client = ctrlClient
}

func defaultOLMHTTPHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	if request.Method == http.MethodGet && request.URL.Query().Has("watch") {
		responseWriter.WriteHeader(http.StatusOK)

		return
	}

	if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodPatch {
		var body map[string]any
		_ = json.NewDecoder(request.Body).Decode(&body)
		_ = json.NewEncoder(responseWriter).Encode(body)

		return
	}

	if request.Method == http.MethodGet && len(request.URL.Path) >= len("catalogsources") && request.URL.Path[len(request.URL.Path)-len("catalogsources"):] == "catalogsources" {
		_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"CatalogSourceList","items":[{"metadata":{"name":"certified-operators"},"status":{"connectionState":{"lastObservedState":"READY"}}},{"metadata":{"name":"community-operators"},"status":{"connectionState":{"lastObservedState":"READY"}}}]}`))

		return
	}

	if request.Method == http.MethodGet && strings.Contains(request.URL.Path, "/catalogsources/") {
		name := request.URL.Path[strings.LastIndex(request.URL.Path, "/")+1:]

		displayName := ""
		if name == "custom-catalog" {
			displayName = "Custom Index"
		}
		_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"CatalogSource","metadata":{"name":"` + name + `","namespace":"openshift-marketplace"},"spec":{"displayName":"` + displayName + `"}}`))

		return
	}

	if request.Method == http.MethodGet && len(request.URL.Path) >= len("clusterversions/version") && request.URL.Path[len(request.URL.Path)-len("clusterversions/version"):] == "clusterversions/version" {
		_, _ = responseWriter.Write([]byte(`{"apiVersion":"config.openshift.io/v1","kind":"ClusterVersion","metadata":{"name":"version"},"status":{"desired":{"version":"4.18.0"}}}`))

		return
	}

	if request.Method == http.MethodGet && (strings.HasSuffix(request.URL.Path, "/installplans") || strings.HasSuffix(request.URL.Path, "/subscriptions") || strings.HasSuffix(request.URL.Path, "/clusterserviceversions")) {
		_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"List","items":[]}`))

		return
	}

	if request.Method == http.MethodGet && request.URL.Path != "" && request.URL.Path[len(request.URL.Path)-1] != '/' {
		_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"Status","status":"Success"}`))

		return
	}
	_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"List","items":[]}`))
}

var _ kubernetes.Interface = (*k8sfake.Clientset)(nil)

package crd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crscaleoperator "github.com/test-network-function/cr-scale-operator/api/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
)

var errInjected = errors.New("injected failure")

const (
	testCRDName   = "memcacheds.cache.example.com"
	testCRName    = "test-cr"
	testNamespace = "test-ns"
	testSelector  = "app=test"
)

var testLabels = map[string]string{"app": "test"}

// withFakeCrds installs fake k8s and apiextensions clients seeded with crds,
// restoring the global client on cleanup. Tests using it must not run in parallel.
func withFakeCrds(t *testing.T, crds ...runtime.Object) *apiextfake.Clientset {
	t.Helper()

	globalhelper.SetTestK8sAPIClient(k8sfake.NewClientset())
	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	apiextensionsClient := apiextfake.NewSimpleClientset(crds...)
	globalhelper.GetAPIClient().ApiextensionsV1Interface = apiextensionsClient.ApiextensionsV1()

	return apiextensionsClient
}

// failVerb makes the fake apiextensions client fail the given verb with errInjected.
func failVerb(client *apiextfake.Clientset, verb string) {
	client.PrependReactor(verb, "customresourcedefinitions",
		func(k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errInjected
		})
}

// withRESTClient points the global client's core client at a test HTTP server, so
// the raw REST call in CreateCustomResourceScale hits it instead of a cluster.
// Tests using it must not run in parallel.
func withRESTClient(t *testing.T, serverURL string) {
	t.Helper()

	globalhelper.SetTestK8sAPIClient(k8sfake.NewClientset())
	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	coreClient, err := typedcorev1.NewForConfig(&rest.Config{Host: serverURL})
	require.NoError(t, err)

	globalhelper.GetAPIClient().CoreV1Interface = coreClient
}

// statusResponse writes a Kubernetes Status body with the given code and reason.
func statusResponse(responseWriter http.ResponseWriter, code int, reason metav1.StatusReason, message string) {
	status := metav1.Status{
		TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
		Status:   metav1.StatusFailure,
		Reason:   reason,
		Message:  message,
		Code:     int32(code),
	}

	encoded, _ := json.Marshal(status)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(code)
	_, _ = responseWriter.Write(encoded)
}

func testCrdNames() apiextv1.CustomResourceDefinitionNames {
	return apiextv1.CustomResourceDefinitionNames{
		Plural:     "memcacheds",
		Singular:   "memcached",
		Kind:       "Memcached",
		ListKind:   "MemcachedList",
		ShortNames: []string{"mc"},
	}
}

func TestDefineCustomResourceDefinition(t *testing.T) {
	names := testCrdNames()

	t.Run("builds the CRD without a status subresource", func(t *testing.T) {
		aCrd := DefineCustomResourceDefinition(names, "cache.example.com", false)

		assert.Equal(t, testCRDName, aCrd.Name)
		assert.Equal(t, "cache.example.com", aCrd.Spec.Group)
		assert.Equal(t, names, aCrd.Spec.Names)
		assert.Equal(t, apiextv1.NamespaceScoped, aCrd.Spec.Scope)
		require.Len(t, aCrd.Spec.Versions, 1)

		version := aCrd.Spec.Versions[0]
		assert.Equal(t, "v1", version.Name)
		assert.True(t, version.Served)
		assert.True(t, version.Storage)
		assert.Nil(t, version.Subresources)
		require.NotNil(t, version.Schema)
		require.NotNil(t, version.Schema.OpenAPIV3Schema)
		assert.Equal(t, "object", version.Schema.OpenAPIV3Schema.Type)

		properties := version.Schema.OpenAPIV3Schema.Properties
		require.Len(t, properties, 1)
		require.Contains(t, properties, "spec")
		assert.Equal(t, "object", properties["spec"].Type)
		assert.Contains(t, properties["spec"].Properties, "specProperty1")
	})

	t.Run("adds the status property and subresource", func(t *testing.T) {
		aCrd := DefineCustomResourceDefinition(names, "cache.example.com", true)
		version := aCrd.Spec.Versions[0]

		require.NotNil(t, version.Subresources)
		require.NotNil(t, version.Subresources.Status)

		properties := version.Schema.OpenAPIV3Schema.Properties
		require.Len(t, properties, 2)
		require.Contains(t, properties, "status")
		assert.Equal(t, "object", properties["status"].Type)
		assert.Contains(t, properties["status"].Properties, "statusProperty1")
	})
}

func TestDefineCustomResource(t *testing.T) {
	aCr := DefineCustomResource(testCRName, testNamespace, testSelector, testLabels)

	assert.Equal(t, crscaleoperator.GroupVersion.String(), aCr.APIVersion)
	assert.Equal(t, "Memcached", aCr.Kind)
	assert.Equal(t, testCRName, aCr.Name)
	assert.Equal(t, testNamespace, aCr.Namespace)
	assert.Equal(t, testLabels, aCr.Labels)
	assert.Equal(t, int32(1), aCr.Spec.Size)
	assert.Equal(t, testSelector, aCr.Status.Selector)
}

func TestRedefineCustomResourceWithReplica(t *testing.T) {
	aCr := DefineCustomResource(testCRName, testNamespace, testSelector, nil)

	RedefineCustomResourceWithReplica(aCr, 3)
	assert.Equal(t, int32(3), aCr.Spec.Size)
}

func TestEnsureCrdExists(t *testing.T) {
	t.Run("reports an existing CRD", func(t *testing.T) {
		withFakeCrds(t, &apiextv1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: testCRDName},
		})

		exists, err := EnsureCrdExists(testCRDName)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("reports a missing CRD", func(t *testing.T) {
		withFakeCrds(t)

		exists, err := EnsureCrdExists(testCRDName)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("returns the client error", func(t *testing.T) {
		client := withFakeCrds(t)
		failVerb(client, "get")

		exists, err := EnsureCrdExists(testCRDName)
		require.ErrorIs(t, err, errInjected)
		assert.False(t, exists)
	})
}

func TestGetAllCustomResourceDefinitions(t *testing.T) {
	t.Run("lists every CRD", func(t *testing.T) {
		withFakeCrds(t,
			&apiextv1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "b.example.com"}},
			&apiextv1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "a.example.com"}},
		)

		crds, err := GetAllCustomResourceDefinitions()
		require.NoError(t, err)
		require.Len(t, crds, 2)

		// The fake client lists in map order, so sort before asserting.
		sort.Slice(crds, func(left, right int) bool { return crds[left].Name < crds[right].Name })
		assert.Equal(t, "a.example.com", crds[0].Name)
		assert.Equal(t, "b.example.com", crds[1].Name)
	})

	t.Run("wraps the list failure", func(t *testing.T) {
		client := withFakeCrds(t)
		failVerb(client, "list")

		_, err := GetAllCustomResourceDefinitions()
		require.ErrorIs(t, err, errInjected)
	})
}

func TestDeleteCustomResourceDefinition(t *testing.T) {
	t.Run("deletes the CRD", func(t *testing.T) {
		client := withFakeCrds(t, &apiextv1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: testCRDName},
		})

		require.NoError(t, DeleteCustomResourceDefinition(testCRDName))

		_, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(
			context.TODO(), testCRDName, metav1.GetOptions{})
		assert.True(t, k8serrors.IsNotFound(err))
	})

	t.Run("tolerates a missing CRD", func(t *testing.T) {
		withFakeCrds(t)

		require.NoError(t, DeleteCustomResourceDefinition(testCRDName))
	})

	t.Run("returns the client error", func(t *testing.T) {
		client := withFakeCrds(t)
		failVerb(client, "delete")

		err := DeleteCustomResourceDefinition(testCRDName)
		require.ErrorIs(t, err, errInjected)
	})
}

func TestCreateCustomResourceScale(t *testing.T) {
	t.Run("posts the custom resource to the memcacheds endpoint", func(t *testing.T) {
		var gotMethod, gotPath string

		var gotBody []byte

		server := httptest.NewServer(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, request *http.Request) {
				gotMethod, gotPath = request.Method, request.URL.Path
				gotBody, _ = io.ReadAll(request.Body)

				responseWriter.WriteHeader(http.StatusCreated)
			}))
		t.Cleanup(server.Close)

		withRESTClient(t, server.URL)

		result, err := CreateCustomResourceScale(testCRName, testNamespace, testSelector, testLabels)
		require.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.Equal(t, http.MethodPost, gotMethod)
		assert.Equal(t, "/apis/cache.example.com/v1/namespaces/"+testNamespace+"/memcacheds", gotPath)

		var posted crscaleoperator.Memcached

		require.NoError(t, json.Unmarshal(gotBody, &posted))
		assert.Equal(t, testCRName, posted.Name)
		assert.Equal(t, testNamespace, posted.Namespace)
		assert.Equal(t, testLabels, posted.Labels)
		assert.Equal(t, int32(1), posted.Spec.Size)
	})

	t.Run("tolerates an already existing custom resource", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, _ *http.Request) {
				statusResponse(responseWriter, http.StatusConflict, metav1.StatusReasonAlreadyExists,
					fmt.Sprintf("%s %q already exists", testCRDName, testCRName))
			}))
		t.Cleanup(server.Close)

		withRESTClient(t, server.URL)

		result, err := CreateCustomResourceScale(testCRName, testNamespace, testSelector, nil)
		require.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("wraps the POST failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, _ *http.Request) {
				statusResponse(responseWriter, http.StatusInternalServerError,
					metav1.StatusReasonInternalError, "boom")
			}))
		t.Cleanup(server.Close)

		withRESTClient(t, server.URL)

		result, err := CreateCustomResourceScale(testCRName, testNamespace, testSelector, nil)
		require.Error(t, err)
		assert.Empty(t, result)

		// DoRaw reports unstructured errors, so the assertion goes by status code.
		var statusErr *k8serrors.StatusError

		require.ErrorAs(t, err, &statusErr)
		assert.Equal(t, http.StatusInternalServerError, int(statusErr.Status().Code))
	})

	t.Run("wraps the marshaling failure", func(t *testing.T) {
		oldMarshal := marshalCustomResource

		t.Cleanup(func() { marshalCustomResource = oldMarshal })

		marshalCustomResource = func(any) ([]byte, error) { return nil, errInjected }

		result, err := CreateCustomResourceScale(testCRName, testNamespace, testSelector, nil)
		require.ErrorIs(t, err, errInjected)
		assert.Empty(t, result)
	})
}

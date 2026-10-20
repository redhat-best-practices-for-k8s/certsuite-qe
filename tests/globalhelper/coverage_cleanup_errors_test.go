//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestCoverageNamespaceCleanupErrors(t *testing.T) {
	k8s, _ := coverageFixture(t, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "clean"}})
	k8s.PrependReactor("delete", "namespaces", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("delete namespace denied")
	})
	assert.Error(t, DeleteNamespaceAndWait("clean", time.Millisecond))

	k8s, _ = coverageFixture(t, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "clean"}})
	k8s.PrependReactor("list", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("PVC list failed") })
	badList := failingControllerClient{Client: egiClient.Client, listErr: errors.New("list failed")}
	egiClient.Client = badList
	assert.Error(t, cleanPVCs("clean", egiClient))

	dyn, ok := egiClient.Interface.(*dynamicfake.FakeDynamicClient)
	if !ok {
		t.Fatal("expected the test dynamic client")
	}

	dyn.PrependReactor("delete-collection", "*", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("clean objects failed")
	})
	assert.Error(t, cleanNamespace("clean", GetAPIClient(), egiClient)) // CleanObjects delete fails

	coverageFixture(t, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "clean"}})
	clientSet := *GetAPIClient()
	clientSet.Client = failingControllerClient{Client: egiClient.Client, listErr: errors.New("NAD list failed")}
	assert.Error(t, cleanNamespace("clean", &clientSet, egiClient))

	k8s, _ = coverageFixture(t, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "clean"}}, &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "claim", Namespace: "clean"}})
	k8s.PrependReactor("delete", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("PVC delete failed")
	})
	assert.Error(t, cleanPVCs("clean", egiClient))

	assertCleanNamespaceHTTPFailure(t, "subscriptions")
	assertCleanNamespaceHTTPFailure(t, "clusterserviceversions")
	assertCleanNamespaceHTTPFailure(t, "installplans")

	k8s, _ = coverageFixture(t, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "clean"}})
	k8s.PrependReactor("delete-collection", "poddisruptionbudgets", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("PDB delete failed")
	})
	assert.Error(t, cleanNamespace("clean", GetAPIClient(), egiClient))
}

func assertCleanNamespaceHTTPFailure(t *testing.T, resource string) {
	t.Helper()
	handler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodDelete && strings.Contains(request.URL.Path, resource) {
			http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"cleanup failed","code":500}`, http.StatusInternalServerError)

			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		_, _ = responseWriter.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Success","code":200}`))
	}
	coverageFixtureWithOLMHandler(t, handler, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "clean"}})
	assert.Error(t, cleanNamespace("clean", GetAPIClient(), egiClient))
}

//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	operatorclient "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stesting "k8s.io/client-go/testing"
)

func TestCoverageOperatorControllerErrors(t *testing.T) {
	k8s, _ := coverageFixture(t)
	operatorGroup := &olmv1.OperatorGroup{ObjectMeta: metav1.ObjectMeta{Name: "og", Namespace: "operators"}}

	k8s.PrependReactor("create", "namespaces", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("namespace create denied")
	})
	assert.Error(t, DeployOperatorGroup("operators", operatorGroup))
	coverageFixture(t)
	broken := failingControllerClient{Client: egiClient.Client, createErr: errors.New("create denied")}
	egiClient.Client, apiclient.Client = broken, broken

	assert.Error(t, DeployOperatorGroup("operators", operatorGroup))
	assert.Error(t, DeployOperator(&operatorv1alpha1.Subscription{ObjectMeta: metav1.ObjectMeta{Name: "sub", Namespace: "operators"}}))
	assert.Error(t, deployRHCatalogSource("redhat-operators", "redhat/index", "Red Hat", "4.18"))
	assert.Error(t, DeployCustomOperatorSource("custom/image"))
	assert.Error(t, CreateAndValidateCatalogSources(true))
	oldVersionLookup := clusterVersionForCatalogSource
	clusterVersionForCatalogSource = func() (string, error) { return "", errors.New("cluster version unavailable") }

	t.Cleanup(func() { clusterVersionForCatalogSource = oldVersionLookup })
	assert.Error(t, deployRHCatalogSource("redhat-operators", "redhat/index", "Red Hat", ""))

	coverageFixture(t)
	alreadyExists := failingControllerClient{Client: egiClient.Client, createErr: k8serrors.NewAlreadyExists(schema.GroupResource{Resource: "operatorgroups"}, "og")}
	egiClient.Client, apiclient.Client = alreadyExists, alreadyExists

	assert.NoError(t, DeployOperatorGroup("operators", operatorGroup))
	assert.NoError(t, DeployOperator(&operatorv1alpha1.Subscription{ObjectMeta: metav1.ObjectMeta{Name: "sub", Namespace: "operators"}}))
	alreadyExists = failingControllerClient{Client: egiClient.Client, createErr: k8serrors.NewAlreadyExists(schema.GroupResource{Resource: "catalogsources"}, "custom-catalog")}
	egiClient.Client, apiclient.Client = alreadyExists, alreadyExists

	assert.NoError(t, DeployCustomOperatorSource("custom/image"))
}

func TestCoverageOperatorRESTErrorAndMatchPaths(t *testing.T) {
	apiError := func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Error(responseWriter, `{"kind":"Status","apiVersion":"v1","status":"Failure","message":"operator API failed","code":500}`, http.StatusInternalServerError)
	}
	_, olm := coverageFixtureWithOLMHandler(t, apiError)
	_, err := GetInstallPlanByCSV("operators", "csv")
	assert.Error(t, err)
	assert.Error(t, updateInstallPlan("operators", &operatorv1alpha1.InstallPlan{ObjectMeta: metav1.ObjectMeta{Name: "plan", Namespace: "operators"}}))
	enabled, err := IsCatalogSourceEnabled("source", "operators", "display")
	assert.NoError(t, err)
	assert.False(t, enabled)
	assert.Error(t, DeleteCatalogSource("source", "operators", "display"))
	assert.Error(t, deleteSubscription("operators", "sub", olm))
	badSubscription := subscriptionGetError{SubscriptionInterface: olm.Subscriptions("operators"), err: errors.New("subscription get failed")}
	assert.Error(t, deleteSubscription("operators", "sub", operatorsWithSubscriptionOverride{OperatorsV1alpha1Interface: olm, subscriptionClient: badSubscription}))
	notFoundSubscription := subscriptionGetError{SubscriptionInterface: olm.Subscriptions("operators"), err: k8serrors.NewNotFound(schema.GroupResource{Resource: "subscriptions"}, "sub")}
	assert.NoError(t, deleteSubscription("operators", "sub", operatorsWithSubscriptionOverride{OperatorsV1alpha1Interface: olm, subscriptionClient: notFoundSubscription}))
	assert.Error(t, createSubscription("operators", &operatorv1alpha1.Subscription{ObjectMeta: metav1.ObjectMeta{Name: "sub", Namespace: "operators"}}, olm))
	assert.Error(t, createInstallPlan(&operatorv1alpha1.InstallPlan{ObjectMeta: metav1.ObjectMeta{Name: "plan", Namespace: "operators"}}, olm))
	assert.Error(t, DisableCatalogSource("community-operators"))

	handler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.URL.Path, "/installplans") {
			responseWriter.Header().Set("Content-Type", "application/json")
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"operators.coreos.com/v1alpha1","kind":"InstallPlanList","items":[{"metadata":{"name":"plan"},"spec":{"clusterServiceVersionNames":["example-operator.v1.0.0"]}}]}`))

			return
		}

		defaultOLMHTTPHandler(responseWriter, request)
	}
	_, olm = coverageFixtureWithOLMHandler(t, handler)
	plan, err := GetInstallPlanByCSV("operators", "operator")
	assert.NoError(t, err)
	assert.Equal(t, "plan", plan.Name)
	plan, err = GetInstallPlanByCSV("operators", "not-present")
	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.NoError(t, createInstallPlan(&operatorv1alpha1.InstallPlan{ObjectMeta: metav1.ObjectMeta{Name: "plan", Namespace: "operators"}}, olm))

	notFoundHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			responseWriter.Header().Set("Content-Type", "application/json")
			responseWriter.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(responseWriter).Encode(map[string]any{"apiVersion": "v1", "kind": "Status", "status": "Failure", "reason": "NotFound", "code": http.StatusNotFound})

			return
		}

		defaultOLMHTTPHandler(responseWriter, request)
	}
	coverageFixtureWithOLMHandler(t, notFoundHandler)
	assert.NoError(t, DeleteCatalogSource("absent", "operators", "display"))

	alreadyExistsHandler := func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost && strings.Contains(request.URL.Path, "/installplans") {
			responseWriter.Header().Set("Content-Type", "application/json")
			responseWriter.WriteHeader(http.StatusConflict)
			_, _ = responseWriter.Write([]byte(`{"apiVersion":"v1","kind":"Status","status":"Failure","reason":"AlreadyExists","code":409}`))

			return
		}

		defaultOLMHTTPHandler(responseWriter, request)
	}
	_, olm = coverageFixtureWithOLMHandler(t, alreadyExistsHandler)
	assert.NoError(t, createInstallPlan(&operatorv1alpha1.InstallPlan{ObjectMeta: metav1.ObjectMeta{Name: "duplicate", Namespace: "operators"}}, olm))
}

type operatorsWithSubscriptionOverride struct {
	operatorclient.OperatorsV1alpha1Interface
	subscriptionClient operatorclient.SubscriptionInterface
}

//nolint:ireturn // This fake must implement the typed client interface.
func (c operatorsWithSubscriptionOverride) Subscriptions(string) operatorclient.SubscriptionInterface {
	return c.subscriptionClient
}

type subscriptionGetError struct {
	operatorclient.SubscriptionInterface
	err error
}

func (c subscriptionGetError) Get(context.Context, string, metav1.GetOptions) (*operatorv1alpha1.Subscription, error) {
	return nil, c.err
}

package helper

import (
	"context"
	"errors"
	"testing"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	olmclient "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
)

const testNamespace = "affiliated-ns"

var errInjected = errors.New("injected failure")

// fakeCSVs implements only the CSV list and update calls the helpers use; any other
// method panics on the nil embedded interface.
type fakeCSVs struct {
	olmclient.ClusterServiceVersionInterface

	csvs      []v1alpha1.ClusterServiceVersion
	listErr   error
	updated   []*v1alpha1.ClusterServiceVersion
	updateErr error
}

func (f *fakeCSVs) List(context.Context, metav1.ListOptions) (*v1alpha1.ClusterServiceVersionList, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}

	return &v1alpha1.ClusterServiceVersionList{Items: f.csvs}, nil
}

func (f *fakeCSVs) Update(_ context.Context, csv *v1alpha1.ClusterServiceVersion,
	_ metav1.UpdateOptions) (*v1alpha1.ClusterServiceVersion, error) {
	f.updated = append(f.updated, csv)

	return csv, f.updateErr
}

// fakeOLM implements only the CSV client accessor the helpers use.
type fakeOLM struct {
	olmclient.OperatorsV1alpha1Interface

	csvs *fakeCSVs
}

//nolint:ireturn // Must match the OLM typed client interface.
func (f *fakeOLM) ClusterServiceVersions(string) olmclient.ClusterServiceVersionInterface {
	return f.csvs
}

// fakeCrds implements only the CRD list and delete calls the helpers use.
type fakeCrds struct {
	apiextclient.CustomResourceDefinitionInterface

	crds      []apiextv1.CustomResourceDefinition
	listErr   error
	deleted   []string
	deleteErr error
}

func (f *fakeCrds) List(context.Context, metav1.ListOptions) (*apiextv1.CustomResourceDefinitionList, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}

	return &apiextv1.CustomResourceDefinitionList{Items: f.crds}, nil
}

func (f *fakeCrds) Delete(_ context.Context, name string, _ metav1.DeleteOptions) error {
	f.deleted = append(f.deleted, name)

	return f.deleteErr
}

// fakeApiextensions implements only the CRD client accessor the helpers use.
type fakeApiextensions struct {
	apiextclient.ApiextensionsV1Interface

	crds *fakeCrds
}

//nolint:ireturn // Must match the apiextensions typed client interface.
func (f *fakeApiextensions) CustomResourceDefinitions() apiextclient.CustomResourceDefinitionInterface {
	return f.crds
}

// withFakeOLM installs fake k8s and OLM CSV clients, restoring the global client on cleanup.
// Tests using it must not run in parallel.
func withFakeOLM(t *testing.T, csvs ...v1alpha1.ClusterServiceVersion) *fakeCSVs {
	t.Helper()

	globalhelper.SetTestK8sAPIClient(k8sfake.NewClientset())

	csvClient := &fakeCSVs{csvs: csvs}
	globalhelper.GetAPIClient().OperatorsV1alpha1Interface = &fakeOLM{csvs: csvClient}

	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	return csvClient
}

// withFakeCrds installs fake k8s and apiextensions clients, restoring the global client on cleanup.
// Tests using it must not run in parallel.
func withFakeCrds(t *testing.T, crds ...apiextv1.CustomResourceDefinition) *fakeCrds {
	t.Helper()

	globalhelper.SetTestK8sAPIClient(k8sfake.NewClientset())

	crdClient := &fakeCrds{crds: crds}
	globalhelper.GetAPIClient().ApiextensionsV1Interface = &fakeApiextensions{crds: crdClient}

	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	return crdClient
}

func testCSV(name string, labels map[string]string) v1alpha1.ClusterServiceVersion {
	return v1alpha1.ClusterServiceVersion{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace, Labels: labels},
	}
}

func testCrd(name string) apiextv1.CustomResourceDefinition {
	return apiextv1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: name}}
}

func TestGetCsvByPrefix(t *testing.T) {
	t.Run("returns the last matching CSV", func(t *testing.T) {
		withFakeOLM(t,
			testCSV("other-operator.v1", nil),
			testCSV("my-operator.v1", nil),
			testCSV("my-operator.v2", nil))

		csv, err := GetCsvByPrefix("my-operator", testNamespace)
		require.NoError(t, err)
		assert.Equal(t, "my-operator.v2", csv.Name)
	})

	t.Run("no matching CSV", func(t *testing.T) {
		withFakeOLM(t, testCSV("other-operator.v1", nil))

		_, err := GetCsvByPrefix("my-operator", testNamespace)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to detect a needed CSV")
	})

	t.Run("list failure", func(t *testing.T) {
		csvClient := withFakeOLM(t)
		csvClient.listErr = errInjected

		_, err := GetCsvByPrefix("my-operator", testNamespace)
		require.ErrorIs(t, err, errInjected)
	})
}

func TestAddLabelToInstalledCSV(t *testing.T) {
	t.Run("merges the label into the CSV", func(t *testing.T) {
		csvClient := withFakeOLM(t, testCSV("my-operator.v1", map[string]string{"keep": "yes"}))

		require.NoError(t, AddLabelToInstalledCSV("my-operator", testNamespace, map[string]string{"new": "label"}))

		require.Len(t, csvClient.updated, 1)
		assert.Equal(t, map[string]string{"keep": "yes", "new": "label"}, csvClient.updated[0].Labels)
	})

	t.Run("no matching CSV", func(t *testing.T) {
		withFakeOLM(t)

		err := AddLabelToInstalledCSV("my-operator", testNamespace, map[string]string{"new": "label"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to detect a needed CSV")
	})

	t.Run("update failure", func(t *testing.T) {
		csvClient := withFakeOLM(t, testCSV("my-operator.v1", nil))
		csvClient.updateErr = errInjected

		err := AddLabelToInstalledCSV("my-operator", testNamespace, map[string]string{"new": "label"})
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "fail to update CSV")
	})
}

func TestDeleteLabelFromInstalledCSV(t *testing.T) {
	t.Run("removes the label from the CSV", func(t *testing.T) {
		csvClient := withFakeOLM(t, testCSV("my-operator.v1",
			map[string]string{"keep": "yes", "drop": "me"}))

		require.NoError(t, DeleteLabelFromInstalledCSV("my-operator", testNamespace, map[string]string{"drop": "ignored"}))

		require.Len(t, csvClient.updated, 1)
		assert.Equal(t, map[string]string{"keep": "yes"}, csvClient.updated[0].Labels)
	})

	t.Run("no matching CSV", func(t *testing.T) {
		withFakeOLM(t)

		err := DeleteLabelFromInstalledCSV("my-operator", testNamespace, map[string]string{"drop": "me"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to detect a needed CSV")
	})

	t.Run("update failure", func(t *testing.T) {
		csvClient := withFakeOLM(t, testCSV("my-operator.v1", nil))
		csvClient.updateErr = errInjected

		err := DeleteLabelFromInstalledCSV("my-operator", testNamespace, map[string]string{"drop": "me"})
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "fail to update CSV")
	})
}

func TestDoesOperatorHaveLabels(t *testing.T) {
	t.Run("all labels match", func(t *testing.T) {
		withFakeOLM(t, testCSV("my-operator.v1", map[string]string{"app": "mine", "tier": "gold"}))

		found, err := DoesOperatorHaveLabels("my-operator", testNamespace, map[string]string{"app": "mine"})
		require.NoError(t, err)
		assert.True(t, found)
	})

	t.Run("missing label", func(t *testing.T) {
		withFakeOLM(t, testCSV("my-operator.v1", map[string]string{"app": "mine"}))

		found, err := DoesOperatorHaveLabels("my-operator", testNamespace, map[string]string{"tier": "gold"})
		require.NoError(t, err)
		assert.False(t, found)
	})

	t.Run("value mismatch", func(t *testing.T) {
		withFakeOLM(t, testCSV("my-operator.v1", map[string]string{"app": "mine"}))

		found, err := DoesOperatorHaveLabels("my-operator", testNamespace, map[string]string{"app": "other"})
		require.NoError(t, err)
		assert.False(t, found)
	})

	t.Run("empty labels match any CSV", func(t *testing.T) {
		withFakeOLM(t, testCSV("my-operator.v1", map[string]string{"app": "mine"}))

		found, err := DoesOperatorHaveLabels("my-operator", testNamespace, map[string]string{})
		require.NoError(t, err)
		assert.True(t, found)
	})

	t.Run("no matching CSV", func(t *testing.T) {
		withFakeOLM(t)

		_, err := DoesOperatorHaveLabels("my-operator", testNamespace, map[string]string{"app": "mine"})
		require.Error(t, err)
	})
}

func TestDeployOperatorSubscription(t *testing.T) {
	oldDeploy := deployOperator

	t.Cleanup(func() { deployOperator = oldDeploy })

	t.Run("deploys the subscription", func(t *testing.T) {
		var deployed *v1alpha1.Subscription

		deployOperator = func(subscription *v1alpha1.Subscription) error {
			deployed = subscription

			return nil
		}

		err := DeployOperatorSubscription("my-sub", "my-package", "stable", testNamespace,
			"my-group", "my-ns", "my-csv.v1", v1alpha1.ApprovalAutomatic)
		require.NoError(t, err)

		require.NotNil(t, deployed)
		assert.Equal(t, "my-sub-subscription", deployed.Name)
		assert.Equal(t, testNamespace, deployed.Namespace)
		require.NotNil(t, deployed.Spec)
		assert.Equal(t, "stable", deployed.Spec.Channel)
		assert.Equal(t, "my-package", deployed.Spec.Package)
		assert.Equal(t, "my-group", deployed.Spec.CatalogSource)
		assert.Equal(t, "my-ns", deployed.Spec.CatalogSourceNamespace)
		assert.Equal(t, "my-csv.v1", deployed.Spec.StartingCSV)
		assert.Equal(t, v1alpha1.ApprovalAutomatic, deployed.Spec.InstallPlanApproval)
	})

	t.Run("deployment failure", func(t *testing.T) {
		deployOperator = func(*v1alpha1.Subscription) error { return errInjected }

		err := DeployOperatorSubscription("my-sub", "my-package", "stable", testNamespace,
			"my-group", "my-ns", "my-csv.v1", v1alpha1.ApprovalManual)
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "Error deploying operator my-package")
	})
}

func TestDeleteIstioCRDs(t *testing.T) {
	t.Run("deletes only the istio CRDs", func(t *testing.T) {
		crdClient := withFakeCrds(t,
			testCrd("gateways.istio.io"),
			testCrd("others.example.com"),
			testCrd("adapters.istio.io"))

		require.NoError(t, DeleteIstioCRDs())
		assert.Equal(t, []string{"gateways.istio.io", "adapters.istio.io"}, crdClient.deleted)
	})

	t.Run("no istio CRDs", func(t *testing.T) {
		crdClient := withFakeCrds(t, testCrd("others.example.com"))

		require.NoError(t, DeleteIstioCRDs())
		assert.Empty(t, crdClient.deleted)
	})

	t.Run("list failure", func(t *testing.T) {
		crdClient := withFakeCrds(t)
		crdClient.listErr = errInjected

		err := DeleteIstioCRDs()
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to list cluster CRDs")
	})

	t.Run("delete failure", func(t *testing.T) {
		crdClient := withFakeCrds(t, testCrd("gateways.istio.io"))
		crdClient.deleteErr = errInjected

		err := DeleteIstioCRDs()
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to delete istio CRD gateways.istio.io")
	})
}

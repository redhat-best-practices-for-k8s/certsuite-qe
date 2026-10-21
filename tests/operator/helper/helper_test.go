package helper

import (
	"context"
	"errors"
	"testing"
	"time"

	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	olmclient "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
)

const testNamespace = "operator-ns"

var errInjected = errors.New("injected failure")

// sequence returns its responses in order and keeps repeating the last one.
type sequence[T any] struct {
	responses []func() (T, error)
	calls     int
}

//nolint:ireturn // T is the concrete list type at every call site.
func (s *sequence[T]) next() (T, error) {
	index := min(s.calls, len(s.responses)-1)
	s.calls++

	return s.responses[index]()
}

func (s *sequence[T]) respond(responses ...func() (T, error)) {
	s.responses = responses
}

func returns[T any](value T) func() (T, error) {
	return func() (T, error) { return value, nil }
}

func fails[T any]() func() (T, error) {
	return func() (T, error) {
		var zero T

		return zero, errInjected
	}
}

// fakeOLM implements only the parts of the OLM typed client the helpers use; any other
// method panics on the nil embedded interface.
type fakeOLM struct {
	olmclient.OperatorsV1alpha1Interface

	csvs *fakeCSVs
	subs *fakeSubscriptions
}

//nolint:ireturn // Must match the OLM typed client interface.
func (f *fakeOLM) ClusterServiceVersions(string) olmclient.ClusterServiceVersionInterface {
	return f.csvs
}

//nolint:ireturn // Must match the OLM typed client interface.
func (f *fakeOLM) Subscriptions(string) olmclient.SubscriptionInterface {
	return f.subs
}

type fakeCSVs struct {
	olmclient.ClusterServiceVersionInterface

	list      sequence[*v1alpha1.ClusterServiceVersionList]
	updated   []*v1alpha1.ClusterServiceVersion
	updateErr error
}

func (f *fakeCSVs) List(context.Context, metav1.ListOptions) (*v1alpha1.ClusterServiceVersionList, error) {
	return f.list.next()
}

func (f *fakeCSVs) Update(_ context.Context, csv *v1alpha1.ClusterServiceVersion,
	_ metav1.UpdateOptions) (*v1alpha1.ClusterServiceVersion, error) {
	f.updated = append(f.updated, csv)

	return csv, f.updateErr
}

type fakeSubscriptions struct {
	olmclient.SubscriptionInterface

	list sequence[*v1alpha1.SubscriptionList]
}

func (f *fakeSubscriptions) List(context.Context, metav1.ListOptions) (*v1alpha1.SubscriptionList, error) {
	return f.list.next()
}

func csvList(csvs ...v1alpha1.ClusterServiceVersion) *v1alpha1.ClusterServiceVersionList {
	return &v1alpha1.ClusterServiceVersionList{Items: csvs}
}

func testCSV(name string, phase v1alpha1.ClusterServiceVersionPhase) v1alpha1.ClusterServiceVersion {
	return v1alpha1.ClusterServiceVersion{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace},
		Status:     v1alpha1.ClusterServiceVersionStatus{Phase: phase},
	}
}

func subscriptionList(installed, current string) *v1alpha1.SubscriptionList {
	return &v1alpha1.SubscriptionList{Items: []v1alpha1.Subscription{{
		ObjectMeta: metav1.ObjectMeta{Name: "sub", Namespace: testNamespace},
		Status:     v1alpha1.SubscriptionStatus{InstalledCSV: installed, CurrentCSV: current},
	}}}
}

// withFakeOLM installs fake k8s and OLM clients, stubs the globalhelper operator seams and
// shortens poll timers, restoring everything on cleanup. Tests using it must not run in parallel.
func withFakeOLM(t *testing.T, objects ...runtime.Object) (*fakeOLM, *k8sfake.Clientset) {
	t.Helper()

	oldInstalled, oldDeployGroup, oldDeploy := isOperatorGroupInstalled, deployOperatorGroup, deployOperator
	oldTimeout, oldInterval, oldSubTimeout := operatorReadyTimeout, operatorPollInterval, subscriptionCheckTimeout

	t.Cleanup(func() {
		isOperatorGroupInstalled, deployOperatorGroup, deployOperator = oldInstalled, oldDeployGroup, oldDeploy
		operatorReadyTimeout, operatorPollInterval, subscriptionCheckTimeout = oldTimeout, oldInterval, oldSubTimeout

		globalhelper.UnsetTestK8sAPIClient()
	})

	client := k8sfake.NewClientset(objects...)
	globalhelper.SetTestK8sAPIClient(client)

	olm := &fakeOLM{
		csvs: &fakeCSVs{list: sequence[*v1alpha1.ClusterServiceVersionList]{
			responses: []func() (*v1alpha1.ClusterServiceVersionList, error){returns(csvList())},
		}},
		subs: &fakeSubscriptions{list: sequence[*v1alpha1.SubscriptionList]{
			responses: []func() (*v1alpha1.SubscriptionList, error){returns(subscriptionList("csv", "csv"))},
		}},
	}
	globalhelper.GetAPIClient().OperatorsV1alpha1Interface = olm

	isOperatorGroupInstalled = func(string, string) error { return errInjected }
	deployOperatorGroup = func(string, *olmv1.OperatorGroup) error { return nil }
	deployOperator = func(*v1alpha1.Subscription) error { return nil }

	// Generous timeouts keep loaded CI from flaking; successful polls return immediately.
	operatorReadyTimeout = 30 * time.Second
	operatorPollInterval = time.Millisecond
	subscriptionCheckTimeout = 30 * time.Second

	return olm, client
}

func TestGetCsvByPrefix(t *testing.T) {
	t.Run("returns the last matching CSV", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		olm.csvs.list.respond(returns(csvList(
			testCSV("other.v1", v1alpha1.CSVPhaseSucceeded),
			testCSV("my-operator.v1", v1alpha1.CSVPhasePending),
			testCSV("my-operator.v2", v1alpha1.CSVPhaseSucceeded),
		)))

		csv, err := GetCsvByPrefix("my-operator", testNamespace)
		require.NoError(t, err)
		assert.Equal(t, "my-operator.v2", csv.Name)
	})

	t.Run("no matching CSV", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		olm.csvs.list.respond(returns(csvList(testCSV("other.v1", v1alpha1.CSVPhaseSucceeded))))

		_, err := GetCsvByPrefix("my-operator", testNamespace)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to detect a needed CSV")
	})

	t.Run("list failure", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		olm.csvs.list.respond(fails[*v1alpha1.ClusterServiceVersionList]())

		_, err := GetCsvByPrefix("my-operator", testNamespace)
		assert.ErrorIs(t, err, errInjected)
	})
}

func TestIsCSVNotSucceeded(t *testing.T) {
	testCases := []struct {
		phase v1alpha1.ClusterServiceVersionPhase
		want  bool
	}{
		{phase: v1alpha1.CSVPhaseSucceeded, want: false},
		{phase: v1alpha1.CSVPhaseFailed, want: true},
		{phase: v1alpha1.CSVPhaseInstalling, want: true},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.phase), func(t *testing.T) {
			olm, _ := withFakeOLM(t)
			olm.csvs.list.respond(returns(csvList(testCSV("my-operator.v1", testCase.phase))))

			notSucceeded, err := IsCSVNotSucceeded("my-operator", testNamespace)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, notSucceeded)
		})
	}

	t.Run("missing CSV", func(t *testing.T) {
		withFakeOLM(t)

		notSucceeded, err := IsCSVNotSucceeded("my-operator", testNamespace)
		require.Error(t, err)
		assert.False(t, notSucceeded)
	})
}

func TestAddLabelToInstalledCSV(t *testing.T) {
	existing := testCSV("my-operator.v1", v1alpha1.CSVPhaseSucceeded)
	existing.Labels = map[string]string{"keep": "me", "override": "old"}

	t.Run("merges labels and updates the CSV", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		olm.csvs.list.respond(returns(csvList(existing)))

		require.NoError(t, AddLabelToInstalledCSV("my-operator", testNamespace, map[string]string{"override": "new", "add": "x"}))
		require.Len(t, olm.csvs.updated, 1)
		assert.Equal(t, map[string]string{"keep": "me", "override": "new", "add": "x"}, olm.csvs.updated[0].Labels)
		assert.Equal(t, "old", existing.Labels["override"], "the listed CSV must not be mutated in place")
	})

	t.Run("update failure is wrapped", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		olm.csvs.list.respond(returns(csvList(existing)))
		olm.csvs.updateErr = errInjected

		err := AddLabelToInstalledCSV("my-operator", testNamespace, map[string]string{"a": "b"})
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "fail to update CSV")
	})

	t.Run("missing CSV", func(t *testing.T) {
		olm, _ := withFakeOLM(t)

		require.Error(t, AddLabelToInstalledCSV("my-operator", testNamespace, map[string]string{"a": "b"}))
		assert.Empty(t, olm.csvs.updated)
	})
}

func TestDeployTestOperatorGroup(t *testing.T) {
	testCases := []struct {
		name        string
		clusterWide bool
		wantTargets []string
	}{
		{name: "namespaced", wantTargets: []string{testNamespace}},
		{name: "cluster wide", clusterWide: true, wantTargets: []string{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			withFakeOLM(t)

			var deployed *olmv1.OperatorGroup

			deployOperatorGroup = func(namespace string, group *olmv1.OperatorGroup) error {
				assert.Equal(t, testNamespace, namespace)
				deployed = group

				return nil
			}

			require.NoError(t, DeployTestOperatorGroup(testNamespace, testCase.clusterWide))
			require.NotNil(t, deployed)
			assert.Equal(t, tsparams.OperatorGroupName, deployed.Name)
			assert.Equal(t, testCase.wantTargets, deployed.Spec.TargetNamespaces)
		})
	}

	t.Run("already installed", func(t *testing.T) {
		withFakeOLM(t)
		isOperatorGroupInstalled = func(string, string) error { return nil }
		deployOperatorGroup = func(string, *olmv1.OperatorGroup) error {
			t.Fatal("operator group must not be redeployed")

			return nil
		}

		assert.NoError(t, DeployTestOperatorGroup(testNamespace, false))
	})

	t.Run("deploy failure", func(t *testing.T) {
		withFakeOLM(t)
		deployOperatorGroup = func(string, *olmv1.OperatorGroup) error { return errInjected }

		assert.ErrorIs(t, DeployTestOperatorGroup(testNamespace, false), errInjected)
	})
}

func TestDeployTestOperatorGroupWithTargetNamespace(t *testing.T) {
	targets := []string{"ns-a", "ns-b"}

	t.Run("deploys when missing", func(t *testing.T) {
		withFakeOLM(t)

		var deployed *olmv1.OperatorGroup

		deployOperatorGroup = func(_ string, group *olmv1.OperatorGroup) error {
			deployed = group

			return nil
		}

		require.NoError(t, DeployTestOperatorGroupWithTargetNamespace("custom-og", testNamespace, targets))
		require.NotNil(t, deployed)
		assert.Equal(t, "custom-og", deployed.Name)
		assert.Equal(t, targets, deployed.Spec.TargetNamespaces)
	})

	t.Run("already installed", func(t *testing.T) {
		withFakeOLM(t)
		isOperatorGroupInstalled = func(name, namespace string) error {
			assert.Equal(t, "custom-og", name)
			assert.Equal(t, testNamespace, namespace)

			return nil
		}
		deployOperatorGroup = func(string, *olmv1.OperatorGroup) error {
			t.Fatal("operator group must not be redeployed")

			return nil
		}

		assert.NoError(t, DeployTestOperatorGroupWithTargetNamespace("custom-og", testNamespace, targets))
	})
}

func TestDeployOperatorSubscription(t *testing.T) {
	t.Run("defines the subscription", func(t *testing.T) {
		withFakeOLM(t)

		var deployed *v1alpha1.Subscription

		deployOperator = func(sub *v1alpha1.Subscription) error {
			deployed = sub

			return nil
		}

		require.NoError(t, DeployOperatorSubscription("my-sub", "my-package", "stable", testNamespace,
			tsparams.CommunityOperatorGroup, tsparams.OperatorSourceNamespace, "my-package.v1", v1alpha1.ApprovalManual))
		require.NotNil(t, deployed)
		assert.Equal(t, "my-sub-subscription", deployed.Name)
		assert.Equal(t, testNamespace, deployed.Namespace)
		assert.Equal(t, v1alpha1.SubscriptionSpec{
			Channel:                "stable",
			Package:                "my-package",
			CatalogSource:          tsparams.CommunityOperatorGroup,
			CatalogSourceNamespace: tsparams.OperatorSourceNamespace,
			StartingCSV:            "my-package.v1",
			InstallPlanApproval:    v1alpha1.ApprovalManual,
		}, *deployed.Spec)
	})

	t.Run("deploy failure is wrapped", func(t *testing.T) {
		withFakeOLM(t)
		deployOperator = func(*v1alpha1.Subscription) error { return errInjected }

		err := DeployOperatorSubscription("my-sub", "my-package", "stable", testNamespace,
			tsparams.CommunityOperatorGroup, tsparams.OperatorSourceNamespace, "", v1alpha1.ApprovalAutomatic)
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "Error deploying operator my-package")
	})
}

func TestDeployOperatorSubscriptionWithNodeSelector(t *testing.T) {
	nodeSelector := map[string]string{"node-role.kubernetes.io/worker": ""}

	t.Run("defines the subscription with node selector", func(t *testing.T) {
		withFakeOLM(t)

		var deployed *v1alpha1.Subscription

		deployOperator = func(sub *v1alpha1.Subscription) error {
			deployed = sub

			return nil
		}

		require.NoError(t, DeployOperatorSubscriptionWithNodeSelector("my-package", "stable", testNamespace,
			tsparams.CommunityOperatorGroup, tsparams.OperatorSourceNamespace, "", v1alpha1.ApprovalAutomatic, nodeSelector))
		require.NotNil(t, deployed)
		assert.Equal(t, "my-package-subscription", deployed.Name)
		assert.Equal(t, "my-package", deployed.Spec.Package)
		assert.Equal(t, nodeSelector, deployed.Spec.Config.NodeSelector)
	})

	t.Run("deploy failure is wrapped", func(t *testing.T) {
		withFakeOLM(t)
		deployOperator = func(*v1alpha1.Subscription) error { return errInjected }

		err := DeployOperatorSubscriptionWithNodeSelector("my-package", "stable", testNamespace,
			tsparams.CommunityOperatorGroup, tsparams.OperatorSourceNamespace, "", v1alpha1.ApprovalAutomatic, nodeSelector)
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "Error deploying operator my-package")
	})
}

func testDeployment(name string, replicas, ready int32) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace},
		Spec:       appsv1.DeploymentSpec{Replicas: ptr.To(replicas)},
		Status: appsv1.DeploymentStatus{
			Replicas:          replicas,
			ReadyReplicas:     ready,
			AvailableReplicas: ready,
			Conditions: []appsv1.DeploymentCondition{
				{Type: appsv1.DeploymentAvailable, Status: corev1.ConditionTrue},
			},
		},
	}
}

func TestWaitUntilOperatorIsReady(t *testing.T) {
	t.Run("retries until the CSV succeeds and the subscription settles", func(t *testing.T) {
		olm, _ := withFakeOLM(t, testDeployment("ready", 1, 1), testDeployment("not-ready", 2, 1))

		pending := testCSV("my-operator.v1", v1alpha1.CSVPhaseInstalling)
		pending.Status.Conditions = []v1alpha1.ClusterServiceVersionCondition{{Phase: v1alpha1.CSVPhaseInstalling}}
		pending.Status.RequirementStatus = []v1alpha1.RequirementStatus{{Kind: "CustomResourceDefinition", Name: "crd"}}

		olm.csvs.list.respond(fails[*v1alpha1.ClusterServiceVersionList](),
			returns(csvList(pending)),
			returns(csvList(testCSV("my-operator.v1", v1alpha1.CSVPhaseSucceeded))))
		olm.subs.list.respond(fails[*v1alpha1.SubscriptionList](),
			returns(&v1alpha1.SubscriptionList{}),
			returns(subscriptionList("", "my-operator.v1")),
			returns(subscriptionList("my-operator.v0", "my-operator.v1")),
			returns(subscriptionList("my-operator.v1", "my-operator.v1")))

		require.NoError(t, WaitUntilOperatorIsReady("my-operator", testNamespace))
		assert.Equal(t, 3, olm.csvs.list.calls)
		assert.Equal(t, 5, olm.subs.list.calls)
	})

	t.Run("times out while the CSV is pending", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		operatorReadyTimeout = 20 * time.Millisecond

		olm.csvs.list.respond(returns(csvList(testCSV("my-operator.v1", v1alpha1.CSVPhasePending))))

		err := WaitUntilOperatorIsReady("my-operator", testNamespace)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timed out")
		assert.Zero(t, olm.subs.list.calls)
	})

	t.Run("operator is ready even if subscription verification times out", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		subscriptionCheckTimeout = 10 * time.Millisecond

		olm.csvs.list.respond(returns(csvList(testCSV("my-operator.v1", v1alpha1.CSVPhaseSucceeded))))
		olm.subs.list.respond(returns(&v1alpha1.SubscriptionList{}))

		require.NoError(t, WaitUntilOperatorIsReady("my-operator", testNamespace))
		assert.Positive(t, olm.subs.list.calls)
	})
}

func TestCheckInstalledCSV(t *testing.T) {
	t.Run("installed CSV matches current CSV", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		olm.subs.list.respond(returns(subscriptionList("my-operator.v1", "my-operator.v1")))

		assert.True(t, checkInstalledCSV(testNamespace))
	})

	t.Run("times out while installed CSV lags current CSV", func(t *testing.T) {
		olm, _ := withFakeOLM(t)
		subscriptionCheckTimeout = 10 * time.Millisecond

		olm.subs.list.respond(returns(subscriptionList("my-operator.v0", "my-operator.v1")))

		assert.False(t, checkInstalledCSV(testNamespace))
	})
}

func TestDisplayDeploymentStatus(t *testing.T) {
	testCases := []struct {
		name        string
		deployments []runtime.Object
		listErr     bool
		wantReady   bool
	}{
		{name: "all deployments ready", deployments: []runtime.Object{testDeployment("ready", 1, 1)}, wantReady: true},
		{name: "deployment not ready", deployments: []runtime.Object{testDeployment("ready", 1, 1), testDeployment("lagging", 3, 2)}},
		{name: "no deployments"},
		{name: "list failure", deployments: []runtime.Object{testDeployment("ready", 1, 1)}, listErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, client := withFakeOLM(t, testCase.deployments...)
			if testCase.listErr {
				client.PrependReactor("list", "deployments", func(k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, errInjected
				})
			}

			assert.Equal(t, testCase.wantReady, displayDeploymentStatus(testNamespace, 0))
		})
	}
}

func TestDisplayCSVStatus(t *testing.T) {
	csv := testCSV("my-operator.v1", v1alpha1.CSVPhaseInstalling)

	// Logging only: a nil CSV must be tolerated.
	assert.NotPanics(t, func() {
		displayCSVStatus(nil, 0)
		displayCSVStatus(&csv, 0)
	})
}

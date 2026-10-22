package helper

import (
	"context"
	"errors"
	"testing"
	"time"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiMco "github.com/openshift-kni/eco-goinfra/pkg/mco"
	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"
	mcv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
)

const (
	masterRoleLabel = "node-role.kubernetes.io/master"
	testMCPName     = "worker"
	testMCName      = "99-worker-realtime"
)

var errInjected = errors.New("injected failure")

// withEgiClient replaces the eco-goinfra client with a fake that knows the MachineConfig types and
// restores the seams on cleanup. Tests using it must not run in parallel.
func withEgiClient(t *testing.T, funcs *interceptor.Funcs, objects ...runtime.Object) *egiClients.Settings {
	t.Helper()

	oldClient, oldTimeout := getEgiClient, mcpStartTimeout

	t.Cleanup(func() { getEgiClient, mcpStartTimeout = oldClient, oldTimeout })

	clientSet, builder := egiClients.GetModifiableTestClients(egiClients.TestClientParams{
		K8sMockObjects:  objects,
		SchemeAttachers: []egiClients.SchemeAttacher{machineconfigv1.Install},
	})

	if funcs != nil {
		builder = builder.WithInterceptorFuncs(*funcs)
	}

	clientSet.Client = builder.Build()
	getEgiClient = func() *egiClients.Settings { return clientSet }
	// Waits check immediately, so a short timeout only matters when a condition is never met.
	mcpStartTimeout = 100 * time.Millisecond

	return clientSet
}

func withoutEgiClient(t *testing.T) {
	t.Helper()

	oldClient := getEgiClient

	t.Cleanup(func() { getEgiClient = oldClient })

	getEgiClient = func() *egiClients.Settings { return nil }
}

func testNode(name string, nodeLabels map[string]string) *corev1.Node {
	return &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: nodeLabels}}
}

func testMCP(name string, machineCount int32, nodeSelector *metav1.LabelSelector) *machineconfigv1.MachineConfigPool {
	return &machineconfigv1.MachineConfigPool{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: machineconfigv1.MachineConfigPoolSpec{
			NodeSelector:          nodeSelector,
			MachineConfigSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"pool": name}},
		},
		Status: machineconfigv1.MachineConfigPoolStatus{MachineCount: machineCount},
	}
}

func failingList() *interceptor.Funcs {
	return &interceptor.Funcs{
		List: func(context.Context, runtimeclient.WithWatch, runtimeclient.ObjectList, ...runtimeclient.ListOption) error {
			return errInjected
		},
	}
}

func matchLabels(key string) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: map[string]string{key: ""}}
}

var invalidSelector = &metav1.LabelSelector{
	MatchExpressions: []metav1.LabelSelectorRequirement{{Key: "pool", Operator: "Bogus"}},
}

func TestDefineDeploymentWithImage(t *testing.T) {
	dep, err := DefineDeploymentWithImage(2, 3, testDeploymentName, parameters.TestAccessControlNameSpace, "quay.io/test/image:1")
	require.NoError(t, err)
	assert.Equal(t, int32(2), *dep.Spec.Replicas)
	require.Len(t, dep.Spec.Template.Spec.Containers, 3)

	for _, container := range dep.Spec.Template.Spec.Containers {
		assert.Equal(t, "quay.io/test/image:1", container.Image)
	}

	_, err = DefineDeploymentWithImage(1, 0, testDeploymentName, parameters.TestAccessControlNameSpace, "image")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid number of containers")
}

func TestDefineDeploymentHelpersRejectZeroContainers(t *testing.T) {
	namespace := parameters.TestAccessControlNameSpace
	builders := map[string]func() error{
		"DefineDeployment": func() error {
			_, err := DefineDeployment(1, 0, testDeploymentName, namespace)

			return err
		},
		"DefineDeploymentWithNamespace": func() error {
			_, err := DefineDeploymentWithNamespace(1, 0, testDeploymentName, namespace)

			return err
		},
		"DefineDeploymentWithClusterRoleBindingWithServiceAccount": func() error {
			_, err := DefineDeploymentWithClusterRoleBindingWithServiceAccount(1, 0, testDeploymentName, namespace, "sa")

			return err
		},
		"DefineDeploymentWithContainerPorts": func() error {
			_, err := DefineDeploymentWithContainerPorts(testDeploymentName, namespace, 1, nil)

			return err
		},
	}

	for name, build := range builders {
		t.Run(name, func(t *testing.T) {
			err := build()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid number of containers")
		})
	}
}

func TestDefineStatefulSet(t *testing.T) {
	sts, err := DefineStatefulSet(3, 2, "test-sts", parameters.TestAccessControlNameSpace)
	require.NoError(t, err)
	assert.Equal(t, "test-sts", sts.Name)
	assert.Equal(t, parameters.TestAccessControlNameSpace, sts.Namespace)
	assert.Equal(t, int32(3), *sts.Spec.Replicas)
	assert.Equal(t, parameters.TestStatefulSetLabels, sts.Spec.Template.Labels)
	require.Len(t, sts.Spec.Template.Spec.Containers, 2)
	assert.Equal(t, parameters.SampleWorkloadImage, sts.Spec.Template.Spec.Containers[1].Image)

	_, err = DefineStatefulSet(1, 0, "test-sts", parameters.TestAccessControlNameSpace)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid number of containers")
}

func withK8sClient(t *testing.T, objects ...runtime.Object) *k8sfake.Clientset {
	t.Helper()

	client := k8sfake.NewClientset(objects...)
	globalhelper.SetTestK8sAPIClient(client)
	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	return client
}

func TestSetServiceAccountAutomountServiceAccountTokenMissingAccount(t *testing.T) {
	withK8sClient(t)

	err := SetServiceAccountAutomountServiceAccountToken(parameters.TestAccessControlNameSpace, "missing", "true")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error getting service account")
}

func TestDefineAndCreateServiceOnClusterErrors(t *testing.T) {
	t.Run("existing service is not an error", func(t *testing.T) {
		withK8sClient(t, &corev1.Service{ObjectMeta: metav1.ObjectMeta{
			Name: "service-name", Namespace: parameters.TestAccessControlNameSpace,
		}})

		assert.NoError(t, DefineAndCreateServiceOnCluster("service-name", parameters.TestAccessControlNameSpace,
			8080, 8080, false, nil, ""))
	})

	t.Run("create failure is wrapped", func(t *testing.T) {
		client := withK8sClient(t)
		client.PrependReactor("create", "services", func(k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errInjected
		})

		err := DefineAndCreateServiceOnCluster("service-name", parameters.TestAccessControlNameSpace,
			8080, 8080, false, nil, "")
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to create service on cluster")
	})

	t.Run("node port redefinition failure is returned", func(t *testing.T) {
		oldRedefine := redefineWithNodePort

		t.Cleanup(func() { redefineWithNodePort = oldRedefine })

		redefineWithNodePort = func(*corev1.Service) error { return errInjected }

		err := DefineAndCreateServiceOnCluster("service-name", parameters.TestAccessControlNameSpace,
			8080, 8080, true, nil, "")
		require.ErrorIs(t, err, errInjected)
	})
}

func TestHasNodeRtKernel(t *testing.T) {
	rtNode := testNode("rt", nil)
	rtNode.Status.NodeInfo.KernelVersion = "5.14.0-284.el9.x86_64+rt"
	stdNode := testNode("std", nil)
	stdNode.Status.NodeInfo.KernelVersion = "5.14.0-284.el9.x86_64"

	withK8sClient(t, rtNode, stdNode)

	hasRt, err := HasNodeRtKernel("rt")
	require.NoError(t, err)
	assert.True(t, hasRt)

	hasRt, err = HasNodeRtKernel("std")
	require.NoError(t, err)
	assert.False(t, hasRt)

	_, err = HasNodeRtKernel("missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get node missing")
}

func TestGetWorkersMCPs(t *testing.T) {
	t.Run("maps pools to the worker nodes they manage", func(t *testing.T) {
		withEgiClient(t, nil,
			testNode("master-0", map[string]string{masterRoleLabel: ""}),
			testNode("worker-0", map[string]string{workerRoleLabel: ""}),
			testNode("worker-1", map[string]string{workerRoleLabel: "", "pool": "mixed"}),
			// The fake lists nodes by name, so this matched non-worker node is scanned before worker-1
			// and must not stop the scan.
			testNode("a-infra-0", map[string]string{"pool": "mixed"}),
			testMCP("master", 1, matchLabels(masterRoleLabel)),
			testMCP("worker", 2, matchLabels(workerRoleLabel)),
			testMCP("mixed", 2, &metav1.LabelSelector{MatchLabels: map[string]string{"pool": "mixed"}}),
			testMCP("empty", 0, matchLabels(workerRoleLabel)),
		)

		mcps, err := GetWorkersMCPs()
		require.NoError(t, err)
		require.Len(t, mcps, 2)
		assert.ElementsMatch(t, []string{"worker-0", "worker-1"}, mcps["worker"])
		assert.Equal(t, []string{"worker-1"}, mcps["mixed"])
	})

	t.Run("invalid node selector", func(t *testing.T) {
		withEgiClient(t, nil, testMCP("broken", 1, invalidSelector))

		_, err := GetWorkersMCPs()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get labelselector from mcp broken")
	})

	t.Run("node list failure", func(t *testing.T) {
		clientSet := withEgiClient(t, nil)
		k8sClient, ok := clientSet.K8sClient.(*k8sfake.Clientset)
		require.True(t, ok)
		k8sClient.PrependReactor("list", "nodes",
			func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errInjected })

		_, err := GetWorkersMCPs()
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to get nodes")
	})

	t.Run("pool list failure", func(t *testing.T) {
		withEgiClient(t, failingList())

		_, err := GetWorkersMCPs()
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to list MachineConfigPools")
	})
}

func testMC(name string, mcLabels map[string]string) *machineconfigv1.MachineConfig {
	return &machineconfigv1.MachineConfig{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: mcLabels}}
}

func TestGetMachineConfigTargetLabels(t *testing.T) {
	t.Run("returns labels of the first matching MachineConfig", func(t *testing.T) {
		withEgiClient(t, nil,
			testMCP(testMCPName, 1, matchLabels(workerRoleLabel)),
			testMC("00-master", map[string]string{"pool": "master"}),
			testMC("00-worker", map[string]string{"pool": testMCPName, "extra": "x"}),
		)

		mcLabels, err := GetMachineConfigTargetLabels(testMCPName)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"pool": testMCPName, "extra": "x"}, mcLabels)
	})

	t.Run("no MachineConfig matches", func(t *testing.T) {
		withEgiClient(t, nil, testMCP(testMCPName, 1, nil), testMC("00-master", map[string]string{"pool": "master"}))

		_, err := GetMachineConfigTargetLabels(testMCPName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no machineconfig matches the label selector")
	})

	t.Run("pool not found", func(t *testing.T) {
		withEgiClient(t, nil)

		_, err := GetMachineConfigTargetLabels(testMCPName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mcp worker not found")
	})

	t.Run("invalid MachineConfig selector", func(t *testing.T) {
		mcp := testMCP(testMCPName, 1, nil)
		mcp.Spec.MachineConfigSelector = invalidSelector
		withEgiClient(t, nil, mcp)

		_, err := GetMachineConfigTargetLabels(testMCPName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create label selector")
	})

	t.Run("MachineConfig list failure", func(t *testing.T) {
		withEgiClient(t, failingList(), testMCP(testMCPName, 1, nil))

		_, err := GetMachineConfigTargetLabels(testMCPName)
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed getting MachineConfigs")
	})

	t.Run("MachineConfigPool builder construction failure", func(t *testing.T) {
		withEgiClient(t, nil)

		oldNewMCPBuilder := newMCPBuilder

		t.Cleanup(func() { newMCPBuilder = oldNewMCPBuilder })

		newMCPBuilder = func(*egiClients.Settings, string) *egiMco.MCPBuilder { return nil }

		_, err := GetMachineConfigTargetLabels(testMCPName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get MachineConfigPool builder")
	})
}

func mcpCondition(condType machineconfigv1.MachineConfigPoolConditionType,
	status corev1.ConditionStatus) machineconfigv1.MachineConfigPoolCondition {
	return machineconfigv1.MachineConfigPoolCondition{Type: condType, Status: status}
}

var (
	updatingTrue  = mcpCondition(machineconfigv1.MachineConfigPoolUpdating, corev1.ConditionTrue)
	updatingFalse = mcpCondition(machineconfigv1.MachineConfigPoolUpdating, corev1.ConditionFalse)
	degradedTrue  = mcpCondition(machineconfigv1.MachineConfigPoolDegraded, corev1.ConditionTrue)
)

// mcpRollout makes the n-th Get of the MachineConfigPool (1-based) report the conditions returned by
// phase, simulating the pool starting and finishing an update.
func mcpRollout(phase func(call int) []machineconfigv1.MachineConfigPoolCondition) *interceptor.Funcs {
	calls := 0

	return &interceptor.Funcs{
		Get: func(ctx context.Context, client runtimeclient.WithWatch, key runtimeclient.ObjectKey,
			obj runtimeclient.Object, opts ...runtimeclient.GetOption) error {
			if err := client.Get(ctx, key, obj, opts...); err != nil {
				return err
			}

			if mcp, ok := obj.(*machineconfigv1.MachineConfigPool); ok {
				calls++
				mcp.Status.Conditions = phase(calls)
			}

			return nil
		},
	}
}

func startsThenFinishes(degraded bool) func(int) []machineconfigv1.MachineConfigPoolCondition {
	return func(call int) []machineconfigv1.MachineConfigPoolCondition {
		if call == 1 {
			return []machineconfigv1.MachineConfigPoolCondition{updatingTrue}
		}

		if degraded {
			return []machineconfigv1.MachineConfigPoolCondition{updatingFalse, degradedTrue}
		}

		return []machineconfigv1.MachineConfigPoolCondition{updatingFalse}
	}
}

func neverStarts(int) []machineconfigv1.MachineConfigPoolCondition {
	return []machineconfigv1.MachineConfigPoolCondition{updatingFalse}
}

func neverFinishes(int) []machineconfigv1.MachineConfigPoolCondition {
	return []machineconfigv1.MachineConfigPoolCondition{updatingTrue}
}

// rolloutFailureCases are shared by DeployRTKernelMachineConfig and RemoveRTKernelMachineConfig, which
// wait on the MachineConfigPool the same way.
var rolloutFailureCases = []struct {
	name    string
	phase   func(int) []machineconfigv1.MachineConfigPoolCondition
	wantErr string
}{
	{name: "pool never starts updating", phase: neverStarts, wantErr: "Updating=True"},
	{name: "pool never finishes updating", phase: neverFinishes, wantErr: "Updating=False"},
	{name: "pool degraded after update", phase: startsThenFinishes(true), wantErr: "Degraded"},
}

func TestDeployRTKernelMachineConfig(t *testing.T) {
	mcpLabels := map[string]string{"pool": testMCPName}

	t.Run("creates the realtime MachineConfig and waits for the rollout", func(t *testing.T) {
		clientSet := withEgiClient(t, mcpRollout(startsThenFinishes(false)), testMCP(testMCPName, 1, nil))

		require.NoError(t, DeployRTKernelMachineConfig(testMCPName, testMCName, mcpLabels, time.Second))

		created := &machineconfigv1.MachineConfig{}
		require.NoError(t, clientSet.Get(context.TODO(), runtimeclient.ObjectKey{Name: testMCName}, created))
		assert.Equal(t, "realtime", created.Spec.KernelType)
		assert.Equal(t, mcpLabels, created.Labels)
	})

	for _, testCase := range rolloutFailureCases {
		t.Run(testCase.name, func(t *testing.T) {
			withEgiClient(t, mcpRollout(testCase.phase), testMCP(testMCPName, 1, nil))

			err := DeployRTKernelMachineConfig(testMCPName, testMCName, mcpLabels, 10*time.Millisecond)
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.wantErr)
		})
	}

	t.Run("invalid MachineConfig label", func(t *testing.T) {
		withEgiClient(t, nil, testMCP(testMCPName, 1, nil))

		err := DeployRTKernelMachineConfig(testMCPName, testMCName, map[string]string{"": "x"}, time.Second)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create MachineConfig")
	})

	t.Run("MachineConfig builder construction failure", func(t *testing.T) {
		withEgiClient(t, nil)

		oldNewMCBuilder := newMCBuilder

		t.Cleanup(func() { newMCBuilder = oldNewMCBuilder })

		newMCBuilder = func(*egiClients.Settings, string) *egiMco.MCBuilder { return nil }

		err := DeployRTKernelMachineConfig(testMCPName, testMCName, mcpLabels, time.Second)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get MachineConfig builder")
	})

	t.Run("rollout MachineConfigPool builder construction failure", func(t *testing.T) {
		withEgiClient(t, nil)

		oldNewMCPBuilder := newMCPBuilder

		t.Cleanup(func() { newMCPBuilder = oldNewMCPBuilder })

		newMCPBuilder = func(*egiClients.Settings, string) *egiMco.MCPBuilder { return nil }

		err := DeployRTKernelMachineConfig(testMCPName, testMCName, mcpLabels, time.Second)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get MachineConfigPool builder")
	})
}

func TestRemoveRTKernelMachineConfig(t *testing.T) {
	t.Run("deletes the MachineConfig and waits for the rollout", func(t *testing.T) {
		clientSet := withEgiClient(t, mcpRollout(startsThenFinishes(false)),
			testMCP(testMCPName, 1, nil), testMC(testMCName, nil))

		require.NoError(t, RemoveRTKernelMachineConfig(testMCPName, testMCName, time.Second))

		err := clientSet.Get(context.TODO(), runtimeclient.ObjectKey{Name: testMCName}, &machineconfigv1.MachineConfig{})
		assert.True(t, k8serrors.IsNotFound(err), "MachineConfig should be deleted, got %v", err)
	})

	for _, testCase := range rolloutFailureCases {
		t.Run(testCase.name, func(t *testing.T) {
			withEgiClient(t, mcpRollout(testCase.phase), testMCP(testMCPName, 1, nil), testMC(testMCName, nil))

			err := RemoveRTKernelMachineConfig(testMCPName, testMCName, 10*time.Millisecond)
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.wantErr)
		})
	}

	t.Run("delete failure", func(t *testing.T) {
		withEgiClient(t, &interceptor.Funcs{
			Delete: func(context.Context, runtimeclient.WithWatch, runtimeclient.Object, ...runtimeclient.DeleteOption) error {
				return errInjected
			},
		}, testMCP(testMCPName, 1, nil), testMC(testMCName, nil))

		err := RemoveRTKernelMachineConfig(testMCPName, testMCName, time.Second)
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to delete machineconfig "+testMCName)
	})
}

func TestIsMCPConditionUpdatingTrue(t *testing.T) {
	testCases := []struct {
		name       string
		conditions []mcv1.MachineConfigPoolCondition
		want       bool
		wantErr    string
	}{
		{name: "no conditions", wantErr: "conditions are empty"},
		{
			name: "updating true",
			conditions: []mcv1.MachineConfigPoolCondition{
				{Type: mcv1.MachineConfigPoolUpdated, Status: corev1.ConditionFalse},
				{Type: mcv1.MachineConfigPoolUpdating, Status: corev1.ConditionTrue},
			},
			want: true,
		},
		{
			name:       "updating false",
			conditions: []mcv1.MachineConfigPoolCondition{{Type: mcv1.MachineConfigPoolUpdating, Status: corev1.ConditionFalse}},
		},
		{
			name:       "updating condition missing",
			conditions: []mcv1.MachineConfigPoolCondition{{Type: mcv1.MachineConfigPoolDegraded, Status: corev1.ConditionFalse}},
			wantErr:    "not found yet",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			updating, err := IsMCPConditionUpdatingTrue(testCase.conditions)
			if testCase.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.want, updating)
		})
	}
}

func TestMCOHelpersRequireEcoGoinfraClient(t *testing.T) {
	calls := map[string]func() error{
		"GetWorkersMCPs": func() error {
			_, err := GetWorkersMCPs()

			return err
		},
		"GetMachineConfigTargetLabels": func() error {
			_, err := GetMachineConfigTargetLabels(testMCPName)

			return err
		},
		"DeployRTKernelMachineConfig": func() error {
			return DeployRTKernelMachineConfig(testMCPName, testMCName, nil, time.Second)
		},
		"RemoveRTKernelMachineConfig": func() error {
			return RemoveRTKernelMachineConfig(testMCPName, testMCName, time.Second)
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			withoutEgiClient(t)

			err := call()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "eco-goinfra client is not initialized")
		})
	}
}

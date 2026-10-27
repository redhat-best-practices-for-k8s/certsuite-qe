package helper

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"
	clientmcv1 "github.com/openshift/client-go/machineconfiguration/clientset/versioned/typed/machineconfiguration/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	testclient "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/client"
)

var errInjected = errors.New("injected failure")

// fakeMachineConfigs and fakeMachineConfigPools implement only the methods the helpers use; any
// other method panics on the nil embedded interface.
type fakeMachineConfigs struct {
	clientmcv1.MachineConfigInterface

	list    machineconfigv1.MachineConfigList
	listErr error

	object     *machineconfigv1.MachineConfig
	getErr     error
	updateErr  error
	updateFunc func(*machineconfigv1.MachineConfig) (*machineconfigv1.MachineConfig, error)
}

func (f *fakeMachineConfigs) List(context.Context, metav1.ListOptions) (*machineconfigv1.MachineConfigList, error) {
	return &f.list, f.listErr
}

func (f *fakeMachineConfigs) Get(_ context.Context, name string, _ metav1.GetOptions) (*machineconfigv1.MachineConfig, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}

	if f.object == nil {
		f.object = &machineconfigv1.MachineConfig{ObjectMeta: metav1.ObjectMeta{Name: name}}
	}

	return f.object, nil
}

func (f *fakeMachineConfigs) Update(_ context.Context, machineConfig *machineconfigv1.MachineConfig,
	_ metav1.UpdateOptions) (*machineconfigv1.MachineConfig, error) {
	if f.updateFunc != nil {
		return f.updateFunc(machineConfig)
	}

	if f.updateErr != nil {
		return nil, f.updateErr
	}

	f.object = machineConfig

	return machineConfig, nil
}

type fakeMachineConfigPools struct {
	clientmcv1.MachineConfigPoolInterface

	list    machineconfigv1.MachineConfigPoolList
	listErr error

	getFunc func() (*machineconfigv1.MachineConfigPool, error)
}

func (f *fakeMachineConfigPools) List(context.Context, metav1.ListOptions) (*machineconfigv1.MachineConfigPoolList, error) {
	return &f.list, f.listErr
}

func (f *fakeMachineConfigPools) Get(context.Context, string, metav1.GetOptions) (*machineconfigv1.MachineConfigPool, error) {
	if f.getFunc != nil {
		return f.getFunc()
	}

	return &machineconfigv1.MachineConfigPool{}, nil
}

// fakeMachineConfigClient implements only the MachineConfiguration getters the helpers use;
// any other method panics on the nil embedded interface.
type fakeMachineConfigClient struct {
	clientmcv1.MachineconfigurationV1Interface

	mcs  *fakeMachineConfigs
	mcps *fakeMachineConfigPools
}

//nolint:ireturn // Must match the MachineconfigurationV1Interface.
func (f *fakeMachineConfigClient) MachineConfigs() clientmcv1.MachineConfigInterface {
	return f.mcs
}

//nolint:ireturn // Must match the MachineconfigurationV1Interface.
func (f *fakeMachineConfigClient) MachineConfigPools() clientmcv1.MachineConfigPoolInterface {
	return f.mcps
}

// withFakeCluster installs fake k8s and MachineConfiguration clients, restoring everything on
// cleanup. Tests using it must not run in parallel.
func withFakeCluster(t *testing.T, nodes ...runtime.Object) (*k8sfake.Clientset, *fakeMachineConfigClient) {
	t.Helper()

	client := k8sfake.NewClientset(nodes...)
	globalhelper.SetTestK8sAPIClient(client)

	machineConfig := &fakeMachineConfigClient{
		mcs:  &fakeMachineConfigs{},
		mcps: &fakeMachineConfigPools{},
	}
	globalhelper.GetAPIClient().MachineconfigurationV1Interface = machineConfig

	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	return client, machineConfig
}

// stubExec replaces the exec seam for the duration of the test.
func stubExec(t *testing.T, exec func(corev1.Pod, []string) (bytes.Buffer, error)) {
	t.Helper()

	old := execPodCommand

	t.Cleanup(func() { execPodCommand = old })

	execPodCommand = exec
}

// stubHugePagesExec answers the hugepages echo command with success and the cat command
// with the given output or error, recording every issued shell command.
func stubHugePagesExec(t *testing.T, catOutput string, catErr error) *[]string {
	t.Helper()

	var commands []string

	stubExec(t, func(_ corev1.Pod, cmd []string) (bytes.Buffer, error) {
		commands = append(commands, strings.Join(cmd, " "))

		if strings.HasPrefix(cmd[2], "echo") {
			return bytes.Buffer{}, nil
		}

		return *bytes.NewBufferString(catOutput), catErr
	})

	return &commands
}

// withFastHugePagesPoll shortens the hugepages poll timers; the defaults (5s interval, 5m
// timeout) would make the failure cases crawl. Successful polls return on the first check.
func withFastHugePagesPoll(t *testing.T) {
	t.Helper()

	oldInterval, oldTimeout := hugePagesRetryInterval, hugePagesWaitTimeout

	t.Cleanup(func() { hugePagesRetryInterval, hugePagesWaitTimeout = oldInterval, oldTimeout })

	hugePagesRetryInterval = time.Millisecond
	hugePagesWaitTimeout = 50 * time.Millisecond
}

func testNode(name string, labels map[string]string, ready corev1.ConditionStatus) *corev1.Node {
	conditions := []corev1.NodeCondition{{Type: corev1.NodeReady, Status: ready}}

	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels},
		Status:     corev1.NodeStatus{Conditions: conditions},
	}
}

func testMachineConfig(kernelArgs []string) *machineconfigv1.MachineConfig {
	return &machineconfigv1.MachineConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "99-worker"},
		Spec:       machineconfigv1.MachineConfigSpec{KernelArguments: kernelArgs},
	}
}

func mcList(configs ...*machineconfigv1.MachineConfig) machineconfigv1.MachineConfigList {
	list := machineconfigv1.MachineConfigList{}

	for _, config := range configs {
		list.Items = append(list.Items, *config)
	}

	return list
}

func mcpList(pools ...string) machineconfigv1.MachineConfigPoolList {
	list := machineconfigv1.MachineConfigPoolList{}

	for _, name := range pools {
		list.Items = append(list.Items, machineconfigv1.MachineConfigPool{
			ObjectMeta: metav1.ObjectMeta{Name: name},
		})
	}

	return list
}

// failNodeList makes the fake clientset fail node listing with the injected error.
func failNodeList(client *k8sfake.Clientset) {
	client.PrependReactor("list", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errInjected
	})
}

func TestIsRealOCPCluster(t *testing.T) {
	testCases := []struct {
		name     string
		pools    machineconfigv1.MachineConfigPoolList
		poolErr  error
		nodes    []runtime.Object
		nodeErr  bool
		wantReal bool
		wantMsg  string
	}{
		{
			name:     "worker-cnf pool marks a real cluster",
			pools:    mcpList("master", "worker-cnf"),
			wantReal: true,
			wantMsg:  "found worker-cnf MachineConfigPool",
		},
		{
			name:    "MachineConfigPool list failure",
			poolErr: errInjected,
			wantMsg: "cannot access MachineConfigPools",
		},
		{
			name: "multi-node cluster",
			nodes: []runtime.Object{
				testNode("master-0", nil, corev1.ConditionTrue),
				testNode("worker-0", nil, corev1.ConditionTrue),
			},
			wantReal: true,
			wantMsg:  "cluster has 2 nodes",
		},
		{
			name:    "node list failure",
			nodeErr: true,
			wantMsg: "cannot list nodes",
		},
		{
			name: "multi-node cluster with a CRC-named node",
			nodes: []runtime.Object{
				testNode("crc-4xk2v", nil, corev1.ConditionTrue),
				testNode("worker-0", nil, corev1.ConditionTrue),
			},
			wantReal: true,
			wantMsg:  "cluster has 2 nodes",
		},
		{
			name:    "single CRC node",
			nodes:   []runtime.Object{testNode("crc-4xk2v", nil, corev1.ConditionTrue)},
			wantMsg: "appears to be a CRC development cluster",
		},
		{
			name: "SNO with worker-cnf role",
			nodes: []runtime.Object{testNode("sno-0",
				map[string]string{"node-role.kubernetes.io/worker-cnf": ""}, corev1.ConditionTrue)},
			wantReal: true,
			wantMsg:  "SNO cluster with worker-cnf role",
		},
		{
			name: "SNO with plain worker role",
			nodes: []runtime.Object{testNode("sno-0",
				map[string]string{"node-role.kubernetes.io/worker": ""}, corev1.ConditionTrue)},
			wantMsg: "appears to be a single-node development cluster",
		},
		{
			name:    "SNO without roles",
			nodes:   []runtime.Object{testNode("sno-0", nil, corev1.ConditionTrue)},
			wantMsg: "appears to be a CRC/SNO development cluster",
		},
		{
			name:    "empty cluster",
			wantMsg: "appears to be a CRC/SNO development cluster",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client, machineConfig := withFakeCluster(t, testCase.nodes...)
			machineConfig.mcps.list = testCase.pools
			machineConfig.mcps.listErr = testCase.poolErr

			if testCase.nodeErr {
				failNodeList(client)
			}

			isReal, reason := IsRealOCPCluster()
			assert.Equal(t, testCase.wantReal, isReal)
			assert.Contains(t, reason, testCase.wantMsg)
		})
	}
}

func TestDetectBootParamsAlterations(t *testing.T) {
	testCases := []struct {
		name      string
		configs   machineconfigv1.MachineConfigList
		listErr   error
		wantFound bool
		wantMsg   string
	}{
		{
			name:      "matching kernel argument",
			configs:   mcList(testMachineConfig([]string{"isolcpus=2-3"})),
			wantFound: true,
			wantMsg:   "MachineConfig 99-worker has kernel arg: isolcpus=2-3",
		},
		{
			name:    "non-matching kernel argument",
			configs: mcList(testMachineConfig([]string{"quiet splash"})),
			wantMsg: "no kernel arg alterations detected",
		},
		{
			name:    "no MachineConfigs",
			wantMsg: "no kernel arg alterations detected",
		},
		{
			name:    "MachineConfig list failure",
			listErr: errInjected,
			wantMsg: "cannot access MachineConfigs",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, machineConfig := withFakeCluster(t)
			machineConfig.mcs.list = testCase.configs
			machineConfig.mcs.listErr = testCase.listErr

			found, details := DetectBootParamsAlterations()
			assert.Equal(t, testCase.wantFound, found)
			assert.Contains(t, details, testCase.wantMsg)
		})
	}
}

func TestHasMachineConfigKernelArguments(t *testing.T) {
	t.Run("MachineConfig with kernel arguments", func(t *testing.T) {
		_, machineConfig := withFakeCluster(t)
		machineConfig.mcs.list = mcList(testMachineConfig([]string{"isolcpus=2-3"}))

		hasArgs, details, err := HasMachineConfigKernelArguments()
		require.NoError(t, err)
		assert.True(t, hasArgs)
		assert.Contains(t, details, "MachineConfig 99-worker has kernelArguments: [isolcpus=2-3]")
	})

	t.Run("no kernel arguments", func(t *testing.T) {
		_, machineConfig := withFakeCluster(t)
		machineConfig.mcs.list = mcList(testMachineConfig(nil))

		hasArgs, details, err := HasMachineConfigKernelArguments()
		require.NoError(t, err)
		assert.False(t, hasArgs)
		assert.Contains(t, details, "No MachineConfigs with custom kernelArguments found")
	})

	t.Run("list failure", func(t *testing.T) {
		_, machineConfig := withFakeCluster(t)
		machineConfig.mcs.listErr = errInjected

		_, _, err := HasMachineConfigKernelArguments()
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to list MachineConfigs")
	})
}

func TestWaitForSpecificNodeCondition(t *testing.T) {
	testCases := []struct {
		name     string
		nodes    []runtime.Object
		nodeName string
		ready    bool
		wantErr  bool
	}{
		{
			name:     "node is ready",
			nodes:    []runtime.Object{testNode("worker-0", nil, corev1.ConditionTrue)},
			nodeName: "worker-0",
			ready:    true,
		},
		{
			name:     "node is not ready",
			nodes:    []runtime.Object{testNode("worker-1", nil, corev1.ConditionFalse)},
			nodeName: "worker-1",
			ready:    false,
		},
		{
			name:     "times out waiting for ready",
			nodes:    []runtime.Object{testNode("worker-1", nil, corev1.ConditionFalse)},
			nodeName: "worker-1",
			ready:    true,
			wantErr:  true,
		},
		{
			name:     "times out when the node is missing",
			nodes:    []runtime.Object{testNode("worker-0", nil, corev1.ConditionTrue)},
			nodeName: "gone-node",
			ready:    true,
			wantErr:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			clients := &testclient.ClientSet{CoreV1Interface: k8sfake.NewClientset(testCase.nodes...).CoreV1()}

			err := WaitForSpecificNodeCondition(clients, 50*time.Millisecond, time.Millisecond,
				testCase.nodeName, testCase.ready)

			if testCase.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}

	t.Run("list failure", func(t *testing.T) {
		fakeClient := k8sfake.NewClientset()
		failNodeList(fakeClient)

		clients := &testclient.ClientSet{CoreV1Interface: fakeClient.CoreV1()}

		err := WaitForSpecificNodeCondition(clients, 50*time.Millisecond, time.Millisecond, "worker-0", true)
		require.ErrorIs(t, err, errInjected)
	})
}

func TestRestoreMachineConfigKernelArguments(t *testing.T) {
	_, machineConfigClient := withFakeCluster(t)
	machineConfigClient.mcs.object = &machineconfigv1.MachineConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "99-worker-cnf-generated"},
		Spec: machineconfigv1.MachineConfigSpec{
			KernelArguments: []string{"skew_tick=1", "nohz=off"},
		},
	}

	require.NoError(t, RestoreMachineConfigKernelArguments("99-worker-cnf-generated", []string{"default_hugepagesz=1G"}))
	assert.Equal(t, []string{"default_hugepagesz=1G"}, machineConfigClient.mcs.object.Spec.KernelArguments)

	// The restored MachineConfig must not alias the caller's slice.
	kernelArguments := []string{"nohz=off"}
	require.NoError(t, RestoreMachineConfigKernelArguments("99-worker-cnf-generated", kernelArguments))
	kernelArguments[0] = "mutated"

	assert.Equal(t, []string{"nohz=off"}, machineConfigClient.mcs.object.Spec.KernelArguments)
}

func TestRestoreMachineConfigKernelArgumentsErrors(t *testing.T) {
	t.Run("get fails", func(t *testing.T) {
		_, machineConfigClient := withFakeCluster(t)
		machineConfigClient.mcs.getErr = errInjected

		err := RestoreMachineConfigKernelArguments("mc", nil)
		require.ErrorIs(t, err, errInjected)
	})

	t.Run("update fails", func(t *testing.T) {
		_, machineConfigClient := withFakeCluster(t)
		machineConfigClient.mcs.updateErr = errInjected

		err := RestoreMachineConfigKernelArguments("mc", nil)
		require.ErrorIs(t, err, errInjected)
	})
}

func TestRestoreMachineConfigKernelArgumentsRetriesOnConflict(t *testing.T) {
	_, machineConfigClient := withFakeCluster(t)
	machineConfigClient.mcs.object = &machineconfigv1.MachineConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "99-worker-cnf-generated"},
		Spec:       machineconfigv1.MachineConfigSpec{KernelArguments: []string{"skew_tick=1"}},
	}

	attempts := 0
	machineConfigClient.mcs.updateFunc = func(
		machineConfig *machineconfigv1.MachineConfig) (*machineconfigv1.MachineConfig, error) {
		attempts++
		if attempts == 1 {
			return nil, apierrors.NewConflict(
				schema.GroupResource{Group: "machineconfiguration.openshift.io", Resource: "machineconfigs"},
				machineConfig.Name, errInjected)
		}

		machineConfigClient.mcs.object = machineConfig

		return machineConfig, nil
	}

	require.NoError(t, RestoreMachineConfigKernelArguments("99-worker-cnf-generated", []string{"nohz=off"}))
	assert.Equal(t, 2, attempts)
	assert.Equal(t, []string{"nohz=off"}, machineConfigClient.mcs.object.Spec.KernelArguments)
}

func updatedMachineConfigPool(conditions ...machineconfigv1.MachineConfigPoolCondition) *machineconfigv1.MachineConfigPool {
	return &machineconfigv1.MachineConfigPool{
		ObjectMeta: metav1.ObjectMeta{Name: "worker-cnf", Generation: 2},
		Status: machineconfigv1.MachineConfigPoolStatus{
			ObservedGeneration: 2,
			Conditions:         conditions,
		},
	}
}

func TestWaitForMachineConfigPoolRollout(t *testing.T) {
	interval := machineConfigPoolRolloutRetryInterval
	machineConfigPoolRolloutRetryInterval = 10 * time.Millisecond

	t.Cleanup(func() { machineConfigPoolRolloutRetryInterval = interval })

	t.Run("returns when the pool finished rolling out", func(t *testing.T) {
		_, machineConfigClient := withFakeCluster(t)
		machineConfigClient.mcps.getFunc = func() (*machineconfigv1.MachineConfigPool, error) {
			return updatedMachineConfigPool(machineconfigv1.MachineConfigPoolCondition{
				Type:   machineconfigv1.MachineConfigPoolUpdated,
				Status: corev1.ConditionTrue,
			}), nil
		}

		require.NoError(t, WaitForMachineConfigPoolRollout("worker-cnf", time.Second))
	})

	t.Run("returns once the pool becomes updated", func(t *testing.T) {
		_, machineConfigClient := withFakeCluster(t)
		calls := 0
		machineConfigClient.mcps.getFunc = func() (*machineconfigv1.MachineConfigPool, error) {
			calls++
			if calls < 3 {
				return updatedMachineConfigPool(), nil
			}

			return updatedMachineConfigPool(machineconfigv1.MachineConfigPoolCondition{
				Type:   machineconfigv1.MachineConfigPoolUpdated,
				Status: corev1.ConditionTrue,
			}), nil
		}

		require.NoError(t, WaitForMachineConfigPoolRollout("worker-cnf", time.Second))
		assert.Equal(t, 3, calls)
	})

	t.Run("fails when the pool is degraded", func(t *testing.T) {
		_, machineConfigClient := withFakeCluster(t)
		machineConfigClient.mcps.getFunc = func() (*machineconfigv1.MachineConfigPool, error) {
			return updatedMachineConfigPool(
				machineconfigv1.MachineConfigPoolCondition{
					Type:   machineconfigv1.MachineConfigPoolUpdated,
					Status: corev1.ConditionTrue,
				},
				machineconfigv1.MachineConfigPoolCondition{
					Type:   machineconfigv1.MachineConfigPoolDegraded,
					Status: corev1.ConditionTrue,
				},
			), nil
		}

		require.Error(t, WaitForMachineConfigPoolRollout("worker-cnf", 100*time.Millisecond))
	})

	t.Run("fails when the rollout does not finish", func(t *testing.T) {
		_, machineConfigClient := withFakeCluster(t)
		machineConfigClient.mcps.getFunc = func() (*machineconfigv1.MachineConfigPool, error) {
			pool := updatedMachineConfigPool(machineconfigv1.MachineConfigPoolCondition{
				Type:   machineconfigv1.MachineConfigPoolUpdated,
				Status: corev1.ConditionTrue,
			})
			pool.Generation = 3

			return pool, nil
		}

		require.Error(t, WaitForMachineConfigPoolRollout("worker-cnf", 100*time.Millisecond))
	})

	t.Run("fails when the pool cannot be read", func(t *testing.T) {
		_, machineConfigClient := withFakeCluster(t)
		machineConfigClient.mcps.getFunc = func() (*machineconfigv1.MachineConfigPool, error) {
			return nil, errInjected
		}

		require.ErrorIs(t, WaitForMachineConfigPoolRollout("worker-cnf", time.Second), errInjected)
	})
}

func TestUpdateAndVerifyHugePagesConfig(t *testing.T) {
	pod := &corev1.Pod{}

	t.Run("updates and verifies the new value", func(t *testing.T) {
		withFastHugePagesPoll(t)

		commands := stubHugePagesExec(t, "20\r\n", nil)

		require.NoError(t, UpdateAndVerifyHugePagesConfig(20, "/sys/nr_hugepages", pod))
		assert.Contains(t, *commands, "/bin/bash -c echo 20 > /sys/nr_hugepages")
		assert.Contains(t, *commands, "/bin/bash -c cat /sys/nr_hugepages")
	})

	t.Run("echo failure is wrapped", func(t *testing.T) {
		stubExec(t, func(corev1.Pod, []string) (bytes.Buffer, error) {
			return bytes.Buffer{}, errInjected
		})

		err := UpdateAndVerifyHugePagesConfig(20, "/sys/nr_hugepages", pod)
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to execute command")
	})

	t.Run("times out when the value never matches", func(t *testing.T) {
		withFastHugePagesPoll(t)

		stubHugePagesExec(t, "10\r\n", nil)

		require.Error(t, UpdateAndVerifyHugePagesConfig(20, "/sys/nr_hugepages", pod))
	})

	t.Run("read failure during the poll is wrapped", func(t *testing.T) {
		withFastHugePagesPoll(t)

		stubHugePagesExec(t, "", errInjected)

		err := UpdateAndVerifyHugePagesConfig(20, "/sys/nr_hugepages", pod)
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to get hugepages number")
	})
}

func TestGetHugePagesConfigNumber(t *testing.T) {
	testCases := []struct {
		name    string
		output  string
		execErr error
		wantNum int
		wantErr bool
	}{
		{name: "parses the number", output: "2048\r\n", wantNum: 2048},
		{name: "non-numeric output", output: "not-a-number\r\n", wantNum: -1, wantErr: true},
		{name: "exec failure", execErr: errInjected, wantNum: -1, wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stubExec(t, func(corev1.Pod, []string) (bytes.Buffer, error) {
				return *bytes.NewBufferString(testCase.output), testCase.execErr
			})

			number, err := GetHugePagesConfigNumber("/sys/nr_hugepages", &corev1.Pod{})
			if testCase.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, testCase.wantNum, number)
		})
	}
}

func TestArgListToMap(t *testing.T) {
	assert.Empty(t, ArgListToMap(nil))
	assert.Equal(t, map[string]string{"isolcpus": "2-3", "nohz": ""},
		ArgListToMap([]string{"isolcpus=2-3", "nohz"}))
	assert.Equal(t, map[string]string{"cpu_affinity": "1 2"},
		ArgListToMap([]string{"cpu_affinity=1 2=3"}))
}

func TestAppendIstioContainerToPod(t *testing.T) {
	pod := &corev1.Pod{Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "app"}}}}

	AppendIstioContainerToPod(pod, "istio/proxyv2:1.31")

	require.Len(t, pod.Spec.Containers, 2)
	assert.Equal(t, "app", pod.Spec.Containers[0].Name)
	assert.Equal(t, corev1.Container{
		Name:    "istio-proxy",
		Image:   "istio/proxyv2:1.31",
		Command: []string{"/bin/bash", "-c", "sleep INF"},
	}, pod.Spec.Containers[1])
}

func assertNonUBIContainer(t *testing.T, template corev1.PodTemplateSpec) {
	t.Helper()

	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, template.Labels)
	require.Len(t, template.Spec.Containers, 1)

	container := template.Spec.Containers[0]
	assert.Equal(t, "test", container.Name)
	assert.Equal(t, tsparams.NotRedHatRelease, container.Image)
	assert.Equal(t, []string{"/bin/sh", "-c", "sleep INF"}, container.Command)
}

func TestDefineDeploymentWithNonUBIContainer(t *testing.T) {
	dep := DefineDeploymentWithNonUBIContainer("test-ns")

	assert.Equal(t, tsparams.TestDeploymentName, dep.Name)
	assert.Equal(t, "test-ns", dep.Namespace)
	assertNonUBIContainer(t, dep.Spec.Template)
}

func TestDefineStatefulSetWithNonUBIContainer(t *testing.T) {
	sts := DefineStatefulSetWithNonUBIContainer("test-ns")

	assert.Equal(t, tsparams.TestStatefulSetName, sts.Name)
	assert.Equal(t, "test-ns", sts.Namespace)
	assertNonUBIContainer(t, sts.Spec.Template)
}

func TestDetectBootParamsAlteration(t *testing.T) {
	testCases := []struct {
		name       string
		cmdline    string
		wantFound  bool
		wantParams []string
	}{
		{
			name:       "isolated CPUs",
			cmdline:    "root=/dev/sda1 isolcpus=2,3 quiet",
			wantFound:  true,
			wantParams: []string{"isolcpus=2,3"},
		},
		{
			name:       "multiple tuning parameters",
			cmdline:    "nohz_full=1 rcu_nocbs=2-3 hugepages=1G",
			wantFound:  true,
			wantParams: []string{"nohz_full=1", "rcu_nocbs=2-3", "hugepages=1G"},
		},
		{
			name:       "parameter without value",
			cmdline:    "quiet nosoftlockup",
			wantFound:  true,
			wantParams: []string{"nosoftlockup"},
		},
		{
			name:      "default kernel cmdline",
			cmdline:   "root=/dev/sda1 quiet rhgb",
			wantFound: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			found, details := DetectBootParamsAlteration(testCase.cmdline)
			assert.Equal(t, testCase.wantFound, found)

			for _, param := range testCase.wantParams {
				assert.Contains(t, details, param)
			}

			if !testCase.wantFound {
				assert.Contains(t, details, "No custom boot parameters detected")
			}
		})
	}
}

func TestIsNodeControlPlane(t *testing.T) {
	testCases := []struct {
		name     string
		nodeName string
		want     bool
	}{
		{name: "master", nodeName: "master-0", want: true},
		{name: "control-plane", nodeName: "control-plane-1", want: true},
		{name: "ctlplane", nodeName: "ctlplane-2", want: true},
		{name: "controlplane uppercase", nodeName: "CONTROLPLANE-0", want: true},
		{name: "worker", nodeName: "worker-0", want: false},
		{name: "empty", nodeName: "", want: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, IsNodeControlPlane(testCase.nodeName))
		})
	}
}

func TestCheckPodsOnControlPlaneNodes(t *testing.T) {
	testPod := func(name, nodeName string) corev1.Pod {
		return corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Spec:       corev1.PodSpec{NodeName: nodeName},
		}
	}

	t.Run("pods on control plane nodes are reported", func(t *testing.T) {
		found, details := CheckPodsOnControlPlaneNodes([]corev1.Pod{
			testPod("worker-pod", "worker-0"), testPod("master-pod", "master-0"),
		})
		assert.True(t, found)
		assert.Contains(t, details, "master-pod (node: master-0)")
		assert.NotContains(t, details, "worker-pod")
	})

	t.Run("no pods on control plane nodes", func(t *testing.T) {
		found, details := CheckPodsOnControlPlaneNodes([]corev1.Pod{testPod("worker-pod", "worker-0")})
		assert.False(t, found)
		assert.Contains(t, details, "No pods running on control plane nodes")
	})
}

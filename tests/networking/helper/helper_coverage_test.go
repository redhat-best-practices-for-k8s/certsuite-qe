package helper

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nad"
)

const (
	testNamespace    = tsparams.NetworkingNamespace
	multusAnnotation = "k8s.v1.cni.cncf.io/networks"
	cnfNodeLabel     = "node-role.kubernetes.io/worker-cnf"
	workerNodeLabel  = "node-role.kubernetes.io/worker"
)

var errInjected = errors.New("injected failure")

var nadScheme = func() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = nadv1.AddToScheme(scheme)

	return scheme
}()

// withFakeCluster wires a fake clientset into globalhelper, stubs every package seam so no
// test can reach a real cluster, and restores everything on cleanup.
func withFakeCluster(t *testing.T, objects ...runtime.Object) *k8sfake.Clientset {
	t.Helper()

	oldDeployment, oldDaemonSet, oldPolicy := createAndWaitDeployment, createAndWaitDaemonSet, createNetworkPolicy
	oldExec, oldConfig, oldCert := execPodCommand, getConfiguration, generateTLSCert

	t.Cleanup(func() {
		createAndWaitDeployment, createAndWaitDaemonSet, createNetworkPolicy = oldDeployment, oldDaemonSet, oldPolicy
		execPodCommand, getConfiguration, generateTLSCert = oldExec, oldConfig, oldCert

		globalhelper.UnsetTestK8sAPIClient()
	})

	client := k8sfake.NewClientset(objects...)
	globalhelper.SetTestK8sAPIClient(client)

	createAndWaitDeployment = func(*appsv1.Deployment, time.Duration) error { return nil }
	createAndWaitDaemonSet = func(*appsv1.DaemonSet, time.Duration) error { return nil }
	createNetworkPolicy = func(*networkingv1.NetworkPolicy, time.Duration) error { return nil }
	execPodCommand = func(corev1.Pod, []string) (bytes.Buffer, error) { return bytes.Buffer{}, nil }
	// Real RSA key generation is slow and already covered by TestGenerateSelfSignedCert.
	generateTLSCert = func() ([]byte, []byte, error) { return []byte("cert"), []byte("key"), nil }

	testConfig := &config.Config{}
	testConfig.General.CnfNodeLabel = cnfNodeLabel
	testConfig.General.WorkerNodeLabel = workerNodeLabel
	getConfiguration = func() *config.Config { return testConfig }

	return client
}

// setNadClient installs a controller-runtime fake client that knows NetworkAttachmentDefinitions.
func setNadClient(t *testing.T, funcs *interceptor.Funcs, objects ...runtimeclient.Object) {
	t.Helper()

	builder := fake.NewClientBuilder().WithScheme(nadScheme).WithObjects(objects...)
	if funcs != nil {
		builder = builder.WithInterceptorFuncs(*funcs)
	}

	globalhelper.GetAPIClient().Client = builder.Build()
}

func failOn(client *k8sfake.Clientset, verb, resource string) {
	client.PrependReactor(verb, resource, func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errInjected
	})
}

func captureDeployment(t *testing.T, err error) **appsv1.Deployment {
	t.Helper()

	var captured *appsv1.Deployment

	createAndWaitDeployment = func(dep *appsv1.Deployment, timeout time.Duration) error {
		assert.Equal(t, tsparams.WaitingTime, timeout)
		captured = dep

		return err
	}

	return &captured
}

func TestDeploymentWrappers(t *testing.T) {
	testCases := []struct {
		name           string
		create         func() error
		wantName       string
		wantReplicas   int32
		wantPrivileged bool
		wantLabel      map[string]string
		wantMultus     bool
	}{
		{
			name:     "default deployment",
			create:   func() error { return DefineAndCreateDeploymentOnCluster(2, testNamespace) },
			wantName: tsparams.TestDeploymentAName, wantReplicas: 2,
		},
		{
			name: "multus deployment",
			create: func() error {
				return DefineAndCreateDeploymentWithMultusOnCluster("multus", testNamespace, []string{"nad-a", "nad-b"}, 3)
			},
			wantName: "multus", wantReplicas: 3, wantPrivileged: true, wantMultus: true,
		},
		{
			name: "multus deployment with skip label",
			create: func() error {
				return DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster("multus-skip", testNamespace, []string{"nad-a"}, 1)
			},
			wantName: "multus-skip", wantReplicas: 1, wantMultus: true, wantLabel: tsparams.NetworkingTestMultusSkipLabel,
		},
		{
			name:     "privileged deployment",
			create:   func() error { return DefineAndCreatePrivilegedDeploymentOnCluster(1, testNamespace) },
			wantName: tsparams.TestDeploymentAName, wantReplicas: 1, wantPrivileged: true,
		},
		{
			name:     "deployment with skipped label",
			create:   func() error { return DefineAndCreateDeploymentWithSkippedLabelOnCluster(1, testNamespace) },
			wantName: tsparams.TestDeploymentAName, wantReplicas: 1, wantPrivileged: true, wantLabel: tsparams.NetworkingTestSkipLabel,
		},
		{
			name:     "deployment B in namespace",
			create:   func() error { return DefineAndCreateDeploymentWithNamespace(testNamespace, 2) },
			wantName: tsparams.TestDeploymentBName, wantReplicas: 2,
		},
		{
			name:     "named deployment",
			create:   func() error { return DefineAndCreateDeployment("named", testNamespace, 4) },
			wantName: "named", wantReplicas: 4,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			withFakeCluster(t)
			captured := captureDeployment(t, nil)

			require.NoError(t, testCase.create())

			dep := *captured
			require.NotNil(t, dep)
			assert.Equal(t, testCase.wantName, dep.Name)
			assert.Equal(t, testNamespace, dep.Namespace)
			assert.Equal(t, testCase.wantReplicas, *dep.Spec.Replicas)
			assert.Equal(t, sampleWorkloadImage, dep.Spec.Template.Spec.Containers[0].Image)

			securityContext := dep.Spec.Template.Spec.Containers[0].SecurityContext
			if testCase.wantPrivileged {
				require.NotNil(t, securityContext)
				assert.True(t, *securityContext.Privileged)
			} else {
				assert.Nil(t, securityContext)
			}

			for _, skipLabel := range []map[string]string{tsparams.NetworkingTestSkipLabel, tsparams.NetworkingTestMultusSkipLabel} {
				for key := range skipLabel {
					_, hasLabel := dep.Spec.Template.Labels[key]
					_, wanted := testCase.wantLabel[key]
					assert.Equal(t, wanted, hasLabel, key)
				}
			}

			_, hasMultus := dep.Spec.Template.Annotations[multusAnnotation]
			assert.Equal(t, testCase.wantMultus, hasMultus)

			// Every wrapper must propagate the create error.
			captureDeployment(t, errInjected)
			assert.ErrorIs(t, testCase.create(), errInjected)
		})
	}

	t.Run("skipped label wrapper adds context to the error", func(t *testing.T) {
		withFakeCluster(t)
		captureDeployment(t, errInjected)

		err := DefineAndCreateDeploymentWithSkippedLabelOnCluster(1, testNamespace)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create deployment")
	})
}

func TestDefineAndCreateDeploymentWithContainerPorts(t *testing.T) {
	withFakeCluster(t)
	captured := captureDeployment(t, nil)

	ports := []corev1.ContainerPort{{ContainerPort: 8080}, {ContainerPort: 9090}}
	require.NoError(t, DefineAndCreateDeploymentWithContainerPorts(2, ports, testNamespace))

	containers := (*captured).Spec.Template.Spec.Containers
	require.Len(t, containers, 2)

	for index, container := range containers {
		assert.Equal(t, []corev1.ContainerPort{ports[index]}, container.Ports)
	}

	*captured = nil
	err := DefineAndCreateDeploymentWithContainerPorts(1, nil, testNamespace)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid containers number")
	assert.Nil(t, *captured)
}

func TestDaemonSetWrappers(t *testing.T) {
	testCases := []struct {
		name      string
		create    func() error
		wantLabel map[string]string
	}{
		{
			name:   "multus daemonset",
			create: func() error { return DefineAndCreateDaemonsetWithMultusOnCluster("nad-a", testNamespace, "ds") },
		},
		{
			name: "multus daemonset with skip label",
			create: func() error {
				return DefineAndCreateDaemonsetWithMultusAndSkipLabelOnCluster("nad-a", testNamespace, "ds")
			},
			wantLabel: tsparams.NetworkingTestMultusSkipLabel,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			withFakeCluster(t)

			var captured *appsv1.DaemonSet

			createAndWaitDaemonSet = func(ds *appsv1.DaemonSet, _ time.Duration) error {
				captured = ds

				return nil
			}

			require.NoError(t, testCase.create())
			require.NotNil(t, captured)
			assert.Equal(t, "ds", captured.Name)
			assert.Equal(t, testNamespace, captured.Namespace)
			assert.Contains(t, captured.Spec.Template.Annotations[multusAnnotation], "nad-a")
			assert.Equal(t, map[string]string{cnfNodeLabel: ""}, captured.Spec.Template.Spec.NodeSelector)

			for key, value := range testCase.wantLabel {
				assert.Equal(t, value, captured.Spec.Template.Labels[key])
			}

			createAndWaitDaemonSet = func(*appsv1.DaemonSet, time.Duration) error { return errInjected }

			assert.ErrorIs(t, testCase.create(), errInjected)
		})
	}
}

func TestDefineAndCreateServiceOnCluster(t *testing.T) {
	testCases := []struct {
		name         string
		ipFamPolicy  string
		withNodePort bool
		headless     bool
	}{
		{name: "cluster IP service"},
		{name: "dual stack policy", ipFamPolicy: string(corev1.IPFamilyPolicyPreferDualStack)},
		{name: "node port service", withNodePort: true},
		{name: "headless service", headless: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client := withFakeCluster(t)

			err := DefineAndCreateServiceOnCluster("svc", testNamespace, 80, 8080, testCase.withNodePort, testCase.headless,
				[]corev1.IPFamily{corev1.IPv4Protocol}, testCase.ipFamPolicy)
			require.NoError(t, err)

			svc, err := client.CoreV1().Services(testNamespace).Get(context.TODO(), "svc", metav1.GetOptions{})
			require.NoError(t, err)

			if testCase.ipFamPolicy == "" {
				assert.Nil(t, svc.Spec.IPFamilyPolicy)
			} else {
				assert.Equal(t, corev1.IPFamilyPolicy(testCase.ipFamPolicy), *svc.Spec.IPFamilyPolicy)
			}

			if testCase.withNodePort {
				assert.Equal(t, corev1.ServiceTypeNodePort, svc.Spec.Type)
				assert.Equal(t, int32(80), svc.Spec.Ports[0].NodePort)
			}

			if testCase.headless {
				assert.Equal(t, corev1.ClusterIPNone, svc.Spec.ClusterIP)
			}
		})
	}

	t.Run("existing service is not an error", func(t *testing.T) {
		withFakeCluster(t, &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc", Namespace: testNamespace}})

		assert.NoError(t, DefineAndCreateServiceOnCluster("svc", testNamespace, 80, 8080, false, false, nil, ""))
	})

	t.Run("create failure is wrapped", func(t *testing.T) {
		client := withFakeCluster(t)
		failOn(client, "create", "services")

		err := DefineAndCreateServiceOnCluster("svc", testNamespace, 80, 8080, false, false, nil, "")
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to create service on cluster")
	})
}

func TestRedefineServiceToHeadless(t *testing.T) {
	svc := &corev1.Service{Spec: corev1.ServiceSpec{ClusterIP: "10.0.0.1"}}
	RedefineServiceToHeadless(svc)
	assert.Equal(t, corev1.ClusterIPNone, svc.Spec.ClusterIP)
}

func TestDefineAndCreateNadOnCluster(t *testing.T) {
	getNad := func(t *testing.T) *nadv1.NetworkAttachmentDefinition {
		t.Helper()

		created := &nadv1.NetworkAttachmentDefinition{}
		require.NoError(t, globalhelper.GetAPIClient().Get(context.TODO(),
			runtimeclient.ObjectKey{Namespace: testNamespace, Name: "nad"}, created))

		return created
	}

	t.Run("without whereabouts IPAM", func(t *testing.T) {
		withFakeCluster(t)
		setNadClient(t, nil)

		require.NoError(t, DefineAndCreateNadOnCluster("nad", testNamespace, ""))
		assert.NotContains(t, getNad(t).Spec.Config, "whereabouts")
	})

	t.Run("with whereabouts IPAM", func(t *testing.T) {
		withFakeCluster(t)
		setNadClient(t, nil)

		require.NoError(t, DefineAndCreateNadOnCluster("nad", testNamespace, "192.168.0.0/24"))
		assert.Contains(t, getNad(t).Spec.Config, `"range": "192.168.0.0/24"`)
	})

	t.Run("existing nad is not an error", func(t *testing.T) {
		withFakeCluster(t)
		setNadClient(t, nil, nad.DefineNad("nad", testNamespace))

		assert.NoError(t, DefineAndCreateNadOnCluster("nad", testNamespace, ""))
	})

	t.Run("create failure is wrapped", func(t *testing.T) {
		withFakeCluster(t)
		setNadClient(t, &interceptor.Funcs{
			Create: func(context.Context, runtimeclient.WithWatch, runtimeclient.Object, ...runtimeclient.CreateOption) error {
				return errInjected
			},
		})

		err := DefineAndCreateNadOnCluster("nad", testNamespace, "")
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to create nad on cluster")
	})
}

func TestDefineAndCreateNetworkPolicy(t *testing.T) {
	withFakeCluster(t)

	var captured *networkingv1.NetworkPolicy

	createNetworkPolicy = func(policy *networkingv1.NetworkPolicy, _ time.Duration) error {
		captured = policy

		return nil
	}

	labels := map[string]string{"app": "test"}
	require.NoError(t, DefineAndCreateNetworkPolicy("deny-all", testNamespace, []string{"Ingress", "Egress"}, labels))
	require.NotNil(t, captured)
	assert.Equal(t, "deny-all", captured.Name)
	assert.Equal(t, testNamespace, captured.Namespace)
	assert.Equal(t, []networkingv1.PolicyType{networkingv1.PolicyTypeIngress, networkingv1.PolicyTypeEgress},
		captured.Spec.PolicyTypes)
	assert.Equal(t, labels, captured.Spec.PodSelector.MatchLabels)

	createNetworkPolicy = func(*networkingv1.NetworkPolicy, time.Duration) error { return errInjected }

	assert.ErrorIs(t, DefineAndCreateNetworkPolicy("deny-all", testNamespace, nil, labels), errInjected)
}

func testPod(name, nodeName string, labels map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace, Labels: labels},
		Spec:       corev1.PodSpec{NodeName: nodeName},
	}
}

func TestExecCmdOnPodsInNamespace(t *testing.T) {
	command := []string{"ping", "-c", "1"}

	countExecs := func(t *testing.T, err error) *[]string {
		t.Helper()

		var podNames []string

		execPodCommand = func(pod corev1.Pod, cmd []string) (bytes.Buffer, error) {
			assert.Equal(t, command, cmd)
			podNames = append(podNames, pod.Name)

			return bytes.Buffer{}, err
		}

		return &podNames
	}

	t.Run("first pod only", func(t *testing.T) {
		withFakeCluster(t, testPod("pod-a", "", nil), testPod("pod-b", "", nil))
		podNames := countExecs(t, nil)

		require.NoError(t, ExecCmdOnOnePodInNamespace(command, testNamespace))
		assert.Equal(t, []string{"pod-a"}, *podNames)
	})

	t.Run("all pods", func(t *testing.T) {
		withFakeCluster(t, testPod("pod-a", "", nil), testPod("pod-b", "", nil))
		podNames := countExecs(t, nil)

		require.NoError(t, ExecCmdOnAllPodInNamespace(command, testNamespace))
		assert.ElementsMatch(t, []string{"pod-a", "pod-b"}, *podNames)
	})

	t.Run("first pod in empty namespace", func(t *testing.T) {
		withFakeCluster(t)
		podNames := countExecs(t, nil)

		err := ExecCmdOnOnePodInNamespace(command, testNamespace)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no pods found")
		assert.Empty(t, *podNames)
	})

	t.Run("invalid target", func(t *testing.T) {
		withFakeCluster(t)

		err := execCmdOnPodsListInNamespace(command, "some", testNamespace)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid parameter some")
	})

	t.Run("exec failure", func(t *testing.T) {
		withFakeCluster(t, testPod("pod-a", "", nil), testPod("pod-b", "", nil))
		podNames := countExecs(t, errInjected)

		assert.ErrorIs(t, ExecCmdOnAllPodInNamespace(command, testNamespace), errInjected)
		assert.Len(t, *podNames, 1)
	})

	t.Run("pod list failure", func(t *testing.T) {
		client := withFakeCluster(t)
		failOn(client, "list", "pods")

		assert.ErrorIs(t, ExecCmdOnAllPodInNamespace(command, testNamespace), errInjected)
	})
}

const ipLinkOutput = `[
	{"ifindex": 1, "ifname": "lo"},
	{"ifindex": 2, "ifname": "ens3"},
	{"ifindex": 3, "ifname": "ens4"},
	{"ifindex": 4, "ifname": "ovn-k8s-mp0"},
	{"ifindex": 5, "ifname": "br-ex"},
	{"ifindex": 6, "ifname": "ovs-system"},
	{"ifindex": 7, "ifname": "veth1", "master": "ovs-system"}
]`

func stubIPLink(output string, err error) {
	execPodCommand = func(corev1.Pod, []string) (bytes.Buffer, error) {
		return *bytes.NewBufferString(output), err
	}
}

func TestGetInterfacesList(t *testing.T) {
	testCases := []struct {
		name    string
		output  string
		execErr error
		want    []string
		wantErr string
	}{
		{name: "filters virtual and enslaved interfaces", output: ipLinkOutput, want: []string{"ens3", "ens4"}},
		{name: "exec failure", execErr: errInjected, wantErr: errInjected.Error()},
		{name: "malformed output", output: "not-json", wantErr: "invalid character"},
		{name: "no multus interfaces", output: `[{"ifname": "lo"}, {"ifname": "br-int"}]`, wantErr: "no multus interfaces"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			withFakeCluster(t)
			stubIPLink(testCase.output, testCase.execErr)

			interfaces, err := getInterfacesList(corev1.Pod{})
			if testCase.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.want, interfaces)
		})
	}
}

func TestGetClusterMultusInterfaces(t *testing.T) {
	workerA := &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "worker-a"}}
	workerB := &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "worker-b"}}
	master := &corev1.Node{ObjectMeta: metav1.ObjectMeta{
		Name:   "master-0",
		Labels: map[string]string{"node-role.kubernetes.io/control-plane": ""},
	}}

	t.Run("intersects interfaces across worker nodes", func(t *testing.T) {
		withFakeCluster(t, workerA, workerB, master,
			testPod("pod-a", "worker-a", nil), testPod("pod-b", "worker-b", nil), testPod("pod-m", "master-0", nil))

		var privilegedDaemonSet *appsv1.DaemonSet

		createAndWaitDaemonSet = func(ds *appsv1.DaemonSet, _ time.Duration) error {
			privilegedDaemonSet = ds

			return nil
		}

		perNode := map[string]string{
			"worker-a": `[{"ifname": "ens3"}, {"ifname": "ens4"}, {"ifname": "ens5"}]`,
			"worker-b": `[{"ifname": "ens4"}, {"ifname": "ens5"}, {"ifname": "ens6"}]`,
		}

		var execNodes []string

		execPodCommand = func(pod corev1.Pod, _ []string) (bytes.Buffer, error) {
			execNodes = append(execNodes, pod.Spec.NodeName)

			return *bytes.NewBufferString(perNode[pod.Spec.NodeName]), nil
		}

		interfaces, err := GetClusterMultusInterfaces(testNamespace)
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"ens4", "ens5"}, interfaces)
		assert.ElementsMatch(t, []string{"worker-a", "worker-b"}, execNodes)

		require.NotNil(t, privilegedDaemonSet)
		assert.Equal(t, map[string]string{workerNodeLabel: ""}, privilegedDaemonSet.Spec.Template.Spec.NodeSelector)
		assert.True(t, privilegedDaemonSet.Spec.Template.Spec.HostNetwork)
	})

	t.Run("control plane only cluster has no interfaces", func(t *testing.T) {
		withFakeCluster(t, master, testPod("pod-m", "master-0", nil))
		stubIPLink("", errInjected)

		interfaces, err := GetClusterMultusInterfaces(testNamespace)
		require.NoError(t, err)
		assert.Empty(t, interfaces)
	})

	t.Run("daemonset creation failure", func(t *testing.T) {
		withFakeCluster(t)
		createAndWaitDaemonSet = func(*appsv1.DaemonSet, time.Duration) error { return errInjected }

		_, err := GetClusterMultusInterfaces(testNamespace)
		assert.ErrorIs(t, err, errInjected)
	})

	t.Run("pod list failure", func(t *testing.T) {
		client := withFakeCluster(t)
		failOn(client, "list", "pods")

		_, err := GetClusterMultusInterfaces(testNamespace)
		assert.ErrorIs(t, err, errInjected)
	})

	t.Run("pod scheduled on unknown node", func(t *testing.T) {
		withFakeCluster(t, testPod("pod-a", "missing", nil))

		_, err := GetClusterMultusInterfaces(testNamespace)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("interface listing failure", func(t *testing.T) {
		withFakeCluster(t, workerA, testPod("pod-a", "worker-a", nil))
		stubIPLink("", errInjected)

		_, err := GetClusterMultusInterfaces(testNamespace)
		assert.ErrorIs(t, err, errInjected)
	})
}

func TestDefineDeploymentWithContainers(t *testing.T) {
	dep, err := DefineDeploymentWithContainers(3, 2, "multi", testNamespace)
	require.NoError(t, err)
	assert.Equal(t, "multi", dep.Name)
	assert.Equal(t, int32(3), *dep.Spec.Replicas)
	assert.Len(t, dep.Spec.Template.Spec.Containers, 2)

	_, err = DefineDeploymentWithContainers(1, 0, "multi", testNamespace)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid containers number")
}

func TestDefineDeploymentWithContainerPorts(t *testing.T) {
	ports := []corev1.ContainerPort{{ContainerPort: 8080}, {ContainerPort: 8443}}

	dep, err := DefineDeploymentWithContainerPorts("ported", testNamespace, 2, ports)
	require.NoError(t, err)
	assert.Equal(t, "ported", dep.Name)
	assert.Equal(t, int32(2), *dep.Spec.Replicas)

	containers := dep.Spec.Template.Spec.Containers
	require.Len(t, containers, 2)

	for index, container := range containers {
		assert.Equal(t, []corev1.ContainerPort{ports[index]}, container.Ports)
	}

	_, err = DefineDeploymentWithContainerPorts("ported", testNamespace, 1, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid number of containers")
}

func TestUniquePodNodeNames(t *testing.T) {
	selected := map[string]string{"app": "nginx"}

	t.Run("lists distinct nodes for matching pods", func(t *testing.T) {
		withFakeCluster(t,
			testPod("pod-a", "worker-a", selected),
			testPod("pod-b", "worker-a", selected),
			testPod("pod-c", "worker-b", selected),
			testPod("pod-d", "worker-c", map[string]string{"app": "other"}))

		names, err := UniquePodNodeNames(testNamespace, "app=nginx")
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"worker-a", "worker-b"}, names)
	})

	t.Run("pod list failure", func(t *testing.T) {
		client := withFakeCluster(t)
		failOn(client, "list", "pods")

		_, err := UniquePodNodeNames(testNamespace, "app=nginx")
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to list pods in "+testNamespace)
	})
}

func TestDefineTLSNginxDeployment(t *testing.T) {
	t.Run("creates TLS materials and mounts them", func(t *testing.T) {
		client := withFakeCluster(t)

		dep, err := DefineTLSNginxDeployment("tls", testNamespace, 2)
		require.NoError(t, err)
		assert.Equal(t, "tls", dep.Name)
		assert.Equal(t, int32(2), *dep.Spec.Replicas)

		container := dep.Spec.Template.Spec.Containers[0]
		assert.Equal(t, tsparams.TLSListenPort, container.Ports[0].ContainerPort)
		assert.Equal(t, tsparams.TLSListenPort, container.ReadinessProbe.TCPSocket.Port.IntVal)
		require.Len(t, container.VolumeMounts, 2)
		assert.Equal(t, tlsCertMountPath, container.VolumeMounts[0].MountPath)
		require.Len(t, dep.Spec.Template.Spec.Volumes, 2)
		assert.Equal(t, tlsSecretName, dep.Spec.Template.Spec.Volumes[0].Secret.SecretName)

		secret, err := client.CoreV1().Secrets(testNamespace).Get(context.TODO(), tlsSecretName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, corev1.SecretTypeTLS, secret.Type)
		assert.NotEmpty(t, secret.Data[corev1.TLSCertKey])
		assert.NotEmpty(t, secret.Data[corev1.TLSPrivateKeyKey])

		configMap, err := client.CoreV1().ConfigMaps(testNamespace).Get(context.TODO(), nginxConfigMapName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, nginxTLSConf, configMap.Data["default.conf"])
	})

	t.Run("existing TLS materials are reused", func(t *testing.T) {
		withFakeCluster(t,
			&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: tlsSecretName, Namespace: testNamespace}},
			&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: nginxConfigMapName, Namespace: testNamespace}})

		_, err := DefineTLSNginxDeployment("tls", testNamespace, 1)
		assert.NoError(t, err)
	})

	failureCases := []struct {
		name    string
		setup   func(*k8sfake.Clientset)
		wantErr string
	}{
		{
			name: "certificate generation failure",
			setup: func(*k8sfake.Clientset) {
				generateTLSCert = func() ([]byte, []byte, error) { return nil, nil, errInjected }
			},
			wantErr: errInjected.Error(),
		},
		{
			name:    "secret creation failure",
			setup:   func(client *k8sfake.Clientset) { failOn(client, "create", "secrets") },
			wantErr: "failed to create TLS secret",
		},
		{
			name:    "configmap creation failure",
			setup:   func(client *k8sfake.Clientset) { failOn(client, "create", "configmaps") },
			wantErr: "failed to create nginx TLS configmap",
		},
	}

	for _, testCase := range failureCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setup(withFakeCluster(t))

			dep, err := DefineTLSNginxDeployment("tls", testNamespace, 1)
			require.ErrorIs(t, err, errInjected)
			assert.Contains(t, err.Error(), testCase.wantErr)
			assert.Nil(t, dep)
		})
	}
}

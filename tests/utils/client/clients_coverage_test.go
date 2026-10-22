package client

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	netattdefv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var errInjected = errors.New("injected failure")

// writeKubeconfig writes a minimal valid kubeconfig pointing at serverURL and
// returns its path.
func writeKubeconfig(t *testing.T, serverURL string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "kubeconfig")

	content := fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- cluster:
    server: %s
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test
current-context: test
users:
- name: test-user
  user:
    token: test-token
`, serverURL)

	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}

// panicsWith calls New and returns the error it panicked with. It fails the
// test when New does not panic or panics with a non-error value.
func panicsWith(t *testing.T, kubeconfig string) (panicked error) {
	t.Helper()

	defer func() {
		rec := recover()
		require.NotNil(t, rec, "expected New to panic")

		var ok bool

		panicked, ok = rec.(error)
		require.True(t, ok, "New panicked with a non-error value: %v", rec)
	}()

	New(kubeconfig)

	return nil
}

// stubSchemeAdders replaces the scheme registration seam with a failing
// function; the previous list is restored on cleanup.
func stubSchemeAdders(t *testing.T) {
	t.Helper()

	oldAdders := schemeAdders

	t.Cleanup(func() { schemeAdders = oldAdders })

	schemeAdders = []func(*runtime.Scheme) error{
		func(*runtime.Scheme) error { return errInjected },
	}
}

// stubRuntimeClient overrides the controller-runtime client seam to fail; the
// previous constructor is restored on cleanup.
func stubRuntimeClient(t *testing.T) {
	t.Helper()

	oldNew := newRuntimeClient

	t.Cleanup(func() { newRuntimeClient = oldNew })

	newRuntimeClient = func(*rest.Config, runtimeclient.Options) (runtimeclient.Client, error) {
		return nil, errInjected
	}
}

func TestNew(t *testing.T) {
	t.Run("wires every client from an explicit kubeconfig path", func(t *testing.T) {
		clientSet := New(writeKubeconfig(t, "https://127.0.0.1:6443"))

		require.NotNil(t, clientSet.Config)
		assert.Equal(t, "https://127.0.0.1:6443", clientSet.Config.Host)

		assert.NotNil(t, clientSet.CoreV1Interface)
		assert.NotNil(t, clientSet.OcpClientInterface)
		assert.NotNil(t, clientSet.NetworkPolicies("default"))
		assert.NotNil(t, clientSet.RbacV1Interface)
		assert.NotNil(t, clientSet.K8sClient)
		assert.NotNil(t, clientSet.AppsV1Interface)
		assert.NotNil(t, clientSet.ApiextensionsV1Interface)
		assert.NotNil(t, clientSet.DiscoveryInterface)
		assert.NotNil(t, clientSet.OperatorsV1alpha1Interface)
		assert.NotNil(t, clientSet.MachineconfigurationV1Interface)
		assert.NotNil(t, clientSet.NodeV1Interface)
		assert.NotNil(t, clientSet.PolicyV1Interface)
		assert.NotNil(t, clientSet.PolicyV1beta1Interface)
		assert.NotNil(t, clientSet.DynamicClient)
		assert.NotNil(t, clientSet.Client)

		// The controller-runtime client must carry the scheme merged from
		// every entry in schemeAdders.
		scheme := clientSet.Scheme()

		knownGVKs := []schema.GroupVersionKind{
			corev1.SchemeGroupVersion.WithKind("Pod"),
			netattdefv1.SchemeGroupVersion.WithKind("NetworkAttachmentDefinition"),
			operatorsv1alpha1.SchemeGroupVersion.WithKind("ClusterServiceVersion"),
		}

		for _, gvk := range knownGVKs {
			assert.True(t, scheme.Recognizes(gvk), "scheme should recognize %v", gvk)
		}
	})

	t.Run("falls back to the KUBECONFIG environment variable", func(t *testing.T) {
		t.Setenv("KUBECONFIG", writeKubeconfig(t, "https://127.0.0.1:6444"))

		clientSet := New("")

		require.NotNil(t, clientSet.Config)
		assert.Equal(t, "https://127.0.0.1:6444", clientSet.Config.Host)
	})

	t.Run("panics when the kubeconfig path does not exist", func(t *testing.T) {
		err := panicsWith(t, filepath.Join(t.TempDir(), "missing-kubeconfig"))
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("panics when the in-cluster config is unavailable", func(t *testing.T) {
		// Empty service host and port make rest.InClusterConfig fail
		// deterministically, even when running inside a pod.
		t.Setenv("KUBECONFIG", "")
		t.Setenv("KUBERNETES_SERVICE_HOST", "")
		t.Setenv("KUBERNETES_PORT", "")

		err := panicsWith(t, "")
		assert.ErrorIs(t, err, rest.ErrNotInCluster)
	})

	t.Run("panics when scheme registration fails", func(t *testing.T) {
		stubSchemeAdders(t)

		err := panicsWith(t, writeKubeconfig(t, "https://127.0.0.1:6443"))
		assert.ErrorIs(t, err, errInjected)
	})

	t.Run("panics when the controller-runtime client cannot be created", func(t *testing.T) {
		stubRuntimeClient(t)

		err := panicsWith(t, writeKubeconfig(t, "https://127.0.0.1:6443"))
		assert.ErrorIs(t, err, errInjected)
	})
}

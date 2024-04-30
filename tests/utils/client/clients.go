package client

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"

	"github.com/golang/glog"
	netattdefv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	ocpclientconfigv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/scheme"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	apiextv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	appsv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	networkv1client "k8s.io/client-go/kubernetes/typed/networking/v1"
	nodev1 "k8s.io/client-go/kubernetes/typed/node/v1"
	policyv1 "k8s.io/client-go/kubernetes/typed/policy/v1"
	policyv1beta1 "k8s.io/client-go/kubernetes/typed/policy/v1beta1"
	rbacv1client "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ClientSet provides the struct to talk with relevant API.
type ClientSet struct {
	corev1client.CoreV1Interface
	OcpClientInterface ocpclientconfigv1.ConfigV1Interface
	networkv1client.NetworkingV1Client
	rbacv1client.RbacV1Interface
	K8sClient kubernetes.Interface
	appsv1client.AppsV1Interface
	apiextv1client.ApiextensionsV1Interface
	discovery.DiscoveryInterface
	Config *rest.Config
	runtimeclient.Client
	v1alpha1.OperatorsV1alpha1Interface
	nodev1.NodeV1Interface
	policyv1.PolicyV1Interface
	policyv1beta1.PolicyV1beta1Interface
	DynamicClient dynamic.Interface
}

// New returns a *ClientBuilder with the given kubeconfig.
func New(kubeconfig string) *ClientSet {
	var (
		config *rest.Config
		err    error
	)

	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	if kubeconfig != "" {
		glog.V(4).Infof("Loading kube client config from path %q", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		glog.V(4).Infof("Using in-cluster kube client config")

		config, err = rest.InClusterConfig()
	}

	if err != nil {
		panic(err)
	}

	clientSet := &ClientSet{}
	clientSet.CoreV1Interface = corev1client.NewForConfigOrDie(config)
	clientSet.OcpClientInterface = ocpclientconfigv1.NewForConfigOrDie(config)
	clientSet.AppsV1Interface = appsv1client.NewForConfigOrDie(config)
	clientSet.RbacV1Interface = rbacv1client.NewForConfigOrDie(config)
	clientSet.DiscoveryInterface = discovery.NewDiscoveryClientForConfigOrDie(config)
	clientSet.NetworkingV1Client = *networkv1client.NewForConfigOrDie(config)
	clientSet.OperatorsV1alpha1Interface = v1alpha1.NewForConfigOrDie(config)
	clientSet.ApiextensionsV1Interface = apiextv1client.NewForConfigOrDie(config)
	clientSet.NodeV1Interface = nodev1.NewForConfigOrDie(config)
	clientSet.PolicyV1Interface = policyv1.NewForConfigOrDie(config)
	clientSet.PolicyV1beta1Interface = policyv1beta1.NewForConfigOrDie(config)
	clientSet.DynamicClient = dynamic.NewForConfigOrDie(config)
	clientSet.K8sClient = kubernetes.NewForConfigOrDie(config)
	clientSet.Config = config

	crScheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(crScheme); err != nil {
		panic(err)
	}

	if err := netattdefv1.SchemeBuilder.AddToScheme(crScheme); err != nil {
		panic(err)
	}

	if err := olm.AddToScheme(crScheme); err != nil {
		panic(err)
	}

	clientSet.Client, err = runtimeclient.New(config, runtimeclient.Options{
		Scheme: crScheme,
	})

	if err != nil {
		panic(err)
	}

	return clientSet
}

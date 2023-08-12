package globalhelper

import (
	"fmt"

	"github.com/golang/glog"
	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"k8s.io/client-go/kubernetes"
)

var (
	apiclient *testclient.ClientSet
	conf      *config.Config
)

func SetTestK8sAPIClient(client kubernetes.Interface) {
	if apiclient == nil {
		apiclient = &testclient.ClientSet{}
	}

	//nolint:godox
	// TODO: Add more client interfaces as needed
	apiclient.CoreV1Interface = client.CoreV1()
	apiclient.AppsV1Interface = client.AppsV1()
	apiclient.RbacV1Interface = client.RbacV1()
	apiclient.NodeV1Interface = client.NodeV1()
	apiclient.K8sClient = client
}

func UnsetTestK8sAPIClient() {
	apiclient = nil
}

func GetAPIClient() *testclient.ClientSet {
	if apiclient != nil {
		return apiclient
	}

	var err error
	apiclient, err = config.DefineClients()

	if err != nil {
		glog.Fatal(fmt.Sprintf("can not load api client. Please check KUBECONFIG env var - %s", err))
	}

	return apiclient
}

func GetConfiguration() *config.Config {
	if conf != nil {
		return conf
	}

	var err error
	conf, err = config.NewConfig()

	if err != nil {
		glog.Fatal(fmt.Sprintf("can not load configuration - %s", err))
	}

	return conf
}

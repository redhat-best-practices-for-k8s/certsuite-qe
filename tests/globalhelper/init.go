package globalhelper

import (
	"fmt"

	"github.com/golang/glog"
	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
)

var (
	apiclient *testclient.ClientSet
	conf      *config.Config
)

func GetAPIClient() *testclient.ClientSet {
	if apiclient != nil {
		return apiclient
	}

	tempClient, err := config.DefineClients()
	if err != nil {
		glog.Fatal(fmt.Sprintf("can not load api client. Please check KUBECONFIG env var - %s", err))
	}

	apiclient = tempClient

	return tempClient
}

func GetConfiguration() *config.Config {
	if conf != nil {
		return conf
	}

	Configuration, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Sprintf("can not load configuration - %s", err))
	}

	conf = Configuration

	return Configuration
}

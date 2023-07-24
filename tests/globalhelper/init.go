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

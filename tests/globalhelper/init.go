package globalhelper

import (
	"fmt"

	"github.com/golang/glog"
	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
)

var (
	APIClient     *testclient.ClientSet
	Configuration *config.Config
)

func init() {
	var err error
	APIClient, err = config.DefineClients()

	if err != nil {
		glog.Fatal(fmt.Errorf("can not load api client. Please check KUBECONFIG env var"))
	}

	Configuration, err = config.NewConfig()

	if err != nil {
		glog.Fatal(fmt.Errorf("can not load configuration"))
	}
}

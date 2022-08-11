package globalhelper

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
)

var (
	APIClient     *testclient.ClientSet
	Configuration *config.Config
)

func init() {
	if os.Getenv("UNIT_TEST") != "" {
		Configuration = &config.Config{}

		return
	}

	var err error
	APIClient, err = config.DefineClients()

	if err != nil {
		glog.Fatal(fmt.Sprintf("can not load api client. Please check KUBECONFIG env var - %s", err))
	}

	Configuration, err = config.NewConfig()

	if err != nil {
		glog.Fatal(fmt.Sprintf("can not load configuration - %s", err))
	}
}

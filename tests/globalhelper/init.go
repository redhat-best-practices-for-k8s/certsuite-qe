package globalhelper

import (
	"fmt"
	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"log"
)

var (
	ApiClient *testclient.ClientSet
)

func init() {
	var err error
	ApiClient, err = config.DefineClients()
	if err != nil {
		log.Fatal(fmt.Errorf("can not load api client. Please check KUBECONFIG env var"))
	}
}

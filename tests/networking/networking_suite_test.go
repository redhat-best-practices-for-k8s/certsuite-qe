package networking

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
	"time"

	"fmt"

	"runtime"
	"testing"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/networking/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/cluster"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func TestNetworking(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	configSuite, err := config.NewConfig()

	if err != nil {
		fmt.Print(err)
		return
	}

	junitPath := configSuite.GetReportPath(currentFile)
	RegisterFailHandler(Fail)
	rr := append([]Reporter{}, reporters.NewJUnitReporter(junitPath))
	RunSpecsWithDefaultAndCustomReporters(t, "CNFCert networking tests", rr)
}

var _ = BeforeSuite(func() {
	By("Validate that cluster is Schedulable")
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.ApiClient)
		Expect(err).ToNot(HaveOccurred())
		return isClusterReady
	}, WaitingTime, 5*time.Second).Should(BeTrue())

	By("Validate that all nodes are Ready")
	err := nodes.WaitForNodesReady(globalhelper.ApiClient, WaitingTime, 5)
	Expect(err).ToNot(HaveOccurred())

	By(fmt.Sprintf("Create %s namespace", TestNetworkingNameSpace))
	err = namespaces.Create(TestNetworkingNameSpace, globalhelper.ApiClient)
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	By(fmt.Sprintf("Remove %s namespace", TestNetworkingNameSpace))
	err := namespaces.DeleteAndWait(globalhelper.ApiClient, TestNetworkingNameSpace, WaitingTime)
	Expect(err).ToNot(HaveOccurred())
})

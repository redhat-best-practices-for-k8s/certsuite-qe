package networking

import (
	"flag"
	"runtime"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/nethelper"

	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	"testing"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/networking/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/cluster"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func TestNetworking(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert networking tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	By("Validate that cluster is Schedulable")
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, netparameters.WaitingTime, netparameters.RetryInterval*time.Second).Should(BeTrue())

	By("Validate that all nodes are Ready")
	err := nodes.WaitForNodesReady(globalhelper.APIClient, netparameters.WaitingTime, netparameters.RetryInterval)
	Expect(err).ToNot(HaveOccurred())

	By(fmt.Sprintf("Create %s namespace", netparameters.TestNetworkingNameSpace))
	err = namespaces.Create(netparameters.TestNetworkingNameSpace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{netparameters.TestNetworkingNameSpace},
		[]string{netparameters.TestPodLabel},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())

	By("Set rbac policy which allows authenticated users to run privileged containers")
	err = nethelper.AllowAuthenticatedUsersRunPrivilegedContainers()
	Expect(err).ToNot(HaveOccurred())

})

var _ = AfterSuite(func() {
	By(fmt.Sprintf("Remove %s namespace", netparameters.TestNetworkingNameSpace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		netparameters.TestNetworkingNameSpace,
		netparameters.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from report directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

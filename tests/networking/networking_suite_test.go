//go:build !utest

package networking

import (
	"flag"
	"fmt"
	"runtime"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/networking/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/cluster"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

func TestNetworking(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert networking tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	By("Validate that cluster is Schedulable")
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, tsparams.WaitingTime, tsparams.RetryInterval*time.Second).Should(BeTrue())

	By("Validate that all nodes are Ready")
	err := nodes.WaitForNodesReady(globalhelper.GetAPIClient(), tsparams.WaitingTime, tsparams.RetryInterval)
	Expect(err).ToNot(HaveOccurred())

	By(fmt.Sprintf("Create %s namespace", tsparams.TestNetworkingNameSpace))
	err = namespaces.Create(tsparams.TestNetworkingNameSpace, globalhelper.GetAPIClient())
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.TestNetworkingNameSpace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())

	By("Set rbac policy which allows authenticated users to run privileged containers")
	err = globalhelper.AllowAuthenticatedUsersRunPrivilegedContainers()
	Expect(err).ToNot(HaveOccurred())

})

var _ = AfterSuite(func() {
	By("Remove networking test namespaces")
	err := namespaces.DeleteAndWait(
		globalhelper.GetAPIClient().CoreV1Interface,
		tsparams.TestNetworkingNameSpace,
		tsparams.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	err = namespaces.DeleteAndWait(
		globalhelper.GetAPIClient().CoreV1Interface,
		tsparams.AdditionalNetworkingNamespace,
		tsparams.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from report directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

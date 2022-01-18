package lifecycle

import (
	"flag"
	"fmt"
	"runtime"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/cluster"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

func TestLifecycle(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.LogLevel)
	junitPath := globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	rr := append([]Reporter{}, reporters.NewJUnitReporter(junitPath))
	RunSpecsWithDefaultAndCustomReporters(t, "CNFCert lifecycle tests", rr)
}

var _ = BeforeSuite(func() {

	By("Validate that the cluster is stable")
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, lifeparameters.WaitingTime, lifeparameters.RetryInterval*time.Second).Should(BeTrue())

	By("Validate that all nodes are ready")
	err := nodes.WaitForNodesReady(globalhelper.APIClient,
		lifeparameters.WaitingTime, lifeparameters.RetryInterval)
	Expect(err).ToNot(HaveOccurred())

	By("Create namespace")
	err = namespaces.Create(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{lifeparameters.LifecycleNamespace},
		[]string{lifeparameters.TestPodLabel},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())

})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", lifeparameters.LifecycleNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		lifeparameters.LifecycleNamespace,
		lifeparameters.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

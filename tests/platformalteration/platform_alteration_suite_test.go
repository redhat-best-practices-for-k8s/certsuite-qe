//go:build !utest

package platformalteration

import (
	"flag"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/cluster"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

func TestPlatformAlteration(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert platform-alteration tests", reporterConfig)
}

var _ = BeforeSuite(func() {
	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	By("Validate that cluster is Schedulable")
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, tsparams.WaitingTime, tsparams.RetryInterval*time.Second).Should(BeTrue())

	By("Validate that all nodes are Ready")
	err = nodes.WaitForNodesReady(globalhelper.GetAPIClient().Nodes(), tsparams.WaitingTime, tsparams.RetryInterval)
	Expect(err).ToNot(HaveOccurred())

	By("Ensure all nodes are labeled with 'worker-cnf' label")
	err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().Nodes(), configSuite.General.CnfNodeLabel)
	Expect(err).ToNot(HaveOccurred())
})

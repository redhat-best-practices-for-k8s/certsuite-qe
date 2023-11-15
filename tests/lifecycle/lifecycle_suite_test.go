//go:build !utest

package lifecycle

import (
	"flag"
	"os"
	"runtime"
	"testing"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
)

func TestLifecycle(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert lifecycle tests", reporterConfig)
}

var _ = SynchronizedBeforeSuite(func() {
	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatalf("can not load config file: %w", err)
	}

	err = tshelper.WaitUntilClusterIsStable()
	Expect(err).ToNot(HaveOccurred())

	By("Ensure all nodes are labeled with 'worker-cnf' label")
	err = nodes.EnsureAllNodesAreLabeled(configSuite.General.CnfNodeLabel)
	Expect(err).ToNot(HaveOccurred())
}, func() {})

var _ = SynchronizedAfterSuite(func() {}, func() {
	err := os.Unsetenv("TNF_NON_INTRUSIVE_ONLY")
	Expect(err).ToNot(HaveOccurred())
})

//go:build !utest

package lifecycle

import (
	"flag"
	"os"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/tests"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"
	klog "k8s.io/klog/v2"

	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/helper"
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
		klog.Fatalf("can not load config file: %v", err)
	}

	err = tshelper.WaitUntilClusterIsStable()
	Expect(err).ToNot(HaveOccurred())

	By("Ensure all nodes are labeled with 'worker-cnf' label")
	err = nodes.EnsureAllNodesAreLabeled(configSuite.General.CnfNodeLabel)
	Expect(err).ToNot(HaveOccurred())
}, func() {})

var _ = SynchronizedAfterSuite(func() {}, func() {
	err := os.Unsetenv("CERTSUITE_NON_INTRUSIVE_ONLY")
	Expect(err).ToNot(HaveOccurred())
})

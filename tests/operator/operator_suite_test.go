//go:build !utest

package operator

import (
	"flag"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/tests"
)

func TestOperator(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert operator tests", reporterConfig)
}

var _ = SynchronizedBeforeSuite(func() {
	if globalhelper.IsKindCluster() {
		Skip("Skipping operator tests on kind cluster")
	}

	// Safeguard against running the operator tests on a cluster without catalog sources
	if !globalhelper.IsKindCluster() {
		By("Create catalog sources and wait for them to become ready")
		err := globalhelper.CreateAndValidateCatalogSources(true)
		Expect(err).ToNot(HaveOccurred(), "All necessary catalog sources are not available")
	}
}, func() {})

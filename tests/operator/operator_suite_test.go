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

	if globalhelper.IsVanillaK8sCluster() {
		Skip("Skipping operator tests on kind cluster")
	}

	// Safeguard against running the operator tests on a cluster without catalog sources
	if !globalhelper.IsVanillaK8sCluster() {
		By("Create community-operators catalog source")
		err := globalhelper.CreateCommunityOperatorsCatalogSource()
		Expect(err).ToNot(HaveOccurred())

		By("Create certified-operators catalog source")
		err = globalhelper.DeployRHCertifiedOperatorSource("")
		Expect(err).ToNot(HaveOccurred())

		By("Create redhat-operators catalog source")
		err = globalhelper.DeployRHOperatorSource("")
		Expect(err).ToNot(HaveOccurred())

		By("Check if catalog sources are available")
		err = globalhelper.ValidateCatalogSources()
		Expect(err).ToNot(HaveOccurred(), "All necessary catalog sources are not available")
	}
}, func() {})

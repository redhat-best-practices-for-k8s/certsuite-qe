//go:build !utest

package operator

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/operator/tests"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/operator/parameters"
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

	By(fmt.Sprintf("Create %s namespace", tsparams.OperatorNamespace))
	err := globalhelper.CreateNamespace(tsparams.OperatorNamespace)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.OperatorNamespace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{},
		[]string{}, globalhelper.GetConfiguration().General.TnfConfigDir)
	Expect(err).ToNot(HaveOccurred())

	// Safeguard against running the operator tests on a cluster without catalog sources
	if !globalhelper.IsKindCluster() {
		By("Create community-operators catalog source")
		err = globalhelper.CreateCommunityOperatorsCatalogSource()
		Expect(err).ToNot(HaveOccurred())

		By("Check if catalog sources are available")
		err = globalhelper.ValidateCatalogSources()
		Expect(err).ToNot(HaveOccurred(), "All necessary catalog sources are not available")
	}
}, func() {})

var _ = SynchronizedAfterSuite(func() {}, func() {
	if globalhelper.IsKindCluster() {
		Skip("Skipping operator tests cleanup on kind cluster")
	}

	By(fmt.Sprintf("Remove %s namespace", tsparams.OperatorNamespace))
	err := globalhelper.DeleteNamespaceAndWait(
		tsparams.OperatorNamespace,
		tsparams.Timeout,
	)
	Expect(err).ToNot(HaveOccurred())
})

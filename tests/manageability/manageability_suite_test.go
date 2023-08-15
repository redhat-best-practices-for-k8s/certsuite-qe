//go:build !utest

package performance

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/manageability/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/manageability/tests"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func TestPerformance(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert performance tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	By("Create namespace")
	err := namespaces.Create(tsparams.ManageabilityNamespace, globalhelper.GetAPIClient())
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.ManageabilityNamespace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", tsparams.ManageabilityNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.GetAPIClient().CoreV1Interface,
		tsparams.ManageabilityNamespace,
		tsparams.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

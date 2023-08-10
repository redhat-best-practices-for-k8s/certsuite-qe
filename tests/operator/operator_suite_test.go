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
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

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

var _ = BeforeSuite(func() {

	By(fmt.Sprintf("Create %s namespace", tsparams.OperatorNamespace))
	err := namespaces.Create(tsparams.OperatorNamespace, globalhelper.GetAPIClient())
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.OperatorNamespace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", tsparams.OperatorNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.GetAPIClient(),
		tsparams.OperatorNamespace,
		tsparams.Timeout,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

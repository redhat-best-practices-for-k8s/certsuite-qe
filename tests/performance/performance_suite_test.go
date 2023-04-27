//go:build !utest

package performance

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/performance/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/performance/tests"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/performance/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func TestPerformance(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert performance tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	By("Create namespace")
	err := namespaces.Create(tsparams.PerformanceNamespace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	// Create service account and roles and roles binding
	err = tshelper.ConfigurePrivilegedServiceAccount(tsparams.PerformanceNamespace)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.PerformanceNamespace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", tsparams.PerformanceNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		tsparams.PerformanceNamespace,
		tsparams.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

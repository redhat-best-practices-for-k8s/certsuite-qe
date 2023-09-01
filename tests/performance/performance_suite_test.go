//go:build !utest

package performance

import (
	"flag"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"

	_ "github.com/test-network-function/cnfcert-tests-verification/tests/performance/tests"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
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

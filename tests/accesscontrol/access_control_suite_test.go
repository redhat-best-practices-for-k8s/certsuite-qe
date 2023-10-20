//go:build !utest

package accesscontrol

import (
	"flag"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"runtime"
	"testing"

	_ "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

func TestAccessControl(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert access-control tests", reporterConfig)
}

var _ = SynchronizedBeforeSuite(func() {
	err := globalhelper.AllowAuthenticatedUsersRunPrivilegedContainers()
	Expect(err).ToNot(HaveOccurred(), "Error creating namespace")
}, func() {})

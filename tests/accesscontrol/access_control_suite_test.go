//go:build !utest

package accesscontrol

import (
	"flag"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"runtime"
	"testing"

	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func TestAccessControl(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert access-control tests", reporterConfig)
}

var _ = BeforeSuite(func() {
	By("Create namespace")
	err := namespaces.Create(parameters.TestAccessControlNameSpace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred(), "Error creating namespace")
})

var _ = AfterSuite(func() {

	By("Remove test namespaces")
	err := helper.DeleteNamespaces(
		[]string{parameters.TestAccessControlNameSpace,
			parameters.AdditionalValidNamespace,
			parameters.InvalidNamespace},
		globalhelper.APIClient,
		parameters.Timeout,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from report directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

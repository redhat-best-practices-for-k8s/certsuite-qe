//go:build !utest

package accesscontrol

import (
	"flag"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"runtime"
	"testing"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
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

var _ = BeforeSuite(func() {
	By("Create namespace")
	err := namespaces.Create(parameters.TestAccessControlNameSpace, globalhelper.GetAPIClient())
	Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

	err = globalhelper.AllowAuthenticatedUsersRunPrivilegedContainers()
	Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

})

var _ = AfterSuite(func() {

	By("Remove test namespaces")
	err := tshelper.DeleteNamespaces(
		[]string{parameters.TestAccessControlNameSpace,
			parameters.AdditionalValidNamespace,
			parameters.InvalidNamespace,
			parameters.TestAnotherNamespace,
		},
		globalhelper.GetAPIClient(),
		parameters.Timeout,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from report directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

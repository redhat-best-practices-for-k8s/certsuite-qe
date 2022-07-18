//go:build !utest

package affiliatedcertification

import (
	"flag"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"

	"runtime"
	"testing"

	_ "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

func TestAffiliatedCertification(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert affiliated-certification tests", reporterConfig)
}

var _ = BeforeSuite(func() {
	By("Create namespace")
	err := namespaces.Create(tsparams.TestCertificationNameSpace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred(), "Error creating namespace")
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", tsparams.TestCertificationNameSpace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		tsparams.TestCertificationNameSpace,
		tsparams.Timeout,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from report directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

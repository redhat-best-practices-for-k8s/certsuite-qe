package observability

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/observability/helper"
	params "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/observability/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func TestObservability(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert observability tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	By(fmt.Sprintf("Create %s namespace", params.QeTestNamespace))
	err := namespaces.Create(params.QeTestNamespace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{params.QeTestNamespace},
		helper.GetTnfTargetPodLabelsSlice(),
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", params.QeTestNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		params.QeTestNamespace,
		params.NsResourcesDeleteTimeoutMins,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

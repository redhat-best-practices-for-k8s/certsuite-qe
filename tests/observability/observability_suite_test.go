package observability

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"

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
	junitPath := globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	rr := append([]Reporter{}, reporters.NewJUnitReporter(junitPath))
	RunSpecsWithDefaultAndCustomReporters(t, "CNFCert observability tests", rr)
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

package lifecycle

import (
	"flag"
	"os"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

func TestLifecycle(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	junitPath := globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	rr := append([]Reporter{}, reporters.NewJUnitReporter(junitPath))
	RunSpecsWithDefaultAndCustomReporters(t, "CNFCert lifecycle tests", rr)
}

var _ = BeforeSuite(func() {

	err := lifehelper.WaitUntilClusterIsStable()
	Expect(err).ToNot(HaveOccurred())

	By("Create namespace")
	err = namespaces.Create(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{lifeparameters.LifecycleNamespace},
		[]string{lifeparameters.TestPodLabel},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	// By(fmt.Sprintf("Remove %s namespace", lifeparameters.LifecycleNamespace))
	// err := namespaces.DeleteAndWait(
	// 	globalhelper.APIClient,
	// 	lifeparameters.LifecycleNamespace,
	// 	lifeparameters.WaitingTime,
	// )
	// Expect(err).ToNot(HaveOccurred())

	// By("Remove reports from reports directory")
	// err = globalhelper.RemoveContentsFromReportDir()
	// Expect(err).ToNot(HaveOccurred())

	By("Remove masters scheduling")
	err := lifehelper.EnableMasterScheduling(false)
	Expect(err).ToNot(HaveOccurred())

	err = os.Unsetenv("TNF_NON_INTRUSIVE_ONLY")
	Expect(err).ToNot(HaveOccurred())

})

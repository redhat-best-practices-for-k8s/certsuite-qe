//go:build !utest

package lifecycle

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

func TestLifecycle(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert lifecycle tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	err := tshelper.WaitUntilClusterIsStable()
	Expect(err).ToNot(HaveOccurred())

	By("Create namespace")
	err = namespaces.Create(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.LifecycleNamespace},
		[]string{tsparams.TestPodLabel},
		[]string{tsparams.TnfTargetOperatorLabels}, // some operator labels are added here
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", tsparams.LifecycleNamespace))
	err := namespaces.DeleteAndWait(globalhelper.GetAPIClient(), tsparams.LifecycleNamespace, tsparams.WaitingTime)
	Expect(err).ToNot(HaveOccurred())

	By(fmt.Sprintf("Remove %s namespace", tsparams.TestCrdNamespace))
	err = namespaces.DeleteAndWait(globalhelper.GetAPIClient(), tsparams.TestCrdNamespace, tsparams.WaitingTime)
	Expect(err).ToNot(HaveOccurred())

	/*By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())*/

	By("Remove masters scheduling")
	err = globalhelper.EnableMasterScheduling(globalhelper.GetAPIClient().CoreV1Interface, false)
	Expect(err).ToNot(HaveOccurred())

	err = os.Unsetenv("TNF_NON_INTRUSIVE_ONLY")
	Expect(err).ToNot(HaveOccurred())
})

package lifecycle

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"

	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

func TestLifecycle(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert lifecycle tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	err := helper.WaitUntilClusterIsStable()
	Expect(err).ToNot(HaveOccurred())

	By("Create namespace")
	err = namespaces.Create(parameters.LifecycleNamespace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{parameters.LifecycleNamespace},
		[]string{parameters.TestPodLabel},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", parameters.LifecycleNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		parameters.LifecycleNamespace,
		parameters.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())

	By("Remove masters scheduling")
	err = helper.EnableMasterScheduling(false)
	Expect(err).ToNot(HaveOccurred())

	err = os.Unsetenv("TNF_NON_INTRUSIVE_ONLY")
	Expect(err).ToNot(HaveOccurred())
})

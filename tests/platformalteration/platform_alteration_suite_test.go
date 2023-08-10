//go:build !utest

package platformalteration

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/tests"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

func TestPlatformAlteration(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert platform-alteration tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	By("Create namespace")
	err := namespaces.Create(tsparams.PlatformAlterationNamespace, globalhelper.GetAPIClient())
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.PlatformAlterationNamespace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", tsparams.PlatformAlterationNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.GetAPIClient(),
		tsparams.PlatformAlterationNamespace,
		tsparams.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())

})

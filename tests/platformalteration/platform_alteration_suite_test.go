package platformalteration

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"

	"github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/platformalterationparameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/tests"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

func TestPlatformAlteration(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert platform-alteration tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	By("Create namespace")
	err := namespaces.Create(platformalterationparameters.PlatformAlterationNamespace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{platformalterationparameters.PlatformAlterationNamespace},
		[]string{platformalterationparameters.TestPodLabel},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", platformalterationparameters.PlatformAlterationNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		platformalterationparameters.PlatformAlterationNamespace,
		platformalterationparameters.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())

})

package affiliated_certification

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	"fmt"

	"runtime"
	"testing"

	_ "github.com/test-network-function/cnfcert-tests-verification/tests/affiliated_certification/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
)

func TestAffiliatedCertification(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	configSuite, err := config.NewConfig()

	if err != nil {
		fmt.Print(err)
		return
	}

	junitPath := configSuite.GetReportPath(currentFile)
	RegisterFailHandler(Fail)
	rr := append([]Reporter{}, reporters.NewJUnitReporter(junitPath))
	RunSpecsWithDefaultAndCustomReporters(t, "CNFCert affiliated-certification tests", rr)
}

var _ = BeforeSuite(func() {

})

var _ = AfterSuite(func() {

})

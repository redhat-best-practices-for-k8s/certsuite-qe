package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
)

var _ = Describe("platform-alteration-ocp-node-os", func() {

	It("Nodes OS should be compatible with OCP version", func() {
		By("Start platform-alteration-ocp-node-os test")
		err := globalhelper.LaunchTests(tsparams.TnfOcpNodeOsName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfOcpNodeOsName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

})

package tests

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-crd-scaling", func() {
	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.LifecycleNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	FIt("Crd deployed, scale in and out", func() {
		// Install crd operator
		By("Install crd-scaling-operator")
		err := installScaleOperator()
		Expect(err).ToNot(HaveOccurred(), "Error installing crd-scaling-operator")

		By("Create a scale custom resource")

		By("Start lifecycle-crd-scaling test")

		err = globalhelper.LaunchTests(tsparams.TnfCrdScaling,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdScaling, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})
})

func installScaleOperator() error {
	command := "git clone https://github.com/bnshr/cnf-operator && cd cnf-operator && " +
		"export IMG=quay.io/rh_ee_bmandal/cnf-collector:latest && export NAMESPACE=bmandal && make deploy"

	// commandStr := fmt.Sprintf("rm -rf cnf-operator && %s", command)

	fmt.Println(command)

	cmd := exec.Command("/bin/bash", "-c", command)
	err := cmd.Run()

	return err
}

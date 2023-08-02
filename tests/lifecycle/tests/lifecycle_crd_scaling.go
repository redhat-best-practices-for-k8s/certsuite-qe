package tests

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

const crdNameSpace = "qe-test" // create this manually

var _ = Describe("lifecycle-crd-scaling", func() {
	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	FIt("Crd deployed, scale in and out", func() {

		By("To be implemented")

		// Install crd operator
		By("Install crd-scaling-operator")

		command := "git clone https://github.com/bnshr/cnf-operator && cd cnf-operator && " +
			"export IMG=quay.io/rh_ee_bmandal/cnf-collector:latest && export NAMESPACE=bmandal && make deploy"

		// commandStr := fmt.Sprintf("rm -rf cnf-operator && %s", command)

		fmt.Println(command)

		cmd := exec.Command("/bin/bash", "-c", command)
		err := cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing crd-scaling-operator")
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})
})

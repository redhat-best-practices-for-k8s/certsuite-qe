package tests

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Affiliated-certification helm chart certification,", func() {
	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.TestHelmChartCertified},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})
	
	BeforeEach(func() {
		By("Create namespace")
		err := namespaces.Create(tsparams.TestHelmChartCertified, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")
	})

	AfterEach(func() {
		By("remove the project")
		err := namespaces.DeleteAndWait(globalhelper.GetAPIClient(), tsparams.TestHelmChartCertified, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred(), "Error delete ns affiliated-certification-helmchart-is-certified")
	})

	It("one helm to test,  are certified", func() {
		By("check if helm is installed")
		cmd := exec.Command("/bin/bash", "-c",
			"helm version")
		err := cmd.Run()
		if err != nil {
			Skip("helm does not exist please install it to run the test.")
		}

		By("Install a helm chart")
		cmd = exec.Command("/bin/bash", "-c",
			"helm repo add openshift-helm-charts https://charts.openshift.io/ "+
				"&& helm repo update && "+
				"helm install example-vault1 openshift-helm-charts/hashicorp-vault -n "+tsparams.TestHelmChartCertified)
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm chart")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestHelmChartCertified,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestHelmChartCertified+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestHelmChartCertified,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one helm to test, chart not certified", func() {
		By("check if helm is installed")
		cmd := exec.Command("/bin/bash", "-c",
			"helm version")
		err := cmd.Run()
		if err != nil {
			Skip("helm does not exist please install it to run the test.")
		}

		By("Create ns istio-system")
		err = namespaces.Create("istio-system", globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Install a helm chart")
		cmd = exec.Command("/bin/bash", "-c",
			"helm repo add istio https://istio-release.storage.googleapis.com/charts "+
				"&& helm repo update &&"+
				"helm install istio-base istio/base --set defaultRevision=default -n "+tsparams.TestHelmChartCertified)
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm chart")

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestHelmChartCertified,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred(), "Error running "+
			tsparams.TestHelmChartCertified+" test")

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestHelmChartCertified,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

		By("remove the istio-system ns and istio chart")
		cmd = exec.Command("/bin/bash", "-c", // uinstall the chart
			"helm uninstall istio-base -n "+tsparams.TestHelmChartCertified)
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm chart")
		err = namespaces.Clean("istio-system", globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error delete ns istio-system")
	})
})

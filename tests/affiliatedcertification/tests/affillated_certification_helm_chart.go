package tests

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
)

var _ = Describe("Affiliated-certification helm chart certification,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestCertificationNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})

	It("One helm to test, are certified", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skipping helm chart test on Kind cluster")
		}

		By("Check if helm is installed")
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
				"helm install example-vault1 openshift-helm-charts/hashicorp-vault -n "+randomNamespace)
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

	It("One helm to test, chart not certified", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skipping helm chart test on Kind cluster")
		}

		By("Check if helm is installed")
		cmd := exec.Command("/bin/bash", "-c",
			"helm version")
		err := cmd.Run()
		if err != nil {
			Skip("helm does not exist please install it to run the test.")
		}

		By("Create ns istio-system")
		err = globalhelper.CreateNamespace("istio-system")
		Expect(err).ToNot(HaveOccurred())

		By("Install a helm chart")
		cmd = exec.Command("/bin/bash", "-c",
			"helm repo add istio https://istio-release.storage.googleapis.com/charts "+
				"&& helm repo update &&"+
				"helm install istio-base istio/base --set defaultRevision=default -n "+randomNamespace)
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

		By("Remove the istio-system ns and istio chart")
		cmd = exec.Command("/bin/bash", "-c", // uinstall the chart
			"helm uninstall istio-base -n "+randomNamespace)
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm chart")
		err = globalhelper.CleanNamespace("istio-system")
		Expect(err).ToNot(HaveOccurred(), "Error delete ns istio-system")
	})
})

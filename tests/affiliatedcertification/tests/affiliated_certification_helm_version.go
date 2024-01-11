package tests

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
)

var _ = Describe("Affiliated-certification helm-version,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestCertificationNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.Timeout)
	})

	// 68120
	It("installed helm version is certified", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skipping helm version test on Kind cluster")
		}

		By("Check if helm is installed")
		cmd := exec.Command("/bin/bash", "-c",
			"helm version")
		err := cmd.Run()
		if err != nil {
			Skip("helm does not exist please install it to run the test.")
		}

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestHelmVersion,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestHelmVersion+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestHelmVersion,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 68121
	It("installed helm version is not certified", func() {
		if globalhelper.IsKindCluster() {
			Skip("Skipping helm version test on Kind cluster")
		}

		By("Remove helm")
		cmd := exec.Command("sudo", "rm", "-rf", "/usr/local/bin/helm")
		err := cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error uninstalling helm")

		By("Install helm v2")
		cmd = exec.Command("/bin/bash", "-c",
			"curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get"+
				" && chmod +x get_helm.sh"+
				" && ./get_helm.sh --version v2.17.0"+
				" && helm init")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm v2")

		DeferCleanup(func() {
			By("Delete tiller from kube-system namespace")
			err := globalhelper.DeleteDeployment("tiller-deploy", "kube-system")
			Expect(err).ToNot(HaveOccurred(), "Error deleting tiller deployment")

			By("Re-install helm v3")
			cmd = exec.Command("/bin/bash", "-c",
				"curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3"+
					" && chmod +x get_helm.sh"+
					" && ./get_helm.sh")
			err = cmd.Run()
			Expect(err).ToNot(HaveOccurred(), "Error installing helm v3")
		})

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestHelmVersion,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestHelmVersion+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestHelmVersion,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

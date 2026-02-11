package tests

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
)

var _ = Describe("Affiliated-certification helm-version,", Serial, Label("affiliatedcertification", "ocp-required"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

	BeforeEach(func() {
		if globalhelper.IsKindCluster() {
			Skip("Skipping helm version test on Kind cluster")
		}

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestCertificationNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 68120
	It("installed helm version is certified", func() {
		By("Check if helm is installed")
		cmd := exec.Command("/bin/bash", "-c",
			"helm version")

		err := cmd.Run()
		if err != nil {
			Skip("helm does not exist please install it to run the test.")
		}

		By("Check that helm version is v3")
		cmd = exec.Command("/bin/bash", "-c",
			"helm version --short | grep v3")

		err = cmd.Run()
		if err != nil {
			Fail("Helm version is not v3")
		}

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestHelmVersion,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestHelmVersion+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestHelmVersion,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

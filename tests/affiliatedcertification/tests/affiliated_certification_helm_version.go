package tests

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
)

var _ = Describe("Affiliated-certification helm-version,", Ordered, Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeAll(func() {
		if globalhelper.IsKindCluster() {
			Skip("Skipping helm version test on Kind cluster")
		}

		By("Install helm v2 to /usr/local/bin/helm2")
		cmd := exec.Command("/bin/bash", "-c",
			"curl -fsSL -o get_helm_v2.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get"+
				" && chmod +x get_helm_v2.sh"+
				" && HELM_INSTALL_DIR=/usr/local/bin ./get_helm_v2.sh --version v2.17.0 --no-sudo"+
				" && mv /usr/local/bin/helm /usr/local/bin/helm2")
		err := cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm v2")

		By("Install helm v3 to /usr/local/bin/helm3")
		cmd = exec.Command("/bin/bash", "-c",
			"curl -fsSL -o get_helm_v3.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3"+
				" && chmod +x get_helm_v3.sh"+
				" && HELM_INSTALL_DIR=/usr/local/bin ./get_helm_v3.sh --no-sudo"+
				" && mv /usr/local/bin/helm /usr/local/bin/helm3")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm v3")

		By("Create default helm symlink to v3")
		cmd = exec.Command("/bin/bash", "-c", "ln -sf /usr/local/bin/helm3 /usr/local/bin/helm")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error creating helm symlink")
	})

	AfterAll(func() {
		By("Cleanup helm installations")
		_ = exec.Command("sudo", "rm", "-f", "/usr/local/bin/helm", "/usr/local/bin/helm2", "/usr/local/bin/helm3").Run()
		_ = exec.Command("rm", "-f", "get_helm_v2.sh", "get_helm_v3.sh").Run()
	})

	BeforeEach(func() {
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
		By("Ensure helm points to v3")
		cmd := exec.Command("/bin/bash", "-c", "ln -sf /usr/local/bin/helm3 /usr/local/bin/helm")
		err := cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error setting helm to v3")

		By("Verify helm version is v3")
		cmd = exec.Command("/bin/bash", "-c", "helm version --short | grep v3")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Helm is not v3")

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
		By("Switch helm to v2")
		cmd := exec.Command("/bin/bash", "-c", "ln -sf /usr/local/bin/helm2 /usr/local/bin/helm")
		err := cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error setting helm to v2")

		By("Initialize tiller for helm v2")
		cmd = exec.Command("/bin/bash", "-c", "helm init")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error initializing helm v2")

		DeferCleanup(func() {
			By("Delete tiller from kube-system namespace")
			err := globalhelper.DeleteDeployment("tiller-deploy", "kube-system")
			Expect(err).ToNot(HaveOccurred(), "Error deleting tiller deployment")

			By("Restore helm to v3")
			cmd = exec.Command("/bin/bash", "-c", "ln -sf /usr/local/bin/helm3 /usr/local/bin/helm")
			err = cmd.Run()
			Expect(err).ToNot(HaveOccurred(), "Error restoring helm to v3")
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

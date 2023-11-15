package tests

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
)

var _ = Describe("Affiliated-certification helm chart certification,", Serial, func() {
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

		By("Add openshift-helm-charts repo")
		cmd = exec.Command("/bin/bash", "-c",
			"helm repo add openshift-helm-charts https://charts.openshift.io/ --force-update "+
				"&& helm repo update")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error adding openshift-helm-carts repo")

		By("Install helm chart")
		cmd = exec.Command("/bin/bash", "-c",
			"helm install example-vault1 openshift-helm-charts/hashicorp-vault -n "+randomNamespace)
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing hashicorp-vault helm chart")

		DeferCleanup(func() {
			By("Remove the example-vault1 helm chart")
			cmd = exec.Command("/bin/bash", "-c", // uninstall the chart
				"helm uninstall example-vault1 --ignore-not-found -n "+randomNamespace)
			err = cmd.Run()
			Expect(err).ToNot(HaveOccurred(), "Error uninstalling helm chart")

			By("Delete clusterrole and clusterrolebinding")
			err = globalhelper.DeleteClusterRoleBindingByName("example-vault1-agent-injector-binding")
			Expect(err).ToNot(HaveOccurred(), "Error deleting clusterrolebinding")
			err = globalhelper.DeleteClusterRoleBindingByName("example-vault1-server-binding")
			Expect(err).ToNot(HaveOccurred(), "Error deleting clusterrolebinding")

			err = globalhelper.DeleteClusterRole("example-vault1-agent-injector-clusterrole")
			Expect(err).ToNot(HaveOccurred(), "Error deleting clusterrole")

			By("Delete mutatingwebhookconfiguration")
			err = globalhelper.DeleteMutatingWebhookConfiguration("example-vault1-agent-injector-cfg")
			Expect(err).ToNot(HaveOccurred(), "Error deleting mutatingwebhookconfiguration")
		})

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestHelmChartCertified,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred(), "Error running "+
			tsparams.TestHelmChartCertified+" test")

		By("Verify test case status in Claim report")
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

		By("Create istio-system namespace")
		err = globalhelper.CreateNamespace("istio-system")
		Expect(err).ToNot(HaveOccurred(), "Error creating istio-system namespace")

		By("Add istio helm chart repo")
		cmd = exec.Command("/bin/bash", "-c",
			"helm repo add istio https://istio-release.storage.googleapis.com/charts --force-update "+
				"&& helm repo update")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error adding istio charts repo")

		By("Install istio-base helm chart")
		cmd = exec.Command("/bin/bash", "-c",
			"helm install istio-base istio/base --set defaultRevision=default -n "+randomNamespace)
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing istio-base helm chart")

		DeferCleanup(func() {
			By("Remove istio-base helm chart")
			cmd = exec.Command("/bin/bash", "-c", // uninstall the chart
				"helm uninstall istio-base --ignore-not-found -n "+randomNamespace)
			err = cmd.Run()
			Expect(err).ToNot(HaveOccurred(), "Error uninstalling helm chart")

			By("Delete istio-system namespace")
			err = globalhelper.DeleteNamespaceAndWait("istio-system", tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting istio-system namespace")

			By("Remove validating webhook configuration")
			err = globalhelper.DeleteValidatingWebhookConfiguration("istiod-default-validator")
			Expect(err).ToNot(HaveOccurred(), "Error deleting validating webhook configuration")
		})

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestHelmChartCertified,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred(), "Error running "+
			tsparams.TestHelmChartCertified+" test")

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestHelmChartCertified,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

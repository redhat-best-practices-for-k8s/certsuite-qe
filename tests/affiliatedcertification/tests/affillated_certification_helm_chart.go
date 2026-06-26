package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
)

var _ = Describe("Affiliated-certification helm chart certification,", Serial,
	Label("affiliatedcertification", "ocp-required"), func() {
		var (
			randomNamespace          string
			randomReportDir          string
			randomCertsuiteConfigDir string
		)

		BeforeEach(func() {
			if globalhelper.IsKindCluster() {
				Skip("Skipping helm chart test on Kind cluster")
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

		It("One helm to test, are certified", func() {
			By("Check if helm is installed")
			_, err := globalhelper.RunShellCommand("helm version")

			if err != nil {
				Skip("helm does not exist please install it to run the test.")
			}

			By("Check that helm version is v3")
			_, err = globalhelper.RunShellCommand("helm version --short | grep v3")

			if err != nil {
				Fail("Helm version is not v3")
			}

			By("Add openshift-helm-charts repo")
			_, err = globalhelper.RunShellCommand(
				"helm repo add hashicorp https://helm.releases.hashicorp.com --force-update" +
					" && helm repo update")
			Expect(err).ToNot(HaveOccurred(), "Error adding openshift-helm-carts repo")

			By("Install helm chart")
			_, err = globalhelper.RunShellCommand(
				"helm install example-vault1 hashicorp/vault -n " + randomNamespace)
			Expect(err).ToNot(HaveOccurred(), "Error installing hashicorp-vault helm chart")

			DeferCleanup(func() {
				By("Remove the example-vault1 helm chart")
				_, err = globalhelper.RunShellCommand(
					"helm uninstall example-vault1 --ignore-not-found -n " + randomNamespace)
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
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				tsparams.TestHelmChartCertified+" test")

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestHelmChartCertified,
				globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("One helm to test, chart not certified", func() {
			By("Check if helm is installed")
			_, err := globalhelper.RunShellCommand("helm version")

			if err != nil {
				Skip("helm does not exist please install it to run the test.")
			}

			By("Check that helm version is v3")
			_, err = globalhelper.RunShellCommand("helm version --short | grep v3")

			if err != nil {
				Fail("Helm version is not v3")
			}

			By("Delete istio-system namespace")
			err = globalhelper.DeleteNamespaceAndWait("istio-system", tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred(), "Error deleting istio-system namespace")

			By("Create istio-system namespace")
			err = globalhelper.CreateNamespace("istio-system")
			Expect(err).ToNot(HaveOccurred(), "Error creating istio-system namespace")

			By("Add istio helm chart repo")
			_, err = globalhelper.RunShellCommand(
				"helm repo add istio https://istio-release.storage.googleapis.com/charts --force-update" +
					" && helm repo update")
			Expect(err).ToNot(HaveOccurred(), "Error adding istio charts repo")

			DeferCleanup(func() {
				By("Remove istio-base helm chart")
				_, err := globalhelper.RunShellCommand(
					"helm uninstall istio-base --ignore-not-found -n " + randomNamespace)
				Expect(err).ToNot(HaveOccurred(), "Error uninstalling helm chart")

				By("Delete istio-system namespace")
				err = globalhelper.DeleteNamespaceAndWait("istio-system", tsparams.Timeout)
				Expect(err).ToNot(HaveOccurred(), "Error deleting istio-system namespace")

				By("Remove validating webhook configuration")
				err = globalhelper.DeleteValidatingWebhookConfiguration("istiod-default-validator")
				Expect(err).ToNot(HaveOccurred(), "Error deleting validating webhook configuration")

				// BugFix: see https://github.com/istio/istio/issues/43204 , which suggests removing CRDs:
				// https://istio.io/latest/docs/setup/install/helm/#optional-deleting-crds-installed-by-istio
				By("Delete istio CRDs")
				err = tshelper.DeleteIstioCRDs()
				Expect(err).ToNot(HaveOccurred(), "Error deleting istio CRDs.")
			})

			By("Install istio-base helm chart")
			_, err = globalhelper.RunShellCommand(
				"helm install istio-base istio/base --set defaultRevision=default -n " + randomNamespace +
					" --set hub=gcr.io/istio-release")
			Expect(err).ToNot(HaveOccurred(), "Error installing istio-base helm chart")

			By("Start test")
			err = globalhelper.LaunchTests(
				tsparams.TestHelmChartCertified,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "Error running "+
				tsparams.TestHelmChartCertified+" test")

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestHelmChartCertified,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})
	})

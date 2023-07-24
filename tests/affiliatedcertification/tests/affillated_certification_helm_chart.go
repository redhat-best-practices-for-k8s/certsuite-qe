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
		By("Installing helm")
		cmd := exec.Command("/bin/bash", "-c",
			"curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3"+
				"&& chmod 700 get_helm.sh"+
				"&& ./get_helm.sh")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm")
	})

	AfterEach(func() {
		By("remove the project")
		err := namespaces.Clean("affiliated-certification-helmchart-is-certified", globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error delete ns affiliated-certification-helmchart-is-certified")
	})

	AfterAll(func() {
		By("remove the istio-system")
		err := namespaces.Clean("istio-system", globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error delete ns istio-system")
	})

	It("one helm to test,  are certified", func() {
		By("Install a helm chart")
		err := namespaces.Create("affiliated-certification-helmchart-is-certified", globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")
		cmd := exec.Command("/bin/bash", "-c",
			"helm repo add openshift-helm-charts https://charts.openshift.io/ "+
				"&& helm repo update && helm install example-vault1 openshift-helm-charts/hashicorp-vault")
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
		By("Create ns istio-system")
		err := namespaces.Create("istio-system", globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
		err = namespaces.Create("affiliated-certification-helmchart-is-certified", globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

		By("Install a helm chart")
		cmd := exec.Command("/bin/bash", "-c",
			"helm repo add istio https://istio-release.storage.googleapis.com/charts "+
				"&& helm repo update && helm install istio-base istio/base --set defaultRevision=default")
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
	})
})

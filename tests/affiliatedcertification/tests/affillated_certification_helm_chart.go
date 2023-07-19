package tests

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
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

	FIt("one helm to test,  are certified", func() {
		By("Install a hellm chart")
		cmd := exec.Command("/bin/bash", "-c", "oc new-project affiliated-certification-helmchart-is-certified"+
			"&& helm repo add openshift-helm-charts https://charts.openshift.io/ "+
			"&& helm repo update && helm install example-vault1 openshift-helm-charts/hashicorp-vault")
		err := cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing hellm chart")

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
		By("remove the project")
		cmd = exec.Command("/bin/bash", "-c", "oc delete ns affiliated-certification-helmchart-is-certified")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error delete ns affiliated-certification-helmchart-is-certified")
	})

	FIt("one helm to test,  are not certified", func() {
		By("Install a hellm chart")
		cmd := exec.Command("/bin/bash", "-c", "oc new-project affiliated-certification-helmchart-is-certified"+
			"&& helm repo add istio https://istio-release.storage.googleapis.com/charts "+
			"&& helm repo update && oc create ns istio-system && helm install istio-base istio/base --set defaultRevision=default")
		err := cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing hellm chart")

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
		By("remove the project")
		cmd = exec.Command("/bin/bash", "-c", "oc delete ns affiliated-certification-helmchart-is-certified")
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error delete ns affiliated-certification-helmchart-is-certified")
	})
})

//go:build !utest

package affiliatedcertification

import (
	"flag"
	"os/exec"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"

	"runtime"
	"testing"

	_ "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/tests"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
)

func TestAffiliatedCertification(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.GetConfiguration().General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.GetConfiguration().GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert affiliated-certification tests", reporterConfig)
}

var isCloudCasaAlreadyLabeled bool

var _ = SynchronizedBeforeSuite(func() {
	// Always install Helm v3 right before running the suite
	By("Install helm v3")
	cmd := exec.Command("/bin/bash", "-c",
		"curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3"+
			" && chmod +x get_helm.sh"+
			" && ./get_helm.sh")
	err := cmd.Run()
	Expect(err).ToNot(HaveOccurred(), "Error installing helm v3")

	By("Preemptively delete tiller-deploy pod if its installed")
	err = globalhelper.DeleteDeployment("tiller-deploy", "kube-system")
	Expect(err).ToNot(HaveOccurred(), "Error deleting tiller deployment")

	By("Preemptively delete clusterrole and clusterrolebinding")
	err = globalhelper.DeleteClusterRoleBindingByName("example-vault1-agent-injector-binding")
	Expect(err).ToNot(HaveOccurred(), "Error deleting clusterrolebinding")
	err = globalhelper.DeleteClusterRoleBindingByName("example-vault1-server-binding")
	Expect(err).ToNot(HaveOccurred(), "Error deleting clusterrolebinding")
	err = globalhelper.DeleteClusterRole("example-vault1-agent-injector-clusterrole")
	Expect(err).ToNot(HaveOccurred(), "Error deleting clusterrole")

	By("Delete mutatingwebhookconfiguration")
	err = globalhelper.DeleteMutatingWebhookConfiguration("example-vault1-agent-injector-cfg")
	Expect(err).ToNot(HaveOccurred(), "Error deleting mutatingwebhookconfiguration")

	By("Remove validating webhook configuration")
	err = globalhelper.DeleteValidatingWebhookConfiguration("istiod-default-validator")
	Expect(err).ToNot(HaveOccurred(), "Error deleting validating webhook configuration")

	By("Create namespace")
	err = globalhelper.CreateNamespace(tsparams.TestCertificationNameSpace)
	Expect(err).ToNot(HaveOccurred(), "Error creating namespace")

	isCloudCasaAlreadyLabeled, err = tshelper.DoesOperatorHaveLabels(tsparams.UnrelatedOperatorPrefixCloudcasa,
		tsparams.UnrelatedNamespace,
		tsparams.OperatorLabel)
	if err != nil {
		glog.Info(tsparams.UnrelatedOperatorPrefixCloudcasa+" not installed or error accessing it: ", err)
	}

	By("Un-label operator used in other suites if labeled")
	if isCloudCasaAlreadyLabeled {
		err = tshelper.DeleteLabelFromInstalledCSV(
			tsparams.UnrelatedOperatorPrefixCloudcasa,
			tsparams.UnrelatedNamespace,
			tsparams.OperatorLabel)
		Expect(err).ToNot(HaveOccurred())
	}

}, func() {})

var _ = SynchronizedAfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", tsparams.TestCertificationNameSpace))
	err := globalhelper.DeleteNamespaceAndWait(tsparams.TestCertificationNameSpace, tsparams.Timeout)
	Expect(err).ToNot(HaveOccurred())

	if isCloudCasaAlreadyLabeled {
		By("Re-label operator used in other suites")
		err = tshelper.AddLabelToInstalledCSV(
			tsparams.UnrelatedOperatorPrefixCloudcasa,
			tsparams.UnrelatedNamespace,
			tsparams.OperatorLabel)
		Expect(err).ToNot(HaveOccurred())
	}
}, func() {})

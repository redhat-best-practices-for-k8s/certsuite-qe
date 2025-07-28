//go:build !utest

package affiliatedcertification

import (
	"flag"
	"fmt"
	"os/exec"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"runtime"
	"testing"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/tests"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
	ophelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
	utils "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/operator"
)

// Suite-level shared variables for operators.
var (
	isCloudCasaAlreadyLabeled bool
	sharedOperatorNamespace   string
	grafanaOperatorName       string
	grafanaChannel            string
	grafanaVersion            string
	grafanaCSVName            string
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

var _ = SynchronizedBeforeSuite(func() {

	if !globalhelper.IsKindCluster() {
		// Always install Helm v3 right before running the suite
		By("Install helm v3")
		cmd := exec.Command("/bin/bash", "-c",
			"curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3"+
				" && chmod +x get_helm.sh"+
				" && ./get_helm.sh")
		out, err := cmd.CombinedOutput()
		Expect(err).ToNot(HaveOccurred(), "Error installing helm v3: "+string(out))
	}

	By("Preemptively delete tiller-deploy pod if its installed")
	err := globalhelper.DeleteDeployment("tiller-deploy", "kube-system")
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

	if !globalhelper.IsKindCluster() {
		By("Create community-operators catalog source")
		err = globalhelper.CreateCommunityOperatorsCatalogSource()
		Expect(err).ToNot(HaveOccurred())

		By("Create certified-operators catalog source")
		err = globalhelper.DeployRHCertifiedOperatorSource("")
		Expect(err).ToNot(HaveOccurred())

		By("Check if catalog sources are available")
		err = globalhelper.ValidateCatalogSources()
		Expect(err).ToNot(HaveOccurred(), "All necessary catalog sources are not available")

		// ==========================================================================
		// 🚀 SUITE-LEVEL OPERATOR DEPLOYMENT FOR PERFORMANCE OPTIMIZATION
		// ==========================================================================
		By("🚀 SUITE OPTIMIZATION: Deploy shared operators for all tests")

		// Create dedicated namespace for shared operators
		sharedOperatorNamespace = tsparams.TestCertificationNameSpace + "-operators"
		By(fmt.Sprintf("Create shared operator namespace: %s", sharedOperatorNamespace))
		err = globalhelper.CreateNamespace(sharedOperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error creating shared operator namespace")

		By("Setup shared operator environment")
		setupSharedOperatorEnvironment(sharedOperatorNamespace)

		By("Deploy cockroachdb operator (uncertified) for shared use")
		err = ophelper.DeployOperatorSubscription(
			"cockroachdb",
			"cockroachdb",
			"stable-v6.x",
			sharedOperatorNamespace,
			tsparams.CommunityOperatorGroup,
			tsparams.OperatorSourceNamespace,
			"",
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying shared cockroachdb operator")

		By("Wait for cockroachdb operator to be ready")
		err = ophelper.WaitUntilOperatorIsReady(tsparams.UncertifiedOperatorPrefixCockroach,
			sharedOperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Shared cockroachdb operator is not ready")

		By("Query packagemanifest for cockroachdb-certified operator")
		channel, err := globalhelper.QueryPackageManifestForDefaultChannel(
			"cockroachdb-certified",
			sharedOperatorNamespace,
		)
		Expect(err).ToNot(HaveOccurred(), "Error querying cockroachdb-certified manifest")

		version, err := globalhelper.QueryPackageManifestForVersion("cockroachdb-certified",
			sharedOperatorNamespace, channel)
		Expect(err).ToNot(HaveOccurred(), "Error querying cockroachdb-certified version")

		By(fmt.Sprintf("Deploy cockroachdb-certified operator %s for shared use", "v"+version))
		err = ophelper.DeployOperatorSubscription(
			"cockroachdb-certified",
			"cockroachdb-certified",
			channel,
			sharedOperatorNamespace,
			tsparams.CertifiedOperatorGroup,
			tsparams.OperatorSourceNamespace,
			tsparams.CertifiedOperatorPrefixCockroachCertified+".v"+version,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying shared cockroachdb-certified operator")

		By("Wait for cockroachdb-certified operator to be ready")
		err = ophelper.WaitUntilOperatorIsReady(tsparams.CertifiedOperatorPrefixCockroachCertified,
			sharedOperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Shared cockroachdb-certified operator is not ready")

		By("Query packagemanifest for Grafana operator")
		var catalogSource string
		grafanaOperatorName, catalogSource = globalhelper.CheckOperatorExistsOrFail("grafana", sharedOperatorNamespace)
		grafanaChannel, grafanaVersion, grafanaCSVName = globalhelper.CheckOperatorChannelAndVersionOrFail(
			grafanaOperatorName, sharedOperatorNamespace)

		By(fmt.Sprintf("Deploy Grafana operator (%s, %s) for shared use", grafanaChannel, grafanaVersion))
		err = ophelper.DeployOperatorSubscription(
			grafanaOperatorName,
			grafanaOperatorName,
			grafanaChannel,
			sharedOperatorNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			grafanaCSVName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying shared Grafana operator")

		By("Wait for Grafana operator to be ready")
		err = ophelper.WaitUntilOperatorIsReady(grafanaOperatorName,
			sharedOperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Shared Grafana operator is not ready")

		By("✅ Suite-level operators deployment completed successfully!")
	}
}, func() {})

// setupSharedOperatorEnvironment configures the shared operator namespace.
func setupSharedOperatorEnvironment(namespace string) {
	By("Clean shared operator namespace")
	err := globalhelper.CleanNamespace(namespace)
	Expect(err).ToNot(HaveOccurred(), "Error cleaning shared operator namespace")

	By("Ensure certified catalog source is enabled")
	catalogEnabled, err := globalhelper.IsCatalogSourceEnabled(
		tsparams.CertifiedOperatorGroup,
		tsparams.OperatorSourceNamespace,
		tsparams.CertifiedOperatorDisplayName)
	Expect(err).ToNot(HaveOccurred(), "Cannot collect catalogSource object")

	if !catalogEnabled {
		Expect(globalhelper.EnableCatalogSource(tsparams.CertifiedOperatorGroup)).ToNot(HaveOccurred())
		Eventually(func() bool {
			catalogEnabled, err = globalhelper.IsCatalogSourceEnabled(
				tsparams.CertifiedOperatorGroup,
				tsparams.OperatorSourceNamespace,
				tsparams.CertifiedOperatorDisplayName)
			Expect(err).ToNot(HaveOccurred())

			return catalogEnabled
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(BeTrue(),
			"Certified catalog source is not enabled")
	}

	By("Deploy OperatorGroup in shared namespace")

	if globalhelper.IsOperatorGroupInstalled(tsparams.OperatorGroupName, namespace) != nil {
		err = globalhelper.DeployOperatorGroup(namespace,
			utils.DefineOperatorGroup(tsparams.OperatorGroupName,
				namespace,
				[]string{namespace}),
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying shared operatorgroup")
	}
}

// GetSharedOperatorNamespace returns the shared operator namespace for tests to use.
func GetSharedOperatorNamespace() string {
	return sharedOperatorNamespace
}

// GetGrafanaOperatorName returns the grafana operator name for tests to use.
func GetGrafanaOperatorName() string {
	return grafanaOperatorName
}

var _ = SynchronizedAfterSuite(func() {}, func() {
	By(fmt.Sprintf("Remove %s namespace", tsparams.TestCertificationNameSpace))
	err := globalhelper.DeleteNamespaceAndWait(tsparams.TestCertificationNameSpace, tsparams.Timeout)
	Expect(err).ToNot(HaveOccurred())

	// Clean up shared operator namespace
	if sharedOperatorNamespace != "" {
		By(fmt.Sprintf("🧹 Clean up shared operator namespace: %s", sharedOperatorNamespace))
		err = globalhelper.DeleteNamespaceAndWait(sharedOperatorNamespace, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())
	}

	if isCloudCasaAlreadyLabeled {
		By("Re-label operator used in other suites")
		err = tshelper.AddLabelToInstalledCSV(
			tsparams.UnrelatedOperatorPrefixCloudcasa,
			tsparams.UnrelatedNamespace,
			tsparams.OperatorLabel)
		Expect(err).ToNot(HaveOccurred())
	}
})

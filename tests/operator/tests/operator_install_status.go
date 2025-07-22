package operator

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/golang/glog"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
)

var _ = Describe("Operator install-source,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string
	var operatorName string
	var catalogSource string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.OperatorNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			tsparams.CertsuiteTargetCrdFilters, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		By("Query the packagemanifest for grafana operator package name and catalog source")
		operatorName, catalogSource = globalhelper.CheckOperatorExistsOrFail("grafana", randomNamespace)

		By("Query the packagemanifest for available channel, version and CSV for " + operatorName)
		channel, _, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(operatorName, randomNamespace)

		By("Deploy grafana operator for testing")
		err = tshelper.DeployOperatorSubscription(
			operatorName,
			operatorName,
			channel,
			randomNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			operatorName)

		err = tshelper.WaitUntilOperatorIsReady(operatorName,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("one operator that reports Succeeded as its installation status", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatus,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("two operators, one does not reports Succeeded as its installation status (quick failure) [negative]", func() {
		By("Query the packagemanifest for postgresql operator package name and catalog source")
		postgresOperatorName, catalogSource, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			"cloud-native-postgresql", randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for postgresql operator")
		Expect(postgresOperatorName).ToNot(Equal("not found"), "PostgreSQL operator package not found")
		Expect(catalogSource).ToNot(Equal("not found"), "PostgreSQL operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + postgresOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			postgresOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+postgresOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		By("Deploy postgresql operator for testing")
		// Deploy PostgreSQL operator with nodeSelector that will cause quick failure
		nodeSelector := map[string]string{"target": "nonexistent-node"}
		err = tshelper.DeployOperatorSubscriptionWithNodeSelector(
			postgresOperatorName,
			channel,
			randomNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
			nodeSelector,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			postgresOperatorName)

		// Do not wait for the PostgreSQL operator to be ready - it should fail due to nodeSelector

		By("Verify that PostgreSQL operator CSV is not in Succeeded phase")
		Eventually(func() bool {
			isNotSucceeded, err := tshelper.IsCSVNotSucceeded(postgresOperatorName, randomNamespace)
			if err != nil {
				fmt.Printf("Error checking CSV status for %s: %v\n", postgresOperatorName, err)

				return false
			}
			fmt.Printf("PostgreSQL operator %s CSV status is not Succeeded: %t\n", postgresOperatorName, isNotSucceeded)

			return isNotSucceeded
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Equal(true),
			"PostgreSQL operator CSV should not be in Succeeded phase for this negative test")

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				postgresOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+postgresOperatorName)

		By("Update certsuite config to include both operators")
		err = globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			tsparams.CertsuiteTargetCrdFilters, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatus,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("two operators, one does not reports Succeeded as its installation status (delayed failure) [negative]", Serial, func() {
		By("Query the packagemanifest for Jaeger operator package name and catalog source")
		jaegerOperatorName, catalogSource2, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			"jaeger", randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for Jaeger operator")
		Expect(jaegerOperatorName).ToNot(Equal("not found"), "Jaeger operator package not found")
		Expect(catalogSource2).ToNot(Equal("not found"), "Jaeger operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + jaegerOperatorName)
		channel, version, csvName, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			jaegerOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+jaegerOperatorName)
		Expect(channel).ToNot(Equal("not found"), "Channel not found")
		Expect(version).ToNot(Equal("not found"), "Version not found")
		Expect(csvName).ToNot(Equal("not found"), "CSV name not found")

		By("Deploy Jaeger operator for testing")
		// The jaeger operator fails to deploy, which creates a delayed failure scenario
		// This allows testing of the CNF Certification Suite timeout mechanism
		// for operator readiness.
		//
		// NOTE: In OCP 4.16 and below, nodeSelector constraints may be enforced differently.
		// We use multiple restrictive constraints to ensure the operator installation fails.
		// Use a more restrictive nodeSelector that is guaranteed to fail
		nodeSelector := map[string]string{
			"kubernetes.io/arch":                     "nonexistent-arch-x999",
			"node.kubernetes.io/unreachable":         "true",
			"node.kubernetes.io/network-unavailable": "true",
			"test.certsuite.fail/operator":           "guaranteed-fail",
		}

		// For OCP 4.16 and below, use additional constraints to ensure failure
		ocpVersion, versionErr := globalhelper.GetClusterVersion()
		if versionErr == nil {
			glog.V(5).Infof("Detected OCP version: %s", ocpVersion)
			// Add even more restrictive constraints for older versions
			nodeSelector["kubernetes.io/os"] = "nonexistent-os-fail"
			nodeSelector["kubernetes.io/hostname"] = "never-exists-hostname-fail"
			nodeSelector["topology.kubernetes.io/zone"] = "fail-zone-nonexistent"
		}

		glog.V(5).Infof("Deploying Jaeger operator with nodeSelector constraints: %v", nodeSelector)
		err = tshelper.DeployOperatorSubscriptionWithNodeSelector(
			jaegerOperatorName,
			channel,
			randomNamespace,
			catalogSource2,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
			nodeSelector,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			jaegerOperatorName)

		// Do not wait until the operator is ready. This time the CNF Certification suite must handle the situation.

		By("Verify that Jaeger operator CSV is not in Succeeded phase")
		glog.V(5).Infof("Starting verification that Jaeger operator %s CSV is not in Succeeded phase", jaegerOperatorName)

		// Debug: Check if there are any nodes that might match our constraints
		nodes, nodeErr := globalhelper.GetAPIClient().Nodes().List(context.TODO(), metav1.ListOptions{})
		if nodeErr == nil {
			glog.V(5).Infof("Found %d nodes in cluster", len(nodes.Items))
			for _, node := range nodes.Items {
				glog.V(5).Infof("Node %s has labels: %v", node.Name, node.Labels)
				// Check if any node could potentially match our constraints
				for key, value := range nodeSelector {
					if nodeValue, exists := node.Labels[key]; exists && nodeValue == value {
						glog.V(5).Infof("WARNING: Node %s matches constraint %s=%s", node.Name, key, value)
					}
				}
			}
		}

		Eventually(func() bool {
			isNotSucceeded, err := tshelper.IsCSVNotSucceeded(jaegerOperatorName, randomNamespace)
			if err != nil {
				glog.V(5).Infof("Error checking CSV status for %s: %v", jaegerOperatorName, err)

				return false
			}

			// Get detailed CSV status for debugging
			csv, csvErr := tshelper.GetCsvByPrefix(jaegerOperatorName, randomNamespace)
			if csvErr == nil {
				glog.V(5).Infof("Jaeger operator %s CSV current phase: %s", jaegerOperatorName, csv.Status.Phase)
				glog.V(5).Infof("Jaeger operator %s CSV conditions: %v", jaegerOperatorName, csv.Status.Conditions)
				if csv.Status.Message != "" {
					glog.V(5).Infof("Jaeger operator %s CSV message: %s", jaegerOperatorName, csv.Status.Message)
				}

				// Check deployment status for more insights
				if csv.Status.Phase == "Succeeded" {
					glog.V(5).Infof("WARNING: CSV succeeded despite nodeSelector constraints!")
					// Log deployment details
					deployments, depErr := globalhelper.GetAPIClient().AppsV1Interface.
						Deployments(randomNamespace).List(context.TODO(), metav1.ListOptions{})
					if depErr == nil {
						for _, dep := range deployments.Items {
							if strings.Contains(dep.Name, "jaeger") {
								glog.V(5).Infof("Jaeger deployment %s status: %+v", dep.Name, dep.Status)
								glog.V(5).Infof("Jaeger deployment %s nodeSelector: %v", dep.Name, dep.Spec.Template.Spec.NodeSelector)
							}
						}
					}
				}
			}

			glog.V(5).Infof("Jaeger operator %s CSV status is not Succeeded: %t", jaegerOperatorName, isNotSucceeded)

			return isNotSucceeded
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Equal(true),
			"Jaeger operator CSV should not be in Succeeded phase for this negative test")

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				jaegerOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+jaegerOperatorName)

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatus,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatus,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

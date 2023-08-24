package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking custom namespace,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		randomNamespace = tsparams.TestNetworkingNameSpace + "-" + globalhelper.GenerateRandomString(10)

		By(fmt.Sprintf("Create %s namespace", randomNamespace))
		err := namespaces.Create(randomNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Override default report directory")
		origReportDir = globalhelper.GetConfiguration().General.TnfReportDir
		reportDir := origReportDir + "/" + randomNamespace
		globalhelper.OverrideReportDir(reportDir)

		By("Override default TNF config directory")
		origTnfConfigDir = globalhelper.GetConfiguration().General.TnfConfigDir
		configDir := origTnfConfigDir + "/" + randomNamespace
		globalhelper.OverrideTnfConfigDir(configDir)

		By("Define TNF config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By(fmt.Sprintf("Remove %s namespace", randomNamespace))
		err := namespaces.DeleteAndWait(
			globalhelper.GetAPIClient().CoreV1Interface,
			randomNamespace,
			tsparams.WaitingTime,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Restore default report directory")
		globalhelper.GetConfiguration().General.TnfReportDir = origReportDir

		By("Restore default TNF config directory")
		globalhelper.GetConfiguration().General.TnfConfigDir = origTnfConfigDir
	})

	// 48328
	It("custom deployment 3 pods, 1 NAD, connectivity via Multus secondary interface", func() {
		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48330
	It("2 custom deployments 3 pods, 1 NAD, connectivity via Multus secondary interface", func() {
		// The NetworkAttachmentDefinition (mcvlan) created for this TC uses the default interface that is connecting
		// all worker/master nodes so that pods have connectivity irrespective of the node they are scheduled on
		// see https://github.com/test-network-function/cnfcert-tests-verification/pull/263

		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define first deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48331
	It("custom deployment and daemonset 3 pods, 2 NADs, connectivity via Multus secondary interfaces", func() {
		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, randomNamespace, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define first deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, randomNamespace, []string{tsparams.TestNadNameB}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48334
	It("custom deployment 3 pods, 1 NAD missing IP, connectivity via Multus secondary interface[skip]", func() {
		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameA, randomNamespace, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48338
	It("custom deployments 3 pods and 1 pod, standalone IP, connectivity via Multus secondary interface[skip]", func() {

		By("Define and create Network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameB, randomNamespace, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment-a and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 1)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment-b and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, randomNamespace, []string{tsparams.TestNadNameB}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48343
	It("custom deployment and daemonset 3 pods, daemonset missing ip, 2 NADs, connectivity via Multus "+
		"secondary interface", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameB, randomNamespace, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")

		err = tshelper.DefineAndCreateDaemonsetWithMultusOnCluster(tsparams.TestNadNameB, randomNamespace, "ds1")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48580
	It("custom daemonset 3 pods with skip label [skip]", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA, randomNamespace, "ds2")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48582
	It("custom deployment and daemonset 3 pods with skip label[skip]", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA, randomNamespace, "ds3")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48582
	It("custom deployment and daemonSet 3 pods, daemonSet has skip label", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA, randomNamespace, "ds4")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48582
	It("custom deployment 3 pods, 2 NADs, multiple Multus interfaces on deployment", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, randomNamespace, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA, tsparams.TestNadNameB}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 48346
	It("custom deployment 3 pods,1 NAD,no connectivity via Multus secondary interface[negative]", func() {

		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Put one deployment's pod interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48347
	It("custom deployment and daemonset 3 pods, 2 NADs, No connectivity on daemonset via Multus secondary "+
		"interface[negative]", func() {

		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, randomNamespace, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Put one deployment's pod interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusOnCluster(tsparams.TestNadNameB, randomNamespace, "ds5")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48590
	It("custom deployment and daemonset 3 pods, 2 NADs, multiple Multus interfaces on deployment no "+
		"connectivity via secondary interface[negative]", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, randomNamespace, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameB, tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Put one deployment's pod interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusOnCluster(tsparams.TestNadNameB, randomNamespace, "ds6")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfMultusIpv4TcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

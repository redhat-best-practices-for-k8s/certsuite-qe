package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking custom namespace,", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	execute.BeforeAll(func() {
		By("Clean namespace before all tests")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all nodes are labeled with 'worker-cnf' label")
		err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().CoreV1Interface, configSuite.General.CnfNodeLabel)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())
	})

	// 48328
	It("custom deployment 3 pods, 1 NAD, connectivity via Multus secondary interface", func() {
		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
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
		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define first deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, []string{tsparams.TestNadNameA}, 3)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define first deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, []string{tsparams.TestNadNameB}, 3)
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
		err := tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameA, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameB, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment-a and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 1)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment-b and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, []string{tsparams.TestNadNameB}, 3)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameB, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")

		err = tshelper.DefineAndCreateDeamonsetWithMultusOnCluster(tsparams.TestNadNameB)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDeamonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDeamonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDeamonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA, tsparams.TestNadNameB}, 3)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Put one deployment's pod  interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"})
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Put one deployment's pod interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"})
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDeamonsetWithMultusOnCluster(tsparams.TestNadNameB)
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
			tsparams.TestNadNameA, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, []string{tsparams.TestNadNameB, tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		By("Put one deployment's pod interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"})
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDeamonsetWithMultusOnCluster(tsparams.TestNadNameB)
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

package tests

import (
	"fmt"
	"os"

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

var _ = Describe("Networking network-policy-deny-all,", Serial, func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	execute.BeforeAll(func() {
		By("Clean namespace before all tests")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, tsparams.TestNetworkingNameSpace)
		Expect(err).ToNot(HaveOccurred())

		By("Create additional namespaces for testing")
		// this namespace will only be used for the networking-network-policy-deny-all tests
		err = namespaces.Create(tsparams.AdditionalNetworkingNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		By("Clean namespaces before each test")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Clean(tsparams.AdditionalNetworkingNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all nodes are labeled with 'worker-cnf' label")
		err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().CoreV1Interface, configSuite.General.CnfNodeLabel)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespaces after each test")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Clean(tsparams.AdditionalNetworkingNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())
	})

	// 59740
	It("one deployment, one pod in a namespace with deny all ingress and egress network policy", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			tsparams.TestNetworkingNameSpace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59741
	It("one deployment, one pod in a namespace with only deny all ingress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			tsparams.TestNetworkingNameSpace, []string{"Ingress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59742
	It("one deployment, one pod in a namespace with only deny all egress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1)
		Expect(err).ToNot(HaveOccurred())

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			tsparams.TestNetworkingNameSpace, []string{"Egress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59743
	It("one deployment, one pod in a namespace with neither deny all ingress or egress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59744
	It("two deployments in different namespaces, one pod each, namespaces have deny all ingress and egress network policy",
		func() {

			By("Define first deployment and create it on cluster")
			err := tshelper.DefineAndCreateDeploymentOnCluster(1)
			Expect(err).ToNot(HaveOccurred())

			By("Define and create first network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
				tsparams.TestNetworkingNameSpace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Define second deployment and create it on cluster")
			err = tshelper.DefineAndCreateDeploymentWithNamespace(tsparams.AdditionalNetworkingNamespace, 1)
			Expect(err).ToNot(HaveOccurred())

			By("Define and create second network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy2",
				tsparams.AdditionalNetworkingNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Define TNF config file")
			err = globalhelper.DefineTnfConfig(
				[]string{tsparams.TestNetworkingNameSpace, tsparams.AdditionalNetworkingNamespace},
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{})
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())
		})

	// 59745
	It("two deployments in different namespaces, one pod each, one namespace has only deny all egress network policy [negative]",
		func() {

			By("Define first deployment and create it on cluster")
			err := tshelper.DefineAndCreateDeploymentOnCluster(1)
			Expect(err).ToNot(HaveOccurred())

			By("Define and create first network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
				tsparams.TestNetworkingNameSpace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Define second deployment and create it on cluster")
			err = tshelper.DefineAndCreateDeploymentWithNamespace(tsparams.AdditionalNetworkingNamespace, 1)
			Expect(err).ToNot(HaveOccurred())

			By("Define and create second network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy2",
				tsparams.AdditionalNetworkingNamespace, []string{"Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).To(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalparameters.TestCaseFailed)
			Expect(err).ToNot(HaveOccurred())

		})
})

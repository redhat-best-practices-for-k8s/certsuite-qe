package tests

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking dual-stack-service,", func() {

	execute.BeforeAll(func() {

		By("Clean namespace before all tests")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, tsparams.TestNetworkingNameSpace)
		Expect(err).ToNot(HaveOccurred())

	})

	BeforeEach(func() {

		By("Clean namespaces before each test")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Clean(tsparams.AdditionalNetworkingNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())

	})

	// 62504
	It("service with ipFamilyPolicy SingleStack and ip version ipv6", func() {

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

	// 62505
	It("service with ipFamilyPolicy RequireDualStack and two ClusterIPs", func() {

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
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNetworkPolicyDenyAllTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 62506
	It("service with ipFamilyPolicy SingleStack and ip version ipv4 [negative]", func() {

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

	// 62507
	It("service with ipFamilyPolicy PreferDualStack and zero ClusterIPs [negative]", func() {

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

	// 62508
	It("service with no ipFamilyPolicy configured [negative]",
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
				[]string{})
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

	// 62509
	It("two services, one with ipFamilyPolicy SingleStack and ip version ipv6 and the other with ipFamilyPolicy PreferDualStack and two ClusterIPs",
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
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfNetworkPolicyDenyAllTcName,
				globalparameters.TestCasePassed)
			Expect(err).ToNot(HaveOccurred())

		})

	// 62510
	It("two services, both with ipFamilyPolicy SingleStack, one with ip version ipv6, the other with ip version ipv4 [negative]",
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

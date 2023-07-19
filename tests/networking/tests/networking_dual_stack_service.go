package tests

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking dual-stack-service,", func() {

	execute.BeforeAll(func() {
		By("Clean namespace before all tests")
		err := namespaces.Clean(tsparams.TestNetworkingNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, tsparams.TestNetworkingNameSpace)
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

	// 62506
	It("service with ipFamilyPolicy SingleStack and ip version ipv4 [negative]", func() {

		By("Define and create service")
		err := tshelper.DefineAndCreateServiceOnCluster("testservice", 3022, 3022, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDualStackServiceTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDualStackServiceTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 62507
	It("service with ipFamilyPolicy PreferDualStack and zero ClusterIPs [negative]", func() {

		By("Define and create service")
		err := tshelper.DefineAndCreateServiceOnCluster("testservice", 3022, 3022, false,
			[]corev1.IPFamily{"IPv4"}, "PreferDualStack")
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDualStackServiceTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDualStackServiceTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 62508
	It("service with no ipFamilyPolicy configured [negative]",
		func() {

			By("Define and create service")
			err := tshelper.DefineAndCreateServiceOnCluster("testservice", 3022, 3022, false, []corev1.IPFamily{"IPv4"}, "")
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				tsparams.TnfDualStackServiceTcName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
			Expect(err).To(HaveOccurred())

			By("Verify test case status in Junit and Claim reports")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TnfDualStackServiceTcName,
				globalparameters.TestCaseFailed)
			Expect(err).ToNot(HaveOccurred())
		})
})

package accesscontrol

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
)

var _ = Describe("Access control custom namespace, custom deployment,", Serial, func() {

	execute.BeforeAll(func() {
		By("Clean namespace before all tests")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, tsparams.TestAccessControlNameSpace)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())
	})

	// 45447
	It("2 custom pods, no service installed, service Should not have type of nodePort", func() {

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNodePortTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45481
	It("2 custom pods, service installed without NodePort, service Should not have type of nodePort", func() {

		By("Define Service")
		err := tshelper.DefineAndCreateServiceOnCluster("testservice", 3022, 3022, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNodePortTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45482
	It("2 custom pods, multiple services installed without NodePort, service Should not have type of nodePort", func() {

		By("Define multiple Services")
		err := tshelper.DefineAndCreateServiceOnCluster("testservicefirst", 3022, 3022, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateServiceOnCluster("testservicesecond", 3023, 3023, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNodePortTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45483
	It("2 custom pods, service installed with NodePort, service Should not have type of nodePort [negative]", func() {

		By("Define Services with NodePort")
		err := tshelper.DefineAndCreateServiceOnCluster("testservice", 30022, 3022, true,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNodePortTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45484
	It("2 custom pods, multiple services installed and one has NodePort, service Should not have type of "+
		"nodePort [negative]", func() {

		By("Define Services")
		err := tshelper.DefineAndCreateServiceOnCluster("testservicefirst", 30022, 3022, true,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())
		err = tshelper.DefineAndCreateServiceOnCluster("testservicesecond", 3022, 3022, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfNodePortTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

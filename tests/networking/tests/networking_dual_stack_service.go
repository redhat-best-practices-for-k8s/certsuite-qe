package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/networking/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
)

var _ = Describe("Networking dual-stack-service,", func() {
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

	// 62506
	It("service with ipFamilyPolicy SingleStack and ip version ipv4 [negative]", func() {

		By("Define and create service")
		err := tshelper.DefineAndCreateServiceOnCluster("testservice", randomNamespace, 3022, 3022, false, false,
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
		err := tshelper.DefineAndCreateServiceOnCluster("testservice", randomNamespace, 3023, 3023,
			false, true, []corev1.IPFamily{"IPv4"}, "PreferDualStack")
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
			err := tshelper.DefineAndCreateServiceOnCluster("testservice", randomNamespace, 3024,
				3024, false, false, []corev1.IPFamily{"IPv4"}, "")
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

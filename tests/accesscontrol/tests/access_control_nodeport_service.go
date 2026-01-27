package accesscontrol

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	corev1 "k8s.io/api/core/v1"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
)

var _ = Describe("Access control custom namespace, custom deployment,", Label("accesscontrol6"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")

		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, randomNamespace)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 45447
	It("2 custom pods, no service installed, service Should not have type of nodePort", func() {
		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert all services in namespace are not nodeport")
		services, err := globalhelper.GetServicesFromNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		for _, service := range services {
			Expect(service.Spec.Type).ToNot(Equal(corev1.ServiceTypeNodePort))
		}

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNodePortTcName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45481
	It("2 custom pods, service installed without NodePort, service Should not have type of nodePort", func() {
		By("Define Service")
		err := tshelper.DefineAndCreateServiceOnCluster("testservice", randomNamespace, 3022, 3022, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert all services in namespace are not nodeport")
		services, err := globalhelper.GetServicesFromNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		for _, service := range services {
			Expect(service.Spec.Type).ToNot(Equal(corev1.ServiceTypeNodePort))
		}

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNodePortTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 45482
	It("2 custom pods, multiple services installed without NodePort, service Should not have type of nodePort", func() {

		By("Define multiple Services")
		err := tshelper.DefineAndCreateServiceOnCluster("testservicefirst", randomNamespace, 3022, 3022, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateServiceOnCluster("testservicesecond", randomNamespace, 3023, 3023, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert all services in namespace are not nodeport")
		services, err := globalhelper.GetServicesFromNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		for _, service := range services {
			Expect(service.Spec.Type).ToNot(Equal(corev1.ServiceTypeNodePort))
		}

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNodePortTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45483
	It("2 custom pods, service installed with NodePort, service Should not have type of nodePort [negative]", func() {

		By("Define Services with NodePort")
		err := tshelper.DefineAndCreateServiceOnCluster("testservice", randomNamespace, 30022, 3022, true,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert testservice in namespace is type nodePort")
		services, err := globalhelper.GetServicesFromNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		for _, service := range services {
			if service.Name == "testservice" {
				Expect(service.Spec.Type).To(Equal(corev1.ServiceTypeNodePort))
			}
		}

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNodePortTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 45484
	It("2 custom pods, multiple services installed and one has NodePort, service Should not have type of "+
		"nodePort [negative]", func() {

		By("Define Services")
		err := tshelper.DefineAndCreateServiceOnCluster("testservicefirst", randomNamespace, 30023, 3023, true,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())
		err = tshelper.DefineAndCreateServiceOnCluster("testservicesecond", randomNamespace, 3023, 3023, false,
			[]corev1.IPFamily{"IPv4"}, "SingleStack")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment")
		dep, err := tshelper.DefineDeployment(3, 1, "acdeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert the services in namespace have corresponding nodeport statuses")
		services, err := globalhelper.GetServicesFromNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		for _, service := range services {
			switch service.Name {
			case "testservicefirst":
				Expect(service.Spec.Type).To(Equal(corev1.ServiceTypeNodePort))
			case "testservicesecond":
				Expect(service.Spec.Type).ToNot(Equal(corev1.ServiceTypeNodePort))
			}
		}

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNodePortTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNodePortTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

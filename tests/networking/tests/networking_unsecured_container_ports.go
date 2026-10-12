package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/parameters"
)

var _ = Describe("Networking unsecured-container-ports,", Serial, Label("networking3"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

	BeforeEach(func() {
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNetworkingNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	It("one deployment, plaintext HTTP listener [negative]", func() {
		By("Define and create plaintext nginx deployment")
		dep := tshelper.DefineHTTPNginxDeployment("plain-http", randomNamespace, 1)
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		checkDetails := mustGetCheckDetails(randomReportDir)
		Expect(globalhelper.CountReasonsContaining(
			checkDetails.NonCompliantObjectsOut, tsparams.PlaintextClaimSubstring)).To(BeNumerically(">=", 1),
			"expected at least one plaintext non-compliant reason")
	})

	It("one deployment, TLS listener", func() {
		By("Define and create TLS nginx deployment")
		dep, err := tshelper.DefineTLSNginxDeployment("tls-http", randomNamespace, 1)
		Expect(err).ToNot(HaveOccurred())
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		checkDetails := mustGetCheckDetails(randomReportDir)
		Expect(globalhelper.CountReasonsContaining(
			checkDetails.CompliantObjectsOut, tsparams.TLSClaimSubstring)).To(BeNumerically(">=", 1),
			"expected at least one TLS compliant reason")
	})

	It("one deployment, plaintext HTTP listener with non-tls-ports annotation", func() {
		By("Define plaintext nginx deployment with exemption annotation")
		dep := tshelper.DefineHTTPNginxDeployment("plain-http-exempt", randomNamespace, 1)
		tshelper.AnnotatePodTemplate(dep, tsparams.NonTLSPortsAnnotationKey, "8080")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		checkDetails := mustGetCheckDetails(randomReportDir)
		Expect(globalhelper.CountReasonsContaining(
			checkDetails.CompliantObjectsOut, tsparams.ExemptClaimSubstring)).To(BeNumerically(">=", 1),
			"expected annotation exemption in compliant reasons")
		Expect(globalhelper.CountReasonsContaining(
			checkDetails.NonCompliantObjectsOut, tsparams.PlaintextClaimSubstring)).To(Equal(0),
			"annotated plaintext port must not be reported as non-compliant")
	})

	It("one deployment, no listening ports", func() {
		By("Define sleep deployment with no listeners")
		dep, err := tshelper.DefineDeploymentWithContainers(1, 1, "no-ports", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		checkDetails := mustGetCheckDetails(randomReportDir)
		Expect(globalhelper.CountReasonsContaining(
			checkDetails.CompliantObjectsOut, tsparams.NoPortsClaimSubstring)).To(BeNumerically(">=", 1),
			"expected no-listening-ports compliant reason")
	})

	It("two replicas with pod anti-affinity, plaintext HTTP [negative]", func() {
		skipIfCannotSpread()

		By("Define spread plaintext nginx deployment")
		dep := tshelper.DefineHTTPNginxDeployment("plain-http-spread", randomNamespace, 2)
		deployment.RedefineWithPodAntiAffinity(dep, tsparams.TestDeploymentLabels)
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		assertSpreadAcrossNodes(randomNamespace)

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		checkDetails := mustGetCheckDetails(randomReportDir)
		Expect(globalhelper.CountReasonsContaining(
			checkDetails.NonCompliantObjectsOut, tsparams.PlaintextClaimSubstring)).To(BeNumerically(">=", 2),
			"node-local openssl must fail plaintext on each replica, not treat cross-node as unreachable")
	})

	It("two replicas with pod anti-affinity, TLS listener", func() {
		skipIfCannotSpread()

		By("Define spread TLS nginx deployment")
		dep, err := tshelper.DefineTLSNginxDeployment("tls-spread", randomNamespace, 2)
		Expect(err).ToNot(HaveOccurred())
		deployment.RedefineWithPodAntiAffinity(dep, tsparams.TestDeploymentLabels)
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		assertSpreadAcrossNodes(randomNamespace)

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteUnsecuredContainerPortsTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		checkDetails := mustGetCheckDetails(randomReportDir)
		Expect(globalhelper.CountReasonsContaining(
			checkDetails.CompliantObjectsOut, tsparams.TLSClaimSubstring)).To(BeNumerically(">=", 2),
			"node-local openssl must see TLS on each replica")
	})
})

func skipIfCannotSpread() {
	GinkgoHelper()

	if globalhelper.IsCRCCluster() {
		Skip("Pod anti-affinity spread is not supported on CRC clusters")
	}

	readyNodes, err := nodes.GetNumOfReadyNodesInCluster(globalhelper.GetAPIClient().Nodes())
	Expect(err).ToNot(HaveOccurred())

	if readyNodes < 2 {
		Skip("Cluster does not have at least 2 ready nodes for anti-affinity spread")
	}
}

func assertSpreadAcrossNodes(namespace string) {
	GinkgoHelper()

	nodeNames, err := tshelper.UniquePodNodeNames(namespace, "")
	Expect(err).ToNot(HaveOccurred())
	Expect(len(nodeNames)).To(BeNumerically(">=", 2),
		"expected replicas scheduled on at least 2 nodes, got %v", nodeNames)
}

func mustGetCheckDetails(reportDir string) *globalhelper.CheckDetails {
	GinkgoHelper()

	checkDetails, checkDetailsErr := globalhelper.GetTestCaseCheckDetails(
		tsparams.CertsuiteUnsecuredContainerPortsTcName, reportDir)
	globalhelper.LogCheckDetails(checkDetails, checkDetailsErr)
	Expect(checkDetailsErr).ToNot(HaveOccurred())

	return checkDetails
}

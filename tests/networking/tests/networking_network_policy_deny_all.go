package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/parameters"
)

var _ = Describe("Networking network-policy-deny-all,", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
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

	// 59740
	It("one deployment, one pod in a namespace with deny all ingress and egress network policy", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has pod running")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(Equal(1))

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			randomNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNetworkPolicyDenyAllTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59741
	It("one deployment, one pod in a namespace with only deny all ingress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has pod running")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(Equal(1))

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			randomNamespace, []string{"Ingress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59742
	It("one deployment, one pod in a namespace with only deny all egress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has pod running")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(Equal(1))

		By("Define and create network policy")
		err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
			randomNamespace, []string{"Egress"}, tsparams.TestDeploymentLabels)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 59743
	It("one deployment, one pod in a namespace with neither deny all ingress or egress network policy [negative]", func() {

		By("Define deployment and create it on cluster")
		err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has pod running")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(Equal(1))

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNetworkPolicyDenyAllTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNetworkPolicyDenyAllTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 59744
	It("two deployments in different namespaces, one pod each, namespaces have deny all ingress and egress network policy",
		func() {

			By("Define first deployment and create it on cluster")
			err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
			Expect(err).ToNot(HaveOccurred())

			By("Assert first deployment has pod running")
			podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(podsList.Items)).To(Equal(1))

			By("Define and create first network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
				randomNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Create additional namespaces for testing")
			randomSecondaryNamespace := tsparams.AdditionalNetworkingNamespace + "-" + globalhelper.GenerateRandomString(5)
			err = globalhelper.CreateNamespace(randomSecondaryNamespace)
			Expect(err).ToNot(HaveOccurred())

			DeferCleanup(func() {
				err = globalhelper.DeleteNamespaceAndWait(randomSecondaryNamespace, tsparams.WaitingTime)
				Expect(err).ToNot(HaveOccurred())
			})

			By("Define second deployment and create it on cluster")
			err = tshelper.DefineAndCreateDeploymentWithNamespace(randomSecondaryNamespace, 1)
			Expect(err).ToNot(HaveOccurred())

			By("Assert second deployment has pod running")
			podsList2, err := globalhelper.GetListOfPodsInNamespace(randomSecondaryNamespace)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(podsList2.Items)).To(Equal(1))

			By("Define and create second network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
				randomSecondaryNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Define certsuite config file")
			err = globalhelper.DefineCertsuiteConfig(
				[]string{randomNamespace, randomSecondaryNamespace},
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				tsparams.CertsuiteNetworkPolicyDenyAllTcName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.CertsuiteNetworkPolicyDenyAllTcName,
				globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})

	// 59745
	It("two deployments in different namespaces, one pod each, one namespace has only deny all egress network policy [negative]",
		func() {

			By("Define first deployment and create it on cluster")
			err := tshelper.DefineAndCreateDeploymentOnCluster(1, randomNamespace)
			Expect(err).ToNot(HaveOccurred())

			By("Assert first deployment has pod running")
			podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(podsList.Items)).To(Equal(1))

			By("Define and create first network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
				randomNamespace, []string{"Ingress", "Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Create additional namespaces for testing")
			randomSecondaryNamespace := tsparams.AdditionalNetworkingNamespace + "-" + globalhelper.GenerateRandomString(5)
			err = globalhelper.CreateNamespace(randomSecondaryNamespace)
			Expect(err).ToNot(HaveOccurred())

			DeferCleanup(func() {
				err = globalhelper.DeleteNamespaceAndWait(randomSecondaryNamespace, tsparams.WaitingTime)
				Expect(err).ToNot(HaveOccurred())
			})

			By("Define second deployment and create it on cluster")
			err = tshelper.DefineAndCreateDeploymentWithNamespace(randomSecondaryNamespace, 1)
			Expect(err).ToNot(HaveOccurred())

			By("Assert second deployment has pod running")
			podsList2, err := globalhelper.GetListOfPodsInNamespace(randomSecondaryNamespace)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(podsList2.Items)).To(Equal(1))

			By("Define and create second network policy")
			err = tshelper.DefineAndCreateNetworkPolicy("netpolicy1",
				randomSecondaryNamespace, []string{"Egress"}, tsparams.TestDeploymentLabels)
			Expect(err).ToNot(HaveOccurred())

			By("Define certsuite config file")
			err = globalhelper.DefineCertsuiteConfig(
				[]string{randomNamespace, randomSecondaryNamespace},
				[]string{tsparams.TestPodLabel},
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Start tests")
			err = globalhelper.LaunchTests(
				tsparams.CertsuiteNetworkPolicyDenyAllTcName,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			By("Verify test case status in Claim report")
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.CertsuiteNetworkPolicyDenyAllTcName,
				globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})
})

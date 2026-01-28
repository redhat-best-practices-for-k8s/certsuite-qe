package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"

	crdutils "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/crd"
)

var _ = Describe("access-control-crd-roles", Serial, Label("accesscontrol2", "ocp-required"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		if globalhelper.IsKindCluster() {
			By("Make masters schedulable")
			err := nodes.EnableMasterScheduling(globalhelper.GetAPIClient().K8sClient.CoreV1().Nodes(), true)
			Expect(err).ToNot(HaveOccurred())
		}

		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.CertsuiteTargetOperatorLabels},
			[]string{},
			[]string{tsparams.CertsuiteTargetCrdFilters}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")

		// We have to pre-install the cr-scale-operator resources prior to running these tests.
		By("Check if cr-scale-operator is installed")
		exists, err := globalhelper.NamespaceExists(tsparams.CertsuiteTargetOperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "error checking if cr-scale-operator is installed")
		if !exists {
			// Skip the test if cr-scaling-operator is not installed
			Skip("cr-scale-operator is not installed, skipping test")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("Custom resource is deployed, proper role defined", func() {
		By("Create a custom resource")
		_, err := crdutils.CreateCustomResourceScale(tsparams.CertsuiteCustomResourceName, randomNamespace,
			tsparams.CertsuiteTargetOperatorLabels, tsparams.CertsuiteTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Create a role for the custom resource")
		testRole := globalhelper.DefineRole("memcached-role", randomNamespace)
		globalhelper.RedefineRoleWithAPIGroups(testRole, []string{tsparams.CertsuiteCustomResourceAPIGroupName})
		globalhelper.RedefineRoleWithResources(testRole, []string{tsparams.CertsuiteCustomResourceResourceName})
		err = globalhelper.CreateRole(testRole)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete role")
			err = globalhelper.DeleteRole(testRole.Name, testRole.Namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdRoles,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdRoles, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Custom resource is deployed, one role defined with multiple api groups [negative]", func() {
		By("Create a scale custom resource")
		_, err := crdutils.CreateCustomResourceScale(tsparams.CertsuiteCustomResourceName, randomNamespace,
			tsparams.CertsuiteTargetOperatorLabels, tsparams.CertsuiteTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Create a role for the custom resource")
		testRole := globalhelper.DefineRole("memcached-role", randomNamespace)
		globalhelper.RedefineRoleWithAPIGroups(testRole, []string{tsparams.CertsuiteCustomResourceAPIGroupName, "rbac.authorization.k8s.io"})
		globalhelper.RedefineRoleWithResources(testRole, []string{tsparams.CertsuiteCustomResourceResourceName})
		err = globalhelper.CreateRole(testRole)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete role")
			err = globalhelper.DeleteRole(testRole.Name, testRole.Namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdRoles,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdRoles, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Custom resource is deployed, one role with multiple resources defined [negative]", func() {
		By("Create a scale custom resource")
		_, err := crdutils.CreateCustomResourceScale(tsparams.CertsuiteCustomResourceName, randomNamespace,
			tsparams.CertsuiteTargetOperatorLabels, tsparams.CertsuiteTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Create a role for the custom resource")
		testRole := globalhelper.DefineRole("memcached-role", randomNamespace)
		globalhelper.RedefineRoleWithAPIGroups(testRole, []string{tsparams.CertsuiteCustomResourceAPIGroupName})
		globalhelper.RedefineRoleWithResources(testRole, []string{tsparams.CertsuiteCustomResourceResourceName, "pods"})
		err = globalhelper.CreateRole(testRole)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete role")
			err = globalhelper.DeleteRole(testRole.Name, testRole.Namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdRoles,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdRoles, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Custom resource is deployed, with improper role [skip]", func() {
		By("Create a scale custom resource")
		_, err := crdutils.CreateCustomResourceScale(tsparams.CertsuiteCustomResourceName, randomNamespace,
			tsparams.CertsuiteTargetOperatorLabels, tsparams.CertsuiteTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Create a role for the custom resource")
		testRole := globalhelper.DefineRole("memcached-role", randomNamespace)
		globalhelper.RedefineRoleWithAPIGroups(testRole, []string{"bad.example.com"})
		globalhelper.RedefineRoleWithResources(testRole, []string{tsparams.CertsuiteCustomResourceResourceName})
		err = globalhelper.CreateRole(testRole)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete role")
			err = globalhelper.DeleteRole(testRole.Name, testRole.Namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteCrdRoles,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteCrdRoles, globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

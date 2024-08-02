package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"

	crdutils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/crd"
)

var _ = Describe("access-control-crd-roles", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		if globalhelper.IsKindCluster() {
			By("Make masters schedulable")
			err := nodes.EnableMasterScheduling(globalhelper.GetAPIClient().K8sClient.CoreV1().Nodes(), true)
			Expect(err).ToNot(HaveOccurred())
		}

		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			[]string{tsparams.TnfTargetCrdFilters}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		// We have to pre-install the cr-scale-operator resources prior to running these tests.
		By("Check if cr-scale-operator is installed")
		exists, err := globalhelper.NamespaceExists(tsparams.TnfTargetOperatorNamespace)
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
		_, err := crdutils.CreateCustomResourceScale(tsparams.TnfCustomResourceName, randomNamespace,
			tsparams.TnfTargetOperatorLabels, tsparams.TnfTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Create a role for the custom resource")
		testRole := globalhelper.DefineRole("memcached-role", randomNamespace)
		globalhelper.RedefineRoleWithAPIGroups(testRole, []string{tsparams.TnfCustomResourceAPIGroupName})
		globalhelper.RedefineRoleWithResources(testRole, []string{tsparams.TnfCustomResourceResourceName})
		err = globalhelper.CreateRole(testRole)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete role")
			err = globalhelper.DeleteRole(testRole.Name, testRole.Namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.TnfCrdRoles,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdRoles, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Custom resource is deployed, one role defined with multiple api groups [negative]", func() {
		By("Create a scale custom resource")
		_, err := crdutils.CreateCustomResourceScale(tsparams.TnfCustomResourceName, randomNamespace,
			tsparams.TnfTargetOperatorLabels, tsparams.TnfTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Create a role for the custom resource")
		testRole := globalhelper.DefineRole("memcached-role", randomNamespace)
		globalhelper.RedefineRoleWithAPIGroups(testRole, []string{tsparams.TnfCustomResourceAPIGroupName, "rbac.authorization.k8s.io"})
		globalhelper.RedefineRoleWithResources(testRole, []string{tsparams.TnfCustomResourceResourceName})
		err = globalhelper.CreateRole(testRole)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete role")
			err = globalhelper.DeleteRole(testRole.Name, testRole.Namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.TnfCrdRoles,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdRoles, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Custom resource is deployed, one role with multiple resources defined [negative]", func() {
		By("Create a scale custom resource")
		_, err := crdutils.CreateCustomResourceScale(tsparams.TnfCustomResourceName, randomNamespace,
			tsparams.TnfTargetOperatorLabels, tsparams.TnfTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Create a role for the custom resource")
		testRole := globalhelper.DefineRole("memcached-role", randomNamespace)
		globalhelper.RedefineRoleWithAPIGroups(testRole, []string{tsparams.TnfCustomResourceAPIGroupName})
		globalhelper.RedefineRoleWithResources(testRole, []string{tsparams.TnfCustomResourceResourceName, "pods"})
		err = globalhelper.CreateRole(testRole)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete role")
			err = globalhelper.DeleteRole(testRole.Name, testRole.Namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.TnfCrdRoles,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdRoles, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Custom resource is deployed, with improper role [skip]", func() {
		By("Create a scale custom resource")
		_, err := crdutils.CreateCustomResourceScale(tsparams.TnfCustomResourceName, randomNamespace,
			tsparams.TnfTargetOperatorLabels, tsparams.TnfTargetOperatorLabelsMap)
		Expect(err).ToNot(HaveOccurred())

		By("Create a role for the custom resource")
		testRole := globalhelper.DefineRole("memcached-role", randomNamespace)
		globalhelper.RedefineRoleWithAPIGroups(testRole, []string{"bad.example.com"})
		globalhelper.RedefineRoleWithResources(testRole, []string{tsparams.TnfCustomResourceResourceName})
		err = globalhelper.CreateRole(testRole)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			By("Delete role")
			err = globalhelper.DeleteRole(testRole.Name, testRole.Namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Start lifecycle-crd-scaling test")
		err = globalhelper.LaunchTests(tsparams.TnfCrdRoles,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdRoles, globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

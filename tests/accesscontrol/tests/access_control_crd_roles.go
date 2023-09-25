package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"

	crdutils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/crd"
)

var _ = Describe("access-control-crd-roles", Serial, func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	execute.BeforeAll(func() {
		// We have to pre-install the cr-scale-operator resources prior to running these tests.
		By("Check if cr-scale-operator is installed")
		exists, err := namespaces.Exists(tsparams.TnfTargetOperatorNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred(), "error checking if cr-scale-operator is installed")
		if !exists {
			// Skip the test if cr-scaling-operator is not installed
			Skip("cr-scale-operator is not installed, skipping test")
		}
	})

	BeforeEach(func() {
		if globalhelper.IsKindCluster() {
			By("Make masters schedulable")
			err := nodes.EnableMasterScheduling(globalhelper.GetAPIClient().CoreV1Interface, true)
			Expect(err).ToNot(HaveOccurred())
		}

		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(
			tsparams.TestAccessControlNameSpace)

		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{tsparams.TnfTargetOperatorLabels},
			[]string{},
			[]string{tsparams.TnfTargetCrdFilters})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdRoles, globalparameters.TestCasePassed)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdRoles, globalparameters.TestCaseFailed)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdRoles, globalparameters.TestCaseFailed)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfCrdRoles, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.Timeout)
	})
})

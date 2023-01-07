package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control namespace-resource-quota,", func() {

	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.TestAccessControlNameSpace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		By("Create additional namespace for testing")
		err = namespaces.Create(tsparams.AdditionalNamespaceForResourceQuotas, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		By("Clean namespaces before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		err = namespaces.Clean(tsparams.AdditionalNamespaceForResourceQuotas, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56469
	It("one deployment, one pod in a namespace with resource quota", func() {
		By("Define deployment")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Apply resource quota to namespace")
		err = tshelper.DefineAndCreateResourceQuota(tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56470
	It("one deployment, one pod in a namespace without resource quota [negative]", func() {
		By("Define deployment")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56471
	It("two deployments, one pod each, both in a namespace with resource quota", func() {
		By("Define deployments")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Apply resource quota to namespace")
		err = tshelper.DefineAndCreateResourceQuota(tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeploymentWithNamespace(1, 1, "accesscontroldeployment2",
			tsparams.AdditionalNamespaceForResourceQuotas)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Apply resource quota to namespace")
		err = tshelper.DefineAndCreateResourceQuota(tsparams.AdditionalNamespaceForResourceQuotas, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56472
	It("two deployments, one pod each, one in a namespace without resource quota [negative]", func() {
		By("Define deployments")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeploymentWithNamespace(1, 1, "accesscontroldeployment2",
			tsparams.AdditionalNamespaceForResourceQuotas)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Apply resource quota to namespace")
		err = tshelper.DefineAndCreateResourceQuota(tsparams.AdditionalNamespaceForResourceQuotas, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

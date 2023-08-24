package accesscontrol

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/accesscontrol/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Access-control security-context-privilege-escalation,", Serial, func() {

	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.TestAccessControlNameSpace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 63992
	It("one deployment, one pod, one container, does not allow privilege escalation", func() {
		By("Define deployment without privilege escalation")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63993
	It("one deployment, one pod, one container, does allow privilege escalation [negative]", func() {
		By("Define deployment with privilege escalation")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextAllowPrivilegeEscalation(dep, true)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63994
	It("two deployments, one pod each, one container each, does not allow privilege escalation", func() {
		By("Define deployments without privilege escalation")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextAllowPrivilegeEscalation(dep, false)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63995
	It("two deployments, one pod each, one container each, one does allow privilege escalation [negative]", func() {
		By("Define deployments with varying privilege escalation capabilities")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithContainersSecurityContextAllowPrivilegeEscalation(dep, true)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPrivilegeEscalation,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

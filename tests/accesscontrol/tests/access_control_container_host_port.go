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

	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Access-control container-host-port,", func() {

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

	// 63884
	It("one deployment, one pod, one container not declaring host port", func() {
		By("Define deployment with container without host port")
		dep, err := tshelper.DefineDeployment(1, 1, "acdeployment")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlContainerHostPort,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlContainerHostPort,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63885
	It("one deployment, one pod, one container declaring host port [negative]", func() {
		ports := []corev1.ContainerPort{{ContainerPort: 22223, HostPort: 22222}}

		By("Define deployment with container declaring host port")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("acdeployment", 1, ports)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlContainerHostPort,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlContainerHostPort,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63886
	It("one deployment, one pod, two containers, neither declaring host port", func() {
		ports := []corev1.ContainerPort{{ContainerPort: 22222}, {ContainerPort: 22223}}

		By("Define deployment with containers not declaring host port")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("acdeployment", 1, ports)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlContainerHostPort,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlContainerHostPort,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 63887
	It("one deployment, one pod, two containers, one declaring host port [negative]", func() {
		ports := []corev1.ContainerPort{{ContainerPort: 22222}, {ContainerPort: 22222, HostPort: 22223}}

		By("Define deployment with one container declaring host port")
		dep, err := tshelper.DefineDeploymentWithContainerPorts("acdeployment", 1, ports)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlContainerHostPort,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlContainerHostPort,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

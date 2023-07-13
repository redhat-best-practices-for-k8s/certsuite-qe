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
)

var _ = Describe("Access-control projected-volume-service-account-token,", func() {

	execute.BeforeAll(func() {
		By("Define tnf config file")
		err := globalhelper.DefineTnfConfig(
			[]string{tsparams.TestAccessControlNameSpace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred(), "error defining tnf config file")

		err = globalhelper.CreateServiceAccount(tsparams.ProjectedVolumeServiceAcctName, tsparams.TestAccessControlNameSpace)
		Expect(err).ToNot(HaveOccurred())

	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		// err := namespaces.Clean(tsparams.TestAccessControlNameSpace, globalhelper.APIClient)
		// Expect(err).ToNot(HaveOccurred())

	})

	// 56427
	It("one deployment, one pod not using a projected volume for service account access", func() {
		By("Define deployment with securityContext RunAsUser not specified")
		dep, err := tshelper.DefineDeployment(1, 1, "projvoldepPos")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithProjectedVolumeSATokenNil(dep, "lla")
		By("deployment spec is " + dep.Spec.Template.Spec.String())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfProjectedVolumeServiceAccountTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfProjectedVolumeServiceAccountTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56428
	It("one deployment, one pod using a projected volume for service account access [negative]", func() {
		By("Define deployment with projected volume service account")
		dep, err := tshelper.DefineDeployment(1, 1, "projvoldepNeg")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithProjectedVolume(dep, "la", tsparams.ProjectedVolumeServiceAcctName)
		By("deployment spec is " + dep.Spec.Template.Spec.String())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfProjectedVolumeServiceAccountTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfProjectedVolumeServiceAccountTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56429
	It("two deployments, one pod each, neither using a projected volume for service account access", func() {
		By("Define deployments with securityContext RunAsUser not specified or not 1337")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodSecurityContextRunAsUser(dep, 1338)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfProjectedVolumeServiceAccountTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfProjectedVolumeServiceAccountTcName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56430
	It("two deployments, one pod each, one using a projected volume for service account access [negative]", func() {
		By("Define deployments with varying securityContext RunAsUser values")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithPodSecurityContextRunAsUser(dep, 1337)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TnfProjectedVolumeServiceAccountTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfProjectedVolumeServiceAccountTcName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})

package accesscontrol

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
)

var _ = Describe("Access-control pod-host-path, ", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomPrivilegedNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 53939
	It("one deployment, one pod, HostPath not set", func() {
		By("Define deployment with hostPath set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateDeploymentNoWait(dep)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment does not have hostPath set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		for _, volume := range runningDeployment.Spec.Template.Spec.Volumes {
			Expect(volume.HostPath).To(BeNil())
		}

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodHostPath,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodHostPath,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53940
	It("one deployment, one pod, HostPath set [negative]", func() {
		By("Define deployment with hostPath set to true")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPath(dep, "volume", "mnt/data")

		By("Create deployment")
		err = globalhelper.CreateDeploymentNoWait(dep)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has hostPath set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		for _, volume := range runningDeployment.Spec.Template.Spec.Volumes {
			if volume.Name == "volume" {
				Expect(volume.HostPath).ToNot(BeNil())
			}
		}

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodHostPath,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodHostPath,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53941
	It("two deployments, one pod each, HostPath not set", func() {
		By("Define deployments with hostPath set to false")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 1")
		err = globalhelper.CreateDeploymentNoWait(dep)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment1 does not have hostPath set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		for _, volume := range runningDeployment.Spec.Template.Spec.Volumes {
			Expect(volume.HostPath).To(BeNil())
		}

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateDeploymentNoWait(dep2)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment2 does not have hostPath set")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		for _, volume := range runningDeployment2.Spec.Template.Spec.Volumes {
			Expect(volume.HostPath).To(BeNil())
		}

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodHostPath,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodHostPath,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 53946
	It("two deployments, one pod each, one HostPath set [negative]", func() {
		By("Define deployments with hostPath set to different values")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 1")
		err = globalhelper.CreateDeploymentNoWait(dep)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment1 does not have hostPath set")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		for _, volume := range runningDeployment.Spec.Template.Spec.Volumes {
			Expect(volume.HostPath).To(BeNil())
		}

		dep2, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithHostPath(dep2, "volume", "mnt/data")

		By("Create deployment 2")
		err = globalhelper.CreateDeploymentNoWait(dep2)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment2 has hostPath set")
		runningDeployment2, err := globalhelper.GetRunningDeployment(dep2.Namespace, dep2.Name)
		Expect(err).ToNot(HaveOccurred())
		for _, volume := range runningDeployment2.Spec.Template.Spec.Volumes {
			if volume.Name == "volume" {
				Expect(volume.HostPath).ToNot(BeNil())
			}
		}

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlPodHostPath,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlPodHostPath,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

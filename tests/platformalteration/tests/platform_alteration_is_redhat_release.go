package tests

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("platform-alteration-is-redhat-release", Label("platformalteration3"), func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.PlatformAlterationNamespace)

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

	// 51319
	It("One deployment, one pod, several containers, all running Red Hat release", func() {

		By("Define deployment")
		deployment := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		globalhelper.AppendContainersToDeployment(deployment, 3, tsparams.SampleWorkloadImage)

		By("Create and wait until deployment is ready")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(deployment.Namespace, deployment.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Assert pods are running with ready containers")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">", 0), "Expected at least one pod")

		for _, p := range podsList.Items {
			Expect(p.Status.Phase).To(Equal(corev1.PodRunning), fmt.Sprintf("Pod %s should be running", p.Name))
			for _, cs := range p.Status.ContainerStatuses {
				Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s in pod %s should be ready", cs.Name, p.Name))
			}
		}

		By("Start platform-alteration-is-redhat-release test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteIsRedHatReleaseName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteIsRedHatReleaseName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51320
	It("One daemonSet that is running Red Hat release", func() {

		By("Define daemonSet")
		daemonSet := daemonset.DefineDaemonSet(randomNamespace,
			tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)

		By("Create and wait until daemonSet is ready")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the daemonSet pods are not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		By("Assert daemonSet has ready pods on nodes")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(daemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDaemonSet.Status.NumberReady).To(BeNumerically(">", 0), "DaemonSet should have ready pods")
		Expect(runningDaemonSet.Status.NumberReady).To(Equal(runningDaemonSet.Status.DesiredNumberScheduled),
			"All scheduled pods should be ready")

		By("Assert pods are running with ready containers")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">", 0), "Expected at least one pod")

		for _, p := range podsList.Items {
			Expect(p.Status.Phase).To(Equal(corev1.PodRunning), fmt.Sprintf("Pod %s should be running", p.Name))
			for _, cs := range p.Status.ContainerStatuses {
				Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s in pod %s should be ready", cs.Name, p.Name))
			}
		}

		By("Start platform-alteration-is-redhat-release test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteIsRedHatReleaseName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteIsRedHatReleaseName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51321
	It("One deployment, one pod, 2 containers, one running Red Hat release, other is not [negative]", func() {

		By("Define deployment")
		dep := tshelper.DefineDeploymentWithNonUBIContainer(randomNamespace)

		// Append UBI-based container.
		globalhelper.AppendContainersToDeployment(dep, 1, tsparams.SampleWorkloadImage)

		By("Create and wait until deployment is ready")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment is ready")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment).ToNot(BeNil())

		By("Assert pods are running with ready containers")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">", 0), "Expected at least one pod")

		for _, p := range podsList.Items {
			Expect(p.Status.Phase).To(Equal(corev1.PodRunning), fmt.Sprintf("Pod %s should be running", p.Name))
			for _, cs := range p.Status.ContainerStatuses {
				Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s in pod %s should be ready", cs.Name, p.Name))
			}
		}

		By("Start platform-alteration-is-redhat-release test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteIsRedHatReleaseName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteIsRedHatReleaseName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51326
	It("One statefulSet, one pod that is not running Red Hat release [negative]", func() {
		By("Define statefulSet")
		statefulSet := tshelper.DefineStatefulSetWithNonUBIContainer(randomNamespace)

		By("Create and wait until statefulSet is ready")
		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}

		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulSet has ready replicas")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(statefulSet.Namespace, statefulSet.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningStatefulSet.Status.ReadyReplicas).To(BeNumerically(">", 0), "StatefulSet should have ready replicas")

		By("Assert pods are running with ready containers")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">", 0), "Expected at least one pod")

		for _, p := range podsList.Items {
			Expect(p.Status.Phase).To(Equal(corev1.PodRunning), fmt.Sprintf("Pod %s should be running", p.Name))
			for _, cs := range p.Status.ContainerStatuses {
				Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s in pod %s should be ready", cs.Name, p.Name))
			}
		}

		By("Start platform-alteration-is-redhat-release test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteIsRedHatReleaseName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteIsRedHatReleaseName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

package tests

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/statefulset"
)

var _ = Describe("platform-alteration-base-image", Label("platformalteration1", "ocp-required"), func() {
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

		if globalhelper.IsKindCluster() {
			// The Certsuite actually proactively skips this test if the cluster is Non-OCP.
			Skip(fmt.Sprintf("%s test is not applicable for Kind cluster", tsparams.CertsuiteBaseImageName))
		}

		By("Verify MCO is healthy and accessible")
		mcoHealthy, err := globalhelper.IsMCOHealthy()
		if err != nil || !mcoHealthy {
			Skip("MCO is not healthy or accessible on this cluster - skipping base image tests")
		}

		By("Verify cluster has worker nodes")
		if !globalhelper.HasWorkerNodes() {
			Skip("Cluster has no worker nodes - skipping base image tests")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 51297
	It("One deployment, one pod, running test image", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentName,
			randomNamespace,
			tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels)

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

		By("Assert pod is running and has containers")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">", 0), "Expected at least one pod")

		// Log pod and container details for debugging
		GinkgoWriter.Printf("Found %d pods in namespace %s\n", len(podsList.Items), randomNamespace)
		for i, pod := range podsList.Items {
			GinkgoWriter.Printf("Pod[%d] name: %s, phase: %s, node: %s\n",
				i, pod.Name, pod.Status.Phase, pod.Spec.NodeName)
			for j, container := range pod.Spec.Containers {
				GinkgoWriter.Printf("  Container[%d] name: %s, image: %s\n",
					j, container.Name, container.Image)
			}
		}

		Expect(podsList.Items[0].Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")
		Expect(len(podsList.Items[0].Status.ContainerStatuses)).To(BeNumerically(">", 0), "Pod should have containers")

		By("Assert all containers are ready")
		for _, cs := range podsList.Items[0].Status.ContainerStatuses {
			GinkgoWriter.Printf("Container %s: ready=%v, image=%s\n", cs.Name, cs.Ready, cs.Image)
			Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s should be ready", cs.Name))
		}

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteBaseImageName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51298
	It("One daemonSet, running test image", func() {
		By("Define daemonSet")
		testDaemonSet := daemonset.DefineDaemonSet(randomNamespace,
			tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels, tsparams.TestDaemonSetName)

		By("Create and wait until daemonSet is ready")
		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(testDaemonSet, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the daemonSet pods are not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert daemonSet is ready")
		runningDaemonSet, err := globalhelper.GetRunningDaemonset(testDaemonSet)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDaemonSet).ToNot(BeNil())

		By("Assert daemonSet has ready pods on nodes")
		GinkgoWriter.Printf("DaemonSet status: NumberReady=%d, DesiredNumberScheduled=%d, CurrentNumberScheduled=%d\n",
			runningDaemonSet.Status.NumberReady,
			runningDaemonSet.Status.DesiredNumberScheduled,
			runningDaemonSet.Status.CurrentNumberScheduled)
		Expect(runningDaemonSet.Status.NumberReady).To(BeNumerically(">", 0), "DaemonSet should have ready pods")
		Expect(runningDaemonSet.Status.NumberReady).To(Equal(runningDaemonSet.Status.DesiredNumberScheduled),
			"All scheduled pods should be ready")

		By("Assert pods are running with ready containers")
		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(podsList.Items)).To(BeNumerically(">", 0), "Expected at least one pod")

		// Log pod and container details for debugging
		GinkgoWriter.Printf("Found %d pods in namespace %s\n", len(podsList.Items), randomNamespace)
		for i, pod := range podsList.Items {
			GinkgoWriter.Printf("Pod[%d] name: %s, phase: %s, node: %s\n",
				i, pod.Name, pod.Status.Phase, pod.Spec.NodeName)
			for j, container := range pod.Spec.Containers {
				GinkgoWriter.Printf("  Container[%d] name: %s, image: %s\n",
					j, container.Name, container.Image)
			}
			Expect(pod.Status.Phase).To(Equal(corev1.PodRunning), fmt.Sprintf("Pod %s should be running", pod.Name))
			for _, cs := range pod.Status.ContainerStatuses {
				GinkgoWriter.Printf("  Container status %s: ready=%v, image=%s\n", cs.Name, cs.Ready, cs.Image)
				Expect(cs.Ready).To(BeTrue(), fmt.Sprintf("Container %s in pod %s should be ready", cs.Name, pod.Name))
			}
		}

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteBaseImageName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51299
	It("Two deployments, one pod each, change container base image by creating a file [negative]", func() {
		By("Define first deployment")
		deploymenta := deployment.DefineDeployment(tsparams.TestDeploymentName, randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		deployment.RedefineWithPrivilegedContainer(deploymenta)

		By("Create first deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Assert there is at least one pod")
		Expect(len(podsList.Items)).NotTo(BeZero())

		// Log pod details for debugging
		GinkgoWriter.Printf("Found %d pods in namespace %s\n", len(podsList.Items), randomNamespace)
		for i, pod := range podsList.Items {
			GinkgoWriter.Printf("Pod[%d] name: %s, phase: %s, node: %s\n",
				i, pod.Name, pod.Status.Phase, pod.Spec.NodeName)
			for j, container := range pod.Spec.Containers {
				GinkgoWriter.Printf("  Container[%d] name: %s, image: %s, privileged: %v\n",
					j, container.Name, container.Image,
					container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged)
			}
		}

		Expect(podsList.Items[0].Status.Phase).To(Equal(corev1.PodRunning), "First pod should be running")

		By("Change container base image")
		GinkgoWriter.Printf("Modifying base image by creating /usr/lib/testfile in pod %s\n", podsList.Items[0].Name)
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "touch /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())

		By("Verify file was created")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "ls -la /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Successfully created /usr/lib/testfile in pod %s\n", podsList.Items[0].Name)

		By("Define second deployment")
		deploymentb := deployment.DefineDeployment("platform-alteration-dpb",
			randomNamespace,
			tsparams.SampleWorkloadImage, tsparams.CertsuiteTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the second pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert second deployment is ready")
		runningDeployment2, err := globalhelper.GetRunningDeployment(deploymentb.Namespace, deploymentb.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningDeployment2).ToNot(BeNil())

		// Log all pods before running the test
		podsList, err = globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Total pods after creating both deployments: %d\n", len(podsList.Items))
		for i, pod := range podsList.Items {
			GinkgoWriter.Printf("Pod[%d] name: %s, phase: %s\n", i, pod.Name, pod.Status.Phase)
		}

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteBaseImageName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One statefulSet, one pod, change container base image by creating a file [negative]", func() {
		By("Define statefulSet")
		sts := statefulset.DefineStatefulSet(tsparams.TestStatefulSetName,
			randomNamespace,
			tsparams.SampleWorkloadImage,
			tsparams.CertsuiteTargetPodLabels)
		statefulset.RedefineWithPrivilegedContainer(sts)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(sts, tshelper.WaitingTime)
		if err != nil && strings.Contains(err.Error(), "not schedulable") {
			Skip("This test cannot run because the pod is not schedulable due to insufficient resources")
		}
		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulSet is ready")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(sts.Namespace, sts.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningStatefulSet).ToNot(BeNil())

		podsList, err := globalhelper.GetListOfPodsInNamespace(randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(podsList.Items)).NotTo(BeZero())

		// Log pod details for debugging
		GinkgoWriter.Printf("Found %d pods in namespace %s\n", len(podsList.Items), randomNamespace)
		for i, pod := range podsList.Items {
			GinkgoWriter.Printf("Pod[%d] name: %s, phase: %s, node: %s\n",
				i, pod.Name, pod.Status.Phase, pod.Spec.NodeName)
			for j, container := range pod.Spec.Containers {
				GinkgoWriter.Printf("  Container[%d] name: %s, image: %s, privileged: %v\n",
					j, container.Name, container.Image,
					container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged)
			}
		}

		Expect(podsList.Items[0].Status.Phase).To(Equal(corev1.PodRunning), "Pod should be running")

		By("Change container base image")
		GinkgoWriter.Printf("Modifying base image by creating /usr/lib/testfile in pod %s\n", podsList.Items[0].Name)
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "touch /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())

		By("Verify file was created")
		_, err = globalhelper.ExecCommand(podsList.Items[0], []string{"/bin/bash", "-c", "ls -la /usr/lib/testfile"})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Successfully created /usr/lib/testfile in pod %s\n", podsList.Items[0].Name)

		By("Start platform-alteration-base-image test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteBaseImageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteBaseImageName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

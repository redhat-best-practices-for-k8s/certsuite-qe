package accesscontrol

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
)

const (
	imageWithSSHDaemon = "quay.io/redhat-best-practices-for-k8s/certsuite-sample-workload"
)

var _ = Describe("Access-control ssh-daemons,", Label("accesscontrol12"), func() {
	var (
		randomNamespace          string
		randomReportDir          string
		randomCertsuiteConfigDir string
	)

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

	It("one pod with no ssh running", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			imageWithSSHDaemon, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is ready")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())

		By("Start ssh-daemons")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNoSSHDaemonsAllowed,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNoSSHDaemonsAllowed,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("one pod with ssh daemon running", func() {
		By("Define pod")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			imageWithSSHDaemon, tsparams.TestDeploymentLabels)

		err := pod.RedefineWithContainerExecCommand(testPod, tsparams.SSHDaemonStartContainerCommand, 0)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Assert pod is ready with ssh daemon command configured")
		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod).ToNot(BeNil())
		Expect(runningPod.Spec.Containers[0].Command).To(Equal(tsparams.SSHDaemonStartContainerCommand))

		By("Start ssh-daemons")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNoSSHDaemonsAllowed,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)

		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNoSSHDaemonsAllowed,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// Mirrors lab: only the worker-local certsuite-probe OOMKills / CrashLoopBackOffs
	// while other probe pods stay Running. The check must FAIL as a probe exec outage,
	// not as evidence that the CNF is running sshd.
	It("one pod with no ssh running, node-local probe OOMKilled", Serial, func() {
		By("Define pod without sshd")

		testPod := pod.DefinePod(tsparams.TestPodName, randomNamespace,
			imageWithSSHDaemon, tsparams.TestDeploymentLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(testPod, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
		Expect(err).ToNot(HaveOccurred())
		Expect(runningPod.Spec.NodeName).ToNot(BeEmpty())

		By("OOM-loop the certsuite-probe on the CNF node only")

		disruptCtx, stopDisrupt := context.WithCancel(context.Background())
		defer stopDisrupt()

		go tshelper.DisruptNodeLocalCertsuiteProbe(disruptCtx, runningPod.Spec.NodeName)

		By("Start ssh-daemons")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteNoSSHDaemonsAllowed,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteNoSSHDaemonsAllowed,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

		checkDetails, checkDetailsErr := globalhelper.GetTestCaseCheckDetails(
			tsparams.CertsuiteNoSSHDaemonsAllowed, randomReportDir)
		globalhelper.LogCheckDetails(checkDetails, checkDetailsErr)
		Expect(checkDetailsErr).ToNot(HaveOccurred())

		probeExecHits := globalhelper.CountReasonsContaining(checkDetails.NonCompliantObjectsOut, tsparams.ProbeExecSshdClaimSubstring)
		sshdHits := globalhelper.CountReasonsContaining(checkDetails.NonCompliantObjectsOut, tsparams.SshdRunningClaimSubstring)

		Expect(probeExecHits).To(BeNumerically(">=", 1),
			"expected probe-exec failure reason, not an sshd finding")
		Expect(sshdHits).To(Equal(0),
			"OOMKilled probe must not be reported as the pod running sshd")
	})
})

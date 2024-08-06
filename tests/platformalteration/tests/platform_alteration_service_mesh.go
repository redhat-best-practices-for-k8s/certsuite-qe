package tests

import (
	"context"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
)

const (
	WaitingTime    = 5 * time.Minute
	istioNamespace = "istio-system"
)

var _ = Describe("platform-alteration-service-mesh-usage-installed", Ordered, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeAll(func() {
		if _, exists := os.LookupEnv("NON_LINUX_ENV"); !exists {
			By("Install istio")
			//nolint:goconst
			cmd := exec.Command("/bin/bash", "-c", "curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.21.2 sh - "+
				"&& istio-1.21.2/bin/istioctl install --set profile=demo -y --set hub=gcr.io/istio-release")
			err := cmd.Run()
			Expect(err).ToNot(HaveOccurred(), "Error installing istio")
		}
	})

	AfterAll(func() {
		if _, exists := os.LookupEnv("NON_LINUX_ENV"); !exists {
			By("Uninstall istio")
			cmd := exec.Command("/bin/bash", "-c", "curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.21.2 sh - "+
				"&& istio-1.21.2/bin/istioctl uninstall -y --purge")
			err := cmd.Run()
			Expect(err).ToNot(HaveOccurred(), "Error uninstalling istio")
		}
	})

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
			Skip("Service mesh test is not applicable for Kind cluster")
		}
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// 56594
	It("istio is installed", func() {
		By("Define a test pod with istio container")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.CertsuiteTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.GetConfiguration().General.TestImage)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteServiceMeshUsageName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56596
	It("istio is installed but proxy containers does not exist [negative]", func() {
		By("Define a test pod without istio container")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.CertsuiteTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteServiceMeshUsageName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56597
	It("istio is installed but proxy container exist on one pod only [negative]", func() {
		By("Define first pod with istio container")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.CertsuiteTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.GetConfiguration().General.TestImage)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		putb := pod.DefinePod("lifecycle-putb", randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.CertsuiteTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(putb, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteServiceMeshUsageName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

var _ = Describe("platform-alteration-service-mesh-usage-uninstalled", Serial, func() {
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

	// 56595
	It("istio is not installed", func() {
		By("Check if Istio resource exists")
		gvr := schema.GroupVersionResource{Group: "install.istio.io", Version: "v1alpha1", Resource: "istiooperators"}

		_, err := globalhelper.GetAPIClient().DynamicClient.Resource(gvr).Namespace(istioNamespace).Get(context.TODO(),
			"installed-state", metav1.GetOptions{})

		if err == nil {
			By("Uninstall istio")
			if _, exists := os.LookupEnv("NON_LINUX_ENV"); !exists {
				cmd := exec.Command("/bin/bash", "-c", "curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.21.2 sh - "+
					"&& istio-1.21.2/bin/istioctl uninstall -y --purge")
				err := cmd.Run()
				Expect(err).ToNot(HaveOccurred(), "Error uninstalling istio")
			}
		}

		By("Define a test pod with istio container")
		put := pod.DefinePod(tsparams.TestPodName, randomNamespace, globalhelper.GetConfiguration().General.TestImage,
			tsparams.CertsuiteTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.GetConfiguration().General.TestImage)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteServiceMeshUsageName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

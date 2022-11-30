package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

var _ = Describe("platform-alteration-service-mesh-usage", func() {

	istioNs := "istio-system"

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56594
	It("istio is installed", func() {
		By("Create istio-system namespace")
		err := namespaces.Create(istioNs, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage, tsparams.TnfTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.Configuration.General.TestImage)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 56595
	It("istio is not installed", func() {
		By("Start platform-alteration-service-mesh-usage test")
		err := globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56596
	It("istio is installed but proxy containers does not exist [negative]", func() {
		By("Create istio-system namespace")
		err := namespaces.Create(istioNs, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 56597
	It("istio is installed but proxy container exist on one pod only [negative]", func() {
		By("Create istio-system namespace")
		err := namespaces.Create(istioNs, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Define first pod with instio container")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.Configuration.General.TestImage)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		putb := pod.DefinePod("lifecycle-putb", tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(putb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

})

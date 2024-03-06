package tests

import (
	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-affinity-required-pods", func() {
	var randomNamespace string
	var randomReportDir string
	var randomTnfConfigDir string

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatalf("can not load config file: %w", err)
	}

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define TNF config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.WaitingTime)
	})

	// 55327
	It("One pod, label is set, Affinity rules are set", func() {
		By("Define and create pod")
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		globalhelper.AppendLabelsToPod(put, tsparams.TestTargetLabels)
		globalhelper.AppendLabelsToPod(put, tsparams.AffinityRequiredPodLabels)
		pod.RedefineWithNodeAffinity(put, configSuite.General.CnfNodeLabel)
		pod.RedefineWithPodAffinity(put, tsparams.TestTargetLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-affinity-required-pods test")
		err = globalhelper.LaunchTests(tsparams.TnfAffinityRequiredPodsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55328
	It("Two pods, labels are set for both, Affinity rules are set", func() {
		By("Define and create pods")
		putA := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		globalhelper.AppendLabelsToPod(putA, tsparams.TestTargetLabels)
		globalhelper.AppendLabelsToPod(putA, tsparams.AffinityRequiredPodLabels)
		pod.RedefineWithNodeAffinity(putA, configSuite.General.CnfNodeLabel)
		pod.RedefineWithPodAffinity(putA, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(putA, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		putB := tshelper.DefinePod("lifecycle-podb", randomNamespace)
		globalhelper.AppendLabelsToPod(putB, tsparams.TestTargetLabels)
		globalhelper.AppendLabelsToPod(putB, tsparams.AffinityRequiredPodLabels)
		pod.RedefineWithNodeAffinity(putB, configSuite.General.CnfNodeLabel)
		pod.RedefineWithPodAffinity(putB, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(putB, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-affinity-required-pods test")
		err = globalhelper.LaunchTests(tsparams.TnfAffinityRequiredPodsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCasePassed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55329
	It("One pod, label is set, affinity rules are not set [negative]", func() {
		By("Define and create pod")
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		globalhelper.AppendLabelsToPod(put, tsparams.TestTargetLabels)
		globalhelper.AppendLabelsToPod(put, tsparams.AffinityRequiredPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-affinity-required-pods test")
		err = globalhelper.LaunchTests(tsparams.TnfAffinityRequiredPodsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55330
	It("One pod, label is set, podantiaffinity is set [negative]", func() {
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		globalhelper.AppendLabelsToPod(put, tsparams.TestTargetLabels)
		globalhelper.AppendLabelsToPod(put, tsparams.AffinityRequiredPodLabels)
		pod.RedefineWithPodAntiAffinity(put, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-affinity-required-pods test")
		err = globalhelper.LaunchTests(
			tsparams.TnfAffinityRequiredPodsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfAffinityRequiredPodsTcName,
			globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55333
	It("Two pods, labels are set for both, affinity rules are not set for one of the pods [negative]", func() {
		By("Define and create pods")
		putA := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		globalhelper.AppendLabelsToPod(putA, tsparams.TestTargetLabels)
		globalhelper.AppendLabelsToPod(putA, tsparams.AffinityRequiredPodLabels)
		pod.RedefineWithNodeAffinity(putA, configSuite.General.CnfNodeLabel)
		pod.RedefineWithPodAffinity(putA, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(putA, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		putB := tshelper.DefinePod("lifecycle-podb", randomNamespace)
		globalhelper.AppendLabelsToPod(putB, tsparams.TestTargetLabels)
		globalhelper.AppendLabelsToPod(putB, tsparams.AffinityRequiredPodLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(putB, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-affinity-required-pods test")
		err = globalhelper.LaunchTests(tsparams.TnfAffinityRequiredPodsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCaseFailed,
			randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

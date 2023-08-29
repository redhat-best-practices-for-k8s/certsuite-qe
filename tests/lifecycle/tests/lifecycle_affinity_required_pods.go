package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-affinity-required-pods", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.LifecycleNamespace)

		By("Define TNF config file")
		err = globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	AfterEach(func() {
		By(fmt.Sprintf("Remove %s namespace", randomNamespace))
		err := namespaces.DeleteAndWait(
			globalhelper.GetAPIClient().CoreV1Interface,
			randomNamespace,
			tsparams.WaitingTime,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Restore default report directory")
		globalhelper.GetConfiguration().General.TnfReportDir = origReportDir

		By("Restore default TNF config directory")
		globalhelper.GetConfiguration().General.TnfConfigDir = origTnfConfigDir
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCasePassed)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCasePassed)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 55330
	It("One pod, label is set, podantiaffinity is set [negative]", func() {
		put := tshelper.DefinePod(tsparams.TestPodName, randomNamespace)
		globalhelper.AppendLabelsToPod(put, tsparams.TestTargetLabels)
		globalhelper.AppendLabelsToPod(put, tsparams.AffinityRequiredPodLabels)
		pod.RedefineWithPodantiAffinity(put, tsparams.TestTargetLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-affinity-required-pods test")
		err = globalhelper.LaunchTests(tsparams.TnfAffinityRequiredPodsTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCaseFailed)
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
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfAffinityRequiredPodsTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

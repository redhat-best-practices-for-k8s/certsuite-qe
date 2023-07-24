package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-liveness", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file"))
	}

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all nodes are labeled with 'worker-cnf' label")
		err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().CoreV1Interface, configSuite.General.CnfNodeLabel)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 50053
	It("One deployment, one pod with a liveness probe", func() {
		By("Define deployment with a liveness probe")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithLivenessProbe(deploymenta)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(tsparams.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfLivenessTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50054
	It("Two deployments, multiple pods each, all have a liveness probe", func() {
		By("Define first deployment with a liveness probe")
		deploymenta, err := tshelper.DefineDeployment(3, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithLivenessProbe(deploymenta)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment with a liveness probe")
		deploymentb, err := tshelper.DefineDeployment(3, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithLivenessProbe(deploymentb)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(tsparams.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfLivenessTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50055
	It("One statefulSet, one pod with a liveness probe", func() {
		By("Define statefulSet with a liveness probe")
		statefulSet := tshelper.DefineStatefulSet(tsparams.TestStatefulSetName)
		statefulset.RedefineWithLivenessProbe(statefulSet)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(tsparams.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfLivenessTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50056
	It("One pod with a liveness probe", func() {
		By("Define pod with a liveness probe")
		put := tshelper.DefinePod(tsparams.TestPodName)
		pod.RedefineWithLivenessProbe(put)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(tsparams.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfLivenessTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50057
	It("One daemonSet without a liveness probe [negative]", func() {
		By("Define daemonSet without a liveness probe")
		daemonSet := daemonset.DefineDaemonSet(tsparams.LifecycleNamespace,
			globalhelper.GetConfiguration().General.TestImage,
			tsparams.TestTargetLabels, tsparams.TestDaemonSetName)

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(tsparams.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfLivenessTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 50058
	It("Two deployments, one pod each, one without a liveness probe [negative]", func() {
		By("Define first deployment with a liveness probe")
		deploymenta, err := tshelper.DefineDeployment(1, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		deployment.RedefineWithLivenessProbe(deploymenta)

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define second deployment without a liveness probe")
		deploymentb, err := tshelper.DefineDeployment(1, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-liveness test")
		err = globalhelper.LaunchTests(tsparams.TnfLivenessTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfLivenessTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})
